#!/bin/sh
set -e

echo "Waiting for database:"

until PGPASSWORD="$POSTGRES_PASSWORD" psql -h db -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1" >/dev/null 2>&1; do
  sleep 1
done

echo "Running migrations:"
goose -dir ./migrations postgres "$DATABASE_URL" up

echo "Starting server:"
exec ./server