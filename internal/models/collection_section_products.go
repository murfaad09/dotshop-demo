package model

type CollectionSectionProduct struct {
	CollectionSectionID uint              `gorm:"primaryKey"`
	ProductID           string            `gorm:"primaryKey"`
	CollectionSection   CollectionSection `gorm:"foreignKey:CollectionSectionID;constraint:OnDelete:CASCADE;"`
	Product             Product           `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}
