package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ProductID         string               `gorm:"primaryKey" json:"id"`
	BrandName         string               `json:"brand_name"`
	BrandImage        string               `json:"brand_image"`
	SupplierName      string               `json:"supplier_name"`
	BrandID           *uint                `json:"brand_id"`
	Brands            *Brand               `gorm:"foreignKey:BrandID" json:"brand"`
	IsActive          bool                 `gorm:"default:true" json:"is_active"`
	Variants          []Variant            `gorm:"foreignKey:ProductID" json:"variants"`
	Notes             string               `json:"notes"`
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	CreatedAt         time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt       `gorm:"index" json:"deleted_at"`
	Tags              *pq.StringArray      `gorm:"type:text[]"`
	Looks             []Look               `gorm:"many2many:look_products;"`
	Collections       []Collection         `gorm:"many2many:collection_products;"`
	CollectionSection []*CollectionSection `gorm:"many2many:collection_section_products;"`
	Curators          []*Curator           `gorm:"many2many:curator_products;"`
	Reviews           []*Review            `gorm:"foreignKey:ProductID" json:"reviews"`
	CategoryID        uint                 `json:"category_id"`
	SubCategoryID     *uint                `json:"sub_category_id"`
	SubCategory       *SubCategory         `gorm:"foreignKey:SubCategoryID" json:"sub_category"`
	PromotionID       *uint                `json:"promotion_id"`
	Promotion         *Promotion           `gorm:"foreignKey:PromotionID" json:"promotion"`
}
