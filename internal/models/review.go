package model

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Rating    int     `gorm:"not null" json:"rating"`
	Comment   *string `gorm:"size:255" json:"comment"`
	UserID    uint    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	ProductID string  `gorm:"type:varchar(255);constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"product_id"`
	Product   Product
	CuratorID uint           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"curator_id"`
	Curator   Curator        `gorm:"foreignKey:CuratorID" json:"curator"`
	Comments  []*Comment     `gorm:"foreignKey:ReviewID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
