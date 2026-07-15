package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type DecisionKind string

const (
	DecisionLLM           DecisionKind = "llm"
	DecisionDeterministic DecisionKind = "deterministic"
	DecisionFallback      DecisionKind = "fallback"
	DecisionReject        DecisionKind = "reject"
)

type OrchestratorResult struct {
	Decision   DecisionKind
	Response   string
	Confidence float64
	Provider   string
	LatencyMs  int64
	Reason     string
}

type Orchestrator struct {
	router    *Router
	policy    *PolicyEvaluator
	fallback  *FallbackManager
	metrics   *MetricsCollector
	providers ProviderRegistry
}

func NewOrchestrator(
	router *Router,
	policy *PolicyEvaluator,
	fallback *FallbackManager,
	metrics *MetricsCollector,
	providers ProviderRegistry,
) *Orchestrator {
	return &Orchestrator{
		router:    router,
		policy:    policy,
		fallback:  fallback,
		metrics:   metrics,
		providers: providers,
	}
}

func (o *Orchestrator) Process(ctx context.Context, msg *entity.Message, session *entity.Session, convCtx *entity.ConversationContext) *OrchestratorResult {
	start := time.Now()

	intent := o.router.Route(ctx, msg.Content, session)

	if !o.policy.ShouldUseLLM(intent, session) {
		result := &OrchestratorResult{
			Decision:  DecisionDeterministic,
			Response:  "",
			Reason:    "policy: intent handled deterministically",
		}
		o.metrics.Record(result)
		return result
	}

	provider, err := o.providers.Select(intent.RequiredCapabilities...)
	if err != nil {
		fallbackResult := o.fallback.Handle(ctx, intent, err)
		o.metrics.Record(fallbackResult)
		return fallbackResult
	}

	llmCtx := BuildContext(session, convCtx)
	prompt := BuildPrompt(intent, llmCtx)

	resp, err := provider.Complete(ctx, prompt)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		o.metrics.RecordProviderError(provider.Name())
		fallbackResult := o.fallback.Handle(ctx, intent, err)
		fallbackResult.LatencyMs = latency
		o.metrics.Record(fallbackResult)
		return fallbackResult
	}

	result := &OrchestratorResult{
		Decision:   DecisionLLM,
		Response:   resp.Content,
		Confidence: resp.Confidence,
		Provider:   provider.Name(),
		LatencyMs:  latency,
		Reason:     "llm: completed successfully",
	}
	o.metrics.Record(result)
	return result
}

func (o *Orchestrator) ProcessDeterministic(ctx context.Context, msg *entity.Message, session *entity.Session) *OrchestratorResult {
	return &OrchestratorResult{
		Decision: DecisionDeterministic,
		Reason:   "bypassed: deterministic path requested",
	}
}
