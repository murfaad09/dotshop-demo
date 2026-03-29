package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/harishash/dotshop-be/internal/config"
	"github.com/harishash/dotshop-be/internal/dto"
)

type APIClient struct {
	baseURL string
	apiKey  string
}

var Client *APIClient

func init() {
	Client = NewAPIClient(config.GetConfig().ConvictionalBaseURL, config.GetConfig().BuyerAPIKey)
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *APIClient) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", c.apiKey)
	return http.DefaultClient.Do(req)
}

func (c *APIClient) CreateOrder(body *dto.CreateOrderRequest_Convictional) (*dto.CreateOrderResponse_Convictional, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, "orders")
	return post[dto.CreateOrderResponse_Convictional](c, url, body)
}

func (c *APIClient) CancelOrder(body *dto.CancelOrderRequest, orderID string) error {
	url := fmt.Sprintf("%s/orders/%s/cancel", c.baseURL, orderID)
	if _, err := post[interface{}](c, url, body); err != nil {
		return err
	}

	return nil
}

func GetVariants(c *APIClient, url string) (*dto.Convictional_Variants_Data, error) {
	data, err := get[dto.Convictional_Variants_Data](c, url)
	if err != nil {
		return nil, err
	}
	return data, nil
}


func GetPendingOrders(c *APIClient, url string) (*dto.Convictional_Order_Status_Data, error) {
	data, err := get[dto.Convictional_Order_Status_Data](c, url)
	if err != nil {
		return nil, err
	}
	return data, nil
}


func post[T any](c *APIClient, url string, body interface{}) (*T, error) {
	var response T

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp *dto.ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("failed to decode response error: %w", err)
		}

		return nil, errors.New(errResp.Errors.General[0])
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &response, nil
}

func (c *APIClient) Get(resource string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, resource)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get resource: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func get[T any](c *APIClient, url string) (*T, error) {
	var response T

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp *dto.ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("failed to decode response error: %w", err)
		}

		return nil, errors.New(errResp.Errors.General[0])
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &response, nil
}

func (c *APIClient) Update(resource string, data interface{}) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, resource)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update resource: %s", resp.Status)
	}

	return nil
}

func (c *APIClient) Delete(resource string) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, resource)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete resource: %s", resp.Status)
	}

	return nil
}

func (c *APIClient) GetAll(resource string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, resource)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get all resources: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
