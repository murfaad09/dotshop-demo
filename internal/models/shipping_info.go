package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type ShippingInfo struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserId         uint           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	FirstName      *string        `json:"first_name"`
	LastName       *string        `json:"last_name"`
	AddressOne     sql.NullString `json:"address_one"`
	AddressTwo     sql.NullString `json:"address_two"`
	City           sql.NullString `json:"city"`
	State          sql.NullString `json:"state"`
	Country        sql.NullString `json:"country"`
	Company        sql.NullString `json:"company"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	Zip            sql.NullString `json:"zip"`
	DefaultAddress bool           `gorm:"default_address; default:false"`
	DefaultBilling bool           `gorm:"default_billing; default:false"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	User           User           `gorm:"foreignKey:UserId"`
}

func (ShippingInfo) TableName() string {
	return "shipping_info"
}
