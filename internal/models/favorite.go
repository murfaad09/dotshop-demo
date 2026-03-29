package model

import (
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
}
