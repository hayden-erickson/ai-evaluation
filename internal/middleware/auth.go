package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
)

// JWTSecret is the secret key for signing JWT tokens
// In production, this should be loaded from environment variables
var JWTSecret = []byte("your-secret-key-change-this-in-production")

// GenerateToken generates a JWT token for a user
func GenerateToken(userID string) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// AuthMiddleware validates JWT tokens and adds user ID to context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("ERROR: Missing Authorization header from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("ERROR: Invalid Authorization header format from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTSecret, nil
		})

		if err != nil {
			log.Printf("ERROR: Failed to parse token from %s: %v", r.RemoteAddr, err)
			respondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Printf("ERROR: Invalid token claims from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		// Extract user ID
		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Printf("ERROR: Missing user_id in token from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "invalid token: missing user_id")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Call next handler
		next(w, r.WithContext(ctx))
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next(w, r)
	}
}

// LoggingMiddleware logs all requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("INFO: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Printf("INFO: %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// In a real application, you'd use json.Marshal here
	// For simplicity with standard library, we'll handle it in handlers
}

// RespondWithErrorJSON sends an error response in JSON format
func RespondWithErrorJSON(w http.ResponseWriter, code int, message string) {
	errResp := models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s", "message": "%s"}`, errResp.Error, errResp.Message)
}
