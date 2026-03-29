package dto

import "time"

type GraphResponse struct {
	Data []GraphData `json:"data"`
}

type GraphData struct {
	Period  string  `json:"period"`
	Sales   float64 `json:"sales"`
	Revenue float64 `json:"revenue"`
}

type ReturnOrderGraphResponse struct {
	TotalOrders uint                    `json:"totalOrders"`
	Data        []ReturnOrdersGraphData `json:"data"`
}

type ReturnOrdersGraphData struct {
	OrderID   string `json:"orderId"`
	CreatedAt string `json:"createdAt"`
}

type OrdersGraphResponse struct {
	TotalOrders uint              `json:"totalOrders"`
	Data        []OrdersGraphData `json:"data"`
}

type OrdersGraphData struct {
	OrderID   string `json:"orderId"`
	CreatedAt string `json:"createdAt"`
}

type AvgOrderValueGraphResponse struct {
	AvgOrderValue float64             `json:"avgOrderValue"`
	Data          []AOVIntervalResult `json:"data"`
}

type AOVIntervalResult struct {
	IntervalStart     time.Time
	AverageOrderValue float64
}

type AOVIntervalResultResponse struct {
	IntervalStart       time.Time
	TotalOrderValue     float64
	TotalNumberOfOrders float64
}

type UnitsSoldPerOrder struct {
	IntervalStart  time.Time
	TotalUnitsSold float64
	TotalOrders    float64
}

type ReturnsGraphResponse struct {
	TotalReturns uint           `json:"totalReturns"`
	ReturnRate   float64        `json:"returnRate"`
	Data         []OrderReturns `json:"data"`
}

type OrderReturns struct {
	IntervalStart         time.Time
	TotalReturnedQuantity uint
}

type SalesGraphResponse struct {
	TotalSales float64               `json:"totalSales"`
	Data       []SalesIntervalResult `json:"data"`
}

type SalesIntervalResult struct {
	IntervalStart  time.Time
	TotalAmountSum float64
}

type OrderGraphResponse struct {
	TotalOrderCount uint         `json:"totalOrderCount"`
	Data            []OrderCount `json:"data"`
}
type OrderCount struct {
	IntervalStart time.Time
	OrderCount    uint
}

type UnitsSold struct {
	IntervalStart time.Time
	UnitsSold     uint
}

type RevenueGraphResponse struct {
	TotalRevenue float64                 `json:"totalRevenue"`
	Data         []RevenueIntervalResult `json:"data"`
}

type RevenueIntervalResult struct {
	IntervalStart time.Time
	TotalRevenue  float64
}

type AvgUnitsPerOrderGraphResponse struct {
	AvgUnitsPerOrder float64             `json:"avgUnitsPerOrder"`
	Data             []AUPIntervalResult `json:"data"`
}

type AUPIntervalResult struct {
	IntervalStart        time.Time
	AverageUnitsPerOrder float64
}

type UnitsSoldGraphResponse struct {
	TotalUnitsSold uint        `json:"totalUnitsSold"`
	Data           []UnitsSold `json:"data"`
}

type GetCuratorTopWishlistResponse struct {
	ProductID    string `json:"productId"`
	ProductName  string `json:"productName"`
	ProductImage string `json:"productImages"`
	ProductBrand string `json:"productBrand"`
	ProductPrice string `json:"productPrice"`
	Count        int64  `json:"count"`
}

type GetCuratorTopSellingProductResponse struct {
	ProductID    string  `json:"productId"`
	CuratorID    uint64  `json:"CuratorId"`
	ProductName  string  `json:"productName"`
	ProductImage string  `json:"productImage"`
	ProductBrand string  `json:"productBrand"`
	ProductPrice float64 `json:"productPrice"`
	Earnings     string  `json:"earnings"`
	Status       string  `json:"status"`
	Sales        uint    `json:"sales"`
}

type GetCuratorTopSellingBrandsResponse struct {
	CuratorID  uint64 `json:"CuratorId"`
	BrandName  string `json:"brandName"`
	BrandImage string `json:"brandImage"`
	Earnings   string `json:"earnings"`
	Sales      uint   `json:"sales"`
}

type CommonProductRequest struct {
	TimeFilter
	PagingParams
}

type GetCuratorTopPurchasersResponse struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	UserID         uint   `json:"userId"`
	CuratorID      uint   `json:"CuratorId"`
	TotalPurchases uint   `json:"totalPurchases"`
	TotalSpent     string `json:"totalSpent"`
}

type SaleByCategoryResponse struct {
	CategoryName string  `json:"category_name"`
	Percentage   float64 `json:"percentage"`
}

type SaleRequest struct {
	TimeFilter
}
