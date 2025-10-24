package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddr string
	DBPath     string
	JWTSecret  string
	JWTIssuer  string
}

// Load reads configuration from environment with sensible defaults.
func Load() *Config {
	cfg := &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		DBPath:     getEnv("SQLITE_PATH", "./app.db"),
		JWTSecret:  getEnv("JWT_SECRET", "dev-secret-change"),
		JWTIssuer:  getEnv("JWT_ISSUER", "habit-api"),
	}
	if len(cfg.JWTSecret) < 16 {
		log.Printf("warning: JWT_SECRET is too short for production")
	}
	return cfg
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
