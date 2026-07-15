# Architecture Decision Records

## ADR-001: Monorepo with Go Modules

**Status:** Accepted
**Date:** Sprint 01

**Context:** Need to organize 14+ microservices with shared foundation.

**Decision:** Single monorepo with `go.work` workspace file. Each engine is a separate Go module. Foundation packages are shared via the root module.

**Consequences:** 
- Simpler dependency management
- Single source of truth for all code
- CI/CD can build all services from one repo
- Go workspace resolves module references automatically

---

## ADR-002: Clean Architecture + DDD

**Status:** Accepted
**Date:** Sprint 02

**Context:** Need consistent architecture across all engines.

**Decision:** Every engine follows `domain → application → infrastructure → api` layers. Domain has zero external dependencies. Application orchestrates use cases. Infrastructure implements repository interfaces.

**Consequences:**
- High testability (domain is pure Go)
- Replaceable infrastructure (DB, message broker, external APIs)
- Steep learning curve for new developers
- More boilerplate than simpler architectures

---

## ADR-003: Event-Driven Communication

**Status:** Accepted
**Date:** Sprint 21

**Context:** Engines need to communicate without direct dependencies.

**Decision:** All cross-engine communication via Domain Events through the EventBus (Integration Layer). No engine imports another engine's domain directly. Saga pattern for multi-step workflows. Outbox/Inbox for guaranteed delivery.

**Consequences:**
- Eventually consistent across engines
- Higher latency for cross-engine operations
- Complex debugging (distributed tracing needed)
- Loose coupling enables independent deployment

---

## ADR-004: API Gateway as Single Entry Point

**Status:** Accepted
**Date:** Sprint 22

**Context:** Frontends should not access engines directly.

**Decision:** API Gateway is the only HTTP entry point. All engines are internal (no public exposure). Gateway handles auth (JWT), rate limiting, CORS, correlation IDs, request routing via reverse proxy.

**Consequences:**
- Single point for security enforcement
- All traffic goes through one service (potential bottleneck)
- Frontends have simpler configuration (one URL)
- Gateway can aggregate responses from multiple engines

---

## ADR-005: Simulated Payment Providers

**Status:** Accepted (temporary)
**Date:** Sprint 28

**Context:** Need payment flow for MVP without real payment provider integration.

**Decision:** Payment Engine uses in-memory processing with simulated gateways. PaymentGateway interface defines the adapter contract. Real providers (Stripe, PayPhone, Kushki) will implement this interface later.

**Consequences:**
- MVP can demonstrate full flow
- No regulatory/compliance concerns
- Real money not at risk
- Provider switch is plug-and-play via adapter

---

## ADR-006: Server-Sent Events for Real-Time Tracking

**Status:** Accepted (MVP)
**Date:** Sprint 27

**Context:** Need real-time updates for live trip tracking.

**Decision:** SSE (Server-Sent Events) over HTTP streaming instead of native WebSocket. Simpler implementation, built-in reconnection, works through standard HTTP proxies and the API Gateway.

**Consequences:**
- Unidirectional (server→client only)
- Client-side reconnection is automatic
- Works through Nginx/API Gateway without special config
- Upgrade to WebSocket for bidirectional needs in future
