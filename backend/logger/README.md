# Logger

Structured logging for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/logger"

log := logger.NewSlog("debug", "json")
log.Info("service started", "port", 8080)
log.Error("something failed", "error", err)

// With fields
log = log.With("service", "conversation")
log.Info("request processed")
```

## How it works

- Implements `logger.Logger` interface using Go's `log/slog`.
- Supports JSON and text output formats.
- Includes context helpers: `ToContext`, `FromContext` for request-scoped logging.
