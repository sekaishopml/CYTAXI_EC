# CYTAXI Operations Guide — v1.0

## 1. Service Level Objectives (SLO)

| Metric | Target | Measurement Window |
|--------|--------|-------------------|
| API Availability | 99.9% | Monthly |
| API Latency (p95) | <500ms | 5 minutes |
| Trip Creation Success | >99.5% | Daily |
| Payment Success Rate | >99% | Daily |
| Matching Time | <10s average | Per request |
| Recovery Time (RTO) | <15 minutes | Per incident |
| Recovery Point (RPO) | <1 hour | Data loss window |

## 2. Incident Severity Levels

| Level | Name | Response | Escalation | Example |
|-------|------|----------|------------|---------|
| SEV1 | Critical | <5 min | CTO + Engineering Lead | Platform completely down |
| SEV2 | Major | <15 min | Engineering Lead + DevOps | Single engine down |
| SEV3 | Minor | <1 hour | DevOps | High latency, degraded |
| SEV4 | Low | <4 hours | Support | Bug, cosmetic issue |

## 3. Incident Response Flow

```
1. Detect (monitoring alert / user report)
2. Acknowledge (<5 min for SEV1)
3. Triage (assign severity, notify team)
4. Mitigate (apply runbook or ad-hoc fix)
5. Resolve (verify recovery)
6. Post-Mortem (<24h after resolution)
7. Track (close incident, update runbooks)
```

## 4. Escalation Matrix

| Role | SEV1 | SEV2 | SEV3 | SEV4 |
|------|------|------|------|------|
| Platform Engineer | ✓ | ✓ | ✓ | ✓ |
| DevOps Engineer | ✓ | ✓ | ✓ | — |
| Engineering Lead | ✓ | ✓ | — | — |
| CTO | ✓ | — | — | — |
| Support Team | Informed | Informed | ✓ | ✓ |

## 5. Operational Procedures

### Daily
- [ ] Check Grafana dashboard for anomalies
- [ ] Review error logs (Loki)
- [ ] Verify backup completion
- [ ] Check disk/memory usage

### Weekly
- [ ] Review incident log
- [ ] Test restore from latest backup
- [ ] Rotate secrets if due
- [ ] Apply security patches (staging first)

### Monthly
- [ ] Performance review (response time trends)
- [ ] Capacity planning (CPU, RAM, disk trends)
- [ ] Security audit (RBAC review, access logs)
- [ ] Runbook review and updates

### Quarterly
- [ ] Full disaster recovery test
- [ ] Load test (validate scaling)
- [ ] Architecture review
- [ ] Compliance audit

## 6. Maintenance Windows

| Type | Schedule | Notice | Max Duration |
|------|----------|--------|-------------|
| Routine | Sunday 02:00-04:00 UTC | 48h | 2 hours |
| Emergency | Any time | As soon as possible | Until resolved |
| Database | Sunday 03:00-04:00 UTC | 72h | 1 hour |

## 7. On-Call Rotation

| Week | Primary | Secondary |
|------|---------|-----------|
| Week 1 | Engineer A | Engineer B |
| Week 2 | Engineer B | Engineer A |

- Primary responds within 5 min (SEV1)
- Secondary escalates if primary unreachable

## 8. Communication Channels

| Channel | Purpose |
|---------|---------|
| WhatsApp Group | Urgent alerts (SEV1/SEV2) |
| Slack #incidents | Incident tracking |
| Email | Post-mortems, weekly reports |
| Status Page | External communication |

## 9. Key Contacts

| Role | Contact |
|------|---------|
| Platform Engineering | support@cytaxi.app |
| DevOps | devops@cytaxi.app |
| Security | security@cytaxi.app |
| Emergency | +593 99 999 9999 |
