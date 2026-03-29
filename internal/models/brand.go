package model

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Image       *string `json:"image"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	Margin      float64 `json:"margin"`
}
