package config

import (
	"errors"
	"os"
	"strconv"
)

// AppConfig holds all application configuration loaded from environment variables
type AppConfig struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	APISecret         string
	TokenHourLifespan int
	CloudinaryURL     string
	GinMode           string
	Port              string
}

// LoadConfig reads all environment variables and validates required fields.
// It returns an error immediately if a required variable is missing.
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     getEnvOrDefault("DB_SSLMODE", "disable"),
		APISecret:     os.Getenv("API_SECRET"),
		CloudinaryURL: os.Getenv("CLOUDINARY_URL"),
		GinMode:       getEnvOrDefault("GIN_MODE", "debug"),
		Port:          getEnvOrDefault("PORT", "8080"),
	}

	// Parse token lifespan
	tokenLifespan, err := strconv.Atoi(getEnvOrDefault("TOKEN_HOUR_LIFESPAN", "24"))
	if err != nil {
		return nil, errors.New("TOKEN_HOUR_LIFESPAN must be a valid integer")
	}
	cfg.TokenHourLifespan = tokenLifespan

	// Validate required fields
	if cfg.DBHost == "" {
		return nil, errors.New("DB_HOST is required")
	}
	if cfg.DBPort == "" {
		return nil, errors.New("DB_PORT is required")
	}
	if cfg.DBUser == "" {
		return nil, errors.New("DB_USER is required")
	}
	if cfg.DBPassword == "" {
		return nil, errors.New("DB_PASSWORD is required")
	}
	if cfg.DBName == "" {
		return nil, errors.New("DB_NAME is required")
	}
	if cfg.APISecret == "" {
		return nil, errors.New("API_SECRET is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
