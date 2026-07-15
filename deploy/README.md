# CYTAXI Deployment Guide

## Prerequisites
- Docker 24+
- Docker Compose 2+
- Go 1.22+ (for local dev)
- PostgreSQL 16 + Redis 7 (via Docker)

## Quick Start

```bash
# Development
docker compose -f docker-compose.dev.yml up -d

# Production
docker compose -f docker-compose.prod.yml up -d

# Verify
curl http://localhost:8000/health
curl http://localhost:8000/api/v1/trip/health
```

## Environment Configuration

| Environment | File | Purpose |
|-------------|------|---------|
| Development | `.env` | Local development, debug=true |
| Testing | `.env.test` | CI/CD, test database |
| Production | `.env.prod` | Secrets via vault/env vars |

## Secrets

```
JWT_SECRET=<256-bit-key>
DB_PASSWORD=<strong-password>
REDIS_PASSWORD=<optional>
PAYMENT_API_KEY=<gateway-key>
```

## Port Map

| Service | Port | Description |
|---------|------|-------------|
| Gateway | 8000 | API entry point |
| Trip | 8087 | Trip lifecycle |
| Pricing | 8088 | Fare calculation |
| Payment | 8091 | Payment processing |
| Customer | 8085 | Customer profiles |
| Driver | 8086 | Driver profiles |
| Notification | 8090 | Notifications |
| Matching | 8089 | Driver matching |
| Admin | 8094 | Administration |
| Analytics | 8093 | Business intelligence |
| Trust | 8092 | Identity verification |

## Health Checks

| Endpoint | Purpose |
|----------|---------|
| GET /health | Liveness (service running) |
| GET /ready | Readiness (dependencies ok) |
| GET /live | Kubernetes liveness probe |

## Logging

All services output JSON structured logs to stdout:
```json
{"time":"...","level":"INFO","msg":"request","method":"GET","path":"/health","status":200,"duration":"2ms"}
```

## Monitoring

- Metrics: Prometheus `/metrics` endpoint (future)
- Tracing: OpenTelemetry (future)
- Logs: stdout JSON → Loki (future)

## Rollback

```bash
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d --build
```

## Backup

```bash
docker exec cytaxi-postgres pg_dump -U cytaxi cytaxi > backup.sql
docker exec cytaxi-redis redis-cli SAVE
```
