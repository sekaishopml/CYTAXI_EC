#!/bin/bash
set -euo pipefail

echo "=== CYTAXI Beta Deployment ==="
echo "IP: 64.176.219.221"
echo ""

COMPOSE_FILE="docker-compose.beta.yml"
TIMEOUT=180

command -v docker >/dev/null 2>&1 || { echo "ERROR: Docker not found"; exit 1; }

echo "[1/5] Building images..."
docker compose -f "$COMPOSE_FILE" build --pull

echo "[2/5] Starting all services..."
docker compose -f "$COMPOSE_FILE" up -d

echo "[3/5] Waiting for services..."
elapsed=0
until curl -sf "http://localhost:80/health" >/dev/null 2>&1; do
    if [ "$elapsed" -ge "$TIMEOUT" ]; then
        echo "ERROR: Timeout"
        docker compose -f "$COMPOSE_FILE" logs --tail=30
        exit 1
    fi
    sleep 5
    elapsed=$((elapsed + 5))
    echo "  ... ${elapsed}s"
done

echo "[4/5] Services healthy!"

echo "[5/5] Verify endpoints..."
echo "  MiniWeb:     http://64.176.219.221/"
echo "  Driver:      http://64.176.219.221/driver"
echo "  Dashboard:   http://64.176.219.221/admin"
echo "  API:         http://64.176.219.221/api/v1/health"
echo "  Swagger:     http://64.176.219.221/docs"
echo "  Health:      http://64.176.219.221/health"
echo ""
echo "=== Beta Deployment Complete ==="
echo "Open http://64.176.219.221 in your browser"
