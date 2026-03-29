package dto

import (
	"encoding/json"
	"time"

	models "github.com/harishash/dotshop-be/integration/vndr/convictional/buyer/dto"
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type VariantOption struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	VariantID string
}

type Variant struct {
	ID              string          `json:"id"`
	ProductID       string          `json:"product_id"`
	SKU             string          `json:"sku"`
	Title           string          `json:"title"`
	InventoryAmount int             `json:"inventoryAmount"`
	Image           string          `json:"image"`
	RetailPrice     float64         `json:"retailPrice"`
	RetailCurrency  string          `json:"retailCurrency"`
	BasePrice       float64         `json:"basePrice"`
	BaseCurrency    string          `json:"baseCurrency"`
	VariantOptions  []VariantOption `json:"variantOptions"`
	Units           string          `json:"units"`
	Attributes      json.RawMessage `json:"attributes"`
}

type Product struct {
	ID           string    `json:"id"`
	BrandName    string    `json:"brand_name"`
	SupplierName string    `json:"supplier_name"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Variants     []Variant `json:"variants"`
}

type CreateFeatureProductRequest struct {
	Products  []CreateProductRequest `json:"products"`
	CuratorID int                    `json:"curator_id"`
}
type CreateFeatureProductResponse struct {
	ProductID    string    `json:"productId"`
	BrandName    string    `json:"brandName"`
	SupplierName string    `json:"supplierNname"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	IsFeature    bool      `json:"is_feature"`
}

type FeatureProductResponse struct {
	Products  []Product `json:"products"`
	CuratorID int       `json:"curator_id"`
}

type FeatureVariantResponse struct {
	Variant        models.Variant         `json:"variant"`
	VariantOptions []models.VariantOption `json:"variant_options"`
}
type GetFeatureProductResponse struct {
	// FeatureID uint `json:"featureId"`
	Products []*ProductResponse
}

type AddProductToFeatureRequest struct {
	Products []ProductRequest `json:"products"`
}

type AddFeatureProductResponse struct {
	ProductID string `json:"productId"`
}

func (r *AddProductToFeatureRequest) Validate() *errors.Error {
	if len(r.Products) <= 0 {
		return errors.New("at least one product is required")
	}
	return nil
}
