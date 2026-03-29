package model

import (
	"time"
)

type WishlistItem struct {
	ID                 uint    `gorm:"primaryKey"`
	CuratorID          uint    `gorm:"not null"`
	UserID             uint    `gorm:"not null"`
	ProductID          string  `gorm:"type:varchar(100);not null"`
	ProductName        string  `gorm:"type:varchar(255);not null"`
	ProductImage       string  `gorm:"type:varchar(255);not null"`
	ProductBrand       string  `gorm:"type:varchar(100);not null"`
	ProductPrice       string  `gorm:"type:varchar(50);not null"`
	VariantID          *string `gorm:"type:varchar(50)"`
	VariantOptionName  string  `gorm:"type:text;" json:"variant_option_name"`
	VariantOptionValue string  `gorm:"type:text;" json:"variant_option_value"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
