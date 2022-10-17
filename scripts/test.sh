#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Running tests with coverage for package internal...'

docker build -t elizabeth-tests -f ./build/test/unit.Dockerfile --build-arg PACKAGES="./internal/..." .

exit 0
