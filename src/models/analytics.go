package models

import "time"

// Analytics represents aggregated analytics data
type Analytics struct {
	TotalClicks      int64            `json:"total_clicks"`
	ClicksByCountry  map[string]int64 `json:"clicks_by_country"`
	ClicksByDevice   map[string]int64 `json:"clicks_by_device"`
	LastClick        time.Time        `json:"last_click"`
} 