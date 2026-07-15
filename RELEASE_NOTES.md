# CYTAXI v1.0.0-rc1 Release Notes

## Release Candidate 1 — Pilot Launch

**Date:** 2026-07-15
**Status:** Ready for controlled pilot

### What's Included

This release candidate includes the complete CYTAXI platform with:

- **11 backend microservices** implementing DDD + Clean Architecture + CQRS + Event Driven
- **API Gateway** as single HTTP entry point with JWT auth, rate limiting, and CORS
- **Integration Layer** with EventBus, Saga, Outbox/Inbox patterns for eventual consistency
- **2 frontend applications** (Customer MiniWeb + Driver Web Portal) built with React/Next.js
- **Full MVP flow** from trip request → driver assignment → live tracking → payment

### System Requirements

- Docker 24+ and Docker Compose 2+
- 4GB RAM minimum (8GB recommended)
- 20GB disk space
- PostgreSQL 16 + Redis 7 (included in Docker compose)

### Quick Start

```bash
git clone https://github.com/sekaishopml/CYTAXI_EC.git
cd CYTAXI_EC
docker compose -f docker-compose.prod.yml up -d
curl http://localhost:8000/health
```

### Services

| Service | Port | Status |
|---------|------|--------|
| API Gateway | 8000 | ✓ |
| Trip Engine | 8087 | ✓ |
| Pricing Engine | 8088 | ✓ |
| Payment Engine | 8091 | ✓ |
| Matching Engine | 8089 | ✓ |
| Customer Engine | 8085 | ✓ |
| Driver Engine | 8086 | ✓ |
| Notification Engine | 8090 | ✓ |
| Admin Engine | 8094 | ✓ |
| Analytics Engine | 8093 | ✓ |
| Trust Engine | 8092 | ✓ |

### Known Limitations

- Payment gateways are simulated (no Stripe/PayPhone integration)
- GPS location is simulated (no real device GPS)
- WhatsApp integration uses placeholder adapter
- No email provider integration
- SSE instead of native WebSocket for live tracking

### Breaking Changes

None. This is the initial release.

### Upgrade Notes

N/A — first release. Follow `deploy/README.md` for fresh installation.
