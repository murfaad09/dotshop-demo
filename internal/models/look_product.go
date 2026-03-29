package model

type LookProduct struct {
	LookID    uint    `gorm:"primaryKey json:look_id"`
	ProductID string  `gorm:"primaryKey json:product_id"`
	Look      Look    `gorm:"foreignKey:LookID;constraint:OnDelete:CASCADE;"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}
