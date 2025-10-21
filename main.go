package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/database"
	"github.com/hayden-erickson/ai-evaluation/internal/handlers"
	"github.com/hayden-erickson/ai-evaluation/internal/middleware"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("INFO: Starting Habit Tracker API...")

	// Initialize database
	db, err := database.InitDB("./data/habits.db")
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	logRepo := repository.NewLogRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	habitService := service.NewHabitService(habitRepo)
	logService := service.NewLogService(logRepo, habitRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)

	// Set up routes
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/users", userHandler.CreateUser)
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/health", healthCheck)

	// Protected routes (authentication required)
	mux.HandleFunc("/users/", chainMiddleware(
		handleUserRoutes(userHandler),
		middleware.AuthMiddleware,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware,
	))

	mux.HandleFunc("/habits", chainMiddleware(
		handleHabitsRoute(habitHandler),
		middleware.AuthMiddleware,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware,
	))

	mux.HandleFunc("/habits/", chainMiddleware(
		handleHabitRoutes(habitHandler, logHandler),
		middleware.AuthMiddleware,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware,
	))

	mux.HandleFunc("/logs/", chainMiddleware(
		handleLogRoutes(logHandler),
		middleware.AuthMiddleware,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware,
	))

	// Create server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("INFO: Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO: Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Server forced to shutdown: %v", err)
	}

	log.Println("INFO: Server stopped gracefully")
}

// healthCheck handles health check requests
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

// handleUserRoutes routes user-specific requests
func handleUserRoutes(handler *handlers.UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Route based on method
		switch r.Method {
		case http.MethodGet:
			handler.GetUser(w, r)
		case http.MethodPut:
			handler.UpdateUser(w, r)
		case http.MethodDelete:
			handler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleHabitsRoute routes /habits requests
func handleHabitsRoute(handler *handlers.HabitHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Route based on method
		switch r.Method {
		case http.MethodPost:
			handler.CreateHabit(w, r)
		case http.MethodGet:
			handler.GetHabits(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleHabitRoutes routes habit-specific requests
func handleHabitRoutes(habitHandler *handlers.HabitHandler, logHandler *handlers.LogHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a log route (contains /logs)
		if len(r.URL.Path) > 8 && r.URL.Path[len(r.URL.Path)-5:] == "/logs" {
			// This is /habits/{id}/logs
			switch r.Method {
			case http.MethodPost:
				logHandler.CreateLog(w, r)
			case http.MethodGet:
				logHandler.GetLogs(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// This is /habits/{id}
		switch r.Method {
		case http.MethodGet:
			habitHandler.GetHabit(w, r)
		case http.MethodPut:
			habitHandler.UpdateHabit(w, r)
		case http.MethodDelete:
			habitHandler.DeleteHabit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleLogRoutes routes log-specific requests
func handleLogRoutes(handler *handlers.LogHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Route based on method
		switch r.Method {
		case http.MethodGet:
			handler.GetLog(w, r)
		case http.MethodPut:
			handler.UpdateLog(w, r)
		case http.MethodDelete:
			handler.DeleteLog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// chainMiddleware chains multiple middleware functions
func chainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	// Apply middlewares in reverse order so they execute in the order specified
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
