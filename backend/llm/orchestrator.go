package llm

import (
	"context"
	"encoding/json"
	"time"
)

type Orchestrator struct {
	provider   Provider
	security   *SecurityFilter
	metrics    *MetricsCollector
}

func NewOrchestrator(provider Provider, security *SecurityFilter, metrics *MetricsCollector) *Orchestrator {
	return &Orchestrator{
		provider: provider,
		security: security,
		metrics:  metrics,
	}
}

type OrchestratorResult struct {
	Response     *LLMResponse
	RawJSON      string
	Confidence   float64
	LatencyMs    int64
	ProviderUsed string
	Error        error
	WasFiltered  bool
}

func (o *Orchestrator) Process(ctx context.Context, rawInput string, contextData map[string]string) *OrchestratorResult {
	start := time.Now()

	sanitized, err := o.security.ValidateAndSanitize(rawInput)
	if err != nil {
		o.metrics.RecordRejected(err)
		return &OrchestratorResult{
			Confidence:   0,
			Error:        err,
			WasFiltered:  true,
			LatencyMs:    time.Since(start).Milliseconds(),
			ProviderUsed: o.provider.Name(),
		}
	}

	req := CompletionRequest{
		RawInput:  sanitized,
		Context:   contextData,
		MaxTokens: o.security.maxTokens,
	}

	resp, err := o.provider.Complete(ctx, req)
	if err != nil {
		o.metrics.RecordProviderError(o.provider.Name())
		return &OrchestratorResult{
			Confidence:   0,
			Error:        err,
			LatencyMs:    time.Since(start).Milliseconds(),
			ProviderUsed: o.provider.Name(),
		}
	}

	var llmResp LLMResponse
	if err := json.Unmarshal([]byte(resp.Content), &llmResp); err != nil {
		o.metrics.RecordParseError()
		return &OrchestratorResult{
			Confidence:   0.3,
			Error:        err,
			RawJSON:      resp.Content,
			LatencyMs:    time.Since(start).Milliseconds(),
			ProviderUsed: o.provider.Name(),
		}
	}

	llmResp.RawInput = rawInput
	latency := time.Since(start).Milliseconds()

	o.metrics.RecordSuccess(llmResp.Intent.Kind, llmResp.Confidence, latency)

	return &OrchestratorResult{
		Response:     &llmResp,
		RawJSON:      resp.Content,
		Confidence:   llmResp.Confidence,
		LatencyMs:    latency,
		ProviderUsed: o.provider.Name(),
	}
}
