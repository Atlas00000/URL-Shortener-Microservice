package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/urlshortener/src/services"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	urlService      *services.URLService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService, urlService *services.URLService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		urlService:      urlService,
	}
}

// GetAnalytics retrieves analytics data for a URL
func (h *AnalyticsHandler) GetAnalytics(c *gin.Context) {
	shortID := c.Query("short_id")
	if shortID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short ID is required"})
		return
	}

	// Get URL record to verify it exists and get its ID
	url, err := h.urlService.GetURLByShortID(shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or expired"})
		return
	}

	// Get analytics data
	analytics, err := h.analyticsService.GetAnalytics(url.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// RecordClick records a click event for a URL
func (h *AnalyticsHandler) RecordClick(c *gin.Context) {
	shortID := c.Query("short_id")
	if shortID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short ID is required"})
		return
	}

	// Get URL record to verify it exists and get its ID
	url, err := h.urlService.GetURLByShortID(shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or expired"})
		return
	}

	// Record the click
	if err := h.analyticsService.RecordClick(url.ID, c.Request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
} 