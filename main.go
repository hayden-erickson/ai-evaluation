package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/clients"
	"github.com/hayden-erickson/ai-evaluation/internal/config"
	"github.com/hayden-erickson/ai-evaluation/internal/handlers"
	"github.com/hayden-erickson/ai-evaluation/internal/middleware"
)

func main() {
	// Load configuration from .env file and environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Initialize database connection
	bank, err := clients.NewBank(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer bank.Close()
	log.Println("Database connection established")

	// Create a new mux for better middleware handling
	mux := http.NewServeMux()

	// Set up HTTP routes with middleware chain
	accessCodeHandler := http.HandlerFunc(handlers.AccessCodeEditHandler)
	middlewareChain := middleware.Chain(
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
		middleware.BankMiddleware(bank),
	)

	mux.Handle("/api/access-code/edit", middlewareChain(accessCodeHandler))

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s with middleware enabled", serverAddr)
	log.Printf("Database: %s@%s:%s/%s", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// Start the server
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
