# Containers

Simple dependency injection container for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/containers"

c := containers.New()
c.Register("logger", myLogger)
c.Register("db", myDB)

logger := c.Resolve("logger").(logger.Logger)
```

Designed for simple wiring. Replace with Google Wire or Uber Fx if complexity grows.
