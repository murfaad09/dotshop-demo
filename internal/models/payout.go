package model

import (
	"time"
)

type Payout struct {
	ID               uint64    `gorm:"type:bigint; not null; primaryKey" json:"id"`
	CuratorID        uint64    `gorm:"type:bigint; not null; unique" json:"curator_id"`
	ReturnAmount     float64   `gorm:"type:double precision;" json:"return_amount"`
	CommissionAmount float64   `gorm:"type:double precision;" json:"commission_amount"`
	PayoutAmount     float64   `gorm:"type:double precision;" json:"payout_amount"`
	Status           string    `gorm:"type:text; default:in progress" json:"status"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PayoutHistory struct {
	ID               uint64    `gorm:"type:bigint; not null; primaryKey" json:"id"`
	CuratorID        uint64    `gorm:"type:bigint; not null;" json:"curator_id"`
	ReturnAmount     float64   `gorm:"type:double precision;" json:"return_amount"`
	CommissionAmount float64   `gorm:"type:double precision;" json:"commission_amount"`
	PayoutAmount     float64   `gorm:"type:double precision;" json:"payout_amount"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
