package dto

import "time"

type PayoutStatus string

const (
	PAYOUT_INPROGRESS PayoutStatus = "IN PROGRESS"
	PAYOUT_SENT       PayoutStatus = "SENT"
	PAYOUT_REJECTED   PayoutStatus = "REJECTED"
)

type PayoutHistoryResponse struct {
	Id               uint64    `json:"id"`
	CuratorId        uint64    `json:"curatorId"`
	ReturnAmount     float64   `json:"returnAmount"`
	CommissionAmount float64   `json:"commissionAmount"`
	PayoutAmount     float64   `json:"payoutAmount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type PayoutResponse struct {
	ID               uint64    `json:"id"`
	CuratorID        uint64    `json:"curatorId"`
	ReturnAmount     float64   `json:"returnAmount"`
	CommissionAmount float64   `json:"commissionAmount"`
	PayoutAmount     float64   `json:"payoutAmount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ProductData struct {
	ProductName                 string  `json:"product_name"`
	ProductImage                string  `json:"product_image"`
	TotalMargin                 float64 `json:"total_margin"`
	TotalSales                  float64 `json:"total_sales"`
	DotshopProfit               float64 `json:"dotshop_profit"`
	CuratorCommission           float64 `json:"curator_commission"`
	DotShopProfitPercentage     float64 `json:"dotshop_profit_percentage"`
	CuratorCommissionPercentage float64 `json:"curator_commission_percentage"`
}
