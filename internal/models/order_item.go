package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	ID        uint            `gorm:"primaryKey" json:"id"`
	OrderID   string          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"order_id"`
	Order     Order           `gorm:"foreignKey:OrderID;reference:ID"`
	ProductID string          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"product_id"`
	Product   Product         `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  sql.NullInt64   `json:"quantity"`
	Price     sql.NullFloat64 `json:"price"`
}
