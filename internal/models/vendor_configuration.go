package model

import (
	"gorm.io/gorm"
)

type VendorConfiguration struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey" json:"id"`
	VendorName string `gorm:"not null;uniqueIndex" json:"vendor_name"`
	ApiKey     string `json:"api_key"`
	Endpoint   string `json:"endpoint"`
	Status     string `json:"status"`
}
