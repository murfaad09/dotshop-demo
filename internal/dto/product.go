package dto

import (
	"encoding/json"
	"time"

	domain "github.com/harishash/dotshop-be/internal/models"
	"github.com/lib/pq"
)

type ProductRequest struct {
	ProductID    string           `json:"id"`
	BrandName    string           `json:"brand_name"`
	SupplierName string           `json:"supplier_name"`
	Variants     []VariantRequest `gorm:"foreignKey:ProductID"`
	Notes        string           `json:"notes"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Category     string           `json:"category"`
}
type VariantRequest struct {
	ID              string                 `json:"id"`
	ProductID       string                 `json:"product_id"`
	SKU             string                 `json:"sku"`
	Title           string                 `json:"title"`
	InventoryAmount int                    `json:"inventoryAmount"`
	Image           string                 `json:"image"`
	RetailPrice     float64                `json:"retailPrice"`
	RetailCurrency  string                 `json:"retailCurrency"`
	BasePrice       float64                `json:"basePrice"`
	BaseCurrency    string                 `json:"baseCurrency"`
	VariantOptions  []VariantOptionRequest `json:"variantOptions"`
	Units           string                 `json:"units"`
	Attributes      json.RawMessage        `json:"attributes"`
}

type VariantOptionRequest struct {
	VariantID string `json:"-"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}
type CreateProductRequest struct {
	ProductID     string                 `json:"id"`
	BrandName     string                 `json:"brand_name"`
	SupplierName  string                 `json:"supplier_name"`
	Variants      []CreateVariantRequest `json:"variants"`
	Notes         string                 `json:"notes"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	CategoryID    uint                   `json:"category_id"`
	SubCategoryID uint                   `json:"sub_category_id"`
	Tags          *pq.StringArray        `json:"tags"`
}
type CreateVariantRequest struct {
	ID              string                       `json:"id"`
	ProductID       string                       `json:"product_id"`
	SKU             string                       `json:"sku"`
	Title           string                       `json:"title"`
	InventoryAmount int                          `json:"inventoryAmount"`
	Image           string                       `json:"image"`
	RetailPrice     float64                      `json:"retailPrice"`
	RetailCurrency  string                       `json:"retailCurrency"`
	BasePrice       float64                      `json:"basePrice"`
	BaseCurrency    string                       `json:"baseCurrency"`
	VariantOptions  []CreateVariantOptionRequest `json:"variantOptions"`
	Units           string                       `json:"units"`
	Attributes      json.RawMessage              `json:"attributes"`
}
type CreateVariantOptionRequest struct {
	VariantID string `json:"-"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type CreateProductResponse struct {
	ProductID     string    `json:"id"`
	BrandName     string    `json:"brandName"`
	SupplierName  string    `json:"supplierName"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CategoryID    uint      `json:"category_id"`
	SubCategoryID *uint     `json:"sub_category_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type FilteredProductResponse struct {
	ProductId     string            `json:"id"`
	Name          string            `json:"name"`
	BrandName     string            `json:"brandName"`
	SupplierName  string            `json:"supplierName"`
	Description   string            `json:"description"`
	Notes         string            `json:"notes"`
	Variants      []VariantResponse `json:"variants"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	CategoryID    *uint             `json:"categoryId"`
	SubCategoryID *uint             `json:"subCategoryId"`
}

type Brands struct {
	BrandName  string `json:"brandName"`
	BrandImage string `json:"brandImage"`
}

type Variants struct {
	ID              string  `json:"id"`
	Sku             string  `json:"sku"`
	Title           string  `json:"title"`
	InventoryAmount int     `json:"inventoryAmount"`
	RetailPrice     int     `json:"retailPrice"`
	RetailCurrency  string  `json:"retailCurrency"`
	BasePrice       float64 `json:"basePrice"`
	BaseCurrency    string  `json:"baseCurrency"`
}

type Convictional_Variants_Data struct {
	Variants []Variants `json:"data"`
}

func RowToFilteredProductResponse(products []domain.Product) []FilteredProductResponse {
	var productResponses []FilteredProductResponse
	for _, product := range products {
		var variants []VariantResponse
		for _, variant := range product.Variants {

			var variantsOption []*VariantOptionResponse
			for _, variantOption := range variant.VariantOptions {
				variantsOption = append(variantsOption, &VariantOptionResponse{
					Name:      variantOption.Name,
					Value:     variantOption.Value,
					VariantID: variantOption.VariantID,
				})
			}

			variants = append(variants, VariantResponse{
				ID:              variant.ID,
				ProductID:       variant.ProductID,
				SKU:             variant.SKU,
				Title:           variant.Title,
				InventoryAmount: variant.InventoryAmount,
				Image:           variant.Image,
				RetailPrice:     variant.RetailPrice,
				RetailCurrency:  variant.RetailCurrency,
				BasePrice:       variant.BasePrice,
				BaseCurrency:    variant.BaseCurrency,
				Units:           variant.Units,
				Attributes:      variant.Attributes,
				VariantOptions:  variantsOption,
			})
		}
		productResponses = append(productResponses, FilteredProductResponse{
			ProductId:     product.ProductID,
			BrandName:     product.BrandName,
			SupplierName:  product.SupplierName,
			Name:          product.Name,
			Notes:         product.Notes,
			CategoryID:    &product.CategoryID,
			SubCategoryID: product.SubCategoryID,
			Description:   product.Description,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			Variants:      variants,
		})
	}
	return productResponses
}

type ListedProducts struct {
	ProductID   string `gorm:"column:product_id" json:"productId"`
	ProductName string `gorm:"column:name" json:"productName"`
	BrandName   string `gorm:"column:brand_name" json:"brandName"`
}

type ListedProductResponse struct {
	ProductID   string                  `gorm:"column:product_id" json:"productId"`
	ProductName string                  `gorm:"column:name" json:"productName"`
	BrandName   string                  `gorm:"column:brand_name" json:"brandName"`
	Variants    []*ListedProductVariant `json:"variants"`
}

type ListedProductVariant struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Image          string  `json:"image"`
	RetailPrice    float64 `json:"retailPrice"`
	RetailCurrency string  `json:"retailCurrency"`
	BasePrice      float64 `json:"basePrice"`
	BaseCurrency   string  `json:"baseCurrency"`
	Units          string  `json:"units"`
}

type ProductWithStats struct {
	ProductID     string  `gorm:"column:product_id" json:"productId"`
	ProductName   string  `gorm:"column:product_name" json:"productName"`
	Image         string  `gorm:"column:image" json:"image"`
	TotalReviews  int     `gorm:"column:total_reviews" json:"totalReviews"`
	AverageRating float64 `gorm:"column:average_rating" json:"averageRating"`
}
