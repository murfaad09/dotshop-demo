package payments

import (
	"time"

	"github.com/harishash/dotshop-be/internal/constants"
)

type Transaction struct {
	ID          *string
	PaymentID   *string
	ChargeID    *string
	Status      *string
	ClientID    *string
	Currency    *string
	Description *string
	Amount      *float64
	CreatedAt   *time.Time
}
type IPaymentGateway interface {
	GetPublishableKey() (*constants.ClientKeyResponse, error)
	Authorize(amount *float64, currency *string) (*Transaction, error)
	Create(amount *float64, currency *string) (*constants.PaymentCreateIntentResponse, error)
	Capture(amount *float64, currency *string, orderID *string) (*constants.PaymentCaptureResponse, error)
	Cancel(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error)
	Refund(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error)
}

func NewPaymentGateway(gateway string) IPaymentGateway {
	switch gateway {
	case "stripe":
		return NewStripe()
	case "paypal":
		return NewPaypal()
	default:
		return nil
	}
}
