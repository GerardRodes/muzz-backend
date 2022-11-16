#!/usr/bin/env bash

basedir=$(dirname "$0")/..

source $basedir/scripts/env.sh
MUZZ_DB_ADDR=localhost:$MUZZ_DEV_DB_PORT MUZZ_KV_ADDR=localhost:$MUZZ_DEV_KV_PORT go run ./cmd/muzz/main.go
