#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Running e2e tests...'

docker compose -f deploy/docker-compose.e2e.yml -p test-elizabeth-e2e up -V --force-recreate --build --abort-on-container-exit --exit-code-from e2e

exit 0
