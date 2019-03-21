#!/usr/bin/env bash
cd `git rev-parse --show-toplevel`
mkdir -p ./pkg/api/v1
protoc --proto_path=./api/proto/v1 --go_out=plugins=grpc:./pkg/api/v1 ./api/proto/v1/account.proto ./api/proto/v1/balance.proto
