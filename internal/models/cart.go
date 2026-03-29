package model

import (
	"time"
)

type Cart struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	User      User        `gorm:"foreignKey:UserID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []*CartItem `gorm:"foreignKey:CartID" json:"items"`
	DeletedAt *time.Time  `json:"deleted_at"`
}
