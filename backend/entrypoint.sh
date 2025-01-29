#!/bin/bash
set -e

echo "Waiting for Postgres to be ready..."
./wait-for-postgres.sh db

migrate -path ./schema -database 'postgres://postgres:qwerty@db:5432/postgres?sslmode=disable' force 1

echo "Applying down migrations..."
migrate -path ./schema -database 'postgres://postgres:qwerty@db:5432/postgres?sslmode=disable' down -all

echo "Applying up migrations..."
migrate -path ./schema -database 'postgres://postgres:qwerty@db:5432/postgres?sslmode=disable' up

echo "Starting the application..."
exec ./main
