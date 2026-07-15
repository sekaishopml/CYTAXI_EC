# Observability

Metrics and health checks for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/observability"

metrics := observability.NewMetrics()
metrics.Counter("http.requests", map[string]string{"method": "GET"})
metrics.Gauge("active.connections", 5, nil)

health := observability.NewHealthChecker()
health.Register("database", func() observability.HealthStatus {
    return observability.HealthStatus{Status: "up"}
})
```

## Health checks

- Register component checks via `health.Register(name, fn)`.
- `health.Health()` returns a map of component → status for `/health` endpoints.
