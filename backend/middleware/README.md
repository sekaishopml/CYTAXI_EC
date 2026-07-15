# Middleware

Reusable HTTP middleware for CYTAXI backend.

## Available middleware

| Middleware | Description |
|-----------|-------------|
| `CorrelationID` | Adds/reads `X-Correlation-ID` header, propagates via context |
| `RequestID` | Adds/reads `X-Request-ID` header, propagates via context |
| `Logging` | Logs method, path, status, duration, correlation_id, request_id |
| `Recovery` | Recovers from panics, logs stack trace, returns 500 |
| `CORS` | Configurable CORS with `Access-Control-*` headers |

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/middleware"

mux := http.NewServeMux()
handler := middleware.Chain(mux,
    middleware.Recovery(logger),
    middleware.CorrelationID,
    middleware.RequestID,
    middleware.Logging(logger),
    middleware.CORS(middleware.DefaultCORSConfig()),
)
```
