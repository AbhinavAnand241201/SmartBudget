package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Port           string
	DatabaseURL    string
	SupabaseKey    string
	SendGridAPIKey string
	HuggingFaceKey string
	Environment    string
	AIServiceURL   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smartbudget?sslmode=disable"),
		SupabaseKey:    getEnv("SUPABASE_KEY", ""),
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		HuggingFaceKey: getEnv("HUGGINGFACE_KEY", ""),
		Environment:    getEnv("ENVIRONMENT", "development"),
		AIServiceURL:   getEnv("AI_SERVICE_URL", "http://localhost:5001"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
