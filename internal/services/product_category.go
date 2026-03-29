package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	dto "github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	repositories "github.com/harishash/dotshop-be/internal/repositories"
)

type CategoryService interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) (*dto.CategoryResponse, error)
	GetAllCategories(page, limit int) ([]*dto.CategoryResponse, int64, error)
	GetCategory(id uint) (*dto.CategoryResponse, error)
	DeleteSubCategory(id uint) (*dto.SubCategoryDTO, error)
	GetCatalogProducts(query *dto.CatalogProducts) (*dto.Response, error)
	UpdateSingleProduct(productID string, request dto.UpdateSingleProductRequest) (*dto.UpdateSingleProductResponse, error)
	AddSubCategory(
		categoryId uint,
		subCategoryName string,
		products []*dto.CreateProductRequest) (
		*dto.CreateSubCategoryResponse,
		error)

	UpdateSubCategory(
		categoryId uint,
		subCategoryId uint,
		products []*dto.CreateProductRequest) (
		*dto.CreateSubCategoryResponse,
		error)

	UpdateProductCategory(
		productIDs []string,
		categories dto.ChangeCategoryRequest) (
		*dto.ChangeCategoryResponse,
		error)

	UpdateSubcategoryName(
		categoryId uint,
		subcategoryId uint,
		body *dto.UpdateSubcategoryRequest) (
		*dto.SubcategoryResponse,
		error)
}

type categoryService struct {
	repo        repositories.CategoryRepository
	productRepo repositories.IProductRepository
}

func NewCategoryService(
	repo repositories.CategoryRepository,
	productRepo repositories.IProductRepository) CategoryService {
	return &categoryService{repo: repo, productRepo: productRepo}
}

func (s *categoryService) CreateCategory(
	req *dto.CreateCategoryRequest) (
	*dto.CategoryResponse, error) {
	category := &domain.Category{
		Name: req.Name,
	}

	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *categoryService) UpdateCategory(
	id uint,
	req *dto.UpdateCategoryRequest) (
	*dto.CategoryResponse,
	error) {
	category := &domain.Category{
		ID:   id,
		Name: req.Name,
	}

	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *categoryService) DeleteCategory(id uint) (*dto.CategoryResponse, error) {
	category, err := s.repo.DeleteCategory(id)
	if err != nil {
		return nil, err
	}

	deletedCategory := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
	return deletedCategory, nil
}

func (s *categoryService) GetAllCategories(page, limit int) ([]*dto.CategoryResponse, int64, error) {
	categories, total, err := s.repo.FindAllCategory(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var categoryResponses []*dto.CategoryResponse
	for _, category := range categories {
		subCategories := make([]*dto.SubCategoryDTO, len(category.SubCategories))
		for i, subCategory := range category.SubCategories {
			count, err := s.productRepo.GetProductCountWithSubCategoryID(subCategory.ID)
			if err != nil {
				log.Warnf("error while getting product count with subcategory id %d: %v", subCategory.ID, err)
			}

			subCategories[i] = &dto.SubCategoryDTO{
				ID:           subCategory.ID,
				Name:         subCategory.Name,
				ProductCount: count,
			}
		}

		categoryResponses = append(categoryResponses, &dto.CategoryResponse{
			ID:                 category.ID,
			Name:               category.Name,
			CountOfSubcategory: uint(len(category.SubCategories)),
			SubCategories:      subCategories,
		})
	}

	return categoryResponses, total, nil
}

func (s *categoryService) GetCategory(id uint) (*dto.CategoryResponse, error) {

	category, err := s.repo.FindCategoryById(id)

	if err != nil {
		return nil, err
	}

	subCategories := getSubCategory(category)

	return &dto.CategoryResponse{
		ID:                 category.ID,
		Name:               category.Name,
		CountOfSubcategory: uint(len(category.SubCategories)),
		SubCategories:      subCategories,
	}, nil
}

func (s *categoryService) AddSubCategory(
	categoryId uint,
	subCategoryName string,
	products []*dto.CreateProductRequest) (
	*dto.CreateSubCategoryResponse,
	error) {

	subCategory, err := s.repo.CreateSubCategory(categoryId, subCategoryName, products)
	if err != nil {
		return nil, err
	}
	return subCategory, nil
}

func (s *categoryService) UpdateSubCategory(
	categoryId uint,
	subCategoryId uint,
	products []*dto.CreateProductRequest) (
	*dto.CreateSubCategoryResponse,
	error) {
	subCategoryProducts, err := s.repo.AddProductsInSubCategory(products, categoryId, subCategoryId)
	if err != nil {
		return nil, err
	}

	subCategory, err := s.repo.FindSubCategoryById(subCategoryId)
	if err != nil {
		return nil, err
	}

	response := &dto.CreateSubCategoryResponse{
		SubCategory: *subCategory,
		Products:    subCategoryProducts,
	}
	return response, nil
}

func (s *categoryService) UpdateProductCategory(
	productIDs []string,
	categories dto.ChangeCategoryRequest) (
	*dto.ChangeCategoryResponse,
	error) {

	var updatedProducts []*dto.CreateProductResponse

	for _, id := range productIDs {
		product, err := s.productRepo.GetProductById(id)
		if err != nil {
			return nil, err
		}

		product.CategoryID = categories.CategoryID
		product.SubCategoryID = categories.SubCategoryID

		updatedProduct, err := s.productRepo.UpdateProduct(product)
		if err != nil {
			return nil, err
		}

		updatedProducts = append(updatedProducts, updatedProduct)
	}

	category, err := s.repo.FindCategoryById(categories.CategoryID)

	if err != nil {
		return nil, err
	}

	subCategories := getSubCategory(category)

	for _, subcategory := range subCategories {
		if subcategory.ID == *categories.SubCategoryID {
			return &dto.ChangeCategoryResponse{
				CategoryName:  category.Name,
				CategoryID:    category.ID,
				SubCategoryID: categories.SubCategoryID,
				SubCategory:   subcategory.Name,
				Products:      updatedProducts,
			}, nil
		}
	}

	// Add a return statement here to handle the case when no subcategory matches
	return nil, fmt.Errorf("no subcategory found with ID %d", *categories.SubCategoryID)
}

func getSubCategory(category *domain.Category) []*dto.SubCategoryDTO {
	subCategories := make([]*dto.SubCategoryDTO, len(category.SubCategories))
	for i, subCategory := range category.SubCategories {
		subCategories[i] = &dto.SubCategoryDTO{
			ID:   subCategory.ID,
			Name: subCategory.Name,
		}
	}
	return subCategories
}

func (s *categoryService) UpdateSubcategoryName(
	categoryId uint,
	subcategoryId uint,
	body *dto.UpdateSubcategoryRequest) (
	*dto.SubcategoryResponse,
	error) {
	subcategoryDB, err := s.repo.FindSubCategoryById(subcategoryId)
	if err != nil {
		return nil, err
	}

	if subcategoryDB == nil {
		return nil, fmt.Errorf("subcategory not found")
	}

	subcategory := &domain.SubCategory{
		ID:         subcategoryId,
		Name:       body.Name,
		CategoryID: categoryId,
	}

	if err := s.repo.UpdateSubcategory(subcategory); err != nil {
		return nil, err
	}

	return &dto.SubcategoryResponse{
		ID:   subcategory.ID,
		Name: subcategory.Name,
	}, nil
}

func (s *categoryService) DeleteSubCategory(id uint) (*dto.SubCategoryDTO, error) {

	subCategory, err := s.repo.DeleteSubcategory(id)
	if err != nil {
		return nil, err
	}

	deletedSubCategory := &dto.SubCategoryDTO{
		ID:   subCategory.ID,
		Name: subCategory.Name,
	}
	return deletedSubCategory, nil
}

func (s *categoryService) GetCatalogProducts(query *dto.CatalogProducts) (*dto.Response, error) {
	data, paging, err := s.repo.GetCatalogProducts(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *categoryService) UpdateSingleProduct(productID string, request dto.UpdateSingleProductRequest) (*dto.UpdateSingleProductResponse, error) {
	product, err := s.productRepo.GetProductById(productID)
	if err != nil {
		return nil, err
	}

	product.Name = request.ProductName
	product.Description = request.Description
	product.CategoryID = request.CategoryID
	product.SubCategoryID = request.SubCategoryID

	updatedProduct, err := s.productRepo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	resp := &dto.UpdateSingleProductResponse{
		ProductId:     updatedProduct.ProductID,
		ProductName:   updatedProduct.Name,
		Description:   updatedProduct.Description,
		CategoryId:    updatedProduct.CategoryID,
		SubCategoryId: updatedProduct.SubCategoryID,
	}

	return resp, nil
}
