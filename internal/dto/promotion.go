package dto

import "time"

type ListPromotionsRequest struct {
	PagingParams
	Status     string  `query:"status"`
	StartValue float64 `query:"startValue"`
	EndValue   float64 `query:"endValue"`
}

type ListPromotionResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	DiscountCode  string    `json:"discountCode"`
	DiscountValue float64   `json:"discountValue"`
	ExpiryDate    time.Time `json:"expiryDate"`
	Status        string    `json:"status"`
	Rule          string    `json:"rule"`
}

type PromotionRequest struct {
	Name          string    `json:"name"`
	DiscountCode  string    `json:"discountCode"`
	DiscountValue float64   `json:"discountValue"`
	ExpiryDate    time.Time `json:"expiryDate" example:"2024-08-21T23:59:59Z"`
	Status        string    `json:"status"`
	Rule          string    `json:"rule"`
}

type PromotionResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	DiscountCode  string  `json:"discountCode"`
	DiscountValue float64 `json:"discountValue"`
	ExpiryDate    string  `json:"expiryDate"`
	Status        string  `json:"status"`
	Rule          string  `json:"rule"`
}

type ApplyBulkDiscountRequest struct {
	ProductIDs  []string `json:"productIds"`
	PromotionID uint     `json:"promotionId"`
}
