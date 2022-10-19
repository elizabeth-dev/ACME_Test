#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Generating mocks for project interfaces using Mockery on Docker...'

docker run -u "$(id -u)" -v "$PWD":/src -w /src vektra/mockery --all --output ./test/mocks

echo -e '\033[0;32m==>\033[0m Done!'

exit 0
