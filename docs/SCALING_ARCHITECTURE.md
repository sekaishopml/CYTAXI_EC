# CYTAXI High Availability & Scaling Architecture

## HA Architecture

```
                    ┌──────────────┐
                    │  DNS / CDN    │
                    └──────┬───────┘
                           │
                    ┌──────┴───────┐
                    │ Load Balancer │ (Nginx / HAProxy)
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         App Server 1  App Server 2  App Server N
         (Gateway)     (Gateway)     (Gateway)
              │            │            │
              └────────────┼────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         PostgreSQL    Redis        RabbitMQ
         (Primary)   (Sentinel)   (Cluster)
         + Replica
```

## Scaling Guidelines

### Stateless Services (scale horizontally)
- API Gateway: 2+ replicas
- Conversation Engine: 1 per 500 concurrent users
- Notification Engine: 1 per 1000 notifications/min
- Matching Engine: 2 minimum

### Stateful Services (scale vertically first)
- PostgreSQL: 4 CPU, 16GB RAM for <10K users
- Redis: 2 CPU, 8GB RAM
- RabbitMQ: 2 CPU, 4GB RAM

### Capacity Planning

| Users | Gateway | Engines | DB | Redis |
|-------|---------|---------|-----|-------|
| <1K | 1 | 1 each | 2CPU/4GB | 1CPU/2GB |
| 1K-5K | 2 | 1 each | 4CPU/8GB | 2CPU/4GB |
| 5K-20K | 4 | 2 each | 8CPU/16GB | 4CPU/8GB |
| 20K+ | 8+ | 4+ each | 16CPU/32GB+ | 8CPU/16GB+ |

## Auto-Recovery Patterns

| Failure | Detection | Recovery | Max Downtime |
|---------|-----------|----------|-------------|
| Gateway crash | Health check 10s | Docker restart | <15s |
| Engine crash | Health check 10s | Docker restart | <15s |
| PostgreSQL crash | pg_isready 5s | Docker restart + WAL | <30s |
| Redis crash | Redis PING 5s | Docker restart + AOF | <10s |
| Disk full | Monitor >80% | Alert + manual cleanup | <1h |

## Deployment Strategy

### Blue/Green (prepared)
```
1. Deploy new version to "green" environment
2. Health check green
3. Switch load balancer: blue → green
4. Monitor for 5 minutes
5. If healthy: decommission blue
6. If unhealthy: switch back to blue
```

### Rollback
```
1. If issue detected within 5 min:
   docker compose up -d --force-recreate {service}
2. If need full rollback:
   docker compose down && docker compose up -d
3. Verify: ./infra/backup/verify_restore.sh
```

## Connection Pooling

| Service | Min | Max | Timeout |
|---------|-----|-----|---------|
| PostgreSQL | 5 | 20 | 30s |
| Redis | 5 | 50 | 5s |
| RabbitMQ | 1 | 10 | 30s |
