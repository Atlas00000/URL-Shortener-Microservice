package storage

import (
	"context"
	"fmt"
	"log"
	"os"
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
	// Get Redis configuration from environment
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "url-shortener-redis.internal"
	}
	
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	log.Printf("Connecting to Redis at %s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Testing Redis connection...")
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s:%s: %v", redisHost, redisPort, err)
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