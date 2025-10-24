package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	DBPath     string
	JWTSecret  string
	ServerPort string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		DBPath:     getEnv("DB_PATH", "./habits.db"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		ServerPort: getEnv("PORT", "8080"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
