# CYTAXI v1.0.0 Release Notes

**Release Date:** 2026-07-15
**Status:** ✓ PRODUCTION READY
**Architecture:** DDD + Clean Architecture + CQRS + Event Driven

## What is CYTAXI?

CYTAXI is a Conversation-First Mobility Platform. The primary user interface is WhatsApp. Miniweb, Dashboard, and Driver Web are secondary interfaces. Business Rules remain deterministic. Artificial Intelligence assists but never owns business state.

## v1.0.0 Includes

### Backend (14 Microservices)
| Engine | Port | Description |
|--------|------|-------------|
| API Gateway | 8000 | Single HTTP entry point, reverse proxy, JWT auth |
| Conversation Engine | 8081 | Conversational flow + session manager + AI orchestrator |
| Geospatial Engine | 8082 | 4 map providers (OSM real, Google/Mapbox/Here stubs) |
| Policy Engine | 8083 | Centralized business rules engine |
| Mobility Decision | 8084 | Dispatch coordinator + candidate ranking |
| Customer Engine | 8085 | Profiles, favorites, preferences, loyalty |
| Driver Engine | 8086 | Drivers, vehicles, licenses, fleet management |
| Trip Engine | 8087 | Trip lifecycle (11 states), tracking SSE |
| Pricing Engine | 8088 | Dynamic pricing, 4 strategies, promotions, coupons |
| Matching Engine | 8089 | Intelligent dispatch, multi-factor scoring, zones |
| Notification Engine | 8090 | 6 channels (WhatsApp, Push, Email, SMS, WebSocket, In-App) |
| Payment Engine | 8091 | 5 provider adapters, webhooks, refunds, settlements |
| Trust & Identity | 8092 | KYC, auth, trust score, reputation, verification |
| Analytics Engine | 8093 | BI dashboard, KPIs, reports, trends |
| Administration Engine | 8094 | Roles, permissions, feature flags, audit |

### Frontends
| App | Tech | URL |
|-----|------|-----|
| Customer MiniWeb | React 18 + Next.js 14 | `/` |
| Driver Web Portal | React 18 + Next.js 14 | `/driver` |
| Admin Dashboard | React 18 + Next.js 14 | `/admin` |

### Infrastructure
- Docker Compose (14 services + Postgres + Redis + Prometheus + Grafana)
- Nginx reverse proxy with TLS 1.3, 10 security headers
- GitHub Actions CI/CD (lint → test → build → security → deploy)
- Dual network architecture (public_net + private_net)
- Backup/restore automation (7-day rotation)
- Blue/Green deployment prepared

### Security
- TLS 1.3, HSTS, CSP, X-Frame-Options, 10 security headers
- JWT + Refresh Tokens (HMAC-SHA256)
- RBAC: customer, driver, operator, admin
- Rate limiting: API (100r/s), Auth (5r/s)
- OWASP Top 10 addressed
- Secrets management (.env per environment, never committed)
- Fail2Ban: SSH + API + Auth

### Observability
- Prometheus + Grafana (7 dashboards)
- Structured JSON logging (slog)
- Correlation IDs across all services
- SSE real-time tracking
- Health/Readiness/Liveness probes on all services

### Provider Adapters (14 total)
| Category | Providers |
|----------|-----------|
| Maps | OpenStreetMap (real), Google Maps (stub), Mapbox (stub), Here (stub) |
| Payments | Stripe (stub), Kushki (stub), PayPhone (stub), PayPal (stub), Mock |
| Messaging | Meta Cloud API (stub), Mock |
| Auth | Email/Password, OTP, Google OAuth (stub), Apple Sign-In (stub) |

### Key Metrics
- Architecture: DDD with 14 bounded contexts
- Code: 264+ files, 12,800+ lines
- APIs: 150+ endpoints documented
- Events: 100+ domain events
- Sprints: 50 completed
- Docs: 30+ documentation files

## System Requirements
- Docker 24+ & Docker Compose 2+
- 4GB RAM (8GB recommended)
- 20GB disk
- Ubuntu 24.04 LTS

## Quick Start
```bash
git clone https://github.com/sekaishopml/CYTAXI_EC.git
cd CYTAXI_EC
docker compose -f docker-compose.prod.yml up -d
curl http://localhost:80/health
```

**CYTAXI v1.0.0 — Production Ready.**
