# Changelog

## v1.0.0-rc1 (2026-07-15)

### Platform Architecture
- Monorepo with 14 Go modules + 2 React/Next.js frontends
- DDD + Clean Architecture + CQRS + Event Driven
- API Gateway as single entry point (port 8000)
- Integration Layer with EventBus, Saga, Outbox/Inbox patterns

### Engines (11 microservices)
- **Conversation Engine** (port 8081) — conversational flow + session manager
- **Geospatial Engine** (port 8082) — 3 provider adapters (Google/OSM/Mapbox)
- **Policy Engine** (port 8083) — centralized business rules
- **Mobility Decision Engine** (port 8084) — dispatch coordinator
- **Customer Engine** (port 8085) — customer profiles + preferences
- **Driver Engine** (port 8086) — driver/vehicle/license management
- **Trip Engine** (port 8087) — trip lifecycle (11 states)
- **Pricing Engine** (port 8088) — fare calculation + promotions
- **Matching Engine** (port 8089) — candidate ranking + selection
- **Notification Engine** (port 8090) — 6-channel notifications
- **Payment Engine** (port 8091) — payment/refund/receipt processing
- **Trust & Identity Engine** (port 8092) — KYC/verification/fraud detection
- **Analytics Engine** (port 8093) — business metrics + BI
- **Administration Engine** (port 8094) — roles/permissions/feature flags

### Frontends
- **Customer MiniWeb** — React/Next.js, 6 pages, 8 components
- **Driver Web Portal** — React/Next.js, 10 pages, sidebar navigation

### Infrastructure
- Docker Compose (dev + prod)
- GitHub Actions CI/CD (lint → test → build → security → deploy)
- PostgreSQL 16 + Redis 7
- SSE real-time tracking
- JWT auth middleware
- Rate limiting, CORS, structured logging

### MVP Flows
- Customer journey: request → pricing → trip creation → notification
- Driver assignment: matching → candidates → accept/reject → reassignment
- Live tracking: SSE streaming → position/ETA → trip completion
- Payment: fare calculation → method selection → payment → receipt

### Production Readiness
- Health/Readiness/Liveness checks on all services
- Deployment guide + runbooks
- Smoke tests
- Disaster recovery plan
