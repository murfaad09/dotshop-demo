package dto

import "time"

type GetTopCuratorsResponse struct {
	CuratorID    uint    `json:"curatorId"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	ProfileImage string  `json:"profileImage"`
	Earnings     string  `json:"earnings"`
	Sales        uint    `json:"sales"`
}

type OrderVariantResponse struct {
	ProductName  string  `json:"product_name"`
	Brand        string  `json:"brand"`
	Size         string  `json:"size"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	VariantImage string  `json:"variant_image"`
}

type GetOrderSalesResponse struct {
	ID                string                 `json:"ID"`
	UserID            uint64                 `json:"userID"`
	CustomerFirstName *string                `json:"customerFirstName"`
	CustomerLastName  *string                `json:"customerLastName"`
	CustomerAddress   *string                `json:"customerAddress"`
	CustomerCity      *string                `json:"customerCity"`
	CustomerState     *string                `json:"customerState"`
	CustomerCountry   *string                `json:"customerCountry"`
	CustomerZip       *string                `json:"customerZip"`
	ShippingAddress   *string                `json:"shippingAddress"`
	ShippingCity      *string                `json:"shippingCity"`
	ShippingState     *string                `json:"shippingState"`
	ShippingCountry   *string                `json:"shippingCountry"`
	ShippingZip       *string                `json:"shippingZip"`
	TotalAmount       float64                `json:"totalAmount"`
	TotalQuantity     uint                   `json:"totalQuantity"`
	Status            string                 `json:"status"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	ShippingMethod    string                 `json:"shippingMethod"`
	PaymentID         string                 `json:"paymentID"`
	Variants          []OrderVariantResponse `json:"variants"`
}

type GetOrderReturnsResponse struct {
	ID                uint               `json:"Id"`
	ReturnId          string             `json:"returnId"`
	UserID            uint64             `json:"userId"`
	OrderVariantId    uint               `json:"orderVariantId"`
	CuratorId         uint               `json:"curatorId"`
	CustomerFirstName *string            `json:"customerFirstName"`
	CustomerLastName  *string            `json:"customerLastName"`
	CustomerAddress   *string            `json:"customerAddress"`
	CustomerCity      *string            `json:"customerCity"`
	CustomerState     *string            `json:"customerState"`
	CustomerCountry   *string            `json:"customerCountry"`
	CustomerZip       *string            `json:"customerZip"`
	ShippingAddress   *string            `json:"shippingAddress"`
	ShippingCity      *string            `json:"shippingCity"`
	ShippingState     *string            `json:"shippingState"`
	ShippingCountry   *string            `json:"shippingCountry"`
	ShippingZip       *string            `json:"shippingZip"`
	Reason            string             `json:"reason"`
	TotalAmount       float64            `json:"totalAmount"`
	TotalQuantity     uint               `json:"totalQuantity"`
	Status            string             `json:"status"`
	OrderReturnDetail OrderReturnDetails `json:"orderReturnDetails"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type CustomerOrdersListResponse struct {
	ID              string                     `json:"ID"`
	UserID          uint64                     `json:"userID"`
	TotalAmount     float64                    `json:"totalAmount"`
	TotalQuantity   uint                       `json:"totalQuantity"`
	Status          string                     `json:"status"`
	ShippingAddress *string                    `json:"shippingAddress"`
	ShippingCity    *string                    `json:"shippingCity"`
	ShippingState   *string                    `json:"shippingState"`
	ShippingCountry *string                    `json:"shippingCountry"`
	ShippingZip     *string                    `json:"shippingZip"`
	BuyerReference  string                     `json:"buyer_reference"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	Note            string                     `json:"note"`
	Variant         []OrderListVariantResponse `json:"variant"`
}

type CuratorOrdersListResponse struct {
	ID                string                     `json:"ID"`
	UserID            uint64                     `json:"userID"`
	CustomerFirstName *string                    `json:"customerFirstName"`
	CustomerLastName  *string                    `json:"customerLastName"`
	TotalAmount       float64                    `json:"totalAmount"`
	TotalQuantity     uint                       `json:"totalQuantity"`
	Status            string                     `json:"status"`
	BuyerReference    string                     `json:"buyer_reference"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	Note              string                     `json:"note"`
	Variant           []OrderListVariantResponse `json:"variant"`
}

type OrderReturnDetails struct {
	OrderVariantId    uint      `json:"orderVariantId"`
	ProductId         string    `json:"productId"`
	BrandName         string    `json:"brandName"`
	ProductName       string    `json:"productName"`
	VariantOptionName string    `json:"variantOptionName"`
	VariantSize       string    `json:"variantSize"`
	VariantImage      string    `json:"variantImage"`
	Reason            string    `json:"reason"`
	TotalAmount       float64   `json:"totalAmount"`
	TotalQuantity     uint      `json:"totalQuantity"`
	CreatedAt         time.Time `json:"created_at"`
}

type GetAllCustomerResponse struct {
	ID            uint       `json:"ID"`
	Email         string     `json:"email"`
	FirstName     *string    `json:"firstName"`
	LastName      *string    `json:"lastName"`
	PhoneNumber   *string    `json:"phoneNumber"`
	LastOrderDate *time.Time `json:"lastOrderDate"`
	LifeTimeSpend float64    `json:"lifeTimeSpend"`
	Orders        int64      `json:"order"`
	Items         uint       `json:"items"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}
