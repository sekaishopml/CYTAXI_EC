# CYTAXI Production Release Checklist v1.0.0-rc2
#
# This checklist must be completed before declaring any release candidate.

## 1. SECURITY HARDENING

- [ ] TLS 1.3 configured and tested (SSL Labs A+)
- [ ] Security headers present: HSTS, CSP, X-Frame-Options, X-Content-Type, X-XSS
- [ ] Rate limiting active on all public endpoints
- [ ] No secrets hardcoded in source code (verify with git grep)
- [ ] All API keys stored in .env files (never committed)
- [ ] JWT secret rotated and at least 256 bits
- [ ] RBAC enforced: customer, driver, operator, admin
- [ ] Audit logging enabled on all admin operations
- [ ] Input validation on all POST/PUT endpoints
- [ ] SQL injection prevention confirmed (parameterized queries)
- [ ] CORS configured restrictively (not wildcard in production)
- [ ] Docker images pinned to specific versions (no :latest)
- [ ] Docker containers run as non-root
- [ ] UFW firewall: only ports 80, 443, 22 open
- [ ] Fail2Ban configured for SSH and API brute force

## 2. INFRASTRUCTURE

- [ ] docker-compose.prod.yml validated with `docker compose config`
- [ ] All services have health checks (10s interval, 3 retries)
- [ ] Restart policy: unless-stopped on all services
- [ ] PostgreSQL: SSL connections, least privilege roles
- [ ] Redis: password protected, protected mode enabled
- [ ] Backup script tested and scheduled (cron)
- [ ] Restore procedure documented and tested
- [ ] Prometheus scraping all targets
- [ ] Grafana dashboard loading correctly
- [ ] Logs shipping to stdout (JSON format)
- [ ] Correlation IDs propagated across all services

## 3. API & CONTRACTS

- [ ] OpenAPI spec matches all implemented endpoints
- [ ] All endpoints return consistent error format: {"error": "message"}
- [ ] /health returns 200 on all engines
- [ ] /ready returns 200 when dependencies available
- [ ] /version returns v1.0.0-rc2
- [ ] All POST endpoints validate input
- [ ] Rate limits documented in OpenAPI

## 4. PERFORMANCE

- [ ] Response time < 200ms on /health (all engines)
- [ ] Database connection pool configured (min 5, max 20 per engine)
- [ ] Redis connection pool configured
- [ ] Asset compression enabled (gzip)
- [ ] Cache headers set on static assets

## 5. RECOVERY

- [ ] Backup restore tested within last 7 days
- [ ] Rollback procedure documented (< 5 minutes)
- [ ] Docker compose down/up cycle verified
- [ ] Database migration tested (forward + rollback)

## 6. DOCUMENTATION

- [ ] API Guide (docs/API_GUIDE.md) up to date
- [ ] Deployment Guide (deploy/README.md) up to date
- [ ] Runbooks (deploy/runbooks.md) cover all known incidents
- [ ] Disaster Recovery (docs/disaster_recovery.md) up to date
- [ ] Architecture decisions (docs/adr/) complete
- [ ] CHANGELOG updated for this release
- [ ] RELEASE_NOTES published

## 7. RELEASE

- [ ] VERSION file updated
- [ ] Git tag created: v1.0.0-rc2
- [ ] CI pipeline passing (lint → test → build → security → deploy)
- [ ] Docker images built and pushed to registry
- [ ] Smoke tests passing on production environment

## Sign-off

- [ ] Security Review: _______
- [ ] Architecture Review: _______
- [ ] Release Manager: _______
- [ ] Date: _______________
