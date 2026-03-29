package dto

import (
	"time"

	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type CartRequest struct {
	UserID   uint              `json:"userid" validate:"required"`
	CartItem []CartItemRequest `json:"cartitem"`
}

type CartItemRequest struct {
	CuratorID   uint    `json:"curatorid" validate:"required"`
	ProductID   string  `json:"prodcutid" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	ProductName string  `json:"productname" validate:"required"`
	BrandName   string  `json:"brandname" validate:"required"`
	VariantID   string  `json:"variantid" validate:"required"`
	ImageURL    string  `json:"imageURL" validate:"required"`
	Color       string  `json:"color" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	Quantity    uint    `json:"quantity" validate:"required"`
}

type CartResponse struct {
	CartID      uint               `json:"cartid"`
	UserID      uint               `json:"userid"`
	CartItem    []CartItemResponse `json:"cartitem"`
	SubTotal    float64            `json:"subtotal"`
	Tax         float64            `json:"tax"`
	TotalAmount float64            `json:"totalamount"`
}

type CartItemResponse struct {
	CartID      uint       `json:"cartid"`
	CuratorID   uint       `json:"curatorid"`
	ProductID   string     `json:"prodcutid"`
	Price       float64    `json:"price"`
	ProductName string     `json:"productname"`
	BrandName   string     `json:"brandname"`
	VariantID   string     `json:"variantid"`
	ImageURL    string     `json:"imageURL"`
	Color       string     `json:"color"`
	Size        string     `json:"size"`
	Quantity    uint       `json:"quantity"`
	DeletedAt   *time.Time `json:"deletedat"`
}
type CartItemQuantityRequest struct {
	Quantity uint `json:"quantity"`
}

type AddCartItemsRequest struct {
	CuratorID   int     `json:"curatorId" validate:"required"`
	ProductID   string  `json:"prodcutid" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	ProductName string  `json:"productname" validate:"required"`
	BrandName   string  `json:"brandname" validate:"required"`
	VariantID   string  `json:"variantid" validate:"required"`
	ImageURL    string  `json:"imageURL" validate:"required"`
	Color       string  `json:"color" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	Quantity    uint    `json:"quantity" validate:"required"`
}

func (r *CartItemQuantityRequest) Validate() *errors.Error {
	if r.Quantity == 0 {
		return errors.New("Quantity Cannot be Zero")
	}
	return nil
}
