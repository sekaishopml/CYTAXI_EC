# CODE FREEZE — v1.0.0

**Effective:** 2026-07-15
**Status:** ACTIVE

## What is frozen
- All new feature development
- Domain model changes
- API contract changes
- Database schema changes
- Architecture modifications to the Blueprint

## What is allowed
- Critical security fixes (with ADR)
- Deployment configuration adjustments
- Documentation updates
- Monitoring/alerting tuning

## Release Gate
Before any change to this codebase:
1. Must pass architecture review against CYTAXI_BLUEPRINT
2. Must not break existing contracts (OpenAPI, Events)
3. Must be documented with ADR if architectural impact

This code freeze remains in effect until the first v1.0 patch release (v1.0.1).
