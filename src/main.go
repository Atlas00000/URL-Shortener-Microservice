package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/urlshortener/src/api"
	"github.com/yourusername/urlshortener/src/logger"
	"github.com/yourusername/urlshortener/src/services"
	"github.com/yourusername/urlshortener/src/storage"
	"github.com/yourusername/urlshortener/config"
	"github.com/yourusername/urlshortener/src/api/handlers"
	"github.com/yourusername/urlshortener/src/geo"
)

func main() {
	// Initialize logger
	if err := logger.Init(true); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("Starting URL shortener service")

	// Load config
	dbConfig, err := config.Load()
	if err != nil {
		logger.LogError(err, "Failed to load config", nil)
		os.Exit(1)
	}
	// Ensure DataDir is set
	if dbConfig.DataDir == "" {
		dbConfig.DataDir = "./data"
	}

	// Initialize database
	db, err := storage.NewDatabase(&dbConfig.Database)
	if err != nil {
		logger.LogError(err, "Failed to initialize database", nil)
		os.Exit(1)
	}

	// Run migrations
	if err := storage.RunMigrations(db.SQLite); err != nil {
		logger.LogError(err, "Failed to run migrations", nil)
		os.Exit(1)
	}

	// Initialize GeoIP service
	geoService, err := geo.NewService(dbConfig.DataDir)
	if err != nil {
		logger.LogError(err, "Failed to initialize geo service", nil)
		logger.LogInfo("Continuing without geo location features", nil)
	} else {
		defer geoService.Close()
		logger.LogInfo("GeoIP service initialized successfully", nil)
	}

	// Initialize services
	urlService := services.NewURLService(db.SQLite)
	analyticsService := services.NewAnalyticsService(db.SQLite, nil)

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(urlService, dbConfig.BaseURL)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, urlService)

	// Initialize server
	server := api.NewServer()
	server.RegisterRoutes(urlHandler, analyticsHandler)

	// Create a channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		serverErrors <- server.Start()
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.LogError(err, "Server error", nil)
		os.Exit(1)
	case <-quit:
		log.Info("Shutting down server...")
	}

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.LogError(err, "Server forced to shutdown", nil)
		os.Exit(1)
	}

	log.Info("Server exiting")
} 