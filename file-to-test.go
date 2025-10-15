package main

import (
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/handlers"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/services"
)

// SetupServer initializes and configures the HTTP server with refactored components
func SetupServer() {
	// Initialize dependencies
	bank := repository.NewBank()
	accessCodeService := services.NewAccessCodeService(bank)
	accessCodeHandler := handlers.NewAccessCodeHandler(accessCodeService)

	// Setup HTTP routes
	http.HandleFunc("/access-code/edit", accessCodeHandler.AccessCodeEditHandler)

	// Start server (example)
	// http.ListenAndServe(":8080", nil)
}
