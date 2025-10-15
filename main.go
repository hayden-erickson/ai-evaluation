package main

import (
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/clients"
	"github.com/hayden-erickson/ai-evaluation/internal/handlers"
)

func main() {
	// Initialize database connection
	bank, err := clients.NewBank()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer bank.Close()

	// Set up HTTP routes
	http.HandleFunc("/api/access-code/edit", handlers.AccessCodeEditHandler)

	log.Println("Server starting on :8080")
	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
