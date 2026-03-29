package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	pagination "github.com/harishash/dotshop-be/internal/utils/pagination"
	"gorm.io/gorm"
)

type IProductRepository interface {
	GetProductsWithFilter(query *dto.Filter) ([]dto.FilteredProductResponse, *dto.Paging, error)
	GetProductById(id string) (*domain.Product, error)
	CreateProduct(product domain.Product) (*domain.Product, error)
	UpdateProduct(product *domain.Product) (*dto.CreateProductResponse, error)
	GetProductWithFirstVariant(productId string) (*domain.Product, *domain.Variant, error)
	DeleteProduct(id string) (*domain.Product, error)
	GetBrandImageByBrandName(brandName string) (string, error)
	GetBrands(searchStr string) ([]*dto.Brands, error)
	GetProductByVariantId(variantID string) (*domain.Product, *domain.Variant, error)
	GetProductIds() []string
	GetBrandName(productId string) string
	GetProductCountWithSubCategoryID(subCategories uint) (int64, error)
	UpdateProductStatusByBrandID(brandID uint, isActive bool) error
	GetVariantsByProductId(productID string) ([]*domain.Variant, error)
	GetAllProductsWithStats(params *dto.ListProductReviewsRequest) ([]*dto.ProductWithStats, *dto.Paging, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) IProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) GetProductsWithFilter(query *dto.Filter) ([]dto.FilteredProductResponse, *dto.Paging, error) {
	validSubCategories := splitIds(query.SubCategoryIds)
	validProductIds := splitIds(query.ProductIds)

	products, paging, err := GetProductsWithFilter(r.db, query.SearchBy, query.SearchByBrandName, query.SortBy, validSubCategories, validProductIds, query.PagingParams)
	if err != nil {
		return nil, nil, err
	}

	return dto.RowToFilteredProductResponse(products), paging, err
}

func splitIds(ids string) []string {
	subCategories := strings.Split(ids, ",")
	var validIds []string
	for _, id := range subCategories {
		if id != "" {
			validIds = append(validIds, id)
		}
	}
	return validIds
}

func GetProductsWithFilter(db *gorm.DB, searchBy, brandName, sortBy string, subCategories, productIds []string, pagingParams dto.PagingParams) ([]domain.Product, *dto.Paging, error) {
	var products []domain.Product
	var totalCount int64

	query := db.Model(&domain.Product{}).Preload("Variants.VariantOptions").Where("is_active = ?", true).Preload("SubCategory")

	if len(productIds) > 0 {
		query = query.Where("products.product_id IN ?", productIds)
	}

	if searchBy != "" {
		query = query.Where("name ILIKE ?", "%"+searchBy+"%")
	}

	if brandName != "" {
		query = query.Where("brand_name ILIKE ?", "%"+brandName+"%")
	}

	if len(subCategories) > 0 {
		query = query.Where("sub_category_id IN ?", subCategories)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	if len(sortBy) > 0 {
		query = pagination.SetProductSort(query, sortBy)
	}

	result := pagination.NewPaginate(query, pagingParams.PageNum, pagingParams.PageSize).Find(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    pagingParams.PageSize,
		CurrentPage: pagingParams.PageNum,
	}

	return products, &paging, nil
}

func (r *productRepo) GetProductById(id string) (*domain.Product, error) {
	product := domain.Product{}
	result := r.db.Where("product_id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *productRepo) CreateProduct(product domain.Product) (*domain.Product, error) {
	results := r.db.Create(&product)
	if results.Error != nil {
		return nil, results.Error
	}
	return &product, nil
}

func (r *productRepo) UpdateProduct(product *domain.Product) (*dto.CreateProductResponse, error) {
	results := r.db.Save(&product)
	if results.Error != nil {
		return nil, results.Error
	}

	productResponse := &dto.CreateProductResponse{
		ProductID:     product.ProductID,
		BrandName:     product.BrandName,
		SupplierName:  product.SupplierName,
		Name:          product.Name,
		Description:   product.Description,
		CreatedAt:     product.CreatedAt,
		CategoryID:    product.CategoryID,
		SubCategoryID: product.SubCategoryID,
	}
	return productResponse, nil
}

func (r *productRepo) deleteVariants(productID string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var variants []domain.Variant
	if err := tx.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		tx.Rollback()
		return err
	}

	var variantIDs []string
	for _, variant := range variants {
		variantIDs = append(variantIDs, variant.ID)
	}

	if err := tx.Where("variant_id IN ?", variantIDs).Delete(&domain.VariantOption{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("variant_id IN ?", variantIDs).Delete(&domain.OrderVariants{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("product_id = ?", productID).Delete(&domain.Variant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepo) DeleteProduct(id string) (*domain.Product, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	product := &domain.Product{}
	result := tx.Where("product_id = ?", id).First(&product)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if product == nil {
		tx.Rollback()
		return nil, errors.New("product not found")
	}

	// if err := r.deleteVariants(id); err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_product_id = ?", id).Delete(&domain.LookProduct{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_id = ?", id).Delete(&domain.CollectionProduct{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_id = ?", id).Delete(&domain.CollectionSectionProduct{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_id = ?", id).Delete(&domain.CuratorProduct{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_id = ?", id).Delete(&domain.WishlistItem{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Where("product_id = ?", id).Delete(&domain.Review{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	if err := tx.Where("product_id = ?", id).Delete(&domain.Product{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return product, nil
}

// InsertProduct begins a transaction to insert a product
func InsertProduct(r *gorm.DB, products []*dto.CreateProductRequest, categoryId, subCategoryId uint) ([]*dto.CreateProductResponse, error) {
	// Begin a transaction
	tx := r.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer tx.Rollback()

	var insertedProducts []*dto.CreateProductResponse

	// Iterate over the list of products
	for _, product := range products {
		var brandID uint
		var err error

		// Check if BrandName is provided and fetch or create brandID
		if product.BrandName != "" {
			brandID, err = getOrCreateBrandID(tx, product.BrandName)
			if err != nil {
				return nil, fmt.Errorf("failed to get brandID: %v", err)
			}
		}

		// Process each product and insert it
		productResponse, err := processProduct(tx, product, categoryId, subCategoryId, brandID)
		if err != nil {
			return nil, err
		}

		insertedProducts = append(insertedProducts, productResponse)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// If commit succeeds, cancel the deferred rollback
	return insertedProducts, nil
}

// processProduct handles the creation or updating of a product
func processProduct(
	tx *gorm.DB,
	product *dto.CreateProductRequest,
	categoryId, subCategoryId, brandID uint) (
	*dto.CreateProductResponse,
	error) {

	var existingProduct *domain.Product

	err := tx.Unscoped().Where("product_id = ?", product.ProductID).Find(&existingProduct).Error

	if len(existingProduct.ProductID) <= 0 {
		return createNewProduct(tx, product, categoryId, subCategoryId, brandID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check product existence: %v", err)
	}

	if len(existingProduct.ProductID) > 0 && !existingProduct.DeletedAt.Time.IsZero() {

		tx.Model(&existingProduct).Update("deleted_at", nil)

		existingProduct.Name = product.Name
		existingProduct.Description = product.Description
		existingProduct.BrandName = product.BrandName
		existingProduct.CategoryID = categoryId
		existingProduct.SubCategoryID = &subCategoryId
		existingProduct.BrandID = &brandID
		existingProduct.UpdatedAt = time.Now()
		existingProduct.Tags = product.Tags

		if err := tx.Save(&existingProduct).Error; err != nil {
			return nil, fmt.Errorf("failed to update product: %v", err)
		}
	} else if len(existingProduct.ProductID) > 0 {
		var categoryName string
		err := tx.Table("categories").
			Select("name").
			Where("id = ?", existingProduct.CategoryID).
			Limit(1).Scan(&categoryName).Error
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve category information: %v", err)
		}

		var subCategoryName string
		err = tx.Table("sub_categories").
			Select("name").
			Where("id = ?", existingProduct.SubCategoryID).
			Limit(1).
			Scan(&subCategoryName).Error
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve subcategory information: %v", err)
		}

		return nil, fmt.Errorf("product '%s' already exists in category '%s' with subcategory '%s'",
			product.Name,
			categoryName,
			subCategoryName)
	}

	return &dto.CreateProductResponse{
		ProductID:     existingProduct.ProductID,
		BrandName:     existingProduct.BrandName,
		SupplierName:  existingProduct.SupplierName,
		Name:          existingProduct.Name,
		Description:   existingProduct.Description,
		CreatedAt:     existingProduct.CreatedAt,
		CategoryID:    existingProduct.CategoryID,
		SubCategoryID: existingProduct.SubCategoryID,
	}, nil
}

// createNewProduct creates a new product and handles its variants
func createNewProduct(
	tx *gorm.DB,
	product *dto.CreateProductRequest,
	categoryId, subCategoryId, brandID uint) (
	*dto.CreateProductResponse,
	error) {

	now := time.Now()
	newProduct := &domain.Product{
		ProductID:     product.ProductID,
		BrandName:     product.BrandName,
		SupplierName:  product.SupplierName,
		Name:          product.Name,
		Description:   product.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
		Tags:          product.Tags,
		CategoryID:    categoryId,
		SubCategoryID: &subCategoryId,
		BrandID:       &brandID,
	}

	if err := tx.Create(&newProduct).Error; err != nil {
		return nil, fmt.Errorf("failed to insert product: %v", err)
	}

	var curatorId uint
	if err := tx.Table("curators").Where("shop_name = ?", "DotShop").Pluck("id", &curatorId).Error; err != nil {
		return nil, err
	}

	featureProduct := &domain.CuratorProduct{
		ProductID: product.ProductID,
		IsFeature: true,
		CuratorID: curatorId,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := tx.Create(featureProduct).Error; err != nil {
		return nil, fmt.Errorf("failed to insert feature product: %v", err)
	}

	if err := handleVariants(tx, product); err != nil {
		return nil, err
	}

	return &dto.CreateProductResponse{
		ProductID:     newProduct.ProductID,
		BrandName:     newProduct.BrandName,
		SupplierName:  newProduct.SupplierName,
		Name:          newProduct.Name,
		Description:   newProduct.Description,
		CreatedAt:     newProduct.CreatedAt,
		CategoryID:    newProduct.CategoryID,
		SubCategoryID: newProduct.SubCategoryID,
	}, nil
}

// handleVariants handles the creation of variants for a product
func handleVariants(tx *gorm.DB, product *dto.CreateProductRequest) error {
	for _, variant := range product.Variants {
		newVariant := &domain.Variant{
			ProductID:       product.ProductID,
			ID:              variant.ID,
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
		}
		if err := tx.Create(&newVariant).Error; err != nil {
			return fmt.Errorf("failed to insert variant: %v", err)
		}

		if err := handleVariantOptions(tx, newVariant.ID, variant.VariantOptions); err != nil {
			return err
		}
	}
	return nil
}

// handleVariantOptions handles the creation of variant options for a variant
func handleVariantOptions(
	tx *gorm.DB,
	variantID string,
	options []dto.CreateVariantOptionRequest,
) error {
	for _, option := range options {
		newOption := &domain.VariantOption{
			VariantID: variantID,
			Name:      option.Name,
			Value:     option.Value,
		}
		if err := tx.Create(&newOption).Error; err != nil {
			return fmt.Errorf("failed to insert variant option: %v", err)
		}
	}
	return nil
}

func (r *productRepo) GetProductWithFirstVariant(productId string) (*domain.Product, *domain.Variant, error) {
	var product domain.Product
	var variant domain.Variant

	productResult := r.db.Where("product_id = ?", productId).First(&product)
	if productResult.Error != nil {
		return nil, nil, productResult.Error
	}

	variantResult := r.db.Where("product_id = ?", productId).Order("id").First(&variant)
	if variantResult.Error != nil {
		return &product, nil, variantResult.Error
	}

	return &product, &variant, nil
}

func (r *productRepo) GetBrandImageByBrandName(brandName string) (string, error) {
	var brandImage string
	brandResult := r.db.Table("products").Where("LOWER(brand_name) = ?", strings.ToLower(brandName)).Pluck("brand_image", &brandImage)
	if brandResult.Error != nil {
		return "", brandResult.Error
	}

	return brandImage, nil
}

func (r *productRepo) GetBrands(searchStr string) ([]*dto.Brands, error) {
	var brands []*dto.Brands
	result := r.db.Table("brands").
		Select("DISTINCT name AS brand_name, image AS brand_image").
		Where("name ILIKE ? AND is_active = true", "%"+searchStr+"%").
		Find(&brands)

	if result.Error != nil {
		return nil, result.Error
	}
	return brands, nil
}

func (r *productRepo) GetProductByVariantId(variantID string) (*domain.Product, *domain.Variant, error) {
	variant := domain.Variant{}
	product := domain.Product{}

	result := r.db.Where("id = ?", variantID).First(&variant)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	result = r.db.Unscoped().Preload("Variants", "id = ?", variantID).
		Preload("Variants.VariantOptions", "variant_id = ?", variantID).
		Where("product_id = ?", variant.ProductID).First(&product)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var loadedVariant *domain.Variant
	if len(product.Variants) > 0 {
		loadedVariant = &product.Variants[0]
	}

	return &product, loadedVariant, nil
}

func (r *productRepo) GetProductIds() []string {

	var productIds []string
	r.db.Table("products").Pluck("product_id", &productIds)
	return productIds
}

func (r *productRepo) GetBrandName(productId string) string {
	var brandName string
	r.db.Table("products").Pluck("brand_name", &brandName).Where("product_id = ?", productId)
	return brandName
}

func (r *productRepo) GetProductCountWithSubCategoryID(subCategories uint) (int64, error) {
	var totalCount int64

	query := r.db.Model(&domain.Product{}).Where("sub_category_id = ? AND is_active = ?", subCategories, true)

	if err := query.Count(&totalCount).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

func getOrCreateBrandID(db *gorm.DB, brandName string) (uint, error) {
	var brand domain.Brand
	brandNameLower := strings.ToLower(brandName)

	result := db.Where("LOWER(name) = ?", brandNameLower).First(&brand)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 0, result.Error
	}

	if result.RowsAffected > 0 {
		return brand.ID, nil
	}

	newBrand := domain.Brand{Name: brandName, IsActive: true}
	if err := db.Create(&newBrand).Error; err != nil {
		return 0, err
	}

	return newBrand.ID, nil
}

func (r *productRepo) UpdateProductStatusByBrandID(brandID uint, isActive bool) error {
	return r.db.Model(&domain.Product{}).Where("brand_id = ?", brandID).Update("is_active", isActive).Error
}

func (r *productRepo) GetVariantsByProductId(productID string) ([]*domain.Variant, error) {
	var variants []*domain.Variant
	result := r.db.Where("product_id = ?", productID).Find(&variants)
	if result.Error != nil {
		return nil, result.Error
	}
	return variants, nil
}

func (r *productRepo) GetAllProductsWithStats(params *dto.ListProductReviewsRequest) ([]*dto.ProductWithStats, *dto.Paging, error) {
	var totalCount int64
	productsWithStats := make([]*dto.ProductWithStats, 0)

	query := r.db.Table("products p").
		Select("p.product_id as product_id, p.name as product_name, COUNT(DISTINCT r.id) AS total_reviews, AVG(r.rating) AS average_rating, MIN(v.image) AS image").
		Joins("LEFT JOIN reviews r ON p.product_id = r.product_id").
		Joins("LEFT JOIN variants v ON p.product_id = v.product_id").
		Group("p.product_id, p.name, r.product_id")

	if params.Stars != nil && len(params.Stars) > 0 {
		query = query.Where("r.rating IN (?)", params.Stars)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	query = query.Offset((params.PageNum - 1) * params.PageSize).
		Limit(params.PageSize)

	if err := query.Find(&productsWithStats).Error; err != nil {
		return nil, nil, err
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		CurrentPage: params.PageNum,
		PageSize:    params.PageSize,
	}

	return productsWithStats, paging, nil
}
