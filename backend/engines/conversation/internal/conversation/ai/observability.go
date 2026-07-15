package ai

import (
	"sync"
	"sync/atomic"
	"time"
)

type MetricsCollector struct {
	mu              sync.RWMutex
	totalRequests   atomic.Int64
	llmRequests     atomic.Int64
	deterministic   atomic.Int64
	fallbackCount   atomic.Int64
	rejectCount     atomic.Int64
	providerErrors  map[string]int64
	latencyBuckets  map[string][]float64
	decisions       map[DecisionKind]int64
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		providerErrors: make(map[string]int64),
		latencyBuckets: make(map[string][]float64),
		decisions:      make(map[DecisionKind]int64),
	}
}

func (m *MetricsCollector) Record(result *OrchestratorResult) {
	m.totalRequests.Add(1)

	m.mu.Lock()
	m.decisions[result.Decision]++
	m.mu.Unlock()

	switch result.Decision {
	case DecisionLLM:
		m.llmRequests.Add(1)
	case DecisionDeterministic:
		m.deterministic.Add(1)
	case DecisionFallback:
		m.fallbackCount.Add(1)
	case DecisionReject:
		m.rejectCount.Add(1)
	}
}

func (m *MetricsCollector) RecordProviderError(provider string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providerErrors[provider]++
}

type MetricSnapshot struct {
	TotalRequests   int64
	LLMRequests     int64
	Deterministic   int64
	FallbackCount   int64
	RejectCount     int64
	ProviderErrors  map[string]int64
	Decisions       map[DecisionKind]int64
	Timestamp       time.Time
}

func (m *MetricsCollector) Snapshot() MetricSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providerErrors := make(map[string]int64)
	for k, v := range m.providerErrors {
		providerErrors[k] = v
	}

	decisions := make(map[DecisionKind]int64)
	for k, v := range m.decisions {
		decisions[k] = v
	}

	return MetricSnapshot{
		TotalRequests:  m.totalRequests.Load(),
		LLMRequests:    m.llmRequests.Load(),
		Deterministic:  m.deterministic.Load(),
		FallbackCount:  m.fallbackCount.Load(),
		RejectCount:    m.rejectCount.Load(),
		ProviderErrors: providerErrors,
		Decisions:      decisions,
		Timestamp:      time.Now(),
	}
}

type AILogger interface {
	LogDecision(ctx Context, intent *Intent, result *OrchestratorResult)
}
