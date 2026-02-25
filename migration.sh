#!/bin/bash
set -e

echo "[migrator] PG_HOST=$PG_HOST PG_PORT=$PG_PORT DB=$PG_DATABASE_NAME USER=$PG_USER"


until pg_isready -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DATABASE_NAME" >/dev/null 2>&1; do
  echo "[migrator] waiting for postgres..."
  sleep 1
done

MIGRATION_DSN="host=$PG_HOST port=$PG_PORT user=$PG_USER password=$PG_PASSWORD dbname=$PG_DATABASE_NAME sslmode=$PG_SSLMODE"

echo "[migrator] running goose up..."
/bin/goose -dir "/migrations" postgres "$MIGRATION_DSN" up -v
echo "[migrator] done"