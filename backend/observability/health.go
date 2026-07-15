package observability

import "sync"

type healthChecker struct {
	mu       sync.RWMutex
	checks   map[string]HealthCheckFunc
}

type HealthCheckFunc func() HealthStatus

func NewHealthChecker() HealthChecker {
	return &healthChecker{
		checks: make(map[string]HealthCheckFunc),
	}
}

func (h *healthChecker) Register(name string, fn HealthCheckFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = fn
}

func (h *healthChecker) Health() map[string]HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	results := make(map[string]HealthStatus, len(h.checks))
	for name, fn := range h.checks {
		results[name] = fn()
	}
	return results
}
