package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/hayden-erickson/ai-evaluation/internal/config"
	"github.com/hayden-erickson/ai-evaluation/internal/handlers"
	"github.com/hayden-erickson/ai-evaluation/internal/middleware"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories (data access layer)
	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	logRepo := repository.NewLogRepository(db)

	// Initialize services (business logic layer)
	userService := service.NewUserService(userRepo)
	habitService := service.NewHabitService(habitRepo, userRepo)
	logService := service.NewLogService(logRepo, habitRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Initialize HTTP handlers
	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)
	authHandler := handlers.NewAuthHandler(authService)

	// Set up HTTP router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Protected routes - User endpoints
	mux.HandleFunc("GET /api/v1/users/{id}", middleware.Auth(cfg.JWTSecret, userHandler.GetUser))
	mux.HandleFunc("PUT /api/v1/users/{id}", middleware.Auth(cfg.JWTSecret, userHandler.UpdateUser))
	mux.HandleFunc("DELETE /api/v1/users/{id}", middleware.Auth(cfg.JWTSecret, userHandler.DeleteUser))

	// Protected routes - Habit endpoints
	mux.HandleFunc("POST /api/v1/habits", middleware.Auth(cfg.JWTSecret, habitHandler.CreateHabit))
	mux.HandleFunc("GET /api/v1/habits", middleware.Auth(cfg.JWTSecret, habitHandler.ListHabits))
	mux.HandleFunc("GET /api/v1/habits/{id}", middleware.Auth(cfg.JWTSecret, habitHandler.GetHabit))
	mux.HandleFunc("PUT /api/v1/habits/{id}", middleware.Auth(cfg.JWTSecret, habitHandler.UpdateHabit))
	mux.HandleFunc("DELETE /api/v1/habits/{id}", middleware.Auth(cfg.JWTSecret, habitHandler.DeleteHabit))

	// Protected routes - Log endpoints
	mux.HandleFunc("POST /api/v1/logs", middleware.Auth(cfg.JWTSecret, logHandler.CreateLog))
	mux.HandleFunc("GET /api/v1/logs", middleware.Auth(cfg.JWTSecret, logHandler.ListLogs))
	mux.HandleFunc("GET /api/v1/logs/{id}", middleware.Auth(cfg.JWTSecret, logHandler.GetLog))
	mux.HandleFunc("PUT /api/v1/logs/{id}", middleware.Auth(cfg.JWTSecret, logHandler.UpdateLog))
	mux.HandleFunc("DELETE /api/v1/logs/{id}", middleware.Auth(cfg.JWTSecret, logHandler.DeleteLog))

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware chain
	handler := middleware.SecurityHeaders(
		middleware.Logging(
			middleware.CORS(mux),
		),
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// initDB initializes the SQLite database connection
func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	log.Printf("Database connection established: %s", cfg.DBPath)
	return db, nil
}

// runMigrations executes SQL migration files
func runMigrations(db *sql.DB) error {
	migrationDir := "migrations"
	
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(migrationDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Split content by semicolon and execute each statement
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}
		}

		log.Printf("Applied migration: %s", file.Name())
	}

	return nil
}
