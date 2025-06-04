package models

import (
	"time"

	"gorm.io/gorm"
)

// URL represents a shortened URL
type URL struct {
	gorm.Model
	LongURL   string     `json:"long_url" gorm:"not null"`
	ShortID   string     `json:"short_id" gorm:"uniqueIndex;not null"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Click represents a click event on a shortened URL
type Click struct {
	ID        uint      `gorm:"primaryKey"`
	URLID     uint      `gorm:"not null"`
	IPAddress string    `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	Country   string    `gorm:"not null"`
	Device    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
} 