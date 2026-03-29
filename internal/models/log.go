package model

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Action    string    `json:"action"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Details   string    `json:"details"`
}
