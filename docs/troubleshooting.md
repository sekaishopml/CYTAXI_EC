# Troubleshooting Guide

## Gateway not responding (port 8000)

```bash
# Check if running
docker compose -f docker-compose.prod.yml ps gateway

# Check logs
docker compose -f docker-compose.prod.yml logs gateway --tail=50

# Restart
docker compose -f docker-compose.prod.yml restart gateway
```

## Engine health check failing

```bash
# Check specific engine
curl http://localhost:{port}/health
docker compose logs {engine} --tail=20

# Common causes:
# 1. PostgreSQL not ready → check `docker exec cytaxi-postgres pg_isready`
# 2. Port conflict → check `netstat -tlnp | grep {port}`
# 3. Out of memory → check `docker stats`
```

## Database connection errors

```bash
# Verify PostgreSQL
docker exec cytaxi-postgres pg_isready -U cytaxi
# Check max connections
docker exec cytaxi-postgres psql -U cytaxi -c "SELECT count(*) FROM pg_stat_activity;"
# Increase pool if needed (update DATABASE_MAX_CONNS env var)
```

## Payment processing fails

```bash
# Payment engine is simulated - no external dependencies
# Check payment engine health
curl http://localhost:8091/health
# Check payment history
curl http://localhost:8000/api/v1/payments/history
```

## Matching not finding drivers

```bash
# Verify driver availability
curl http://localhost:8000/api/v1/driver/requests

# Check matching engine
curl http://localhost:8089/health

# If no drivers: drivers must set status to "online"
```

## SSE tracking not streaming

```bash
# SSE uses HTTP streaming - ensure proxy doesn't buffer
# Check Nginx config has: proxy_buffering off;
# Test directly:
curl -N "http://localhost:8000/api/v1/trip/ws?trip_id=test"
```

## Disk space full

```bash
# Check Docker disk usage
docker system df

# Clean up
docker system prune -a --volumes -f

# Verify backup rotation
ls -lh /backups/
```
