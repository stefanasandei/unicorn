package config

import (
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Port        string
	LogLevel    string

	JWTSecret       string
	TokenExpiration time.Duration

	// Debug API token for Swagger UI testing
	DebugAPIToken string

	// Service URLs
	LambdaURL string
}

// New creates a new Config instance
func New() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		TokenExpiration: time.Duration(time.Hour * 24),
		JWTSecret:       getEnv("JWTSecret", "lmao"),

		// Service URLs
		LambdaURL: getEnv("LAMBDA_URL", "http://localhost:8081"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
