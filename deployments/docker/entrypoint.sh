#!/bin/sh

set -e

DB_DSN="host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE"

wait_for_database() {
    echo ""
    echo "====================================="
    echo "Waiting for PostgreSQL..."
    echo "====================================="

    until pg_isready \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER"
    do
        sleep 2
    done

    echo "✓ PostgreSQL is ready."
}

run_migrations() {
    echo ""
    echo "====================================="
    echo "Running database migrations..."
    echo "====================================="

    goose \
        -dir /app/migrations \
        postgres \
        "$DB_DSN" \
        up

    echo "✓ Database migrations completed."
}

run_seeders() {
    echo ""
    echo "====================================="
    echo "Running database seeders..."
    echo "====================================="

    ./seed

    echo "✓ Database seeding completed."
}

start_api() {
    echo ""
    echo "====================================="
    echo "Starting API..."
    echo "====================================="

    exec ./api
}

wait_for_database
run_migrations
run_seeders
start_api