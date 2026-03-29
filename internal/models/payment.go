package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       string    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"order_id"`
	Order         Order     `gorm:"foreignKey:OrderID;reference:ID"`
	Amount        float64   `gorm:"type:real; default:0.0" json:"amount"`
	Status        string    `gorm:"type:text;" json:"status"`
	PaymentMethod string    `gorm:"type:text;" json:"payment_method"`
	TransactionID string    `gorm:"type:text;" json:"transaction_id"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
