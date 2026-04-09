#!/bin/bash
set -e

export APP_VERSION=$(git describe --tags --always)
export APP_GIT_COMMIT=$(git rev-parse --short HEAD)
export APP_BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "--- Starting deployment script ---"

DEPLOY_COMPOSE_PATH_FILE="${1:-docker-compose.yml}"

docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" down || true

echo "--- Starting MariaDB & Redis ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up -d mariadb redis

echo "--- Running database migrations ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" run --build --rm resik-migrate migrate up

echo "--- Starting application services ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up --build -d resik-rest resik-grpc

echo "--- Deployment script finished ---"
