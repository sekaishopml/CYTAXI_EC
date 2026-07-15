# Telemetry

Distributed tracing for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/telemetry"

ctx, span := telemetry.WithSpan(r.Context(), "processRequest")
defer span.End()

span.SetAttribute("user_id", userID)
// ... do work ...
```

## How it works

- Context-based span propagation.
- `WithSpan` creates a span and stores it in context.
- `SpanFromContext` retrieves the current span for child operations.
- Designed to be replaced with OpenTelemetry in production.
