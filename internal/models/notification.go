package model

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Message   string    `json:"message"`
	Title     string    `json:"title"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeleteAt  time.Time `json:"delete_at"`
}
