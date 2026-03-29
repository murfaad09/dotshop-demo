package dto

import (
	"time"

	domain "github.com/harishash/dotshop-be/internal/models"
)

type CreateOrderResponse_Convictional struct {
	ID                string         `json:"id"`
	BuyerReference    string         `json:"buyerReference"`
	CustomerReference string         `json:"customerReference"`
	OrderedDate       time.Time      `json:"orderedDate"`
	Created           time.Time      `json:"created"`
	Updated           time.Time      `json:"updated"`
	Address           Address        `json:"address"`
	CustomerEmail     string         `json:"customerEmail"`
	Items             []Item         `json:"items"`
	Fulfillments      []FulFillments `json:"fulfillments"`
	ItemsNotAccepted  []interface{}  `json:"itemsNotAccepted"`
	SellerOrders      []SellerOrder  `json:"sellerOrders"`
	Metafields        struct{}       `json:"metafields"`
	Attributes        struct{}       `json:"attributes"`
	QuoteID           string         `json:"quoteId"`
	IsTest            bool           `json:"isTest"`
	Note              string         `json:"note"`
}

type Address struct {
	Name        string `json:"name"`
	AddressOne  string `json:"addressOne"`
	AddressTwo  string `json:"addressTwo"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Zip         string `json:"zip"`
	Company     string `json:"company"`
	PhoneNumber string `json:"phoneNumber"`
}

type Item struct {
	ID                string `json:"id"`
	ProductID         string `json:"productId"`
	VariantID         string `json:"variantId"`
	BuyerReference    string `json:"buyerReference"`
	SellerOrderID     string `json:"sellerOrderId"`
	SellerOrderItemID string `json:"sellerOrderItemId"`
	Quantity          uint   `json:"quantity"`
	RetailPrice       int    `json:"retailPrice"`
	BuyerVariantCode  string `json:"buyerVariantCode"`
}

type SellerItem struct {
	ID                string      `json:"id"`
	VariantID         string      `json:"variantId"`
	BuyerReference    string      `json:"buyerReference"`
	SellerOrderID     interface{} `json:"sellerOrderId"`
	SellerOrderItemID interface{} `json:"sellerOrderItemId"`
	Quantity          int         `json:"quantity"`
	BasePrice         float64     `json:"basePrice"`
	RetailPrice       int         `json:"retailPrice"`
	Cancelled         bool        `json:"cancelled"`
	CancelledReason   string      `json:"cancelledReason"`
	CancelledBy       string      `json:"cancelledBy"`
	CancelledDate     interface{} `json:"cancelledDate"`
}

type SellerOrder struct {
	ID                string        `json:"id"`
	BuyerOrderID      string        `json:"buyerOrderID"`
	BuyerReference    string        `json:"buyerReference"`
	SellerReference   string        `json:"sellerReference"`
	CustomerReference string        `json:"customerReference"`
	CompanyID         string        `json:"companyId"`
	BaseCurrency      string        `json:"baseCurrency"`
	PackingSlipURL    string        `json:"packingSlipUrl"`
	InvoiceID         string        `json:"invoiceId"`
	Posted            bool          `json:"posted"`
	PostedDate        interface{}   `json:"postedDate"`
	Fulfilled         bool          `json:"fulfilled"`
	FulfilledDate     interface{}   `json:"fulfilledDate"`
	Invoiced          bool          `json:"invoiced"`
	InvoicedDate      interface{}   `json:"invoicedDate"`
	Created           time.Time     `json:"created"`
	Updated           time.Time     `json:"updated"`
	Address           Address       `json:"address"`
	Items             []SellerItem  `json:"items"`
	Fulfillments      []interface{} `json:"fulfillments"`
	HasCancellations  bool          `json:"hasCancellations"`
	Attributes        struct{}      `json:"attributes"`
	IsTest            bool          `json:"isTest"`
	Flagged           bool          `json:"flagged"`
	FlaggedDate       time.Time     `json:"flaggedDate"`
	FlaggedMessage    string        `json:"flaggedMessage"`
	FlaggedAt         string        `json:"flaggedAt"`
}

type CreateOrderRequest_Convictional struct {
	Address        Address            `json:"address"`
	Items          []OrderItemRequest `json:"items"`
	BuyerReference string             `json:"buyerReference"`
	CustomerEmail  string             `json:"customerEmail"`
	IsTest         bool               `json:"isTest"`
	Note           string             `json:"note"`
	OrderedDate    time.Time          `json:"orderedDate"`
}

type ItemRequest struct {
	BuyerReference string `json:"buyerReference"`
	Quantity       uint   `json:"quantity"`
	VariantID      string `json:"variantId"`
}

type FulFillments struct {
	ID           string             `json:"id"`
	Posted       bool               `json:"posted"`
	PostedDate   interface{}        `json:"postedDate"`
	Created      time.Time          `json:"created"`
	Updated      time.Time          `json:"updated"`
	Carrier      string             `json:"carrier"`
	TrackingCode string             `json:"trackingCode"`
	TrackingUrls []interface{}      `json:"trackingUrls"`
	Items        []FulFillmentsItem `json:"items"`
}

type FulFillmentsItem struct {
	ID          string `json:"id"`
	OrderItemID string `json:"orderItemId"`
	Quantity    int    `json:"quantity"`
}

// TODO: UserId will remove from the request when auth workflow is completed
type OrderRequest struct {
	UserId         uint64             `json:"userId"`
	AddressID      uint64             `json:"addressId"`
	PaymentID      string             `json:"paymentId"`
	ShippingMethod string             `json:"shippingMethod"`
	Items          []OrderItemRequest `json:"items"`
	BuyerReference string             `json:"buyerReference"`
	Note           string             `json:"note"`
	IsTest         bool               `json:"isTest"`
}

type OrderItemRequest struct {
	CuratorID         uint64 `json:"curatorId"`
	BuyerReference    string `json:"buyerReference"`
	ProductID         string `json:"productId"`
	Quantity          uint   `json:"quantity"`
	VariantID         string `json:"variantId"`
	VariantOptionName string `json:"variantOptionName"`
	VariantSize       string `json:"variantSize"`
}

func ConvicationalOrderRequest(body *OrderRequest, userAddress *domain.ShippingInfo) *CreateOrderRequest_Convictional {
	address := NewAddress(userAddress)
	return &CreateOrderRequest_Convictional{
		Address:        *address,
		Items:          body.Items,
		BuyerReference: body.BuyerReference,
		Note:           body.Note,
		IsTest:         body.IsTest,
	}
}

func NewAddress(address *domain.ShippingInfo) *Address {
	return &Address{
		AddressOne:  address.AddressOne.String,
		AddressTwo:  address.AddressTwo.String,
		City:        address.City.String,
		State:       address.State.String,
		Country:     address.Country.String,
		Zip:         address.Zip.String,
		Company:     address.Company.String,
		PhoneNumber: address.PhoneNumber.String,
	}
}

func NewFilfullments(address *domain.ShippingInfo) *Address {
	return &Address{
		AddressOne:  address.AddressOne.String,
		AddressTwo:  address.AddressTwo.String,
		City:        address.City.String,
		State:       address.State.String,
		Country:     address.Country.String,
		Zip:         address.Zip.String,
		Company:     address.Company.String,
		PhoneNumber: address.PhoneNumber.String,
	}
}

type OrderResponse struct {
	UserId           uint64                      `json:"userId"`
	OrderId          string                      `json:"orderId"`
	Address          *AddConsumerAddressResponse `json:"address"`
	Items            []OrderItemRequest          `json:"items"`
	BuyerReference   string                      `json:"buyerReference"`
	Note             string                      `json:"note"`
	IsTest           bool                        `json:"isTest"`
	FulFillments     []FulFillments              `json:"fulfillments"`
	HasCancellations bool                        `json:"hasCancellations"`
	CreatedAt        time.Time                   `json:"created_at"`
}

type OrdersListResponse struct {
	ID             string                     `json:"ID"`
	UserID         uint64                     `json:"userID"`
	TotalAmount    float64                    `json:"totalAmount"`
	TotalQuantity  uint                       `json:"totalQuantity"`
	Status         string                     `json:"status"`
	BuyerReference string                     `json:"buyer_reference"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
	Note           string                     `json:"note"`
	Variant        []OrderListVariantResponse `json:"variant"`
}

type OrderListVariantResponse struct {
	ID                string  `gorm:"primaryKey" json:"id"`
	ProductID         string  `json:"product_id"`
	CuratorID         uint64  `json:"curatorID"`
	ProductName       string  `json:"productName"`
	Quantity          uint    `json:"quantity"`
	BrandName         string  `json:"brandName"`
	Price             float64 `json:"price"`
	Description       string  `json:"description"`
	SKU               string  `json:"sku"`
	Title             string  `json:"title"`
	Image             string  `json:"image"`
	RetailPrice       float64 `json:"retailPrice"`
	RetailCurrency    string  `json:"retailCurrency"`
	VariantOptionName string  `json:"variant_option_name"`
	VariantSize       string  `json:"variant_size"`
}

type OrderListVariantOptionResponse struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	VariantID string `json:"VariantID"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type CancelOrderResponse_Convictional struct {
	Error []interface{} `json:"error"`
}

type ReturnRequest struct {
	OrderId        string                      `json:"orderId"`
	Reason         string                      `json:"reason"`
	ReturnVariants []ReturnOrderVariantRequest `json:"returnVariants"`
}

type ReturnOrderVariantRequest struct {
	VariantId         string `json:"variantId"`
	SellerOrderId     string `json:"sellerOrderId"`
	SellerOrderItemId string `json:"sellerOrderItemId"`
	BuyerCode         string `json:"buyerCode"`
	Quantity          uint   `json:"qunatity"`
}

type ReturnResponse struct {
	Id                uint64    `json:"id"`
	UserId            uint      `json:"userId"`
	OrderId           string    `json:"orderId"`
	VariantId         string    `json:"variantId"`
	OrderVariantId    uint64    `json:"orderVariantId"`
	SellerOrderId     string    `json:"sellerOrderId"`
	SellerOrderItemId string    `json:"sellerOrderItemId"`
	BuyerCode         string    `json:"buyerCode"`
	Quantity          uint      `json:"qunatity"`
	Status            string    `json:"status"`
	Reason            string    `json:"reason"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type Convictional_Order_Status_Data struct {
	Data Data `json:"data"`
}

// Data represents a data block containing seller orders.
type Data struct {
	ID           string               `json:"id"`
	SellerOrders []SellerOrder_Status `json:"sellerOrders"`
}

// SellerOrder represents the details of a seller order.
type SellerOrder_Status struct {
	InvoiceID        string    `json:"invoiceId"`
	Posted           bool      `json:"posted"`
	PostedDate       time.Time `json:"postedDate"`
	Fulfilled        bool      `json:"fulfilled"`
	FulfilledDate    time.Time `json:"fulfilledDate"`
	Invoiced         bool      `json:"invoiced"`
	InvoicedDate     time.Time `json:"invoicedDate"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	HasCancellations bool      `json:"hasCancellations"`
}
