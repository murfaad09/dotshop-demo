package models

type Product struct {
	ID       string  `gorm:"primaryKey"`
	Quantity float64 `gorm:"not null"`
}

type Order struct {
	ID     string `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}
