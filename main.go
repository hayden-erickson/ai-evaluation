package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/handler"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
	"github.com/hayden-erickson/ai-evaluation/pkg/database"
)

func main() {
	// Initialize database
	db, err := database.InitDB("./data/habits.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, "./migrations"); err != nil {
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
	userHandler := handler.NewUserHandler(userService)
	habitHandler := handler.NewHabitHandler(habitService)
	logHandler := handler.NewLogHandler(logService)

	// Setup router
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/api/users", userHandler.HandleUsers)
	mux.HandleFunc("/api/users/", userHandler.HandleUsers)
	mux.HandleFunc("/api/login", userHandler.HandleLogin)

	// Protected routes (authentication required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/habits", habitHandler.HandleHabits)
	protectedMux.HandleFunc("/api/habits/", habitHandler.HandleHabits)
	protectedMux.HandleFunc("/api/logs", logHandler.HandleLogs)
	protectedMux.HandleFunc("/api/logs/", logHandler.HandleLogs)

	// Apply middleware
	var mainHandler http.Handler = mux
	mainHandler = handler.SecurityHeadersMiddleware(mainHandler)
	mainHandler = handler.LoggingMiddleware(mainHandler)

	// Wrap protected routes with auth middleware
	protectedHandler := handler.AuthMiddleware(protectedMux)

	// Combine handlers
	finalMux := http.NewServeMux()
	finalMux.Handle("/api/habits", protectedHandler)
	finalMux.Handle("/api/habits/", protectedHandler)
	finalMux.Handle("/api/logs", protectedHandler)
	finalMux.Handle("/api/logs/", protectedHandler)
	finalMux.Handle("/api/users", mainHandler)
	finalMux.Handle("/api/users/", mainHandler)
	finalMux.Handle("/api/login", mainHandler)

	// Apply global middleware to final mux
	var finalHandler http.Handler = finalMux
	finalHandler = handler.SecurityHeadersMiddleware(finalHandler)
	finalHandler = handler.LoggingMiddleware(finalHandler)

	// Create server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
