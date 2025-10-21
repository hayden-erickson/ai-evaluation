package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "habits.db"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret-change-in-production"),
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
