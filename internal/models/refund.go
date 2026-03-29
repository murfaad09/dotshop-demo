package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Refund struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   string         `gorm:"type:text;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"order_id"`
	Order     Order          `gorm:"foreignKey:OrderID;references:ID"`
	Reason    sql.NullString `json:"reason"`
	Status    sql.NullString `json:"status"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
