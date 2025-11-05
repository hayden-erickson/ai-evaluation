package config

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/go-sql-driver/mysql"
)

// Database wraps the database connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection and runs migrations
func NewDatabase(dsn string) (*Database, error) {
	// Connect to MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

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
