# Errors

Typed error handling for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/errors"

func FindDriver(id string) (*Driver, error) {
    if id == "" {
        return nil, errors.New(errors.KindValidation, "FindDriver", "id is required")
    }
    // ...
    return nil, errors.Wrap(errors.KindNotFound, "FindDriver", "driver not found", err)
}
```

## Error kinds

| Kind | HTTP Status | Use case |
|------|-------------|----------|
| `KindInternal` | 500 | Unexpected server errors |
| `KindValidation` | 400 | Invalid input |
| `KindNotFound` | 404 | Resource not found |
| `KindUnauthorized` | 401 | Missing/invalid auth |
| `KindForbidden` | 403 | Insufficient permissions |
| `KindConflict` | 409 | Duplicate or conflict |
| `KindTimeout` | 504 | Operation timed out |
| `KindExternal` | 502 | External service error |

## HTTP helpers

```go
status, body := errors.EncodeHTTP(err)
// Returns (500, {"error": "..."})
```
