# CYTAXI Incident Response Plan

## Quick Reference Card

```
┌─────────────────────────────────────────────────┐
│ SEV1 - CRITICAL (Platform Down)                 │
│ → Call Engineer on-call: +593 99 999 9999      │
│ → Slack: #incidents @channel                   │
│ → Check: /health, docker ps, nginx status      │
│ → Runbook: service_down (6 steps)              │
│ → Target: <15 min resolution                   │
└─────────────────────────────────────────────────┘
```

## Incident Categories

| Category | Symptoms | Initial Action |
|----------|----------|---------------|
| Platform Down | All /health fail | Restart Gateway → Check Docker |
| DB Failure | PostgreSQL connection errors | Check pg_isready → Restart |
| High Latency | p95 > 500ms | Check CPU/RAM → Scale |
| Payment Issue | Payment failures >5% | Check provider → Disable auto-capture |
| Matching Issue | No drivers assigned | Check Driver Engine → Restart |
| Security Breach | Unauthorized access detected | Isolate service → Rotate secrets |

## Response Templates

### Initial communication (SEV1/SEV2)
```
@channel INCIDENT: [Title]
Severity: SEV1
Impact: [Description]
Started: [Time]
Engineer: [Name]
Status: Investigating
Updates in #incidents
```

### Resolution communication
```
RESOLVED: [Title]
Duration: [X minutes]
Root Cause: [Brief description]
Fix: [What was done]
Prevention: [How to prevent recurrence]
Post-mortem: [Link]
```

## Post-Mortem Template

```
# Incident Post-Mortem: [Title]

**Date:** [Date]
**Duration:** [X minutes]
**Severity:** SEV1/2/3/4

## Timeline
- [Time] — Incident detected
- [Time] — Engineer acknowledged
- [Time] — Root cause identified
- [Time] — Fix applied
- [Time] — Verified resolved

## Root Cause
[Technical explanation]

## Impact
- Users affected: [N]
- Revenue impact: [$X]
- Data loss: [Yes/No]

## Action Items
1. [Preventive measure]
2. [Monitoring improvement]
3. [Runbook update]
```
