# Disaster Recovery Plan

## RPO (Recovery Point Objective): 1 hour
## RTO (Recovery Time Objective): 2 hours

## Backup Strategy

### PostgreSQL
```bash
# Automated hourly backup (cron)
0 * * * * docker exec cytaxi-postgres pg_dump -U cytaxi cytaxi | gzip > /backups/cytaxi_$(date +\%Y\%m\%d_\%H).sql.gz

# Retain 7 daily + 4 weekly + 3 monthly
```

### Redis
```bash
# Every 30 minutes
*/30 * * * * docker exec cytaxi-redis redis-cli BGSAVE
```

## Recovery Procedure

### Full System Recovery
1. Stop all services: `docker compose -f docker-compose.prod.yml down`
2. Restore PostgreSQL: `gunzip -c /backups/cytaxi_20260715_14.sql.gz | docker exec -i cytaxi-postgres psql -U cytaxi`
3. Start services: `docker compose -f docker-compose.prod.yml up -d`
4. Verify: `./scripts/deploy.sh`

### Single Engine Recovery
1. Restart engine: `docker compose -f docker-compose.prod.yml restart {engine}`
2. Check health: `curl http://localhost:{port}/health`
3. If persistent: `docker compose -f docker-compose.prod.yml up -d --force-recreate {engine}`

### Database Corruption
1. Stop engines (keep postgres): `docker compose stop`
2. Restore from latest backup
3. Start all: `docker compose start`
4. Verify data: `docker exec cytaxi-postgres psql -U cytaxi -c "SELECT count(*) FROM trips;"`

### Region Failure (Cloud)
1. Deploy to DR region using same docker-compose.prod.yml
2. Restore latest backup
3. Update DNS: `cytaxi.app → DR region IP`
4. Verify all endpoints

## Verification Checklist
- [ ] API Gateway health: `curl https://api.cytaxi.app/health`
- [ ] All engines responding: check each port
- [ ] Trip creation: POST /api/v1/trip/request
- [ ] Payment: POST /api/v1/payments
- [ ] Driver matching: POST /api/v1/matching/start
- [ ] SSE tracking: GET /api/v1/trip/ws?trip_id=test
