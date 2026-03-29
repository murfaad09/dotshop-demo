package model

import "time"

type Collection struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	TileColor   string     `json:"tile_color"`
	CuratorID   uint       `json:"curator_id"`
	Products    []Product  `gorm:"many2many:collection_products;" json:"products"`
	CreatedAt   *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`

	CollectionSection []*CollectionSection `gorm:"foreignKey:CollectionID" json:"collection_sections"`
}
