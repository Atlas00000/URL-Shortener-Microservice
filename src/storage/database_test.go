package storage

import (
	"testing"

	"github.com/yourusername/urlshortener/config"
)

func TestDatabaseConnections(t *testing.T) {
	// Load configuration
	cfg := config.LoadDatabaseConfig()

	// Create database connections
	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create database connections: %v", err)
	}
	defer db.Close()

	// Test PostgreSQL connection
	sqlDB, err := db.Postgres.DB()
	if err != nil {
		t.Fatalf("Failed to get PostgreSQL connection: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Test Redis connection
	ctx := db.Redis.Context()
	if err := db.Redis.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}
} 