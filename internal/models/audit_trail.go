package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// type AuditTrail struct {
// 	gorm.Model
// 	ID            uint           `gorm:"primaryKey" json:"id"`
// 	UserID        int            `gorm:"column:user_id" json:"user_id"`
// 	User          User           `gorm:"foreignKey:UserID;references:ID"`
// 	Action        sql.NullString `json:"action"`
// 	Timestamp     time.Time      `gorm:"default:CURRENT_TIMESTAMP()" json:"timestamp"`
// 	Details       sql.NullString `json:"details"`
// 	TableAffected sql.NullString `json:"table_affected"`
// 	RecordID      sql.NullInt64  `json:"record_id"`
// }

type AuditTrail struct {
	gorm.Model
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	Action        sql.NullString `json:"action"`
	Timestamp     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Details       sql.NullString `json:"details"`
	TableAffected sql.NullString `gorm:"size:255" json:"table_affected"`
	RecordID      sql.NullInt64  `json:"record_id"`
	User          User           `gorm:"foreignKey:UserID;references:ID"`
}

func (AuditTrail) TableName() string {
	return "audit_trail"
}
