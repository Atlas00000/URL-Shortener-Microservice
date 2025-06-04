package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourusername/urlshortener/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database connections
type Database struct {
	SQLite *gorm.DB
	Redis  *redis.Client
}

// NewDatabase creates new database connections
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	// Initialize SQLite connection
	sqliteDB, err := initSQLite(cfg.SQLite)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite: %v", err)
	}

	// Initialize Redis connection
	redisClient, err := initRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %v", err)
	}

	return &Database{
		SQLite: sqliteDB,
		Redis:  redisClient,
	}, nil
}

// initSQLite initializes SQLite connection
func initSQLite(cfg config.SQLiteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// initRedis initializes Redis connection
func initRedis(cfg config.RedisConfig) (*redis.Client, error) {
	// Get Redis URL from environment
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	// Log the raw Redis URL for debugging
	log.Printf("Raw Redis URL from environment: %s", redisURL)

	// Ensure URL has proper scheme
	if !strings.HasPrefix(redisURL, "redis://") && !strings.HasPrefix(redisURL, "rediss://") {
		redisURL = "redis://" + redisURL
		log.Printf("Added redis:// scheme to URL: %s", redisURL)
	}

	// Parse Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	// Create Redis client
	client := redis.NewClient(opt)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Testing Redis connection to %s...", redisURL)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	log.Printf("Redis connection successful")

	return client, nil
}

// Close closes all database connections
func (db *Database) Close() error {
	// Close SQLite connection
	sqlDB, err := db.SQLite.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	// Close Redis connection
	if err := db.Redis.Close(); err != nil {
		return err
	}

	return nil
} 