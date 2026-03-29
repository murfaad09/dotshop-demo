package model

type CartItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	CartID      uint    `gorm:"not null" json:"cart_id"`
	ProductID   string  `gorm:"not null" json:"product_id"`
	Price       float64 `gorm:"price" json:"price"`
	ProductName string  `gorm:"product_name" json:"product_name"`
	BrandName   string  `gorm:"brand_name" json:"brand_name"`
	VariantID   string  `gorm:"variant_id" json:"variant_id"`
	ImageURL    string  `gorm:"image_URL" json:"image_URL"`
	CuratorID   int     `gorm:"curator_id" json:"curator_id"`
	Color       string  `gorm:"color" json:"color"`
	Size        string  `gorm:"size" json:"size"`
	Quantity    uint    `gorm:"quantity" json:"quantity"`
	Cart        Cart    `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart"`
}
