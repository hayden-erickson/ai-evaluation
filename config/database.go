package config

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the database connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection and runs migrations
func NewDatabase(dbPath string) (*Database, error) {
	// Create database file if it doesn't exist
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	database := &Database{DB: db}

	// Run migrations
	if err := database.runMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// runMigrations runs all SQL migration files in the migrations directory
func (d *Database) runMigrations() error {
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Println("No migrations directory found, skipping migrations")
		return nil
	}

	// Read all migration files
	var migrationFiles []string
	err := filepath.WalkDir(migrationsDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only process .sql files
		if !entry.IsDir() && filepath.Ext(path) == ".sql" {
			migrationFiles = append(migrationFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migration files to ensure they run in order
	sort.Strings(migrationFiles)

	// Execute each migration
	for _, file := range migrationFiles {
		log.Printf("Running migration: %s", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute the migration
		if _, err := d.DB.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
