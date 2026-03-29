package constants

type PaymentAuthoriseRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentCreateRequest struct {
	PaymentAuthoriseRequest
	ClientID  int `json:"clientId"`
	ProductID int `json:"productId"`
}

type PaymentCaptureRequest struct {
	PaymentAuthoriseRequest
	OrderID string `json:"orderId"`
}
