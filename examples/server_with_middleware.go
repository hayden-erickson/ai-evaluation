package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/config"
	"github.com/hayden-erickson/ai-evaluation/handlers"
	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/services"
)

// SetupServerWithFullMiddleware demonstrates using a complete middleware chain
func SetupServerWithFullMiddleware() (*config.Config, error) {
	// Load environment variables from .env file
	if err := config.LoadEnvFile(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Initialize dependencies
	bank := repository.NewBank()
	accessCodeService := services.NewAccessCodeService(bank)
	accessCodeHandler := handlers.NewAccessCodeHandler(accessCodeService)

	// Create a new ServeMux for better control
	mux := http.NewServeMux()

	// Setup routes
	mux.HandleFunc("/access-code/edit", accessCodeHandler.AccessCodeEditHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Build middleware chain (applied in reverse order)
	var handler http.Handler = mux

	// Add bank context middleware
	handler = middleware.BankMiddleware(bank)(handler)

	// Add CORS middleware
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000", "https://yourdomain.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
	handler = middleware.CORSWithConfig(corsConfig)(handler)

	// Add logging middleware (outermost layer)
	handler = middleware.LoggingMiddleware(handler)

	// Set the handler for the default ServeMux
	http.Handle("/", handler)

	return cfg, nil
}

// main demonstrates the full middleware chain approach
func main() {
	// Setup server with full middleware chain
	cfg, err := SetupServerWithFullMiddleware()
	if err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	// Log configuration (excluding sensitive data)
	log.Printf("Starting server with middleware on %s", cfg.GetServerAddress())
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Available endpoints:")
	log.Printf("  POST /access-code/edit - Access code management")
	log.Printf("  GET  /health          - Health check")

	// Start the HTTP server
	log.Printf("Server listening on http://%s", cfg.GetServerAddress())
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
