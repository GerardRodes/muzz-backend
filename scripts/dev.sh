#!/usr/bin/env bash

basedir=$(dirname "$0")/..

source $basedir/scripts/env.sh

MUZZ_DB_ADDR=localhost:$MUZZ_HOST_DB_PORT
MUZZ_KV_ADDR=localhost:$MUZZ_HOST_KV_PORT
MUZZ_HTTP_PORT=$MUZZ_HOST_HTTP_PORT

go run ./cmd/muzz/main.go
