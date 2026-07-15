# CYTAXI Operational Dashboard — Runbook

## Grafana Dashboard URLs

| Dashboard | URL | Refresh |
|-----------|-----|---------|
| Platform Overview | `/d/cytaxi-overview` | 10s |
| Engine Health | `/d/engine-health` | 10s |
| Business Metrics | `/d/business-metrics` | 30s |
| Security & Auth | `/d/security` | 30s |

## Critical Alerts (Prometheus)

| Alert | Condition | Severity | Action |
|-------|-----------|----------|--------|
| ServiceDown | `up == 0` for >30s | SEV1 | Restart service |
| HighLatency | `p95 > 1s` for 5min | SEV2 | Check CPU, scale |
| DBDown | `pg_up == 0` for >10s | SEV1 | Check pg_isready |
| RedisDown | `redis_up == 0` for >10s | SEV2 | Restart Redis |
| HighErrorRate | `error_rate > 5%` for 5min | SEV2 | Check logs |
| DiskFull | `disk_used > 85%` | SEV3 | Cleanup logs/backups |
| CertificateExpiry | `cert_expiry < 7 days` | SEV3 | Renew cert |

## Daily Health Check Commands

```bash
# Gateway
curl -s http://localhost:8000/health

# All engines
for p in 8087 8088 8089 8091 8085 8086 8090 8094 8093 8092; do
  curl -s -o /dev/null -w "Port $p: %{http_code}\n" http://localhost:$p/health
done

# Database
docker exec cytaxi-postgres pg_isready -U cytaxi

# Redis
docker exec cytaxi-redis redis-cli PING

# Disk
df -h | grep -E "/$|/data"
```

## Weekly Report Template

```
CYTAXI Weekly Operations Report
Week: [Date Range]

### Availability
- Platform uptime: 99.X%
- Incidents: N (SEV1: 0, SEV2: 0, SEV3: N)
- MTTR: X minutes

### Performance
- Avg response time: Xms
- p95 response time: Xms
- Peak RPS: X

### Capacity
- CPU avg: X%
- Memory avg: X%
- Disk: X% used
- DB size: XGB

### Business (last 7 days)
- Trips: N
- Revenue: $X
- New users: N
- Active drivers: N
```
