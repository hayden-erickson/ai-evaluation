package db

import (
	"database/sql"
	"embed"
	"log"
	"sort"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate applies all SQL files in the migrations directory in lex order.
func Migrate(database *sql.DB) error {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}
	// Sort by name to apply in order
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		b, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		log.Printf("applying migration %s", name)
		if _, err := database.Exec(string(b)); err != nil {
			return err
		}
	}
	return nil
}
