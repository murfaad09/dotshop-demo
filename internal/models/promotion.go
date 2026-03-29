package model

import (
	"time"

	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	DiscountCode  string    `gorm:"not null;uniqueIndex" json:"discount_code"`
	DiscountValue float64   `gorm:"not null" json:"discount_value"`
	ExpiryDate    time.Time `gorm:"not null" json:"expiry_date"`
	Status        string    `gorm:"default:'active'" json:"status"`
	Rule          string    `gorm:"not null" json:"rule"`
}
