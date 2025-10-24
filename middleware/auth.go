package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/utils"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// UserIDKey is the context key for the user ID
const UserIDKey contextKey = "user_id"

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtManager *utils.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Missing Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Println("Invalid Authorization header format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add the user ID to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
