package model

type CollectionProduct struct {
	CollectionID uint       `gorm:"primaryKey"`
	ProductID    string     `gorm:"primaryKey"`
	Collection   Collection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE;"`
	Product      Product    `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}
