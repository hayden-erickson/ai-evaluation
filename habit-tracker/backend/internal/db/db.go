package db

import (
	"log"
	"os"

	"github.com/hayden-erickson/habit-tracker/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database instance
var DB *gorm.DB

// Initialize sets up the database connection and runs migrations
func Initialize() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./habit-tracker.db"
	}

	var err error
	DB, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Enable foreign key constraints
	DB.Exec("PRAGMA foreign_keys = ON")

	// Enable WAL mode for better concurrent access
	DB.Exec("PRAGMA journal_mode = WAL")

	// Set connection pool settings
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	// Auto-migrate models
	DB.AutoMigrate(&models.User{}, &models.Habit{}, &models.Log{}, &models.Tag{})
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
