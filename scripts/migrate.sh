#!/usr/bin/env bash

basedir=$(dirname "$0")/..

source $basedir/scripts/env.sh

migrate -path $basedir/internal/mariadb/migrations -database "mysql://root:$MUZZ_DB_ROOT_PASSWORD@tcp(localhost:13306)/$MUZZ_DB_NAME" $@
