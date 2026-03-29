package model

import "time"

type Look struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `json:"name"`
	ImageURL         string     `json:"imageurl"`
	EmbedLink        string     `json:"embed_link"`
	SocialID         string     `json:"social_id"`
	SocialTitle      string     `json:"social_title"`
	VideoDescription string     `json:"video_description"`
	CuratorID        uint       `json:"curator_id"`
	SocialType       string     `json:"social_type"`
	CreatedAt        *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`

	Products []Product `gorm:"many2many:look_products;"`
}
