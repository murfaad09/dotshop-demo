package dto

import (
	"time"

	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type CreateReviewRequest struct {
	Rating    int     `json:"rating" binding:"required"`
	Comment   *string `json:"comment"`
	UserID    uint    `json:"userId" binding:"required"`
	ProductID string  `json:"productId" binding:"required"`
	CuratorID uint    `json:"curatorId" binding:"required"`
}

func (r *CreateReviewRequest) Validate() *errors.Error {
	if r.ProductID == "" {
		return errors.New("productID is required")
	}
	if r.UserID == 0 {
		return errors.New("userID is required")
	}
	if r.CuratorID == 0 {
		return errors.New("curatorID is required")
	}
	if !(0 <= r.Rating && r.Rating <= 5) {
		return errors.New("rating must be greater than 0 and less than or equal to 5")
	}
	return nil
}

type UpdateReviewRequest struct {
	Rating  int     `json:"rating"`
	Comment *string `json:"comment"`
}

type ReviewResponse struct {
	ID        uint       `json:"id"`
	Rating    int        `json:"rating"`
	Comment   *string    `json:"comment"`
	UserID    uint       `json:"user_id"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	ProductID string     `json:"product_id"`
	CuratorID uint       `json:"curator_id"`
	Comments  []*Comment `json:"comments"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type DetailedReviewResponse struct {
	ID        uint             `json:"id"`
	Rating    int              `json:"rating"`
	Comment   *string          `json:"comment"`
	User      UserResponse     `json:"user"`
	Product   ProductResponses `json:"product"`
	Curator   CuratorResponse  `json:"curator"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ProductResponses struct {
	ProductID    string `json:"id"`
	BrandName    string `json:"brandName"`
	SupplierName string `json:"supplierName"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
}

type CuratorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ListProductReviewsRequest struct {
	PagingParams
	Stars []int `query:"stars"`
}
