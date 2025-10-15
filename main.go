package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hayden-erickson/ai-evaluation/contextutil"
	"github.com/hayden-erickson/ai-evaluation/database"
	"github.com/hayden-erickson/ai-evaluation/handlers"
)

// SetupServer sets up the HTTP server with routes and middleware
func SetupServer(bank *database.Bank) *http.ServeMux {
	mux := http.NewServeMux()
	
	// Add middleware to inject bank into context
	bankMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contextutil.NewBankContext(r.Context(), bank)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	
	// Register handlers with middleware
	mux.Handle("/api/access-code/edit", bankMiddleware(http.HandlerFunc(handlers.AccessCodeEditHandler)))
	
	return mux
}

func main() {
	// Initialize database connection
	config := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     3306,
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_DATABASE", "ai_evaluation"),
	}
	
	bank, err := database.NewBank(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer bank.Close()
	
	// Test database connection
	if err := bank.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	log.Println("Database connection established successfully")
	
	server := SetupServer(bank)
	
	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on :%s\n", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
