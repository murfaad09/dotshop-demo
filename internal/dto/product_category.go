package dto

import (
	"time"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateSubCategoryRequest struct {
	Name     string                  `json:"name"`
	Products []*CreateProductRequest `json:"products"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID                 uint              `json:"id"`
	Name               string            `json:"name"`
	CountOfSubcategory uint              `json:"countOfSubcategory"`
	SubCategories      []*SubCategoryDTO `json:"subcategories"`
}

type SubCategoryDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"productCount"`
}

type CreateSubCategoryResponse struct {
	SubCategory SubCategoryDTO           `json:"subcategory"`
	Products    []*CreateProductResponse `json:"products"`
}

type ChangeCategoryRequest struct {
	CategoryID    uint  `json:"category_id"`
	SubCategoryID *uint `json:"sub_category_id"`
}

type ChangeCategoryResponse struct {
	CategoryID    uint                     `json:"category_id"`
	CategoryName  string                   `json:"category_name"`
	SubCategoryID *uint                    `json:"sub_category_id"`
	SubCategory   string                   `json:"sub_category"`
	Products      []*CreateProductResponse `json:"products"`
}

type UpdateSingleProductRequest struct {
	ProductName   string `json:"product_name"`
	Description   string `json:"description"`
	CategoryID    uint   `json:"category_id"`
	SubCategoryID *uint  `json:"sub_category_id"`
}

type UpdateSubcategoryRequest struct {
	UpdateCategoryRequest
}

type SubcategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FilteredCatalogProductResponse struct {
	ProductId       string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	BrandName       string            `json:"brandName"`
	Price           float64           `json:"price"`
	PromoName       string            `json:"promoName"`
	DiscountValue   int               `json:"discountValue"`
	Variants        []VariantResponse `json:"variants"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	CategoryID      *uint             `json:"categoryId"`
	CategoryName    string            `json:"categoryName"`
	SubCategoryID   *uint             `json:"subCategoryId"`
	SubCategoryName string            `json:"subCategoryName"`
}

type UpdateSingleProductResponse struct {
	ProductId     string `json:"productID"`
	ProductName   string `json:"productName"`
	Description   string `json:"description"`
	CategoryId    uint   `json:"category_id"`
	SubCategoryId *uint  `json:"sub_category_id"`
}
