package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"unique;not null" json:"name"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	SubCategories []*SubCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subcategories"`
}

type SubCategory struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"unique;not null" json:"name"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CategoryID uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category_id"`
}
