package dto

type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	General []string `json:"general"`
}
