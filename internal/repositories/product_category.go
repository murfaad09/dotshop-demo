package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	pagination "github.com/harishash/dotshop-be/internal/utils/pagination"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id uint) (*domain.Category, error)
	FindAllCategory(page, limit int) ([]domain.Category, int64, error)
	FindCategoryById(id uint) (*domain.Category, error)
	FindSubCategoryById(id uint) (*dto.SubCategoryDTO, error)
	UpdateSubcategory(subcategory *domain.SubCategory) error
	DeleteSubcategory(id uint) (*domain.SubCategory, error)
	GetCatalogProducts(query *dto.CatalogProducts) ([]dto.FilteredCatalogProductResponse, *dto.Paging, error)
	CreateSubCategory(
		categoryId uint,
		subCategoryName string,
		products []*dto.CreateProductRequest) (
		*dto.CreateSubCategoryResponse,
		error)

	AddProductsInSubCategory(
		products []*dto.CreateProductRequest,
		categoryId, subCategoryId uint) (
		[]*dto.CreateProductResponse,
		error)
}

type categoryRepository struct {
	DB          *gorm.DB
	productRepo IProductRepository
}

func NewCategoryRepository(db *gorm.DB, productRepo IProductRepository) CategoryRepository {
	return &categoryRepository{DB: db, productRepo: productRepo}
}

func (r *categoryRepository) CreateCategory(category *domain.Category) error {
	dbCategory := domain.Category{}
	err := r.DB.Model(&dbCategory).Unscoped().Where("name = ?", category.Name).Find(&dbCategory).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if dbCategory.ID <= 0 {
		return r.DB.Create(category).Error
	}

	if dbCategory.ID > 0 && !dbCategory.DeletedAt.Time.IsZero() {
		r.DB.Model(&dbCategory).Update("deleted_at", nil)
		dbCategory.Name = category.Name

		if err := r.DB.Save(&dbCategory).Error; err != nil {
			return fmt.Errorf("failed to create category: %v", err)
		}
	} else {
		return errors.New("category already exists")
	}

	return nil
}

func (r *categoryRepository) UpdateCategory(category *domain.Category) error {
	return r.DB.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(id uint) (*domain.Category, error) {
	var subCategoryIDs []*uint
	if err := r.DB.Model(&domain.SubCategory{}).Unscoped().Where("category_id = ?", id).
		Pluck("id", &subCategoryIDs).
		Error; err != nil {
		return nil, err
	}

	for _, id := range subCategoryIDs {
		r.DeleteSubcategory(*id)
	}

	var category domain.Category
	if err := r.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Where("id = ?", id).Delete(&domain.Category{}).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAllCategory(page, limit int) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	r.DB.Model(&domain.Category{}).Count(&total)

	query := r.DB.Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
		return db.Order("sub_categories.name ASC")
	}).Order("categories.name ASC")
	result := pagination.Paginate(query, page, limit).Find(&categories)

	return categories, total, result.Error
}

func (r *categoryRepository) FindCategoryById(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.DB.Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
		return db.Order("sub_categories.name ASC")
	}).Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *categoryRepository) CreateSubCategory(
	categoryId uint,
	subCategoryName string,
	products []*dto.CreateProductRequest) (
	*dto.CreateSubCategoryResponse,
	error) {

	subCategory, err := createSubCategory(r.DB, categoryId, subCategoryName)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return &dto.CreateSubCategoryResponse{
			SubCategory: *subCategory,
		}, nil
	}

	insertedProducts, err := InsertProduct(r.DB, products, categoryId, subCategory.ID)
	if err != nil {
		r.DB.Table("sub_categories").Delete(subCategory).Where("id = ?", subCategory.ID)
		return nil, err
	}

	return &dto.CreateSubCategoryResponse{
		SubCategory: *subCategory,
		Products:    insertedProducts,
	}, nil

}

func (r *categoryRepository) AddProductsInSubCategory(
	products []*dto.CreateProductRequest,
	categoryId, subCategoryId uint) (
	[]*dto.CreateProductResponse,
	error) {

	insertedProducts, err := InsertProduct(r.DB, products, categoryId, subCategoryId)
	if err != nil {
		return nil, err
	}
	return insertedProducts, nil
}

func (r *categoryRepository) FindSubCategoryById(id uint) (*dto.SubCategoryDTO, error) {
	var subCategory domain.SubCategory
	err := r.DB.Where("id = ?", id).First(&subCategory).Error

	return &dto.SubCategoryDTO{
		ID:   subCategory.ID,
		Name: subCategory.Name,
	}, err
}

func createSubCategory(db *gorm.DB, categoryId uint, subCategoryName string) (*dto.SubCategoryDTO, error) {

	subCategory := &domain.SubCategory{
		CategoryID: categoryId,
		Name:       subCategoryName,
	}

	dbSubCategory := &domain.SubCategory{}
	err := db.Model(&dbSubCategory).Unscoped().Where("name = ?", subCategoryName).Find(&dbSubCategory).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if dbSubCategory.ID <= 0 {
		if err := db.Create(subCategory).Error; err != nil {
			return nil, err
		}

		return &dto.SubCategoryDTO{
			ID:   subCategory.ID,
			Name: subCategory.Name,
		}, nil
	}

	if dbSubCategory.ID > 0 && !dbSubCategory.DeletedAt.Time.IsZero() {
		db.Model(&dbSubCategory).Update("deleted_at", nil)
		dbSubCategory.CategoryID = categoryId
		dbSubCategory.Name = subCategoryName

		if err := db.Save(&dbSubCategory).Error; err != nil {
			return nil, fmt.Errorf("failed to create category: %v", err)
		}

		return &dto.SubCategoryDTO{
			ID:   dbSubCategory.ID,
			Name: dbSubCategory.Name,
		}, nil
	} else {
		return nil, errors.New("sub-category already exists")
	}
}

func (r *categoryRepository) UpdateSubcategory(subcategory *domain.SubCategory) error {
	result := r.DB.Model(&domain.SubCategory{}).
		Where("id = ? AND category_id = ?", subcategory.ID, subcategory.CategoryID).
		Update("name", subcategory.Name)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("subcategory not found")
	}

	return nil
}

func (r *categoryRepository) DeleteSubcategory(id uint) (*domain.SubCategory, error) {

	var subCategory *domain.SubCategory

	if err := r.DB.Where("id = ?", id).First(&subCategory).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var productIDs []*string
	// Product_ids
	if err := r.DB.Model(&domain.Product{}).Where("sub_category_id = ?", id).
		Pluck("product_id", &productIDs).
		Error; err != nil {
		return nil, err
	}

	for _, id := range productIDs {
		r.productRepo.DeleteProduct(*id)
	}

	if err := r.DB.Where("id = ?", id).Delete(&domain.SubCategory{}).Error; err != nil {
		return nil, err
	}

	return subCategory, nil
}

func (r *categoryRepository) GetCatalogProducts(query *dto.CatalogProducts) ([]dto.FilteredCatalogProductResponse, *dto.Paging, error) {
	validSubCategories := splitIds(query.SubCategoryIds)
	filter := func(db *gorm.DB) *gorm.DB {
		if len(query.SortBy) > 0 {
			db = pagination.SetCatalogProductsSort(db, query.SortBy)
		} else {
			db = db.Order("created_at DESC")
		}

		if len(query.Product) > 0 {
			db = db.Where("products.name ILIKE ?", "%"+query.Product+"%")
		}

		db.Where("products.is_active = ?", true)

		return db
	}

	return getCatalogProducts(r, r.DB, query, validSubCategories, filter)
}

func getCatalogProducts(r *categoryRepository, db *gorm.DB, query *dto.CatalogProducts, subCategories []string, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.FilteredCatalogProductResponse, *dto.Paging, error) {
	var totalCount int64
	var products []domain.Product

	dbQuery := db.Model(&domain.Product{}).Preload("Variants.VariantOptions").Preload("SubCategory")

	dbQuery = filterFunc(dbQuery)

	if len(query.SearchByBrandName) > 0 {
		dbQuery = dbQuery.Where("brand_name ILIKE ?", "%"+query.SearchByBrandName+"%")
	}

	if len(subCategories) > 0 {
		dbQuery = dbQuery.Where("sub_category_id IN ?", subCategories)
	}

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	result := pagination.NewPaginate(dbQuery, query.PageNum, query.PageSize).Find(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var resp []dto.FilteredCatalogProductResponse
	for _, product := range products {
		var category *domain.Category
		if product.SubCategoryID != nil {
			category, err = r.FindCategoryById(product.CategoryID)
			if err != nil {
				log.Printf("error finding category for id: %v error: %v", product.CategoryID, err)
			}
		}

		var variants []dto.VariantResponse
		for _, variant := range product.Variants {
			var variantsOption []*dto.VariantOptionResponse
			for _, variantOption := range variant.VariantOptions {
				variantsOption = append(variantsOption, &dto.VariantOptionResponse{
					Name:      variantOption.Name,
					Value:     variantOption.Value,
					VariantID: variantOption.VariantID,
				})
			}

			variants = append(variants, dto.VariantResponse{
				ID:              variant.ID,
				ProductID:       variant.ProductID,
				SKU:             variant.SKU,
				Title:           variant.Title,
				InventoryAmount: variant.InventoryAmount,
				Image:           variant.Image,
				RetailPrice:     variant.RetailPrice,
				RetailCurrency:  variant.RetailCurrency,
				BasePrice:       variant.BasePrice,
				BaseCurrency:    variant.BaseCurrency,
				Units:           variant.Units,
				Attributes:      variant.Attributes,
				VariantOptions:  variantsOption,
			})
		}

		var price float64
		if len(variants) > 0 {
			price = variants[0].RetailPrice
		}

		resp = append(resp, dto.FilteredCatalogProductResponse{
			ProductId:       product.ProductID,
			Name:            product.Name,
			Description:     product.Description,
			BrandName:       product.BrandName,
			Price:           price,
			Variants:        variants,
			CreatedAt:       product.CreatedAt,
			UpdatedAt:       product.UpdatedAt,
			CategoryID:      &product.CategoryID,
			CategoryName:    category.Name,
			SubCategoryID:   product.SubCategoryID,
			SubCategoryName: product.SubCategory.Name,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}
