package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration values for the application
type Config struct {
	// Server Configuration
	Port string
	Host string

	// Database Configuration
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	// External Services
	CommandCenterURL    string
	CommandCenterAPIKey string

	// Security
	JWTSecret     string
	EncryptionKey string

	// Environment
	Environment string
	LogLevel    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		// Server defaults
		Port: getEnvWithDefault("PORT", "8080"),
		Host: getEnvWithDefault("HOST", "localhost"),

		// Database defaults
		DBHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DBName:     getEnvWithDefault("DB_NAME", "access_control"),
		DBUser:     getEnvWithDefault("DB_USER", ""),
		DBPassword: getEnvWithDefault("DB_PASSWORD", ""),

		// External services
		CommandCenterURL:    getEnvWithDefault("COMMAND_CENTER_URL", ""),
		CommandCenterAPIKey: getEnvWithDefault("COMMAND_CENTER_API_KEY", ""),

		// Security
		JWTSecret:     getEnvWithDefault("JWT_SECRET", ""),
		EncryptionKey: getEnvWithDefault("ENCRYPTION_KEY", ""),

		// Environment
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
	}

	// Parse DB port
	dbPortStr := getEnvWithDefault("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}
	config.DBPort = dbPort

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks that required configuration values are present
func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	return nil
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Host + ":" + c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
