package models

import "time"

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
	ID                string     `json:"id"`
	ProductID         string     `json:"productId"`
	VariantID         string     `json:"variantId"`
	BuyerReference    string     `json:"buyerReference"`
	SellerOrderID     *string    `json:"sellerOrderId"`
	SellerOrderItemID *string    `json:"sellerOrderItemId"`
	Quantity          int        `json:"quantity"`
	BasePrice         float32    `json:"basePrice"`
	RetailPrice       int        `json:"retailPrice"`
	BuyerVariantCode  string     `json:"buyerVariantCode"`
	Cancelled         bool       `json:"cancelled"`
	CancelledReason   string     `json:"cancelledReason"`
	CancelledBy       string     `json:"cancelledBy"`
	CancelledDate     *time.Time `json:"cancelledDate"`
}

type FulfillmentItem struct {
	ID          string `json:"id"`
	OrderItemID string `json:"orderItemId"`
	Quantity    int    `json:"quantity"`
}

type Fulfillment struct {
	ID           string            `json:"id"`
	Posted       bool              `json:"posted"`
	PostedDate   *time.Time        `json:"postedDate"`
	Created      time.Time         `json:"created"`
	Updated      time.Time         `json:"updated"`
	Carrier      string            `json:"carrier"`
	TrackingCode string            `json:"trackingCode"`
	TrackingUrls []string          `json:"trackingUrls"`
	Items        []FulfillmentItem `json:"items"`
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
	PostedDate        time.Time     `json:"postedDate"`
	Fulfilled         bool          `json:"fulfilled"`
	FulfilledDate     *time.Time    `json:"fulfilledDate"`
	Invoiced          bool          `json:"invoiced"`
	InvoicedDate      *time.Time    `json:"invoicedDate"`
	Created           time.Time     `json:"created"`
	Updated           time.Time     `json:"updated"`
	Address           Address       `json:"address"`
	Items             []Item        `json:"items"`
	Fulfillments      []Fulfillment `json:"fulfillments"`
	HasCancellations  bool          `json:"hasCancellations"`
	Attributes        struct{}      `json:"attributes"`
	IsTest            bool          `json:"isTest"`
	Flagged           bool          `json:"flagged"`
	FlaggedDate       time.Time     `json:"flaggedDate"`
	FlaggedMessage    string        `json:"flaggedMessage"`
	FlaggedAt         string        `json:"flaggedAt"`
}

type Order struct {
	HasMore  bool   `json:"hasMore"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Data     []struct {
		ID                string        `json:"id"`
		BuyerReference    string        `json:"buyerReference"`
		CustomerReference string        `json:"customerReference"`
		OrderedDate       time.Time     `json:"orderedDate"`
		Created           time.Time     `json:"created"`
		Updated           time.Time     `json:"updated"`
		Address           Address       `json:"address"`
		CustomerEmail     string        `json:"customerEmail"`
		Items             []Item        `json:"items"`
		ItemsNotAccepted  []interface{} `json:"itemsNotAccepted"`
		SellerOrders      []SellerOrder `json:"sellerOrders"`
		Metafields        struct{}      `json:"metafields"`
		Attributes        struct{}      `json:"attributes"`
		QuoteID           string        `json:"quoteId"`
		IsTest            bool          `json:"isTest"`
		Note              string        `json:"note"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

// Create Order dto
type CreateOrder struct {
	Address           Address          `json:"address"`
	BuyerReference    string           `json:"buyerReference"`
	CustomerEmail     string           `json:"customerEmail"`
	CustomerReference string           `json:"customerReference"`
	IsTest            bool             `json:"isTest"`
	Items             []Item           `json:"items"`
	Note              string           `json:"note"`
	OrderedDate       time.Time        `json:"orderedDate"`
	ShippingMethods   []ShippingMethod `json:"shippingMethods"`
}
type ShippingMethod struct {
	Price map[string]string `json:"price"`
	Code  string            `json:"code"`
	Title string            `json:"title"`
}
