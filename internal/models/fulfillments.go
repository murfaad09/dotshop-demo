package model

import (
	"encoding/json"
	"time"
)

type Fulfillments struct {
	ID           string          `gorm:"primaryKey" json:"id"`
	OrderID      string          `gorm:"not null" json:"order_id"`
	Posted       bool            `json:"posted"`
	PostedDate   *time.Time      `json:"postedDate"`
	Carrier      string          `json:"carrier"`
	TrackingCode string          `json:"trackingCode"`
	TrackingUrls json.RawMessage `gorm:"type:json" json:"trackingUrls"`
	CreatedAt    time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
