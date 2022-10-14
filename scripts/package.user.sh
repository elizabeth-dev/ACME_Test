#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Packaging User microservice in a container...'

docker build -t elizabeth-user . -f ./build/package/user.Dockerfile

echo -e '\033[0;32m==>\033[0m Done! The image is tagged as elizabeth-user'

exit 0
