package handlers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/urlshortener/src/logger"
	"github.com/yourusername/urlshortener/src/models"
	"go.uber.org/zap"
)

// URLService defines the interface for URL operations
type URLService interface {
	CreateShortURL(longURL string, expiresAt *time.Time) (*models.URL, error)
	GetLongURL(shortID string) (string, error)
}

// URLHandler handles URL-related HTTP requests
type URLHandler struct {
	urlService URLService
	logger     *zap.Logger
	baseURL    string
}

// NewURLHandler creates a new URL handler
func NewURLHandler(urlService URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		logger:     logger.Get(),
		baseURL:    baseURL,
	}
}

// ShortenURL handles requests to create a shortened URL
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var input struct {
		URL           string `json:"url" binding:"required"`
		ExpirationDays int    `json:"expiration_days"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid input for URL shortening",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate URL
	if _, err := url.ParseRequestURI(input.URL); err != nil {
		h.logger.Warn("Invalid URL format",
			zap.Error(err),
			zap.String("url", input.URL))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// Set expiration time if specified
	var expiresAt *time.Time
	if input.ExpirationDays > 0 {
		t := time.Now().AddDate(0, 0, input.ExpirationDays)
		expiresAt = &t
	}

	// Create shortened URL
	shortURL, err := h.urlService.CreateShortURL(input.URL, expiresAt)
	if err != nil {
		h.logger.Error("Failed to create short URL",
			zap.Error(err),
			zap.String("long_url", input.URL))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("URL shortened successfully",
		zap.String("short_id", shortURL.ShortID),
		zap.String("long_url", shortURL.LongURL),
		zap.Timep("expires_at", shortURL.ExpiresAt))

	// Construct the full shortened URL using the base URL
	shortenedURL := h.baseURL + "/" + shortURL.ShortID

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortenedURL,
		"long_url":  shortURL.LongURL,
		"expires_at": shortURL.ExpiresAt,
	})
}

// RedirectToLongURL handles requests to redirect to the original URL
func (h *URLHandler) RedirectToLongURL(c *gin.Context) {
	shortID := c.Param("shortID")
	if shortID == "" {
		h.logger.Warn("Empty short ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short ID is required"})
		return
	}

	// Get the original URL
	longURL, err := h.urlService.GetLongURL(shortID)
	if err != nil {
		h.logger.Warn("URL not found or expired",
			zap.Error(err),
			zap.String("short_id", shortID))
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or expired"})
		return
	}

	h.logger.Info("Redirecting to long URL",
		zap.String("short_id", shortID),
		zap.String("long_url", longURL))

	c.Redirect(http.StatusMovedPermanently, longURL)
} 