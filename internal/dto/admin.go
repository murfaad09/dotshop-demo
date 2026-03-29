package dto

import (
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type Status string

const (
	ACTIVE   Status = "active"
	REJECTED Status = "rejected"
	PENDING  Status = "pending"
	BLOCKED  Status = "blocked"
)

type ChangeCuratorStatusRequest struct {
	Status Status `json:"status"`
}

func (r *ChangeCuratorStatusRequest) Validate() error {
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

type ReturnOrderStatusRequest struct {
	Status bool `json:"status"`
}

type BlockCustomerRequest struct {
	IsBlock bool `json:"isBlock"`
}

type BlockCuratorRequest struct {
	IsBlock bool `json:"isBlock"`
}
type OrderSalesRequest struct {
	TimeFilter
	PagingParams
	SortBy string `query:"sort"`
}

func (f *OrderSalesRequest) ValidateSortParam() error {
	switch f.SortBy {
	case "customer_name_asc":
		f.SortBy = "customer_name_asc"
	case "customer_name_desc":
		f.SortBy = "customer_name_desc"
	case "date_asc":
		f.SortBy = "date_asc"
	case "date_desc":
		f.SortBy = "date_desc"
	case "items_low_to_high":
		f.SortBy = "items_low_to_high"
	case "items_high_to_low":
		f.SortBy = "items_high_to_low"
	case "amount_low_to_high":
		f.SortBy = "amount_low_to_high"
	case "amount_high_to_low":
		f.SortBy = "amount_high_to_low"
	default:
		return errors.New("invalid sort parameter")
	}

	return nil
}

type CatalogProducts struct {
	PagingParams
	SubCategoryFilterParams
	SearchByBrandName string `query:"searchByBrandName"`
	Product           string `query:"product"`
	SortBy            string `query:"sort"`
}

func (f *CatalogProducts) ValidateSortParam() error {
	switch f.SortBy {
	case "product_name_asc":
		f.SortBy = "product_name_asc"
	case "product_name_desc":
		f.SortBy = "product_name_desc"
	case "brand_name_asc":
		f.SortBy = "brand_name_asc"
	case "brand_name_desc":
		f.SortBy = "brand_name_desc"
	case "discount_low_to_high":
		f.SortBy = "discount_low_to_high"
	case "discount_high_to_low":
		f.SortBy = "discount_high_to_low"
	case "price_low_to_high":
		f.SortBy = "price_low_to_high"
	case "price_high_to_low":
		f.SortBy = "price_high_to_low"
	default:
		return errors.New("invalid sort parameter")
	}

	return nil
}

type CustomersRequest struct {
	// TimeFilter
	PagingParams
	// SortBy string `query:"sort"`
}

type CuratorRequest struct {
	// TimeFilter
	PagingParams
	// SortBy string `query:"sort"`
}

type GetAllCuratorResponse struct {
	CuratorID    uint    `json:"curator_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	ShopName     string  `json:"shop_name"`
	Image        string  `json:"image"`
	Email        string  `json:"email"`
	Status       string  `json:"status"`
	TotalRevenue float64 `json:"total_revenue"`
	NoOfOrders   int     `json:"no_of_orders"`
	ItemsSold    int     `json:"items_sold"`
}

type CustomersOrderListRequest struct {
	PagingParams
}

type CuratorsOrderListRequest struct {
	PagingParams
}

type PaymentDistributionRequest struct {
	PagingParams
	CommissionType string `query:"commission_type"`
}
