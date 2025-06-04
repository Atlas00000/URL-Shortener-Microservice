package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourusername/urlshortener/src/logger"
	"github.com/yourusername/urlshortener/src/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// URLService handles URL shortening operations
type URLService struct {
	db                 *gorm.DB
	logger             *zap.Logger
	// Rate limiting
	rateLimiter        map[string][]time.Time
	// Geo-fencing
	restrictedCountries map[string]bool
}

// NewURLService creates a new URL service
func NewURLService(db *gorm.DB) *URLService {
	// Initialize restricted countries
	restrictedCountries := map[string]bool{
		"RU": true, // Example: Restrict Russia
		"CN": true, // Example: Restrict China
	}

	return &URLService{
		db:                 db,
		logger:             logger.Get(),
		rateLimiter:        make(map[string][]time.Time),
		restrictedCountries: restrictedCountries,
	}
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(longURL string, expiresAt *time.Time) (*models.URL, error) {
	// Generate a random short ID
	shortID, err := generateShortID()
	if err != nil {
		s.logger.Error("Failed to generate short ID",
			zap.Error(err))
		return nil, err
	}

	// Create URL record
	url := &models.URL{
		ShortID:   shortID,
		LongURL:   longURL,
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(url).Error; err != nil {
		s.logger.Error("Failed to create URL record",
			zap.Error(err),
			zap.String("short_id", shortID),
			zap.String("long_url", longURL))
		return nil, err
	}

	s.logger.Info("Created new short URL",
		zap.String("short_id", shortID),
		zap.String("long_url", longURL),
		zap.Timep("expires_at", expiresAt))

	return url, nil
	}

// GetURLByShortID retrieves a URL by its short ID
func (s *URLService) GetURLByShortID(shortID string) (*models.URL, error) {
	var url models.URL
	if err := s.db.Where("short_id = ? AND (expires_at IS NULL OR expires_at > ?)", shortID, time.Now()).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL not found or expired")
		}
		return nil, err
	}
	return &url, nil
}

// GetLongURL retrieves the original URL for a given short ID
func (s *URLService) GetLongURL(shortID string) (string, error) {
	var url models.URL
	if err := s.db.Where("short_id = ?", shortID).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("URL not found",
				zap.String("short_id", shortID))
			return "", errors.New("URL not found")
		}
		s.logger.Error("Database error while retrieving URL",
			zap.Error(err),
			zap.String("short_id", shortID))
		return "", err
	}

	// Check if URL has expired
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		s.logger.Warn("URL has expired",
			zap.String("short_id", shortID),
			zap.Time("expires_at", *url.ExpiresAt))
		return "", errors.New("URL has expired")
	}

	s.logger.Info("Retrieved long URL",
		zap.String("short_id", shortID),
		zap.String("long_url", url.LongURL))

	return url.LongURL, nil
}

// validateURL checks if the given URL is valid
func validateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("URL must include scheme and host")
	}
	return nil
}

// generateShortID generates a random short ID
func generateShortID() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}

// DetectDeviceType detects the device type from the user agent
func (s *URLService) DetectDeviceType(r *http.Request) string {
	userAgent := strings.ToLower(r.UserAgent())

	switch {
	case strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "android"):
		return "mobile"
	case strings.Contains(userAgent, "ipad"):
		return "tablet"
	default:
		return "desktop"
	}
}

// CheckGeoFencing checks if the request is allowed based on geo-fencing rules
func (s *URLService) CheckGeoFencing(r *http.Request) (string, bool) {
	// In a real implementation, you would use a geo-IP database
	// For testing, we'll use a simple mapping
	ipToCountry := map[string]string{
		"8.8.8.8":         "US",
		"1.1.1.1":         "AU",
		"185.143.223.12":  "RU",
	}

	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	country := ipToCountry[ip]
	if country == "" {
		country = "US" // Default to US for unknown IPs
	}

	return country, !s.restrictedCountries[country]
}

// CheckRateLimit checks if the request is within rate limits
func (s *URLService) CheckRateLimit(r *http.Request) bool {
	ip := r.RemoteAddr
	now := time.Now()

	// Clean up old timestamps
	if timestamps, exists := s.rateLimiter[ip]; exists {
		var validTimestamps []time.Time
		for _, ts := range timestamps {
			if now.Sub(ts) < time.Minute {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		s.rateLimiter[ip] = validTimestamps
	}

	// Check if rate limit exceeded
	if timestamps, exists := s.rateLimiter[ip]; exists {
		if len(timestamps) >= 60 { // 60 requests per minute
			return false
		}
	}

	// Add new timestamp
	s.rateLimiter[ip] = append(s.rateLimiter[ip], now)
	return true
}

// ForceExpireURL forces a URL to expire (for testing)
func (s *URLService) ForceExpireURL(shortID string) error {
	return s.db.Model(&models.URL{}).
		Where("short_id = ?", shortID).
		Update("expires_at", time.Now().Add(-time.Hour)).
		Error
} 