package dto

import "github.com/harishash/dotshop-be/internal/utils/errors"

type SearchLookParams struct {
	PagingParams
	Look string `query:"look"`
}

type SearchProductParams struct {
	PagingParams
	Product string `query:"product"`
}

type SearchCollectionParams struct {
	PagingParams
	Collection string `query:"collection"`
}

type SearchSectionParams struct {
	PagingParams
	Section string `query:"section"`
}

type Filter struct {
	SubCategoryFilterParams
	ProductFilterParams
	PagingParams
	SearchByBrandName string `query:"searchByBrandName"`
	SearchBy          string `query:"searchBy"`
	SortBy            string `query:"sort"`
}

type GlobalSearchResponse struct {
	Data        interface{} `json:"data"`
	Suggestions []string    `json:"suggestions"`
	Paging      Paging      `json:"paging"`
}

func (f *Filter) ValidateSortParam() error {
	switch f.SortBy {
	case "new_in":
		f.SortBy = "new_in"
	case "price_low_to_high":
		f.SortBy = "price_low_to_high"
	case "price_high_to_low":
		f.SortBy = "price_high_to_low"
	default:
		return errors.New("invalid sort parameter")
	}

	return nil
}
