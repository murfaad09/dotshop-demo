package dto

type BrandsResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsActive         bool   `json:"is_active"`
	NumberOfProducts int    `json:"number_of_products"`
}

type BrandsRequest struct {
	PagingParams
}

type UpdateBrandStatusRequest struct {
	IsActive bool `json:"isActive"`
}
