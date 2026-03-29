package model

type CollectionSection struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	ImageURL     *string    `json:"image_url"`
	CollectionID uint       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"collection_id"`
	Products     []*Product `gorm:"many2many:collection_section_products;"`
}
