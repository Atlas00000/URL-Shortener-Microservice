package services

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

// GeoService handles IP-based location lookups
type GeoService struct {
	reader *geoip2.Reader
}

// Location represents a geographical location
type Location struct {
	Country     string  `json:"country"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	CountryCode string  `json:"country_code"`
}

// NewGeoService creates a new GeoService instance
func NewGeoService() (*GeoService, error) {
	dbPath := filepath.Join("data", "geoip", "GeoLite2-City.mmdb")
	reader, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoLite2 database: %v", err)
	}

	return &GeoService{
		reader: reader,
	}, nil
}

// GetLocation looks up the location for an IP address
func (s *GeoService) GetLocation(ipStr string) (*Location, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := s.reader.City(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IP: %v", err)
	}

	location := &Location{
		Country:     record.Country.Names["en"],
		City:        record.City.Names["en"],
		Latitude:    record.Location.Latitude,
		Longitude:   record.Location.Longitude,
		Timezone:    record.Location.TimeZone,
		CountryCode: record.Country.IsoCode,
	}

	return location, nil
}

// Close closes the GeoLite2 database reader
func (s *GeoService) Close() error {
	return s.reader.Close()
} 