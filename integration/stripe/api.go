package stripe

import (
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/tax/calculation"
)

const (
	key = ""
)

type TaxReq struct {
	Amount     int64
	Currency   string
	Country    string
	PostalCode string
}

func PerformTaxCalculation(body TaxReq) (*int64, error) {

	stripe.Key = key

	params := &stripe.TaxCalculationParams{

		Currency: stripe.String(body.Currency),
		LineItems: []*stripe.TaxCalculationLineItemParams{
			{
				Amount:    stripe.Int64(body.Amount),
				TaxCode:   stripe.String("txcd_99999999"),
				Reference: stripe.String("L1"),
			},
		},
		CustomerDetails: &stripe.TaxCalculationCustomerDetailsParams{
			Address: &stripe.AddressParams{

				PostalCode: stripe.String(body.PostalCode),
				Country:    stripe.String(body.Country),
			},
			AddressSource: stripe.String(string(stripe.TaxCalculationCustomerDetailsAddressSourceShipping)),
		},
	}
	params.AddExpand("line_items")

	result, err := calculation.New(params)
	if err != nil {
		return nil, err
	}

	return &result.TaxAmountExclusive, nil
}
