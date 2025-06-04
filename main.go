package main

import (
	"log"
	"os"

	"github.com/yourusername/urlshortener/src/api"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Create and start server
	server, err := api.NewServer(port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 