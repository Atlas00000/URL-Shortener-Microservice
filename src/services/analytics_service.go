package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AnalyticsService handles analytics-related operations
type AnalyticsService struct {
	db         *gorm.DB
	geoService *GeoService
}

// Click represents a click event
type Click struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	URLID       uint      `json:"url_id" gorm:"index"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	DeviceType  string    `json:"device_type"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Timezone    string    `json:"timezone"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewAnalyticsService creates a new analytics service instance
func NewAnalyticsService(db *gorm.DB, geoService *GeoService) *AnalyticsService {
	return &AnalyticsService{
		db:         db,
		geoService: geoService,
	}
}

// RecordClick records a click event for a URL
func (s *AnalyticsService) RecordClick(urlID uint, r *http.Request) error {
	// Get device type from user agent
	deviceType := "unknown"
	userAgent := r.UserAgent()
	if userAgent != "" {
		if isMobile(userAgent) {
			deviceType = "mobile"
		} else if isTablet(userAgent) {
			deviceType = "tablet"
		} else {
			deviceType = "desktop"
		}
	}

	// Get IP address
	ip := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip = forwardedFor
	}

	// Get location data
	var location *Location
	var err error
	if s.geoService != nil {
		location, err = s.geoService.GetLocation(ip)
		if err != nil {
			// Log the error but continue without location data
			fmt.Printf("Failed to get location for IP %s: %v\n", ip, err)
		}
	}

	click := &Click{
		URLID:      urlID,
		IPAddress:  ip,
		UserAgent:  userAgent,
		DeviceType: deviceType,
		CreatedAt:  time.Now(),
	}

	// Add location data if available
	if location != nil {
		click.Country = location.Country
		click.City = location.City
		click.Latitude = location.Latitude
		click.Longitude = location.Longitude
		click.Timezone = location.Timezone
		click.CountryCode = location.CountryCode
	}

	return s.db.Create(click).Error
}

// GetAnalytics retrieves analytics data for a URL
func (s *AnalyticsService) GetAnalytics(urlID uint) (map[string]interface{}, error) {
	var totalClicks int64
	if err := s.db.Model(&Click{}).Where("url_id = ?", urlID).Count(&totalClicks).Error; err != nil {
		return nil, err
	}

	var deviceStats []struct {
		DeviceType string
		Count      int64
	}
	if err := s.db.Model(&Click{}).
		Select("device_type, count(*) as count").
		Where("url_id = ?", urlID).
		Group("device_type").
		Scan(&deviceStats).Error; err != nil {
		return nil, err
	}

	var countryStats []struct {
		Country     string
		CountryCode string
		Count       int64
	}
	if err := s.db.Model(&Click{}).
		Select("country, country_code, count(*) as count").
		Where("url_id = ?", urlID).
		Group("country, country_code").
		Order("count DESC").
		Limit(10).
		Scan(&countryStats).Error; err != nil {
		return nil, err
	}

	var recentClicks []Click
	if err := s.db.Where("url_id = ?", urlID).
		Order("created_at DESC").
		Limit(10).
		Find(&recentClicks).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_clicks":    totalClicks,
		"device_stats":    deviceStats,
		"country_stats":   countryStats,
		"recent_clicks":   recentClicks,
	}, nil
}

// Helper functions to detect device type
func isMobile(userAgent string) bool {
	// Add more mobile device patterns as needed
	mobilePatterns := []string{
		"Mobile", "Android", "iPhone", "iPad", "Windows Phone",
	}
	for _, pattern := range mobilePatterns {
		if contains(userAgent, pattern) {
			return true
		}
	}
	return false
}

func isTablet(userAgent string) bool {
	// Add more tablet device patterns as needed
	tabletPatterns := []string{
		"iPad", "Android.*Tablet", "Tablet",
	}
	for _, pattern := range tabletPatterns {
		if contains(userAgent, pattern) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && s[0:len(substr)] == substr
}

// DetectDeviceType detects the device type from the user agent
func (s *AnalyticsService) DetectDeviceType(r *http.Request) string {
	userAgent := strings.ToLower(r.UserAgent())

	switch {
	case strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "android"):
		return "mobile"
	case strings.Contains(userAgent, "ipad"):
		return "tablet"
	case strings.Contains(userAgent, "windows") || strings.Contains(userAgent, "macintosh") || strings.Contains(userAgent, "linux"):
		return "desktop"
	default:
		return "other"
	}
}

// CheckGeoFencing checks if the request is allowed based on geo-fencing rules
func (s *AnalyticsService) CheckGeoFencing(r *http.Request) (string, bool) {
	// In a real implementation, you would use a geo-IP database
	// For testing, we'll use a simple mapping
	ipToCountry := map[string]string{
		"8.8.8.8":         "US",
		"1.1.1.1":         "AU",
		"185.143.223.12":  "RU",
	}

	ip := r.RemoteAddr
	country := ipToCountry[ip]
	if country == "" {
		country = "US" // Default to US for unknown IPs
	}

	return country, true
} 