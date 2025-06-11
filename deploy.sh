#!/bin/bash
set -e

echo "--- Starting deployment script ---"

DEPLOY_COMPOSE_PATH_FILE="${1:-docker-compose.yml}"

docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" down || true

echo "--- Starting MariaDB ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up -d mariadb

echo "--- Running database migrations ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" run --build --rm resik-migrate migrate up

echo "--- Starting application services ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up --build -d resik-rest resik-grpc

echo "--- Deployment script finished ---"
