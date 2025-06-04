package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/urlshortener/src/api/handlers"
	"github.com/yourusername/urlshortener/src/logger"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	port   string
	server *http.Server
}

// NewServer creates a new server instance
func NewServer() *Server {
	router := gin.Default()
	
	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Serve static files
	router.Static("/static", "./static")
	
	// Serve index.html for root path
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return &Server{
		router: router,
		port:   "8080",
	}
}

// RegisterRoutes registers all API routes
func (s *Server) RegisterRoutes(urlHandler *handlers.URLHandler, analyticsHandler *handlers.AnalyticsHandler) {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// URL routes
	s.router.POST("/shorten", urlHandler.ShortenURL)
	s.router.GET("/:shortID", urlHandler.RedirectToLongURL)

	// Analytics routes
	s.router.GET("/analytics", analyticsHandler.GetAnalytics)
	s.router.POST("/analytics/click", analyticsHandler.RecordClick)
}

// Start starts the server
func (s *Server) Start() error {
	// Create HTTP server with timeouts
	s.server = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	logger.LogInfo("Server starting on port "+s.port, nil)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
} 