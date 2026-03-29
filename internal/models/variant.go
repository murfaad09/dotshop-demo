package model

import (
	"encoding/json"
)

type Variant struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	ProductID       string          `json:"product_id"`
	SKU             string          `json:"sku"`
	Title           string          `json:"title"`
	InventoryAmount int             `json:"inventoryAmount"`
	Image           string          `json:"image"`
	RetailPrice     float64         `json:"retailPrice"`
	RetailCurrency  string          `json:"retailCurrency"`
	BasePrice       float64         `json:"basePrice"`
	BaseCurrency    string          `json:"baseCurrency"`
	VariantOptions  []VariantOption `gorm:"foreignKey:VariantID"`
	Units           string          `json:"units"`
	Attributes      json.RawMessage `gorm:"type:jsonb" json:"attributes"`
	Variant         []*Order        `gorm:"many2many:order_variants;"`
}

type VariantOption struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	VariantID string `json:"VariantID"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}
