package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	CuratorID uint           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"curator_id"`
	Curator   Curator        `gorm:"foreignKey:CuratorID;references:ID"`
	MediaType sql.NullString `json:"media_type"`
	MediaURL  sql.NullString `json:"media_url"`
	Caption   sql.NullString `json:"caption"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Media) TableName() string {
	return "media"
}
