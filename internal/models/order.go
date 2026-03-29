package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              string           `gorm:"primaryKey;type:text" json:"id"`
	UserID          uint64           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	PaymentID       string           `json:"payment_id"`
	ShippingMethod  string           `json:"shipping_method"`
	User            User             `gorm:"foreignKey:UserID"`
	TotalAmount     float64          `gorm:"type:real; default:0.0" json:"total_amount"`
	TotalQuantity   uint             `gorm:"default:1" json:"total_quantity"`
	Status          string           `gorm:"type:text; default:pending" json:"status"`
	BuyerReference  string           `gorm:"type:text; default:''" json:"buyer_reference"`
	CreatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Note            string           `gorm:"type:text" json:"note"`
	Cancelled       bool             `gorm:"default:false" json:"cancelled"`
	CancelledReason string           `gorm:"type:text" json:"cancelled_reason"`
	CancelledBy     uint             `gorm:"type=bigint" json:"cancelled_by"`
	CancelledData   *time.Time       `gorm:"cancelled_data; default:null" json:"cancelled_data"`
	RawData         json.RawMessage  `gorm:"type:json" json:"raw_data"`
	StatusUpdatedAt time.Time        `gorm:"status_updated_at; default:null" json:"status_updated_at"`
	BuyerCancelled  bool             `gorm:"buyer_cancelled; default:false" json:"buyer_cancelled"`
	IsTest          bool             `gorm:"is_test; default:false" json:"is_test"`
	SellerCancelled bool             `gorm:"seller_cancelled; default:false" json:"seller_cancelled"`
	OrderVariants   []*OrderVariants `gorm:"foreignKey:OrderID"`
	ShippingAddress *string          `gorm:"type:text" json:"shippingAddress"`
	ShippingCity    *string          `gorm:"type:text" json:"shippingCity"`
	ShippingState   *string          `gorm:"type:text" json:"shippingState"`
	ShippingCountry *string          `gorm:"type:text" json:"shippingCountry"`
	ShippingZip     *string          `gorm:"type:text" json:"shippingZip"`
	Variant         []*Variant       `gorm:"many2many:order_variants;"`
	FulFillments    []*Fulfillments  `gorm:"foreignKey:OrderID"`
}

type OrderVariants struct {
	gorm.Model
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID           string     `gorm:"type=text; not null;"`
	VariantID         string     `gorm:"type=text; not null;"`
	ProductID         string     `gorm:"type:text;"`
	BrandName         string     `gorm:"type:text;" json:"brand_name"`
	CuratorID         uint64     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"curator_id"`
	SellerOrderId     string     `gorm:"type:text;" json:"seller_order_id"`
	SellerOrderItemId string     `gorm:"type:text;" json:"seller_order_item_id"`
	BuyerReference    string     `gorm:"type:text;" json:"buyer_reference"`
	CreatedAt         time.Time  `gorm:"created_at; default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"updated_at; default:CURRENT_TIMESTAMP" json:"updated_at"`
	Price             float64    `gorm:"type:double precision; default:0.0" json:"price"`
	VariantOptionName string     `gorm:"type:text;" json:"variant_option_name"`
	VariantSize       string     `gorm:"type:text;" json:"variant_size"`
	Quantity          uint       `gorm:"default:1" json:"quantity"`
	Cancelled         bool       `gorm:"default:false" json:"cancelled"`
	CancelledReason   string     `gorm:"type:text" json:"cancelled_reason"`
	CancelledBy       uint       `gorm:"type=bigint" json:"cancelled_by"`
	CancelledData     *time.Time `gorm:"cancelled_data; default:null" json:"cancelled_data"`
	Order             *Order
	Variant           *Variant
}
