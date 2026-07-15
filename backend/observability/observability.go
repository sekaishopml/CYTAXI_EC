package observability

type Metrics interface {
	Counter(name string, tags map[string]string)
	Gauge(name string, value float64, tags map[string]string)
	Histogram(name string, value float64, tags map[string]string)
}

type HealthChecker interface {
	Health() map[string]HealthStatus
}

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
