package llm

import (
	"sync"
	"sync/atomic"
	"time"
)

type MetricsCollector struct {
	mu              sync.RWMutex
	totalRequests   atomic.Int64
	successCount    atomic.Int64
	rejectedCount   atomic.Int64
	parseErrors     atomic.Int64
	providerErrors  map[string]int64
	intentCounts    map[IntentKind]int64
	latencies       []float64
	confidences     []float64
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		providerErrors: make(map[string]int64),
		intentCounts:   make(map[IntentKind]int64),
	}
}

func (m *MetricsCollector) RecordSuccess(intent IntentKind, confidence float64, latencyMs int64) {
	m.totalRequests.Add(1)
	m.successCount.Add(1)

	m.mu.Lock()
	m.intentCounts[intent]++
	m.latencies = append(m.latencies, float64(latencyMs))
	if len(m.latencies) > 1000 {
		m.latencies = m.latencies[len(m.latencies)-1000:]
	}
	m.confidences = append(m.confidences, confidence)
	if len(m.confidences) > 1000 {
		m.confidences = m.confidences[len(m.confidences)-1000:]
	}
	m.mu.Unlock()
}

func (m *MetricsCollector) RecordRejected(err error) {
	m.totalRequests.Add(1)
	m.rejectedCount.Add(1)
}

func (m *MetricsCollector) RecordParseError() {
	m.totalRequests.Add(1)
	m.parseErrors.Add(1)
}

func (m *MetricsCollector) RecordProviderError(provider string) {
	m.mu.Lock()
	m.providerErrors[provider]++
	m.mu.Unlock()
}

type MetricSnapshot struct {
	TotalRequests   int64               `json:"total_requests"`
	SuccessCount    int64               `json:"success_count"`
	RejectedCount   int64               `json:"rejected_count"`
	ParseErrors     int64               `json:"parse_errors"`
	SuccessRate     float64             `json:"success_rate"`
	AvgLatencyMs    float64             `json:"avg_latency_ms"`
	AvgConfidence   float64             `json:"avg_confidence"`
	ProviderErrors  map[string]int64    `json:"provider_errors"`
	IntentCounts    map[IntentKind]int64 `json:"intent_counts"`
	Timestamp       time.Time           `json:"timestamp"`
}

func (m *MetricsCollector) Snapshot() MetricSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := m.totalRequests.Load()
	providerErrors := make(map[string]int64)
	for k, v := range m.providerErrors {
		providerErrors[k] = v
	}
	intentCounts := make(map[IntentKind]int64)
	for k, v := range m.intentCounts {
		intentCounts[k] = v
	}

	avgLat := 0.0
	if len(m.latencies) > 0 {
		sum := 0.0
		for _, v := range m.latencies {
			sum += v
		}
		avgLat = sum / float64(len(m.latencies))
	}

	avgConf := 0.0
	if len(m.confidences) > 0 {
		sum := 0.0
		for _, v := range m.confidences {
			sum += v
		}
		avgConf = sum / float64(len(m.confidences))
	}

	successRate := 0.0
	if total > 0 {
		successRate = float64(m.successCount.Load()) / float64(total) * 100
	}

	return MetricSnapshot{
		TotalRequests:  total,
		SuccessCount:   m.successCount.Load(),
		RejectedCount:  m.rejectedCount.Load(),
		ParseErrors:    m.parseErrors.Load(),
		SuccessRate:    successRate,
		AvgLatencyMs:   avgLat,
		AvgConfidence:  avgConf,
		ProviderErrors: providerErrors,
		IntentCounts:   intentCounts,
		Timestamp:      time.Now(),
	}
}
