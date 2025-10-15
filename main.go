package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/contextutil"
	"github.com/hayden-erickson/ai-evaluation/database"
	"github.com/hayden-erickson/ai-evaluation/handlers"
)

// SetupServer sets up the HTTP server with routes and middleware
func SetupServer() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Add middleware to inject bank into context
	bankMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bank := &database.Bank{}
			ctx := contextutil.NewBankContext(r.Context(), bank)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	
	// Register handlers with middleware
	mux.Handle("/api/access-code/edit", bankMiddleware(http.HandlerFunc(handlers.AccessCodeEditHandler)))
	
	return mux
}

func main() {
	server := SetupServer()
	
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
