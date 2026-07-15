#!/bin/bash
set -euo pipefail

echo "=== CYTAXI v1.0.0-rc1 Deployment ==="
echo ""

VERSION="v1.0.0-rc1"
COMPOSE_FILE="docker-compose.prod.yml"
HEALTH_URL="http://localhost:8000/health"
TIMEOUT=120

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "ERROR: Docker not found"; exit 1; }
command -v docker compose >/dev/null 2>&1 || { echo "ERROR: Docker Compose not found"; exit 1; }

# Backup existing deployment
if docker compose -f "$COMPOSE_FILE" ps --services 2>/dev/null | grep -q .; then
    echo "[1/6] Backing up database..."
    docker exec cytaxi-postgres pg_dump -U cytaxi cytaxi > "backup_$(date +%Y%m%d_%H%M%S).sql" 2>/dev/null || true
fi

# Pull latest images
echo "[2/6] Building Docker images..."
docker compose -f "$COMPOSE_FILE" build --pull

# Start services
echo "[3/6] Starting services..."
docker compose -f "$COMPOSE_FILE" up -d

# Wait for health
echo "[4/6] Waiting for services to be healthy..."
elapsed=0
until curl -sf "$HEALTH_URL" >/dev/null 2>&1; do
    if [ "$elapsed" -ge "$TIMEOUT" ]; then
        echo "ERROR: Health check timeout after ${TIMEOUT}s"
        docker compose -f "$COMPOSE_FILE" logs --tail=50
        exit 1
    fi
    sleep 3
    elapsed=$((elapsed + 3))
    echo "  ... ${elapsed}s"
done

echo "[5/6] All services healthy!"

# Run smoke tests
echo "[6/6] Running smoke tests..."
ENGINES=(8087 8088 8089 8091 8085 8086 8090 8094 8093 8092)
for port in "${ENGINES[@]}"; do
    if curl -sf "http://localhost:$port/health" >/dev/null 2>&1; then
        echo "  ✓ localhost:$port"
    else
        echo "  ✗ localhost:$port FAILED"
    fi
done

echo ""
echo "=== Deployment Complete ==="
echo "API Gateway:  http://localhost:8000"
echo "Health:       http://localhost:8000/health"
echo "MiniWeb:      http://localhost:3000"
echo "Driver Portal: http://localhost:3001"
echo "=== CYTAXI $VERSION ==="
