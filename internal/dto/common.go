package dto

import (
	"time"

	"gorm.io/gorm"
)

type PagingParams struct {
	PageNum  int `query:"pageNum" validate:"omitempty,min=1"`
	PageSize int `query:"pageSize" validate:"omitempty,min=1,max=100"`
}

type Paging struct {
	TotalCount  int64 `json:"total_count" example:"100"`
	CurrentPage int   `json:"current_page" example:"1"`
	PageSize    int   `json:"page_size" example:"10"`
}

type Response struct {
	Data   interface{} `json:"data"`
	Paging Paging      `json:"paging"`
}

type SubCategoryFilterParams struct {
	SubCategoryIds string `query:"subCategoryIds"`
}

type ProductFilterParams struct {
	ProductIds string `query:"productIds"`
}

func QueryToPagingParams(query *PagingParams) PagingParams {
	if query.PageNum < 1 {
		query.PageNum = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	return PagingParams{
		PageNum:  query.PageNum,
		PageSize: query.PageSize,
	}
}

type TimeFilter struct {
	From *time.Time `query:"from" example:"2006-01-02T15:04:05Z"`
	To   *time.Time `query:"to" example:"2006-01-02T15:04:05Z"`
}

func (c *TimeFilter) ApplyTimeFilterToQuery(query *gorm.DB) *gorm.DB {
	return setDefault(c, query)
}

func (c *TimeFilter) ApplyTimeFilterToSpecificTableQuery(query *gorm.DB, tableName string) *gorm.DB {
	return setDefaultForSpecificTable(c, query, tableName)
}

func setDefault(c *TimeFilter, query *gorm.DB) *gorm.DB {
	if query == nil {
		return query
	}

	if c.From != nil && c.To != nil {
		query = query.Where("created_at BETWEEN ? AND ?", c.From, c.To)
	} else if c.From != nil {
		query = query.Where("created_at >= ?", c.From)
	} else if c.To != nil {
		query = query.Where("created_at <= ?", c.To)
	}

	return query
}

func setDefaultForSpecificTable(c *TimeFilter, query *gorm.DB, tableName string) *gorm.DB {
	if query == nil {
		return query
	}

	if c.From != nil && c.To != nil {
		query = query.Where(tableName+".created_at BETWEEN ? AND ?", c.From, c.To)
	} else if c.From != nil {
		query = query.Where(tableName+".created_at >= ?", c.From)
	} else if c.To != nil {
		query = query.Where(tableName+".created_at <= ?", c.To)
	}

	return query
}
