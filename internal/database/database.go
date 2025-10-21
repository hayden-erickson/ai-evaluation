package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database connection and runs migrations
func InitDB(dbPath string) (*sql.DB, error) {
	// Create database file if it doesn't exist
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("INFO: Database connection established at %s", dbPath)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations executes SQL migration files
func runMigrations(db *sql.DB) error {
	// Enable foreign keys in SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	log.Println("INFO: Running database migrations...")

	// Migration files in order
	migrations := []string{
		"migrations/001_create_users_table.sql",
		"migrations/002_create_habits_table.sql",
		"migrations/003_create_logs_table.sql",
	}

	// Execute each migration
	for _, migrationFile := range migrations {
		log.Printf("INFO: Running migration: %s", migrationFile)

		// Read migration file
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationFile, err)
		}

		log.Printf("INFO: Migration %s completed successfully", migrationFile)
	}

	log.Println("INFO: All migrations completed successfully")
	return nil
}
