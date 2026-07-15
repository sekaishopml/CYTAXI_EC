#!/bin/bash
set -euo pipefail

BACKUP_DIR="/backups/cytaxi"
RETENTION_DAYS=7
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

echo "=== CYTAXI Backup $(date) ==="

# PostgreSQL
echo "[1/2] Backing up PostgreSQL..."
docker exec cytaxi-postgres pg_dump -U cytaxi cytaxi | gzip > "$BACKUP_DIR/postgres_${TIMESTAMP}.sql.gz"

# Redis
echo "[2/2] Backing up Redis..."
docker exec cytaxi-redis redis-cli BGSAVE

# Cleanup old backups
find "$BACKUP_DIR" -name "*.sql.gz" -mtime "+${RETENTION_DAYS}" -delete

echo "=== Backup Complete: $BACKUP_DIR ==="
ls -lh "$BACKUP_DIR/postgres_${TIMESTAMP}.sql.gz"
