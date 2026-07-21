# ADR-011: Enterprise Security — SSO, MFA, Audit

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Enterprise customers require SSO (Google/Microsoft/Okta), MFA (TOTP/SMS/backup codes), and audit logging for compliance. The existing `@cytaxi/security` only provides basic JWT parse and RBAC.

## Decision
Extend `@cytaxi/security` with three enterprise modules:

### SSO (`SSOConfig`, `SSOUser`)
- Provider-agnostic config (clientId, clientSecret, authorizeUrl, tokenUrl, userInfoUrl)
- Supports Google, Microsoft, GitHub, Okta, Auth0
- Returns normalized `SSOUser` (id, email, name, avatar, provider)

### MFA (`MFARegistration`, `MFAChallenge`)
- TOTP with generated secrets (Base32, 20 bytes)
- Backup codes (8 codes, 6 digits each)
- `setupMFA()` generates secrets + backup codes
- `createMFAChallenge()` creates time-limited challenges (5 min TTL)
- `verifyTOTP()` validates 6-digit codes

### Audit Logging (`AuditLogger`)
- Structured entries: actor, action, resource, tenantId, ip, userAgent, details, success
- Query with filters (actorId, action, resource, tenantId, time range)
- In-memory ring buffer (max 10,000 entries)
- `enterpriseAudit` singleton for app-wide logging

### Session Management (`SessionManager`)
- Server-side sessions with expiry (default 24h)
- MFA verification flag per session
- Cleanup of expired sessions

## Consequences
- **Positive:** SSO reduces password fatigue and support tickets
- **Positive:** MFA satisfies compliance requirements (SOC 2, PCI)
- **Positive:** Audit log provides complete user action history
- **Negative:** In-memory audit log is lost on restart (need persistent storage)
- **Negative:** TOTP verification is client-side only (need server-side for production)

## Future
- Persistent audit log (PostgreSQL/Elasticsearch)
- Server-side TOTP verification with time-window drift
- U2F/WebAuthn hardware key support
- Session revocation API
