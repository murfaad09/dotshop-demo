package dto

import (
	"github.com/go-playground/validator/v10"
)

type AddWishlistItemRequest struct {
	CuratorID          uint    `json:"curatorId" validate:"required"`
	ProductID          string  `json:"productId" validate:"required"`
	ProductName        string  `json:"productName" validate:"required"`
	ProductImage       string  `json:"productImages" validate:"required,url"`
	ProductBrand       string  `json:"productBrand" validate:"required"`
	ProductPrice       string  `json:"productPrice" validate:"required,numeric"`
	VariantID          *string `json:"variantId"`
	VariantOptionName  string  `json:"variantOptionName"`
	VariantOptionValue string  `json:"variantOptionValue"`
}

type AddWishlistItemResponse struct {
	Success bool
	Message string
}

type GetWishlistResponse struct {
	Success bool
	Data    []*AddWishlistItemRequest
}

func ValidateWishlistItemRequest(request AddWishlistItemRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}
