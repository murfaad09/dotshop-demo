package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type SupportTicket struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	Subject   sql.NullString `json:"subject"`
	Message   sql.NullString `json:"message"`
	Status    sql.NullString `json:"status"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	User      User           `gorm:"foreignKey:UserID"`
}
