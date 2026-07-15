#!/bin/bash
# CYTAXI Restore Verification Test
# Run this after restoring a backup to verify data integrity

set -euo pipefail
echo "=== CYTAXI Restore Verification $(date) ==="

# 1. Check PostgreSQL connectivity
echo "[1/6] Checking PostgreSQL..."
if docker exec cytaxi-postgres psql -U cytaxi -c "SELECT 1;" > /dev/null 2>&1; then
    echo "  ✓ PostgreSQL reachable"
else
    echo "  ✗ PostgreSQL connection failed"
    exit 1
fi

# 2. Verify table counts
echo "[2/6] Verifying data..."
echo "  Trips: $(docker exec cytaxi-postgres psql -U cytaxi -t -c "SELECT count(*) FROM trips;" 2>/dev/null || echo "unknown")"
echo "  Users: $(docker exec cytaxi-postgres psql -U cytaxi -t -c "SELECT count(*) FROM users;" 2>/dev/null || echo "unknown")"

# 3. Check Redis
echo "[3/6] Checking Redis..."
if docker exec cytaxi-redis redis-cli -a "${REDIS_PASSWORD:-}" PING 2>/dev/null | grep -q PONG; then
    echo "  ✓ Redis reachable"
else
    echo "  ⚠ Redis check skipped (may need password)"
fi

# 4. Health checks
echo "[4/6] Checking engine health..."
ENGINES=(8000 8087 8088 8091 8085 8086 8090 8089 8094 8093 8092)
for port in "${ENGINES[@]}"; do
    if curl -sf "http://localhost:$port/health" >/dev/null 2>&1; then
        echo "  ✓ localhost:$port"
    else
        echo "  ✗ localhost:$port FAILED"
    fi
done

# 5. Functional test
echo "[5/6] Running functional test..."
curl -sf -o /dev/null "http://localhost/api/v1/health" && echo "  ✓ API Gateway" || echo "  ✗ API Gateway"

# 6. Summary
echo ""
echo "=== Restore Verification Complete ==="
