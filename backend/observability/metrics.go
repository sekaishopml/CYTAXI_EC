package observability

import (
	"sync"
	"sync/atomic"
)

type metricType int

const (
	metricCounter metricType = iota
	metricGauge
	metricHistogram
)

type metric struct {
	kind  metricType
	value atomic.Int64
}

type metricsStore struct {
	mu      sync.RWMutex
	metrics map[string]*metric
}

var store = &metricsStore{
	metrics: make(map[string]*metric),
}

type basicMetrics struct{}

func NewMetrics() Metrics {
	return &basicMetrics{}
}

func (m *basicMetrics) Counter(name string, tags map[string]string) {
	key := metricKey(name, tags)
	store.mu.Lock()
	defer store.mu.Unlock()

	met, ok := store.metrics[key]
	if !ok {
		met = &metric{kind: metricCounter}
		store.metrics[key] = met
	}
	met.value.Add(1)
}

func (m *basicMetrics) Gauge(name string, value float64, tags map[string]string) {
	key := metricKey(name, tags)
	store.mu.Lock()
	defer store.mu.Unlock()

	met, ok := store.metrics[key]
	if !ok {
		met = &metric{kind: metricGauge}
		store.metrics[key] = met
	}
	met.value.Store(int64(value * 100))
}

func (m *basicMetrics) Histogram(name string, value float64, tags map[string]string) {
	key := metricKey(name, tags)
	store.mu.Lock()
	defer store.mu.Unlock()

	met, ok := store.metrics[key]
	if !ok {
		met = &metric{kind: metricHistogram}
		store.metrics[key] = met
	}
	met.value.Add(int64(value))
}

func metricKey(name string, tags map[string]string) string {
	key := name
	for k, v := range tags {
		key += "|" + k + "=" + v
	}
	return key
}
