package model

import "time"

type ReturnOrder struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ReturnId        string     `gorm:"type:text;" json:"return_id"`
	UserId          uint       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	CuratorId       uint       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"curator_id"`
	OrderId         string     `gorm:"type:text;" json:"order_id"`
	OrderVariantId  uint64     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"order_variant_id"`
	Status          string     `gorm:"type:text;default:pending" json:"status"`
	Reason          string     `gorm:"type:text;" json:"reason"`
	RejectionReason string     `gorm:"type:text;" json:"rejection_reason"`
	Quantity        uint       `gorm:"type:bigint;default:1" json:"quantity"`
	Amount          float64    `gorm:"type:real; default:0.0" json:"amount"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	RejectedAt      *time.Time `json:"rejected_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	CancelledBy     string     `gorm:"type:text;" json:"cancelled_by"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	User          User          `gorm:"foreignKey:UserId"`
	OrderVariants OrderVariants `gorm:"foreignKey:OrderVariantId"`
}
