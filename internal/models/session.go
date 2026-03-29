package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	SessionToken string    `gorm:"not null" json:"session_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User      `gorm:"foreignKey:UserID"`
}
