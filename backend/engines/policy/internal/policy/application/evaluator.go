package application

import (
	"context"
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/domain"
)

type RuleEvaluator struct {
	conditionEvaluator *ConditionEvaluator
}

func NewRuleEvaluator(conditionEvaluator *ConditionEvaluator) *RuleEvaluator {
	return &RuleEvaluator{
		conditionEvaluator: conditionEvaluator,
	}
}

func (e *RuleEvaluator) EvaluatePolicy(ctx context.Context, policy domain.Policy, decisionCtx domain.DecisionContext) (*domain.EvaluationResult, error) {
	result := domain.NewEvaluationResult()

	for _, rule := range policy.Rules {
		if !rule.Enabled {
			continue
		}

		matched, err := e.evaluateRule(ctx, rule, decisionCtx)
		if err != nil {
			result.AddError(fmt.Errorf("rule %s: %w", rule.ID, err))
			continue
		}

		decision := domain.Decision{
			PolicyID:    policy.ID,
			RuleID:      rule.ID,
			Actions:     rule.Actions,
			Matched:     matched,
			EvaluatedAt: time.Now(),
		}
		result.AddDecision(decision)

		if matched {
			break
		}
	}

	return result, nil
}

func (e *RuleEvaluator) evaluateRule(ctx context.Context, rule domain.Rule, decisionCtx domain.DecisionContext) (bool, error) {
	if len(rule.Conditions) == 0 {
		return true, nil
	}

	for _, condition := range rule.Conditions {
		matches, err := e.conditionEvaluator.Evaluate(ctx, condition, decisionCtx)
		if err != nil {
			return false, fmt.Errorf("evaluate condition: %w", err)
		}
		if !matches {
			return false, nil
		}
	}

	return true, nil
}
