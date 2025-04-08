#!/bin/bash
source .env

export MIGRATIION_DSN="host=pg port=5432 dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATIION_DSN}" up -v