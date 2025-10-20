#!/bin/bash

# Database Migration Script
# Run this script to apply database migrations

set -e

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-habits_user}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME:-habits_db}
MIGRATIONS_DIR=${MIGRATIONS_DIR:-./migrations}

echo "=== Running Database Migrations ==="
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo ""

# Check if password is set
if [ -z "$DB_PASSWORD" ]; then
    echo "Error: DB_PASSWORD environment variable is not set"
    exit 1
fi

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Run each migration file in order
for migration_file in $(ls -1 $MIGRATIONS_DIR/*.sql | sort); do
    echo "Applying migration: $(basename $migration_file)"
    mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME < $migration_file
    
    if [ $? -eq 0 ]; then
        echo "✓ Success: $(basename $migration_file)"
    else
        echo "✗ Failed: $(basename $migration_file)"
        exit 1
    fi
done

echo ""
echo "=== All migrations applied successfully ==="
