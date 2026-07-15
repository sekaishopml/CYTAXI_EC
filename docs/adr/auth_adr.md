# ADR-010: Centralized Authentication Platform

**Status:** Accepted
**Date:** Sprint 35

**Context:** Need centralized authentication for all actors (customer, driver, operator, admin) without duplicating auth logic across engines.

**Decision:** Trust & Identity Engine hosts the auth platform:
- Auth providers via adapter pattern (email/password, OTP, Google OAuth, Apple Sign-In)
- JWT + refresh tokens with HMAC-SHA256
- RBAC with 4 roles (customer, driver, operator, admin)
- Session management with revocation
- All auth goes through API Gateway

**Roles:**
- customer: trip create/read, profile read/write, payment read/create
- driver: trip read/accept/reject/start/finish, profile read, vehicle manage, availability manage
- operator: trip read/cancel, driver read, payment read, refund create
- admin: * (all permissions)

**Consequences:**
- Single auth flow for all frontends
- Tokens validatable by any engine (shared secret)
- Provider switch requires zero code changes in consuming engines
