# CYTAXI

Conversation-First Mobility Platform.

## Architecture

- **Primary Channel:** WhatsApp
- **Secondary Channels:** Miniweb, Driver Web App, Admin Dashboard
- **Backend:** Go
- **Architecture:** Event-Driven Microservices, DDD, Clean Architecture, CQRS

## Repository Structure

```
CYTAXI_EC/
├── backend/
│   ├── foundation/       # Shared foundation packages
│   ├── config/           # Configuration management
│   ├── logger/           # Logging infrastructure
│   ├── errors/           # Error handling
│   ├── http/             # HTTP server and client
│   ├── middleware/       # HTTP middleware
│   ├── validation/       # Validation utilities
│   ├── observability/    # Observability infrastructure
│   ├── telemetry/        # Distributed tracing
│   ├── auth/             # Authentication and authorization
│   ├── events/           # Event bus interfaces
│   ├── utils/            # General utilities
│   ├── testing/          # Testing utilities
│   ├── containers/       # DI container utilities
│   ├── engine-template/  # Standard Engine template (Clean Architecture + DDD)
│   └── engines/          # Future Engine implementations (TBD)
├── docs/
│   └── BACKEND_ARCHITECTURE.md  # Backend architecture specification
├── go.work
├── .env.example
├── .editorconfig
├── .golangci.yml
├── Makefile
└── docker-compose.dev.yml
```

## Development

```bash
make lint    # Run linters
make test    # Run tests
make build   # Build project
```

## Principles

- Documentation before implementation
- Architecture before optimization
- Business before technology
- Consistency over speed
- Replaceability
- Long-term thinking

## Source of Truth

The [CYTAXI-BLUEPRINT](https://github.com/sekaishopml/cydigital-blueprint) repository is the single source of truth for architecture, business rules, and engineering standards.
