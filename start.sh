#!/bin/bash

set -e

echo "Starting DB Migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_URL" -verbose up

echo "Starting API"
exec "$@"
