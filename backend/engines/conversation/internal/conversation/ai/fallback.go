package ai

import (
	"context"
	"fmt"
	"time"
)

type FallbackStrategy int

const (
	FallbackRetry    FallbackStrategy = iota
	FallbackFallbackProvider
	FallbackStaticResponse
)

type FallbackRule struct {
	MaxRetries  int
	Strategy    FallbackStrategy
	StaticMsg   string
	FallbackTo  ProviderKind
}

type FallbackManager struct {
	rules []FallbackRule
}

func NewFallbackManager(rules []FallbackRule) *FallbackManager {
	return &FallbackManager{rules: rules}
}

func NewDefaultFallback() *FallbackManager {
	return NewFallbackManager([]FallbackRule{
		{
			Strategy:   FallbackRetry,
			MaxRetries: 2,
		},
		{
			Strategy:   FallbackStaticResponse,
			StaticMsg:  "Lo siento, no puedo procesar tu solicitud en este momento. Por favor intenta más tarde.",
		},
	})
}

func (m *FallbackManager) Handle(ctx context.Context, intent *Intent, err error) *OrchestratorResult {
	for _, rule := range m.rules {
		switch rule.Strategy {
		case FallbackRetry:
			for i := 0; i < rule.MaxRetries; i++ {
				time.Sleep(500 * time.Millisecond)
			}
			return &OrchestratorResult{
				Decision: DecisionFallback,
				Response: rule.StaticMsg,
				Reason:   fmt.Sprintf("fallback: retries exhausted, error: %v", err),
			}

		case FallbackStaticResponse:
			return &OrchestratorResult{
				Decision: DecisionFallback,
				Response: rule.StaticMsg,
				Reason:   fmt.Sprintf("fallback: static response, error: %v", err),
			}

		case FallbackFallbackProvider:
			return &OrchestratorResult{
				Decision: DecisionFallback,
				Provider: string(rule.FallbackTo),
				Reason:   fmt.Sprintf("fallback: switching to %s, error: %v", rule.FallbackTo, err),
			}
		}
	}

	return &OrchestratorResult{
		Decision: DecisionReject,
		Reason:   fmt.Sprintf("reject: all fallbacks exhausted, error: %v", err),
	}
}
