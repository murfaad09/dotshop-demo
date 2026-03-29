package klaviyo

type ProfileDTO struct {
	Data DataDTO `json:"data"`
}

type DataDTO struct {
	Type       string        `json:"type"`
	Attributes AttributesDTO `json:"attributes"`
}

type AttributesDTO struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	General []string `json:"general"`
}

type CreateProfileResponse struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}
