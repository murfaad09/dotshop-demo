package klaviyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseURL          = "https://a.klaviyo.com/api/"
	CreateProfileURL = BaseURL + "profiles/"
	RevisionDate     = "2024-06-15"
)

type Klaviyo interface {
	CreateProfile(email, firstName, lastName string) (*CreateProfileResponse, error)
	AddProfileToList(listID, profileID string) error
}

type klaviyo struct {
	pk string
}

func NewKlaviyoAPI(pk string) Klaviyo {
	return &klaviyo{pk}
}

func (k *klaviyo) CallAPI(url, method string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", fmt.Sprintf("Klaviyo-API-Key %s", k.pk))
	req.Header.Add("REVISION", RevisionDate)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}

func (k *klaviyo) CreateProfile(email, firstName, lastName string) (*CreateProfileResponse, error) {
	data := ProfileDTO{
		Data: DataDTO{
			Type: "profile",
			Attributes: AttributesDTO{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
		},
	}

	resp, err := k.CallAPI(CreateProfileURL, "POST", data)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(respBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %v", err)
		}
		return nil, fmt.Errorf("API error: %v", errorResponse.Errors.General)
	}

	var createResponse CreateProfileResponse

	if err := json.Unmarshal(respBody, &createResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create profile response: %v", err)
	}

	return &createResponse, nil
}

func (k *klaviyo) AddProfileToList(listID, profileID string) error {
	url := fmt.Sprintf("%slists/%s/relationships/profiles/", BaseURL, listID)

	payload := map[string]interface{}{
		"data": []map[string]string{
			{
				"type": "profile",
				"id":   profileID,
			},
		},
	}

	resp, err := k.CallAPI(url, "POST", payload)
	if err != nil {
		return fmt.Errorf("failed to add profile to list: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}
