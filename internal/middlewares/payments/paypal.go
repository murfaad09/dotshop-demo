package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/harishash/dotshop-be/internal/config"
	"github.com/harishash/dotshop-be/internal/constants"
)

type Paypal struct {
	clientID     *string
	clientSecret *string
	httpClient   *httpclient.Client
	apiURL       *string
}

func NewPaypal() *Paypal {
	paypalClientID := config.GetConfig().PaypalClientID
	paypalClientSecret := config.GetConfig().PaypalClientSecret
	paypalApiURL := config.GetConfig().PaypalApiURL
	timeout := 10000 * time.Millisecond

	return &Paypal{
		clientID:     &paypalClientID,
		clientSecret: &paypalClientSecret,
		httpClient:   httpclient.NewClient(httpclient.WithHTTPTimeout(timeout)),
		apiURL:       &paypalApiURL,
	}
}

func (p *Paypal) GetPublishableKey() (*constants.ClientKeyResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	body := bytes.NewBufferString(form.Encode())
	url := fmt.Sprintf("%s/v1/oauth2/token", *p.apiURL)

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(*p.clientID, *p.clientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token constants.PaypalAuthResponse
	err = json.Unmarshal(bodyBytes, &token)
	if err != nil {
		return nil, err
	}

	return &constants.ClientKeyResponse{
		PaypalAccessToken: &token.AccessToken,
		PaypalAppID:       &token.AppID,
		PaypalClientID:    p.clientID,
	}, nil
}

func (p *Paypal) Authorize(amount *float64, currency *string) (*Transaction, error) {
	return nil, nil
}

func (p *Paypal) Create(amount *float64, currency *string) (*constants.PaymentCreateIntentResponse, error) {
	getToken, err := p.GetPublishableKey()
	if err != nil {
		return nil, err
	}
	accessToken := *getToken.PaypalAccessToken
	ordersUrl := fmt.Sprintf("%s/v2/checkout/orders", *p.apiURL)

	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"value":         amount,
					"currency_code": currency,
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(payloadBytes))

	body := bytes.NewBuffer(payloadBytes)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	req, err := http.NewRequest(http.MethodPost, ordersUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var jsonResponse *constants.PaypalOrderResponse
	err = json.Unmarshal(resBody, &jsonResponse)
	if err != nil {
		return nil, err
	}

	return &constants.PaymentCreateIntentResponse{
		StripeClientSecret: nil,
		PaypalData:         jsonResponse,
	}, nil
}

func (p *Paypal) Capture(amount *float64, currency *string, orderID *string) (*constants.PaymentCaptureResponse, error) {
	getToken, err := p.GetPublishableKey()
	if err != nil {
		return nil, err
	}
	accessToken := *getToken.PaypalAccessToken
	ordersUrl := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", *p.apiURL, *orderID)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	req, err := http.NewRequest(http.MethodPost, ordersUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var jsonResponse *constants.PaypalCaptureResponse
	err = json.Unmarshal(resBody, &jsonResponse)
	if err != nil {
		return nil, err
	}
	return &constants.PaymentCaptureResponse{
		StripeResponse: nil,
		PaypalResponse: jsonResponse,
	}, nil
}

func (p *Paypal) Cancel(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error) {
	return nil, nil
}

func (p *Paypal) Refund(clientID *string, paymentID *string, amount *float64, currency *string) (*Transaction, error) {
	return nil, nil
}
