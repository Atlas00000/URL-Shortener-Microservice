package tests

import (
	"testing"
	"time"
	"net/http/httptest"

	"github.com/yourusername/urlshortener/src/services"
	"github.com/yourusername/urlshortener/src/storage"
	"github.com/yourusername/urlshortener/config"
)

// setupTestDB creates a test database instance
func setupTestDB(t *testing.T) *storage.Database {
	cfg := &config.DatabaseConfig{
		SQLite: config.SQLiteConfig{
			Path: ":memory:", // Use in-memory SQLite for tests
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	db, err := storage.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	if err := storage.RunMigrations(db.SQLite); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// TestURLShortening tests the basic URL shortening functionality
func TestURLShortening(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := services.NewURLService(db)

	// Test valid URL
	longURL := "https://example.com/test"
	urlRecord, err := service.ShortenURL(longURL, 30)
	if err != nil {
		t.Errorf("Failed to shorten URL: %v", err)
	}

	if len(urlRecord.ShortID) != 7 {
		t.Errorf("Expected short ID length of 7, got %d", len(urlRecord.ShortID))
	}

	// Test URL retrieval
	retrievedURL, err := service.GetLongURL(urlRecord.ShortID)
	if err != nil {
		t.Errorf("Failed to retrieve URL: %v", err)
	}

	if retrievedURL != longURL {
		t.Errorf("Expected URL %s, got %s", longURL, retrievedURL)
	}
}

// TestDeviceDetection tests the device detection functionality
func TestDeviceDetection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := services.NewURLService(db)

	// Test cases for different user agents
	testCases := []struct {
		userAgent string
		expected  string
	}{
		{
			"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15",
			"mobile",
		},
		{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"desktop",
		},
		{
			"Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15",
			"tablet",
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("User-Agent", tc.userAgent)

		deviceType := service.DetectDeviceType(req)
		if deviceType != tc.expected {
			t.Errorf("Expected device type %s for user agent %s, got %s",
				tc.expected, tc.userAgent, deviceType)
		}
	}
}

// TestGeoFencing tests the geo-fencing functionality
func TestGeoFencing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := services.NewURLService(db)

	// Test cases for different IP addresses
	testCases := []struct {
		ip       string
		country  string
		allowed  bool
	}{
		{
			"8.8.8.8",    // Google DNS (US)
			"US",
			true,
		},
		{
			"1.1.1.1",    // Cloudflare DNS (AU)
			"AU",
			true,
		},
		{
			"185.143.223.12", // Example IP
			"RU",
			false, // Assuming RU is in restricted countries
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = tc.ip

		country, allowed := service.CheckGeoFencing(req)
		if country != tc.country {
			t.Errorf("Expected country %s for IP %s, got %s",
				tc.country, tc.ip, country)
		}
		if allowed != tc.allowed {
			t.Errorf("Expected allowed=%v for IP %s, got %v",
				tc.allowed, tc.ip, allowed)
		}
	}
}

// TestURLExpiration tests URL expiration functionality
func TestURLExpiration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := services.NewURLService(db)

	// Test URL with 1-day expiration
	longURL := "https://example.com/expiring"
	urlRecord, err := service.ShortenURL(longURL, 1)
	if err != nil {
		t.Errorf("Failed to shorten URL: %v", err)
	}

	// Verify expiration time
	expectedExpiration := time.Now().AddDate(0, 0, 1)
	if urlRecord.ExpiresAt.Sub(expectedExpiration) > time.Hour {
		t.Errorf("Expiration time not set correctly")
	}

	// Test expired URL
	service.ForceExpireURL(urlRecord.ShortID)
	_, err = service.GetLongURL(urlRecord.ShortID)
	if err == nil {
		t.Error("Expected error for expired URL, got nil")
	}
}

// TestRateLimiting tests the rate limiting functionality
func TestRateLimiting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := services.NewURLService(db)

	// Test rate limiting for same IP
	ip := "192.168.1.1"
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = ip

		allowed := service.CheckRateLimit(req)
		if i >= 60 && allowed { // Assuming 60 requests per minute limit
			t.Errorf("Rate limit not enforced at request %d", i+1)
		}
	}
} 