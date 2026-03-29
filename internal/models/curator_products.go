package model

import "time"

type CuratorProduct struct {
	FeatureID *uint      `gorm:"autoIncrement" json:"feature_id"`
	CuratorID uint       `gorm:"primaryKey" json:"curator_id"`
	ProductID string     `gorm:"primaryKey" json:"product_id"`
	Curator   *Curator   `gorm:"foreignKey:CuratorID;constraint:OnDelete:CASCADE;"`
	Product   *Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	IsFeature bool       `json:"is_feature"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
