package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type PolicyEngine struct {
	registry  *PolicyRegistry
	evaluator *RuleEvaluator
	logger    *slog.Logger
}

func NewPolicyEngine(registry *PolicyRegistry, evaluator *RuleEvaluator, logger *slog.Logger) *PolicyEngine {
	return &PolicyEngine{
		registry:  registry,
		evaluator: evaluator,
		logger:    logger,
	}
}

func (e *PolicyEngine) Evaluate(ctx context.Context, decisionCtx domain.DecisionContext, domains ...string) (*domain.EvaluationResult, error) {
	start := time.Now()
	result := domain.NewEvaluationResult()

	policies := e.registry.GetPolicies(domains...)
	if len(policies) == 0 {
		e.logger.Warn("no policies found for domains", "domains", domains)
		return result, nil
	}

	for _, policy := range policies {
		if !policy.Enabled {
			continue
		}

		policyResult, err := e.evaluator.EvaluatePolicy(ctx, policy, decisionCtx)
		if err != nil {
			e.logger.Error("policy evaluation failed",
				"policy", policy.ID,
				"error", err,
			)
			result.AddError(fmt.Errorf("policy %s: %w", policy.ID, err))
			continue
		}

		for _, d := range policyResult.Decisions {
			result.AddDecision(d)
		}
	}

	e.logger.Info("policy evaluation completed",
		"domains", domains,
		"policies", len(policies),
		"matched", result.Matched,
		"duration", time.Since(start).String(),
	)

	return result, nil
}

func (e *PolicyEngine) EvaluatePolicy(ctx context.Context, policyID domain.PolicyID, decisionCtx domain.DecisionContext) (*domain.EvaluationResult, error) {
	policy, err := e.registry.GetPolicy(policyID)
	if err != nil {
		return nil, fmt.Errorf("policy %s: %w", policyID, err)
	}

	if !policy.Enabled {
		return nil, fmt.Errorf("policy %s: %w", policyID, domain.ErrPolicyDisabled)
	}

	return e.evaluator.EvaluatePolicy(ctx, *policy, decisionCtx)
}

func (e *PolicyEngine) ReloadPolicies(ctx context.Context) error {
	return e.registry.Reload(ctx)
}
