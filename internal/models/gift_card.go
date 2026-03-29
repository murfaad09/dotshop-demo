package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type GiftCard struct {
	gorm.Model
	ID         uint            `gorm:"primaryKey" json:"id"`
	GiftCardID string          `gorm:"uniqueIndex" json:"gift_card_id"`
	Balance    sql.NullFloat64 `json:"balance"`
	IsRedeemed bool            `gorm:"default:false" json:"is_redeemed"`
	CreatedAt  time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
