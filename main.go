package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hayden-erickson/ai-evaluation/config"
	"github.com/hayden-erickson/ai-evaluation/handlers"
	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/service"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

func main() {
	// Get environment variables with defaults
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "habits.db"
	}

	// Initialize database
	database, err := config.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized successfully")

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(jwtSecret)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	habitRepo := repository.NewHabitRepository(database.DB)
	logRepo := repository.NewLogRepository(database.DB)

	// Initialize services
	userService := service.NewUserService(userRepo, jwtManager)
	habitService := service.NewHabitService(habitRepo)
	logService := service.NewLogService(logRepo, habitRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/users/register", userHandler.Register)
	mux.HandleFunc("/users/login", userHandler.Login)

	// Protected user routes
	mux.Handle("/users/", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route to the appropriate handler based on the method
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Protected habit routes
	mux.Handle("/habits", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle /habits endpoint for listing user's habits or creating a new habit
		switch r.Method {
		case http.MethodGet:
			habitHandler.GetUserHabits(w, r)
		case http.MethodPost:
			habitHandler.CreateHabit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/habits/", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a logs endpoint
		if len(r.URL.Path) > 7 && r.URL.Path[len(r.URL.Path)-5:] == "/logs" {
			// This is a habit logs endpoint: /habits/{habit_id}/logs
			switch r.Method {
			case http.MethodGet:
				logHandler.GetHabitLogs(w, r)
			case http.MethodPost:
				logHandler.CreateLog(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			// This is a single habit endpoint: /habits/{id}
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
	})))

	// Protected log routes
	mux.Handle("/logs/", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route to the appropriate handler based on the method
		switch r.Method {
		case http.MethodGet:
			logHandler.GetLog(w, r)
		case http.MethodPut:
			logHandler.UpdateLog(w, r)
		case http.MethodDelete:
			logHandler.DeleteLog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware to the mux (order: logging -> CORS -> security headers)
	handler := middleware.LoggingMiddleware(middleware.CorsMiddleware(middleware.SecurityHeadersMiddleware(mux)))

	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
