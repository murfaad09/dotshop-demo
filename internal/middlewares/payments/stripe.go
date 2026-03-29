package payments

import (
	"time"

	"github.com/harishash/dotshop-be/internal/config"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
)

type Stripe struct {
	paymentIntent *stripe.PaymentIntent
}

func NewStripe() IPaymentGateway {
	stripe.Key = config.GetConfig().StripeSecretKey
	return &Stripe{
		paymentIntent: nil,
	}
}

func (s *Stripe) GetPublishableKey() (*constants.ClientKeyResponse, error) {
	publishableKey := config.GetConfig().StripePublicKey
	stripeSecretKey := config.GetConfig().StripeSecretKey

	if publishableKey == "" {
		return nil, nil
	}
	return &constants.ClientKeyResponse{
			StripeClientKey:    &publishableKey,
			StripeClientSecret: &stripeSecretKey},
		nil
}
func (s *Stripe) Authorize(amount *float64, currency *string) (*Transaction, error) {
	params := &stripe.ChargeParams{}
	sc := &client.API{}
	sc.Init(config.GetConfig().StripeSecretKey, nil)

	//TODO: This requires updated database schema
	//Get the payment intent from the database
	//ChargeID shoud come from the database
	//This can be called after the payment intent is captured

	chargeID := "ch_3Ln3j02eZvKYlo2C0d5IZWuG"

	charge, err := sc.Charges.Get(chargeID, params)
	if err != nil {
		return nil, err
	}
	created := time.Unix(charge.Created, 0)
	return &Transaction{
		ID:          &charge.ID,
		Status:      &charge.LastResponse.Status,
		ClientID:    nil,
		PaymentID:   nil,
		ChargeID:    &chargeID,
		Currency:    currency,
		Amount:      amount,
		Description: &charge.Description,
		CreatedAt:   &created,
	}, nil
}

// Create implements IPaymentGateway.
func (s *Stripe) Create(amount *float64, currency *string) (*constants.PaymentCreateIntentResponse, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(*amount)),
		Currency: stripe.String(getStripeCurrency(currency)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	//TODO: Save the payment intent in the database
	//TODO: This requires updated database schema

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return &constants.PaymentCreateIntentResponse{
		StripeClientSecret: &pi.ClientSecret,
		PaypalData:         nil,
	}, nil
}

// Capture implements IPaymentGateway.
func (s *Stripe) Capture(amount *float64, currency *string, orrderID *string) (*constants.PaymentCaptureResponse, error) {
	params := &stripe.PaymentIntentCaptureParams{}

	//TODO: This requires updated database schema
	//Get the payment intent from the database
	//save the latest_charge": "ch_1EXUPv2eZvKYlo2CStIqOmbY",

	pi := "pi_3Ln3j02eZvKYlo2C0d5IZWuG"

	result, err := paymentintent.Capture(pi, params)
	if err != nil {
		return nil, err
	}
	return &constants.PaymentCaptureResponse{
		StripeResponse: result,
		PaypalResponse: nil,
	}, nil
} // Capture implements IPaymentGateway.

// Refund implements IPaymentGateway.
func (s *Stripe) Refund(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error) {
	//TODO: This requires updated database schema
	//Get the payment intent from the database
	chargeID := "ch_3Ln3j02eZvKYlo2C0d5IZWuG"
	params := &stripe.RefundParams{
		Charge:        stripe.String(chargeID),
		PaymentIntent: stripe.String(*paymentID),
	}
	result, err := refund.New(params)
	if err != nil {
		return nil, err
	}
	created := time.Unix(result.Created, 0)
	return &Transaction{
		ID:          &result.ID,
		Status:      &result.LastResponse.Status,
		ClientID:    nil,
		PaymentID:   nil,
		ChargeID:    &chargeID,
		Currency:    currency,
		Amount:      amount,
		Description: &result.Description,
		CreatedAt:   &created,
	}, nil
}

// Cancel implements IPaymentGateway.
func (s *Stripe) Cancel(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error) {
	params := &stripe.PaymentIntentCancelParams{}
	//TODO: This requires updated database schema
	//Get the payment intent from the database
	pi := "pi_3Ln3j02eZvKYlo2C0d5IZWuG"
	result, err := paymentintent.Cancel(pi, params)
	if err != nil {
		return nil, err
	}
	created := time.Unix(result.Created, 0)
	return &Transaction{
		ID:          &result.ID,
		Status:      &result.LastResponse.Status,
		ClientID:    nil,
		Currency:    currency,
		Amount:      amount,
		Description: &result.Description,
		CreatedAt:   &created,
	}, nil
}

func getStripeCurrency(currency *string) string {
	if currency == nil {
		return string(stripe.CurrencyUSD)
	}
	return string(stripe.CurrencyUSD)
}
