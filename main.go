package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/config"
	"github.com/hayden-erickson/ai-evaluation/handlers"
	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := repository.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("Database initialized successfully at %s", cfg.DBPath)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB())
	habitRepo := repository.NewHabitRepository(db.DB())
	logRepo := repository.NewLogRepository(db.DB())

	// Initialize services
	userService := service.NewUserService(userRepo)
	habitService := service.NewHabitService(habitRepo)
	logService := service.NewLogService(logRepo, habitRepo)

	// Initialize JWT manager
	jwtManager := middleware.NewJWTManager(cfg.JWTSecret)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, jwtManager)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)

	// Setup routes
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/api/users", userHandler.Register)
	mux.HandleFunc("/api/login", userHandler.Login)

	// Protected routes (authentication required)
	mux.Handle("/api/users/", jwtManager.AuthMiddleware(http.HandlerFunc(userHandler.HandleUser)))
	mux.Handle("/api/habits", jwtManager.AuthMiddleware(http.HandlerFunc(habitHandler.HandleHabits)))
	mux.Handle("/api/habits/", jwtManager.AuthMiddleware(http.HandlerFunc(handleHabitRoutes(habitHandler, logHandler))))
	mux.Handle("/api/logs", jwtManager.AuthMiddleware(http.HandlerFunc(logHandler.HandleLogs)))
	mux.Handle("/api/logs/", jwtManager.AuthMiddleware(http.HandlerFunc(logHandler.HandleLog)))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware chain
	handler := middleware.LoggingMiddleware(
		middleware.SecurityHeadersMiddleware(mux),
	)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleHabitRoutes routes requests for /api/habits/* to appropriate handlers
func handleHabitRoutes(habitHandler *handlers.HabitHandler, logHandler *handlers.LogHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a logs endpoint: /api/habits/{id}/logs
		if strings.Contains(r.URL.Path, "/logs") {
			logHandler.HandleLogs(w, r)
			return
		}

		// Otherwise, handle as a habit endpoint
		habitHandler.HandleHabit(w, r)
	}
}
