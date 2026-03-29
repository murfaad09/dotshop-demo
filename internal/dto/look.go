package dto

import (
	"time"

	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type CreateLookResponse struct {
	LookID           int                         `json:"lookId"`
	Name             string                      `json:"name"`
	ImageURL         string                      `json:"imageURL"`
	CuratorID        int                         `json:"curatorId"`
	EmbedLink        string                      `json:"embedLink"`
	SocialID         string                      `json:"socialId"`
	SocialTitle      string                      `json:"socialTitle"`
	VideoDescription string                      `json:"videoDescription"`
	SocialType       string                      `json:"socialType"`
	Product          []CreateLookProductResponse `json:"product"`
}
type CreateLookProductResponse struct {
	ID           string    `json:"id"`
	BrandName    string    `json:"brandName"`
	SupplierName string    `json:"supplierName"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateLookRequest struct {
	Name             string           `json:"name"`
	ImageURL         string           `json:"imageurl"`
	CuratorID        uint             `json:"curator_id"`
	EmbedLink        string           `json:"embedLink"`
	SocialID         string           `json:"socialId"`
	SocialTitle      string           `json:"socialTitle"`
	VideoDescription string           `json:"videoDescription"`
	SocialType       string           `json:"socialType"`
	Products         []ProductRequest `json:"products"`
}

func (r *CreateLookRequest) Validate() *errors.Error {
	if r.Name == "" {
		return errors.New("look name is required")
	}
	if r.ImageURL == "" {
		return errors.New("look image url is required")
	}

	return nil
}

type LooksResponse struct {
	ID               uint               `json:"id"`
	Name             string             `json:"name"`
	ImageURL         string             `json:"imageurl"`
	CuratorID        uint               `json:"curator_id"`
	EmbedLink        string             `json:"embedLink"`
	SocialID         string             `json:"socialId"`
	SocialTitle      string             `json:"socialTitle"`
	VideoDescription string             `json:"videoDescription"`
	SocialType       string             `json:"socialType"`
	Products         []*ProductResponse `json:"products"`
	CreatedAt        *time.Time         `json:"created_at"`
	UpdatedAt        *time.Time         `json:"updated_at"`
}

type AddProductToLookRequest struct {
	Products []ProductRequest `json:"products"`
}

func (r *AddProductToLookRequest) Validate() *errors.Error {
	if len(r.Products) <= 0 {
		return errors.New("at least one product is required")
	}
	return nil
}
