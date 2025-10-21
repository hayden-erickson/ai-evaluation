package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var DB *sql.DB

// Connect initializes the database connection
func Connect(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	return DB.Ping()
}

// Migrate runs the database migrations
func Migrate(migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(migrationsPath, file.Name())
		migration, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("could not read migration file %s: %w", file.Name(), err)
		}

		if _, err := DB.Exec(string(migration)); err != nil {
			return fmt.Errorf("could not execute migration %s: %w", file.Name(), err)
		}
		log.Printf("Successfully applied migration %s", file.Name())
	}

	return nil
}

// InitDB initializes the database and runs migrations
func InitDB() {
	dbPath := "./habits.db"
	if err := Connect(dbPath); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	log.Println("Database connected.")

	migrationsPath := "internal/database/migrations"
	if err := Migrate(migrationsPath); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	log.Println("Migrations completed.")
}
