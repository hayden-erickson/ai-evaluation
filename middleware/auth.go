package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/utils"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
)

// AuthMiddleware validates JWT tokens and adds user ID to context
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			
			// Extract token from "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
			
			tokenString := parts[1]
			
			// Validate token
			claims, err := utils.ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				log.Printf("JWT validation error: %v", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			
			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
