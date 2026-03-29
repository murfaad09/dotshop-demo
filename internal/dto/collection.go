package dto

import (
	"encoding/json"
	"time"

	"github.com/harishash/dotshop-be/internal/utils/errors"

	domain "github.com/harishash/dotshop-be/internal/models"
)

type VariantResponse struct {
	ID              string                   `json:"id"`
	ProductID       string                   `json:"productId"`
	SKU             string                   `json:"sku"`
	Title           string                   `json:"title"`
	InventoryAmount int                      `json:"inventoryAmount"`
	Image           string                   `json:"image"`
	RetailPrice     float64                  `json:"retailPrice"`
	RetailCurrency  string                   `json:"retailCurrency"`
	BasePrice       float64                  `json:"basePrice"`
	BaseCurrency    string                   `json:"baseCurrency"`
	VariantOptions  []*VariantOptionResponse `json:"variantOptions"`
	Units           string                   `json:"units"`
	Attributes      json.RawMessage          `json:"attributes"`
}

type VariantOptionResponse struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	VariantID string `json:"variantId"`
}
type ProductResponse struct {
	ID           string             `json:"id"`
	BrandName    string             `json:"brandName"`
	SupplierName string             `json:"supplierName"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Category     string             `json:"category"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Variants     []*VariantResponse `json:"variants"`
}

type AddCollectionProductResponse struct {
	ProductID string `json:"productId"`
}
type AddLookProductResponse struct {
	ProductID string `json:"productId"`
}
type CollectionResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	TileColor   string             `json:"tile_color"`
	CreatedAt   *time.Time         `json:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at"`
	Products    []*ProductResponse `json:"products"`
	Sections    []*SectionResponse `json:"sections"`
}

type SectionResponse struct {
	ID           uint               `json:"id"`
	Name         *string            `json:"name"`
	ImageURL     *string            `json:"imageurl"`
	Description  *string            `json:"description"`
	CollectionID uint               `json:"collectionId"`
	Products     []*ProductResponse `json:"products"`
}
type CreateCollectionRequest struct {
	CuratorID   uint                   `json:"curator_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TileColor   string                 `json:"tile_color"`
	Products    []CreateProductRequest `json:"products"`
}

type CreateCollectionSectionRequest struct {
	Name         *string           `json:"name"`
	ImageURL     *string           `json:"imageurl"`
	Description  *string           `json:"description"`
	CollectionID uint              `json:"collectionId"`
	Products     []*ProductRequest `json:"products"`
}

type CreateCollectionSectionResponse struct {
	CollectionSectionID int                                `json:"collectionSectionId"`
	Name                *string                            `json:"name"`
	ImageURL            *string                            `json:"imageURL"`
	CollectionID        int                                `json:"collectionId"`
	Description         *string                            `json:"description"`
	Product             []*CreateCollectionProductResponse `json:"product"`
}

type UpdateCollectionSectionRequest struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageurl"`
	Description *string `json:"description"`
}

type UpdateCollectionSectionResponse struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageurl"`
	Description *string `json:"description"`
}

type CreateCollectionResponse struct {
	ID          uint                              `json:"id"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	TileColor   string                            `json:"tile_color"`
	Products    []CreateCollectionProductResponse `json:"products"`
}
type CreateCollectionProductResponse struct {
	ID           string    `json:"id"`
	BrandName    string    `json:"brandName"`
	SupplierName string    `json:"supplierName"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *CreateCollectionRequest) Validate() *errors.Error {
	if r.Name == "" {
		return errors.New("look name is required")
	}
	return nil
}

type AddProductToCollectionRequest struct {
	Products []ProductRequest `json:"products"`
}

type AddProductToCollectionResponse struct {
	CollectionId uint   `json:"collectionId"`
	ProductId    string `json:"productIds"`
}

func (r *AddProductToCollectionRequest) Validate() *errors.Error {
	if len(r.Products) <= 0 {
		return errors.New("at least one product is required")
	}
	return nil
}

func NewProduct(product ProductRequest) *domain.Product {
	return &domain.Product{
		ProductID:    product.ProductID,
		BrandName:    product.BrandName,
		SupplierName: product.SupplierName,
		Name:         product.Name,
		Description:  product.Description,
	}
}

func NewVariant(variant VariantRequest) *domain.Variant {
	return &domain.Variant{
		ID:              variant.ID,
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
	}
}

type UpdateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TileColor   string `json:"tileColor"`
}

type UpdateCollectionResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TileColor   string `json:"tileColor"`
}

type UpdateLookRequest struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageURL"`
}

type UpdateLookResponse struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	ImageURL         string `json:"imageURL"`
	EmbedLink        string `json:"embedLink"`
	SocialID         string `json:"socialId"`
	SocialTitle      string `json:"socialTitle"`
	SocialType       string `json:"socialType"`
	VideoDescription string `json:"videoDescription"`
}

type AddProductToSectionRequest struct {
	Products []*CreateProductRequest `json:"products"`
}

type AddProductToSectionResponse struct {
	SectionId uint                     `json:"sectionId"`
	Products  []*CreateProductResponse `json:"products"`
}

func ProductToResponse(products []domain.Product) []*ProductResponse {
	var productResponses []*ProductResponse
	for _, product := range products {
		var variants []*VariantResponse
		for _, variant := range product.Variants {

			var variantsOption []*VariantOptionResponse
			for _, variantOption := range variant.VariantOptions {
				variantsOption = append(variantsOption, &VariantOptionResponse{
					Name:      variantOption.Name,
					Value:     variantOption.Value,
					VariantID: variantOption.VariantID,
				})
			}

			variants = append(variants, &VariantResponse{
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
		productResponses = append(productResponses, &ProductResponse{
			ID:           product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
			Variants:     variants,
		})
	}
	return productResponses
}

type CollectionWithProducts struct {
	domain.Collection
	Products []domain.Product `json:"products"`
}

func CollectionToResponse(collection *domain.Collection) *CollectionResponse {
	return &CollectionResponse{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		TileColor:   collection.TileColor,
		Products:    ProductToResponse(collection.Products),
	}
}

func CollectionSectionsToResponse(collectionSections []domain.CollectionSection) []SectionResponse {
	var sectionResponses []SectionResponse

	for _, cs := range collectionSections {
		products := make([]domain.Product, len(cs.Products))
		for i, product := range cs.Products {
			products[i] = *product
		}

		sectionResponse := SectionResponse{
			ID:           cs.ID,
			Name:         cs.Name,
			ImageURL:     cs.ImageURL,
			Description:  cs.Description,
			CollectionID: cs.CollectionID,
			Products:     ProductToResponse(products),
		}

		sectionResponses = append(sectionResponses, sectionResponse)
	}

	return sectionResponses
}
