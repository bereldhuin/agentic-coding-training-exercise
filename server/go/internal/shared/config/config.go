package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Port         int
	DatabasePath string
	Environment  string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	port := getEnvAsInt("PORT", 8081)
	dbPath := getEnv("DATABASE_PATH", "../database/db.sqlite")
	env := getEnv("ENV", "development")

	// Resolve database path to absolute path
	absDbPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve database path: %w", err)
	}

	// Validate configuration
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d (must be between 1 and 65535)", port)
	}

	if absDbPath == "" {
		return nil, fmt.Errorf("database path cannot be empty")
	}

	if env != "development" && env != "production" && env != "test" {
		return nil, fmt.Errorf("invalid environment: %s (must be development, production, or test)", env)
	}

	return &Config{
		Port:         port,
		DatabasePath: absDbPath,
		Environment:  env,
	}, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

// GetAddr returns the address to bind the server to
func (c *Config) GetAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}
