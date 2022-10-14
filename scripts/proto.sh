#!/usr/bin/env bash

echo -e '\033[0;34m==>\033[0m Compiling protocol buffers with Docker...'

docker run --rm -u "$(id -u)"  \
    -v"${PWD}"/api/proto/v1:/source  \
    -v"${PWD}"/pkg/api/v1:/output  \
    -w/source jaegertracing/protobuf  \
      --experimental_allow_proto3_optional \
      --proto_path=/source \
      --go_out=paths=source_relative,plugins=grpc:/output \
      -I/usr/include/google/protobuf  \
      /source/*

echo -e '\033[0;32m==>\033[0m Done!'

exit 0
