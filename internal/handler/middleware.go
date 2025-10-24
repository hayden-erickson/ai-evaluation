package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/pkg/auth"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for storing user ID
	UserIDKey contextKey = "userID"
)

// AuthMiddleware validates JWT tokens and adds user ID to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Error: Missing Authorization header from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Error: Invalid Authorization header format from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		userID, err := auth.ValidateJWT(token)
		if err != nil {
			log.Printf("Error: Invalid token from %s: %v", r.RemoteAddr, err)
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
