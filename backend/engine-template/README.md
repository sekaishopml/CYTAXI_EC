# Engine Template

Standard CYTAXI Engine structure following Clean Architecture + DDD.

## Structure

```
cmd/                    # Entrypoint
domain/                 # Business logic (zero dependencies)
application/            # Use cases, ports, DTOs
infrastructure/         # DB, cache, broker, external APIs
api/                    # HTTP handlers, routes, middleware
events/                 # Domain event definitions and handlers
config/                 # Engine configuration
tests/                  # Integration tests
```

## Dependency Rules

- `domain` → nothing
- `application` → `domain`, `foundation`
- `infrastructure` → `domain`, `foundation`
- `api` → `application`, `foundation`
- `events` → `domain`, `foundation`
- `cmd` → all layers

## Adding a New Engine

1. Copy this template to `backend/engines/<engine-name>/`
2. Create `go.mod` for the engine
3. Update `go.work` to include the new module
4. Implement domain entities and events
5. Implement use cases
6. Implement infrastructure
7. Wire dependencies in `cmd/`
