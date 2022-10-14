#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Running all microservices locally with Docker Compose...'

docker compose -f deploy/docker-compose.yml -p test-elizabeth up --build -d

echo -e '\033[0;32m==>\033[0m Running! Compose project name: test-elizabeth'
echo -e '\033[0;32m==>\033[0m Ephemeral MongoDB is available at mongodb://root:example@localhost:27017/test?authSource=admin'
echo -e '\033[0;32m==>\033[0m User microservice gRPC is available at localhost:8080'

exit 0
