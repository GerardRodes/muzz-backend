#!/usr/bin/env bash

basedir=$(dirname "$0")/..

source $basedir/scripts/env.sh

migrate -path $basedir/internal/mariadb/migrations -database "mysql://root:$DB_ROOT_PASSWORD@tcp(localhost:13306)/$DB_DATABASE" $@
