# HTTP

HTTP server and response helpers for CYTAXI backend.

## Usage

```go
import cytaxihttp "github.com/sekaishopml/cytaxi/backend/http"

mux := http.NewServeMux()
srv := cytaxihttp.New(mux, cytaxihttp.DefaultConfig())

go srv.Start()
// ... wait for signal ...
srv.Shutdown()
```

## Response helpers

```go
cytaxihttp.OK(w, data)           // 200
cytaxihttp.Created(w, data)      // 201
cytaxihttp.NoContent(w)          // 204
cytaxihttp.WriteError(w, err)     // Maps typed errors to HTTP
cytaxihttp.WriteValidationError(w, fieldErrors) // 400
```
