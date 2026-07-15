# Backend Architecture

## Overview

CYTAXI backend follows **Clean Architecture** + **Domain-Driven Design** (DDD) + **Event-Driven Architecture**.

Every Engine in the system shares the same architecture template. This document defines that standard.

---

## Architecture Principles

1. **Domain is the center** — Business logic lives in `domain/`, has zero external dependencies.
2. **Dependency inversion** — High-level modules (domain) never import low-level modules (infrastructure).
3. **Boundaries by layer** — Each layer communicates through ports (interfaces).
4. **Events are first-class citizens** — Business facts are expressed as Domain Events.
5. **Infrastructure is replaceable** — DB, message broker, external APIs are all behind interfaces.

---

## Layer Architecture

```
┌──────────────────────────────────────────────────┐
│                    cmd/                           │
│            Composition Root (DI)                  │
├──────────────────────────────────────────────────┤
│                  api/                             │
│         HTTP Handlers, Routes, DTOs               │
├──────────────────────────────────────────────────┤
│               application/                        │
│          Use Cases, Ports, DTOs                   │
├──────────────────────────────────────────────────┤
│                domain/                            │
│   Entities, Value Objects, Aggregates, Events     │
├──────────────────────────────────────────────────┤
│            infrastructure/                        │
│  DB, Cache, Message Broker, External APIs         │
├──────────────────────────────────────────────────┤
│           foundation (shared/pkg)                 │
│  Logger, Config, Errors, HTTP, Events, Telemetry  │
└──────────────────────────────────────────────────┘
```

---

## Standard Engine Structure

Every Engine MUST follow this structure:

```
engine-name/
├── cmd/
│   └── engine-name/
│       └── main.go              # Entrypoint
├── internal/
│   └── engine-name/
│       ├── domain/
│       │   ├── entity/          # Entities & Aggregates
│       │   ├── valueobject/     # Value Objects
│       │   ├── repository/      # Repository interfaces
│       │   └── event/           # Domain events
│       ├── application/
│       │   ├── usecase/         # Use cases / Interactors
│       │   ├── dto/             # Data Transfer Objects
│       │   └── port/            # Inbound/Outbound ports
│       ├── infrastructure/
│       │   ├── database/        # DB implementations
│       │   ├── cache/           # Cache implementations
│       │   ├── messagebroker/   # Message broker implementations
│       │   └── externalapi/     # External API clients
│       ├── api/
│       │   ├── handler/         # HTTP handlers
│       │   ├── router/          # Route definitions
│       │   └── middleware/      # Engine-specific middleware
│       ├── events/
│       │   ├── definition/      # Event type definitions
│       │   └── handler/         # Event handlers
│       ├── config/              # Engine configuration
│       └── tests/               # Integration & e2e tests
├── README.md
└── Dockerfile
```

---

## Folder Purpose Reference

| Folder | Purpose | Can import |
|--------|---------|------------|
| `cmd/` | Application entrypoint (main func, DI wiring) | All layers |
| `domain/entity/` | Business entities, aggregates, invariants | Nothing |
| `domain/valueobject/` | Immutable value objects | Domain types |
| `domain/repository/` | Repository interfaces (ports) | Domain types |
| `domain/event/` | Domain event structs | Domain types |
| `application/usecase/` | Orchestration of business rules | `domain/` |
| `application/dto/` | Input/output DTOs for use cases | Domain types |
| `application/port/` | Inbound/outbound port interfaces | `domain/` |
| `infrastructure/database/` | DB implementations (PostgreSQL, Redis) | `domain/repository/`, `application/port/` |
| `infrastructure/cache/` | Cache implementations | `application/port/` |
| `infrastructure/messagebroker/` | Message broker implementations (NATS, RabbitMQ) | `events/` |
| `infrastructure/externalapi/` | External HTTP/gRPC clients | `application/port/` |
| `api/handler/` | HTTP handlers, request/response mapping | `application/usecase/` |
| `api/router/` | Route registration, middleware chains | `api/handler/` |
| `api/middleware/` | Engine-specific HTTP middleware | Foundation |
| `events/definition/` | Event schemas, versioning | `domain/event/` |
| `events/handler/` | Event consumption logic | `application/usecase/` |
| `config/` | Env loading, config structs | Foundation |
| `tests/` | Integration, e2e, contract tests | All layers |

---

## Dependency Rules

### Strict dependency direction

```
domain  ←  application  ←  infrastructure
domain  ←  application  ←  api
domain  ←  events
```

### Allowed imports

- **domain/** imports NOTHING (not even foundation). Pure Go.
- **application/** imports `domain/` and `foundation` types.
- **infrastructure/** imports `domain/` (interfaces) and `foundation`.
- **api/** imports `application/` and `foundation`.
- **events/** imports `domain/` and `foundation`.
- **cmd/** imports all layers and foundation.
- **tests/** imports whatever is under test.

### Foundation packages available to all layers

- `backend/foundation/` — Shared base types
- `backend/logger/` — Logging interface
- `backend/errors/` — Error handling
- `backend/config/` — Configuration loader
- `backend/http/` — HTTP server & response helpers
- `backend/middleware/` — Common middleware
- `backend/validation/` — Validation utilities
- `backend/observability/` — Metrics & health
- `backend/telemetry/` — Distributed tracing
- `backend/auth/` — Auth interfaces
- `backend/events/` — Event bus interface
- `backend/containers/` — DI container

### Forbidden

- `domain/` importing `infrastructure/`, `api/`, `cmd/`, external frameworks.
- `application/` importing `infrastructure/`, `api/`, external frameworks.
- Circular imports between any engine packages.
- Direct DB calls from `api/handler/`.
- Business logic in `api/handler/`.

---

## Naming Conventions

### Packages

- All lowercase, no underscores, no mixed case.
- Single-word names preferred: `entity`, `usecase`, `handler`, `dto`.
- Avoid `common`, `util`, `misc` — be specific.
- Repository package name matches the domain concept: `driverrepository`, `triprepository`.

### Files

- Snake case: `driver_repository.go`, `create_trip.go`.
- One primary type per file.

### Interfaces

- Single method interfaces use the method name + `er` suffix: `Finder`, `Saver`, `Dispatcher`.
- Repository interfaces: `DriverRepository`, `TripRepository`.

### Errors

- Defined in the package they belong to.
- Use `errors.Kind` for classification.
- Prefix with context: `ErrDriverNotFound`, `ErrInvalidTripState`.

---

## Import Conventions

```go
// Standard library
import (
    "context"
    "fmt"
)

// Foundation packages
import (
    "github.com/sekaishopml/cytaxi/backend/errors"
    "github.com/sekaishopml/cytaxi/backend/logger"
)

// Same-engine packages
import (
    "github.com/sekaishopml/cytaxi/backend/engines/conversation/domain/entity"
    "github.com/sekaishopml/cytaxi/backend/engines/conversation/application/usecase"
)

// Cross-engine (only through events or shared kernel)
import (
    "github.com/sekaishopml/cytaxi/backend/engines/driver/domain/event"
)
```

---

## Module Strategy

Each Engine is a separate Go module within the monorepo.

```
go.work                        # workspace root
backend/go.mod                 # foundation module
backend/engines/conversation/go.mod
backend/engines/driver/go.mod
backend/engines/revenue/go.mod
...
```

Engine modules depend on the foundation module and, when necessary, on other engines' event packages.

---

## Transaction Boundaries

- **Application layer** opens and commits transactions.
- **Domain layer** never knows about transactions.
- **Infrastructure** implements transaction injection via `context.Context`.
- Use `Unit of Work` pattern across aggregates when needed.

---

## Testing Strategy

| Layer | Test type | Scope |
|-------|-----------|-------|
| `domain/` | Unit tests | Pure logic, no mocks |
| `application/` | Unit + integration | Mock repositories |
| `infrastructure/` | Integration | Real DB, real broker |
| `api/` | Integration | HTTP tests, mock use cases |
| `events/` | Integration | Pub/sub with test broker |

All tests use the shared `backend/testing/` utilities.

---

## Health & Observability

Each Engine exposes:

- `/health` — Liveness check
- `/ready` — Readiness check (DB, broker, cache)
- `/metrics` — Prometheus metrics (if applicable)
- Structured logs via `backend/logger`
- Distributed tracing via `backend/telemetry`

---

## Configuration

Each Engine loads config via `backend/config`:

```go
type Config struct {
    Engine  EngineConfig
    DB      config.DBConfig
    Redis   config.RedisConfig
    Log     config.LogConfig
    Auth    config.AuthConfig
}
```

All environment variables follow the pattern `ENGINE_NAME_VARIABLE_NAME`.

---

## Event-Driven Communication

- Engines communicate through Domain Events (not direct service calls).
- `backend/events.Bus` is the transport abstraction.
- Event schemas live in each Engine's `events/definition/`.
- Events are versioned and immutable.
- Cross-engine events use `nats` or `rabbitmq` in production.

---

## Anti-Patterns to Avoid

- ❌ Business logic in handlers
- ❌ Domain entities in API responses directly
- ❌ Skip use case layer (handler → repo directly)
- ❌ Transaction logic spread across handlers
- ❌ Circular imports between engines
- ❌ Shared database across engines
- ❌ Business logic in infrastructure
- ❌ Silent errors (always log + return)
