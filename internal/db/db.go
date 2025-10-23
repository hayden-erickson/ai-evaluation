package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"

	"github.com/hayden-erickson/ai-evaluation/internal/config"
)

// Open opens a SQLite database using modernc.org/sqlite driver.
func Open(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", cfg.DBPath)
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	log.Printf("connected to sqlite at %s", cfg.DBPath)
	return database, nil
}
