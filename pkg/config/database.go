package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes a database connection using the provided AppConfig.
// It prefers DATABASE_URL if available, otherwise builds DSN from individual fields.
// Returns an error instead of panicking.
func ConnectDatabase(cfg *AppConfig) (*gorm.DB, error) {
	var dsn string

	if cfg.DatabaseURL != "" {
		// Use the full connection string (e.g. from Render)
		dsn = cfg.DatabaseURL
	} else {
		// Build DSN from individual fields (local development)
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
			cfg.DBSSLMode,
		)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Keep global reference for backward compatibility (AutoMigrate, etc.)
	DB = database
	return database, nil
}
