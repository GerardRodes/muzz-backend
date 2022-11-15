#!/usr/bin/env bash
# Requires https://github.com/golang-migrate/migrate

basedir=$(dirname "$0")/..

migrate create -dir $basedir/internal/sql/migrations -ext sql $@
