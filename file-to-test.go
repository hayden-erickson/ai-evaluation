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

// SetupServer initializes and configures the HTTP server with refactored components
func SetupServer() (*config.Config, error) {
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
	bank, err := repository.NewBank(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}
	accessCodeService := services.NewAccessCodeService(bank)
	accessCodeHandler := handlers.NewAccessCodeHandler(accessCodeService)

	// Setup HTTP routes with middleware
	// Method 1: Using BankHandlerFunc wrapper
	http.HandleFunc("/access-code/edit", middleware.BankHandlerFunc(bank, accessCodeHandler.AccessCodeEditHandler))

	// Method 2: Using a mux with middleware chain (alternative approach)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/access-code/edit", accessCodeHandler.AccessCodeEditHandler)
	//
	// // Apply middleware chain
	// handler := middleware.LoggingMiddleware(
	// 	middleware.CORSMiddleware(
	// 		middleware.BankMiddleware(bank)(mux),
	// 	),
	// )
	// http.Handle("/", handler)

	return cfg, nil
}

// RunServer demonstrates how to use the refactored server with environment configuration
// Uncomment and rename to main() to run this server instead of the original
func RunServer() {
	// Setup server and load configuration
	cfg, err := SetupServer()
	if err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	// Log configuration (excluding sensitive data)
	log.Printf("Starting server on %s", cfg.GetServerAddress())
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Log Level: %s", cfg.LogLevel)

	// Start the HTTP server
	log.Printf("Server listening on http://%s", cfg.GetServerAddress())
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
