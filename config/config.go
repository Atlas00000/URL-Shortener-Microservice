package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig
	BaseURL  string
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	SQLite SQLiteConfig
	Redis  RedisConfig
}

// SQLiteConfig represents SQLite configuration
type SQLiteConfig struct {
	Path string
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// GetDSN returns the SQLite connection string
func (c SQLiteConfig) GetDSN() string {
	return c.Path
}

// GetRedisURL returns the Redis URL
func (c RedisConfig) GetRedisURL() string {
	log.Printf("Redis URL from config: %s", c.URL)
	return c.URL
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Ensure data directory exists
	dataDir := getEnv("DATA_DIR", "./data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	redisURL := getEnv("REDIS_URL", "redis://redis:6379/0")
	log.Printf("Loading Redis URL from environment: %s", redisURL)

	return &Config{
		Database: DatabaseConfig{
			SQLite: SQLiteConfig{
				Path: filepath.Join(dataDir, "urlshortener.db"),
			},
			Redis: RedisConfig{
				URL:      redisURL,
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       0,
			},
		},
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Environment variable %s found: %s", key, value)
		return value
	}
	log.Printf("Environment variable %s not found, using default: %s", key, defaultValue)
	return defaultValue
} 