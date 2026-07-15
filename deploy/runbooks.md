# CYTAXI Runbooks

## Incident: Service Down

1. Check Gateway: `curl http://localhost:8000/health`
2. Check specific engine: `curl http://localhost:{port}/health`
3. Check Docker: `docker compose -f docker-compose.prod.yml ps`
4. Check logs: `docker compose -f docker-compose.prod.yml logs -f {service}`
5. Restart: `docker compose -f docker-compose.prod.yml restart {service}`

## Incident: High Latency

1. Check rate limiter: `curl http://localhost:8000/health`
2. Check PostgreSQL: `docker exec cytaxi-postgres pg_isready`
3. Check Redis: `docker exec cytaxi-redis redis-cli PING`
4. Check CPU/Memory: `docker stats`
5. Scale if needed: update replicas in compose

## Incident: Payment Failure

1. Verify Payment Engine: `curl http://localhost:8091/health`
2. Check payment gateway status (simulated)
3. Check recent payments: `curl http://localhost:8000/api/v1/payments/history`
4. Trigger retry if needed

## Incident: Matching Failure

1. Verify Matching Engine: `curl http://localhost:8089/health`
2. Check Driver Engine: `curl http://localhost:8086/health`
3. Restart matching: `docker compose restart matching`
4. Monitor candidates

## Database Recovery

```bash
# List backups
ls -la /backups/

# Restore PostgreSQL
docker exec -i cytaxi-postgres psql -U cytaxi cytaxi < backup.sql

# Verify
docker exec cytaxi-postgres psql -U cytaxi -c "SELECT count(*) FROM trips;"
```

## Emergency Contacts

- Platform: WhatsApp +593 99 999 9999
- Database: DBA on-call
- Payment: Finance team
