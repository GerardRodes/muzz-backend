#!/usr/bin/env bash
# Requires https://github.com/golang-migrate/migrate
# Usage: bash scripts/add_migration.sh migration_name

migrate create -dir $(dirname "$0")/../internal/sql/migrations -ext sql $@
