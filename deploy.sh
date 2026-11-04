#!/bin/bash
set -e

export APP_VERSION=$(git describe --tags --always)
export APP_GIT_COMMIT=$(git rev-parse --short HEAD)
export APP_BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

ENV_FILE="${ENV_FILE:-.env}"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "Error: Environment file '${ENV_FILE}' not found"
  exit 1
fi

DATABASE_HOST=$(
  sed -n 's/^[[:space:]]*DATABASE_HOST[[:space:]]*=[[:space:]]*//p' "${ENV_FILE}" |
  tail -n 1 |
  sed -E 's/[[:space:]]+#.*$//; s/^["'\'']|["'\'']$//g; s/[[:space:]]+$//'
)

case "${DATABASE_HOST}" in
  mariadb)
    DB_SERVICE="mariadb"
    ;;
  postgres)
    DB_SERVICE="postgres"
    ;;
  *)
    echo "Error: Unsupported DATABASE_HOST '${DATABASE_HOST}'"
    exit 1
    ;;
esac

DEPLOY_COMPOSE_PATH_FILE="docker-compose.${DB_SERVICE}.yml"

echo "--- Starting deployment script ---"

docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" down || true

echo "--- Starting ${DB_SERVICE} & redis ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up -d "${DB_SERVICE}" redis

echo "--- Running database migrations ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" run --build --rm resik-migrate migrate up

echo "--- Starting application services ---"
docker compose -f "${DEPLOY_COMPOSE_PATH_FILE}" up --build -d resik-rest resik-grpc

echo "--- Deployment script finished ---"
