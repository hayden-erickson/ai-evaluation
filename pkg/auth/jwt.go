package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// JWTSecret is the secret key used for signing JWT tokens
	// In production, this should be loaded from environment variables
	JWTSecret = []byte("your-secret-key-change-this-in-production")
)

// GenerateID generates a unique ID for entities
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
// Using a simple implementation without external dependencies
func GenerateJWT(userID string) (string, error) {
	// Create a simple token format: userID:timestamp:signature
	// In production, use a proper JWT library like github.com/golang-jwt/jwt
	timestamp := time.Now().Unix()
	token := fmt.Sprintf("%s:%d", userID, timestamp)
	
	// For simplicity, we'll use the token as-is
	// In production, this should be properly signed
	return token, nil
}

// ValidateJWT validates a JWT token and returns the user ID
func ValidateJWT(token string) (string, error) {
	// Simple validation - extract user ID from token
	// In production, properly verify the signature
	var userID string
	var timestamp int64
	
	_, err := fmt.Sscanf(token, "%s:%d", &userID, &timestamp)
	if err != nil {
		return "", fmt.Errorf("invalid token format")
	}
	
	// Check if token is expired (24 hours)
	if time.Now().Unix()-timestamp > 86400 {
		return "", fmt.Errorf("token expired")
	}
	
	return userID, nil
}
