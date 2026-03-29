package constants

import (
	"time"

	domain "github.com/harishash/dotshop-be/internal/models"
	"github.com/stripe/stripe-go/v72"
)

type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	UserId      int     `json:"user_id"`
	Email       string  `json:"email"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	RoleId      uint    `json:"role_id"`
	CuratorId   *uint   `json:"curator_id,omitempty"`
}

type ProductsResponse struct {
	Products []*domain.Product `json:"products"`
}

type ClientKeyResponse struct {
	PaypalAccessToken  *string `json:"paypalAccessToken,omitempty"`
	PaypalAppID        *string `json:"paypalAppId,omitempty"`
	PaypalClientID     *string `json:"paypalClientId,omitempty"`
	PaypalApiURL       *string `json:"paypalApiUrl,omitempty"`
	StripeClientKey    *string `json:"stripeClientKey,omitempty"`
	StripeClientSecret *string `json:"stripeClientSecret,omitempty"`
}

type PaymentCreateIntentResponse struct {
	StripeClientSecret *string              `json:"clientSecret"`
	PaypalData         *PaypalOrderResponse `json:"paypalData"`
}

type PaypalOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

type PaypalAuthResponse struct {
	Scope                 string         `json:"scope"`
	AccessToken           string         `json:"access_token"`
	TokenType             string         `json:"token_type"`
	AppID                 string         `json:"app_id"`
	ExpiresIn             int            `json:"expires_in"`
	SupportedAuthnSchemes []string       `json:"supported_authn_schemes"`
	Nonce                 string         `json:"nonce"`
	ClientMetadata        ClientMetadata `json:"client_metadata"`
}
type ClientMetadata struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	LogoURI     string   `json:"logo_uri"`
	Scopes      []string `json:"scopes"`
	UIType      string   `json:"ui_type"`
}

type PaymentCaptureResponse struct {
	StripeResponse *stripe.PaymentIntent  `json:"stripeResponse"`
	PaypalResponse *PaypalCaptureResponse `json:"paypalResponse"`
}

type PaypalCaptureResponse struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentSource struct {
		Paypal struct {
			Name struct {
				GivenName string `json:"given_name"`
				Surname   string `json:"surname"`
			} `json:"name"`
			EmailAddress string `json:"email_address"`
			AccountID    string `json:"account_id"`
		} `json:"paypal"`
	} `json:"payment_source"`
	PurchaseUnits []struct {
		ReferenceID string `json:"reference_id"`
		Shipping    struct {
			Address struct {
				AddressLine1 string `json:"address_line_1"`
				AddressLine2 string `json:"address_line_2"`
				AdminArea2   string `json:"admin_area_2"`
				AdminArea1   string `json:"admin_area_1"`
				PostalCode   string `json:"postal_code"`
				CountryCode  string `json:"country_code"`
			} `json:"address"`
		} `json:"shipping"`
		Payments struct {
			Captures []struct {
				ID     string `json:"id"`
				Status string `json:"status"`
				Amount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"amount"`
				SellerProtection struct {
					Status            string   `json:"status"`
					DisputeCategories []string `json:"dispute_categories"`
				} `json:"seller_protection"`
				FinalCapture              bool   `json:"final_capture"`
				DisbursementMode          string `json:"disbursement_mode"`
				SellerReceivableBreakdown struct {
					GrossAmount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"gross_amount"`
					PaypalFee struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"paypal_fee"`
					NetAmount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"net_amount"`
				} `json:"seller_receivable_breakdown"`
				CreateTime time.Time `json:"create_time"`
				UpdateTime time.Time `json:"update_time"`
				Links      []struct {
					Href   string `json:"href"`
					Rel    string `json:"rel"`
					Method string `json:"method"`
				} `json:"links"`
			} `json:"captures"`
		} `json:"payments"`
	} `json:"purchase_units"`
	Payer struct {
		Name struct {
			GivenName string `json:"given_name"`
			Surname   string `json:"surname"`
		} `json:"name"`
		EmailAddress string `json:"email_address"`
		PayerID      string `json:"payer_id"`
	} `json:"payer"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}
