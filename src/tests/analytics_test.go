package tests

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/urlshortener/config"
	"github.com/yourusername/urlshortener/src/models"
	"github.com/yourusername/urlshortener/src/services"
	"github.com/yourusername/urlshortener/src/storage"
)

func TestAnalytics(t *testing.T) {
	// Initialize database config for in-memory SQLite and test Redis
	cfg := &config.DatabaseConfig{
		SQLite: config.SQLiteConfig{Path: ":memory:"},
		Redis:  config.RedisConfig{Host: "localhost", Port: "6379", Password: "", DB: 1},
	}
	// Initialize database
	db, err := storage.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := storage.RunMigrations(db.SQLite); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create analytics service
	analyticsService := services.NewAnalyticsService(db)

	// Create test URL
	url := &models.URL{
		OriginalURL: "https://example.com",
		ShortID:     "test123",
		CreatedAt:   time.Now(),
	}
	if err := db.SQLite.Create(url).Error; err != nil {
		t.Fatalf("Failed to create test URL: %v", err)
	}

	// Test recording clicks
	t.Run("RecordClicks", func(t *testing.T) {
		// Create test request
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "8.8.8.8" // US IP
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)")

		// Record click
		if err := analyticsService.RecordClick(url.ID, req); err != nil {
			t.Errorf("Failed to record click: %v", err)
		}

		// Verify click was recorded
		var click models.Click
		if err := db.SQLite.Where("url_id = ?", url.ID).First(&click).Error; err != nil {
			t.Errorf("Failed to find recorded click: %v", err)
		}

		// Verify click details
		if click.IPAddress != "8.8.8.8" {
			t.Errorf("Expected IP 8.8.8.8, got %s", click.IPAddress)
		}
		if click.Country != "US" {
			t.Errorf("Expected country US, got %s", click.Country)
		}
		if click.Device != "mobile" {
			t.Errorf("Expected device mobile, got %s", click.Device)
		}
	})

	// Test getting analytics
	t.Run("GetAnalytics", func(t *testing.T) {
		// Record more clicks with different devices and countries
		clicks := []struct {
			ip      string
			ua      string
			country string
			device  string
		}{
			{"1.1.1.1", "Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X)", "AU", "tablet"},
			{"185.143.223.12", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", "RU", "desktop"},
		}

		for _, c := range clicks {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = c.ip
			req.Header.Set("User-Agent", c.ua)
			if err := analyticsService.RecordClick(url.ID, req); err != nil {
				t.Errorf("Failed to record click: %v", err)
			}
		}

		// Get analytics
		analytics, err := analyticsService.GetAnalytics(url.ID)
		if err != nil {
			t.Errorf("Failed to get analytics: %v", err)
		}

		// Verify analytics
		if analytics.TotalClicks != 3 {
			t.Errorf("Expected 3 total clicks, got %d", analytics.TotalClicks)
		}

		// Verify clicks by country
		if analytics.ClicksByCountry["US"] != 1 {
			t.Errorf("Expected 1 click from US, got %d", analytics.ClicksByCountry["US"])
		}
		if analytics.ClicksByCountry["AU"] != 1 {
			t.Errorf("Expected 1 click from AU, got %d", analytics.ClicksByCountry["AU"])
		}
		if analytics.ClicksByCountry["RU"] != 1 {
			t.Errorf("Expected 1 click from RU, got %d", analytics.ClicksByCountry["RU"])
		}

		// Verify clicks by device
		if analytics.ClicksByDevice["mobile"] != 1 {
			t.Errorf("Expected 1 mobile click, got %d", analytics.ClicksByDevice["mobile"])
		}
		if analytics.ClicksByDevice["tablet"] != 1 {
			t.Errorf("Expected 1 tablet click, got %d", analytics.ClicksByDevice["tablet"])
		}
		if analytics.ClicksByDevice["desktop"] != 1 {
			t.Errorf("Expected 1 desktop click, got %d", analytics.ClicksByDevice["desktop"])
		}

		// Verify last click time
		if analytics.LastClick.IsZero() {
			t.Error("Expected non-zero last click time")
		}
	})
} 