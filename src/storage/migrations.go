package storage

import (
	"github.com/yourusername/urlshortener/src/models"
	"gorm.io/gorm"
)

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	// Auto migrate the schema
	if err := db.AutoMigrate(&models.URL{}, &models.Click{}); err != nil {
		return err
	}

	// Create indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_urls_short_id ON urls(short_id)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id)").Error; err != nil {
		return err
	}

	return nil
} 