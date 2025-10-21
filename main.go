package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/hayden-erickson/ai-evaluation/handlers"
	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/service"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

const (
	defaultDBPath = "./habits.db"
	defaultPort   = "8080"
)

func main() {
	// Get configuration from environment or use defaults
	dbPath := getEnv("DB_PATH", defaultDBPath)
	port := getEnv("PORT", defaultPort)
	
	// Generate or get JWT secret key
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		var err error
		jwtSecret, err = utils.GenerateSecretKey()
		if err != nil {
			log.Fatalf("Failed to generate JWT secret: %v", err)
		}
		log.Println("Generated JWT secret key. Set JWT_SECRET environment variable for production.")
	}
	
	// Initialize database
	db, err := initDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	logRepo := repository.NewLogRepository(db)
	
	// Initialize services
	userService := service.NewUserService(userRepo)
	habitService := service.NewHabitService(habitRepo)
	logService := service.NewLogService(logRepo, habitRepo)
	
	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, jwtSecret)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)
	
	// Setup routes
	mux := http.NewServeMux()
	
	// Public routes (no authentication required)
	mux.HandleFunc("/api/register", userHandler.Register)
	mux.HandleFunc("/api/login", userHandler.Login)
	
	// Protected routes (authentication required)
	authMux := http.NewServeMux()
	authMux.HandleFunc("/api/users/", userHandler.HandleUser)
	authMux.HandleFunc("/api/habits", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/habits (list and create)
		if r.URL.Path == "/api/habits" {
			if r.Method == http.MethodGet {
				habitHandler.List(w, r)
			} else if r.Method == http.MethodPost {
				habitHandler.Create(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	authMux.HandleFunc("/api/habits/", habitHandler.HandleHabit)
	authMux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/logs (list and create)
		if r.URL.Path == "/api/logs" {
			if r.Method == http.MethodGet {
				logHandler.ListByHabit(w, r)
			} else if r.Method == http.MethodPost {
				logHandler.Create(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	authMux.HandleFunc("/api/logs/", logHandler.HandleLog)
	
	// Wrap authenticated routes with auth middleware
	mux.Handle("/api/users/", middleware.AuthMiddleware(jwtSecret)(authMux))
	mux.Handle("/api/habits", middleware.AuthMiddleware(jwtSecret)(authMux))
	mux.Handle("/api/habits/", middleware.AuthMiddleware(jwtSecret)(authMux))
	mux.Handle("/api/logs", middleware.AuthMiddleware(jwtSecret)(authMux))
	mux.Handle("/api/logs/", middleware.AuthMiddleware(jwtSecret)(authMux))
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply middleware chain
	handler := middleware.LoggingMiddleware(
		middleware.SecurityHeadersMiddleware(
			middleware.CORSMiddleware(mux),
		),
	)
	
	// Start server
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	log.Printf("Database: %s", dbPath)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// initDB initializes the database connection
func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	
	return db, nil
}

// runMigrations executes all migration files
func runMigrations(db *sql.DB) error {
	migrationDir := "./migrations"
	
	// Get all migration files
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}
	
	// Execute migrations in order
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		
		migrationPath := filepath.Join(migrationDir, file.Name())
		log.Printf("Running migration: %s", file.Name())
		
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}
		
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}
	
	log.Println("Migrations completed successfully")
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
