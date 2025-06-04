package geo

import (
	"fmt"
	"os"
	"path/filepath"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type Service struct {
	reader *geoip2.Reader
}

func NewService(dataDir string) (*Service, error) {
	dbPath := filepath.Join(dataDir, "geoip", "GeoLite2-City.mmdb")
	
	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("GeoLite2 database file not found at %s", dbPath)
	}

	reader, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoLite2 database: %v", err)
	}

	return &Service{
		reader: reader,
	}, nil
}

func (s *Service) GetLocation(ip string) (string, error) {
	if s.reader == nil {
		return "Unknown", nil
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "Unknown", fmt.Errorf("invalid IP address: %s", ip)
	}

	record, err := s.reader.City(parsedIP)
	if err != nil {
		return "Unknown", err
	}

	return record.Country.Names["en"], nil
}

func (s *Service) Close() {
	if s.reader != nil {
		s.reader.Close()
	}
} 