package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	dto "github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	pagination "github.com/harishash/dotshop-be/internal/utils/pagination"

	"gorm.io/gorm"
)

// CommonRepo interface
type CommonRepo interface {
	FetchAllProducts(curatorID uint, subCategories string, isFeature bool, page, pageSize int) (*dto.GetFeatureProductResponse, error)
	GetTotalProducts(curatorID uint, subCategories string, isFeature bool) (int, error)
	FetchAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error)
	FetchSectionByID(sectionID uint) (*dto.SectionResponse, error)
	GetTotalProductsCount(curatorID uint) (int64, error)
	GetTotalCollectionProductsCount(collectionID uint) (int64, error)
	FetchCuratorAllLooks(curatorid uint, pageNum, pageSize int) ([]*domain.Look, error)
	GetTotalCuratorLooksCount(curatorID uint) (int64, error)
	FetchProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error)
	FetchProductsByLookID(lookID uint, page, pageSize int) (*domain.Look, error)
	GetTotalLookProductsCount(lookID uint) (int64, error)
	SearchFeatureProductsByName(curatorID uint, searchQuery *dto.SearchProductParams) ([]domain.Product, *dto.Paging, error)
	SearchLooksByName(queryParams *dto.SearchLookParams) ([]domain.Look, *dto.Paging, error)
	SearchProductsWithinCuratorLooks(curatorID uint, searchQuery *dto.SearchProductParams) ([]domain.Product, *dto.Paging, error)
	SearchCollectionsByName(searchQuery *dto.SearchCollectionParams, curatorId uint) ([]dto.CollectionWithProducts, *dto.Paging, error)
	SearchCollectionProductByName(collectionId uint64, searchQuery *dto.SearchProductParams) (*domain.Collection, *dto.Paging, error)
	SearchCollectionSectionsByName(collectionId uint64, queryParams *dto.SearchSectionParams) ([]domain.CollectionSection, *dto.Paging, error)
	GlobalSearch(searchQuery *dto.SearchProductParams) ([]domain.Product, []string, *dto.Paging, error)
	FetchAllLooks(pageNum, pageSize int) ([]*domain.Look, error)
	GetTotalLooksCount() (int64, error)
	GetCuratorByEmail(email string) (*domain.Curator, error)
}

// commonRepo struct
type commonRepo struct {
	db *gorm.DB
}

// NewCommonRepo creates a new CommonRepo
func NewCommonRepo() CommonRepo {
	instance := GetDatabaseConnection()
	return &commonRepo{
		db: instance.Connection}
}

func (r *commonRepo) GetFeatureIDByCuratorID(curatorID uint) (uint, error) {
	var featureID uint
	err := r.db.Table("curator_products").
		Select("feature_id").
		Where("curator_id = ?", curatorID).
		Row().
		Scan(&featureID)

	if err != nil {
		return 0, err
	}

	return featureID, nil
}

func (r *commonRepo) FetchAllProducts(curatorID uint, subCategories string, isFeature bool, page, pageSize int) (*dto.GetFeatureProductResponse, error) {
	validSubCategories := splitIds(subCategories)

	var products []*domain.Product
	query := r.db.Preload("Variants.VariantOptions").Preload("SubCategory").
		Joins("INNER JOIN curator_products ON products.product_id = curator_products.product_id").
		Where("curator_products.curator_id = ? AND curator_products.is_feature = ? AND products.is_active = ?", curatorID, isFeature, true).
		Order("products.created_at DESC")

	if len(validSubCategories) > 0 {
		query = query.Where("products.sub_category_id IN ?", validSubCategories)
	}

	result := pagination.Paginate(query, page, pageSize).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	var productResponses []*dto.ProductResponse
	for _, collectionProduct := range products {
		productResponse, err := r.fetchProductDetails(collectionProduct.ProductID)
		if productResponse == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		productResponses = append(productResponses, productResponse)
	}
	response := &dto.GetFeatureProductResponse{
		Products: productResponses,
	}

	return response, nil
}

func (r *commonRepo) GetTotalProducts(curatorID uint, subCategories string, isFeature bool) (int, error) {
	validSubCategories := splitIds(subCategories)

	var count int64
	query := r.db.Model(&domain.Product{}).Preload("Variants.VariantOptions").Preload("SubCategory").
		Joins("INNER JOIN curator_products ON products.product_id = curator_products.product_id").
		Where("curator_products.curator_id = ? AND curator_products.is_feature = ? AND products.is_active = ?", curatorID, isFeature, true)

	if len(validSubCategories) > 0 {
		query = query.Where("products.sub_category_id IN ?", validSubCategories)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *commonRepo) FetchAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error) {
	var collections []*domain.Collection
	// Fetch collections for the given curator ID with pagination using Paginate utility function
	if err := r.db.Where("curator_id = ?", curatorID).Find(&collections).Error; err != nil {
		return nil, err
	}

	var response []*dto.CollectionResponse

	for _, collection := range collections {
		productResponses, err := r.fetchProductsForCollection(collection.ID)
		if err != nil {
			return nil, err
		}

		sectionResponses, err := r.fetchSectionsForCollection(collection.ID)
		if err != nil {
			return nil, err
		}

		response = append(response, &dto.CollectionResponse{
			ID:          collection.ID,
			Name:        collection.Name,
			Description: collection.Description,
			TileColor:   collection.TileColor,
			CreatedAt:   collection.CreatedAt,
			UpdatedAt:   collection.UpdatedAt,
			Products:    productResponses,
			Sections:    sectionResponses,
		})
	}
	return response, nil
}

func (r *commonRepo) FetchSectionByID(sectionID uint) (*dto.SectionResponse, error) {
	// Fetch the section based on section ID
	var section domain.CollectionSection
	if err := r.db.Where("id = ?", sectionID).First(&section).Error; err != nil {
		return nil, err
	}

	// Fetch products associated with the section
	sectionProducts, err := r.fetchSectionProducts(section.ID)
	if err != nil {
		return nil, err
	}

	// Build the response
	sectionResponse := &dto.SectionResponse{
		ID:           section.ID,
		Name:         section.Name,
		ImageURL:     section.ImageURL,
		Description:  section.Description,
		CollectionID: section.CollectionID,
		Products:     sectionProducts,
	}

	return sectionResponse, nil
}
func (r *commonRepo) fetchProductsForCollection(collectionID uint) ([]*dto.ProductResponse, error) {
	var collectionProducts []*domain.CollectionProduct
	if err := r.db.Joins("JOIN products ON products.product_id = collection_products.product_id").Where("collection_id = ? AND is_active = ?", collectionID, true).Find(&collectionProducts).Error; err != nil {
		return nil, err
	}

	var productResponses []*dto.ProductResponse
	for _, collectionProduct := range collectionProducts {
		productResponse, err := r.fetchProductDetails(collectionProduct.ProductID)
		if productResponse == nil {
			continue
		}

		if err != nil {
			return nil, err
		}
		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}

func (r *commonRepo) fetchSectionsForCollection(collectionID uint) ([]*dto.SectionResponse, error) {
	var collectionSections []*domain.CollectionSection
	if err := r.db.Where("collection_id = ?", collectionID).Find(&collectionSections).Error; err != nil {
		return nil, err
	}

	var sectionResponses []*dto.SectionResponse
	for _, section := range collectionSections {
		sectionProducts, err := r.fetchSectionProducts(section.ID)
		if sectionProducts == nil {
			continue
		}
		if err != nil {
			return nil, err
		}

		sectionResponses = append(sectionResponses, &dto.SectionResponse{
			ID:           section.ID,
			Name:         section.Name,
			ImageURL:     section.ImageURL,
			Description:  section.Description,
			CollectionID: section.CollectionID,
			Products:     sectionProducts,
		})
	}

	return sectionResponses, nil
}

func (r *commonRepo) fetchSectionProducts(sectionID uint) ([]*dto.ProductResponse, error) {
	var sectionProducts []*domain.CollectionSectionProduct
	if err := r.db.Joins("JOIN products ON products.product_id = collection_section_products.product_id").Where("collection_section_id = ? AND is_active = ?", sectionID, true).Find(&sectionProducts).Error; err != nil {
		return nil, err
	}

	var productResponses []*dto.ProductResponse
	for _, sectionProduct := range sectionProducts {
		productResponse, err := r.fetchProductDetails(sectionProduct.ProductID)
		if productResponse == nil {
			continue
		}

		if err != nil {
			return nil, err
		}
		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}
func (r *commonRepo) fetchProductDetails(productID string) (*dto.ProductResponse, error) {
	var product *domain.Product
	if err := r.db.Where("product_id = ? AND is_active = ?", productID, true).Preload("Variants.VariantOptions").First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	var variantResponses []*dto.VariantResponse
	for _, variant := range product.Variants {
		var variantOptionResponses []*dto.VariantOptionResponse
		for _, option := range variant.VariantOptions {
			variantOptionResponses = append(variantOptionResponses, &dto.VariantOptionResponse{
				VariantID: option.VariantID,
				Name:      option.Name,
				Value:     option.Value,
			})
		}

		variantResponses = append(variantResponses, &dto.VariantResponse{
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
			VariantOptions:  variantOptionResponses,
		})
	}

	return &dto.ProductResponse{
		ID:           product.ProductID,
		BrandName:    product.BrandName,
		SupplierName: product.SupplierName,
		Name:         product.Name,
		Description:  product.Description,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		Variants:     variantResponses,
	}, nil
}

func (r *commonRepo) GetTotalProductsCount(curatorID uint) (int64, error) {
	var totalProducts int64
	// Fetch total number of products for the given curator ID
	if err := r.db.Model(&domain.Collection{}).Where("curator_id = ?", curatorID).Count(&totalProducts).Error; err != nil {
		return 0, err
	}

	return totalProducts, nil
}

func (r *commonRepo) FetchCuratorAllLooks(curatorid uint, pageNum, pageSize int) ([]*domain.Look, error) {
	var looks []*domain.Look

	// Apply the pagination to the query
	query := r.db.Preload("Products", "is_active = ?", true).Preload("Products.Variants.VariantOptions").Order("created_at DESC").Where("curator_id = ?", curatorid)
	paginatedQuery := pagination.Paginate(query, pageNum, pageSize).Find(&looks)

	if paginatedQuery.Error != nil {
		return nil, fmt.Errorf("failed to fetch looks: %w", paginatedQuery.Error)
	}

	// Return looks with pagination info
	return looks, nil
}

func (r *commonRepo) FetchAllLooks(pageNum, pageSize int) ([]*domain.Look, error) {
	var looks []*domain.Look

	// Start the query
	query := r.db.Preload("Products", "is_active = ?", true).Preload("Products.Variants.VariantOptions").Order("created_at DESC")


		

	// Apply pagination
	paginatedQuery := pagination.Paginate(query, pageNum, pageSize).Find(&looks)

	if paginatedQuery.Error != nil {
		return nil, fmt.Errorf("failed to fetch looks: %w", paginatedQuery.Error)
	}

	// Return looks with pagination info
	return looks, nil
}

func (r *commonRepo) GetTotalCuratorLooksCount(curatorID uint) (int64, error) {
	var count int64

	// Fetch total number of looks for the given curator ID
	if err := r.db.Model(&domain.Look{}).Where("curator_id = ?", curatorID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commonRepo) GetTotalLooksCount() (int64, error) {
	var count int64

	query := r.db.Preload("Products", "is_active = ?", true).Preload("Products.Variants.VariantOptions")

	// Fetch total number of looks for the given curator ID
	if err := query.Model(&domain.Look{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Fetch products by collection ID along with sections
func (r *commonRepo) FetchProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error) {

	var collection domain.Collection

	// Fetch the collection based on collection ID
	if err := r.db.Where("id = ?", collectionID).First(&collection).Error; err != nil {
		return &dto.CollectionResponse{}, err
	}

	// Fetch products associated with the collection with pagination
	var collectionProducts []domain.CollectionProduct
	query := r.db.Where("collection_id = ?", collectionID)
	paginatedQuery := pagination.Paginate(query, page, pageSize).Find(&collectionProducts)
	if paginatedQuery.Error != nil {
		return &dto.CollectionResponse{}, paginatedQuery.Error
	}

	var productResponses []*dto.ProductResponse
	for _, collectionProduct := range collectionProducts {
		productResponse, err := r.fetchProductDetails(collectionProduct.ProductID)
		if productResponse == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		productResponses = append(productResponses, productResponse)
	}

	sectionResponses, err := r.fetchSectionsForCollection(collectionID)
	if err != nil {
		return nil, err
	}

	response := &dto.CollectionResponse{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		TileColor:   collection.TileColor,
		Products:    productResponses,
		Sections:    sectionResponses,
	}

	return response, nil
}

// func (r *commonRepo) FetchProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error) {
// 	var collection domain.Collection

// 	// Fetch the collection based on collection ID
// 	if err := r.db.Where("id = ?", collectionID).First(&collection).Error; err != nil {
// 		return &dto.CollectionResponse{}, err
// 	}

// 	// Fetch products associated with the collection with pagination
// 	var collectionProducts []domain.CollectionProduct
// 	query := r.db.Where("collection_id = ?", collectionID)
// 	paginatedQuery := pagination.Paginate(query, page, pageSize).Find(&collectionProducts)
// 	if paginatedQuery.Error != nil {
// 		return &dto.CollectionResponse{}, paginatedQuery.Error
// 	}

// 	var productResponses []*dto.ProductResponse
// 	for _, collectionProduct := range collectionProducts {
// 		productResponse, err := r.fetchProductDetails(collectionProduct.ProductID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		productResponses = append(productResponses, productResponse)
// 	}

// 	response := &dto.CollectionResponse{
// 		ID:          collection.ID,
// 		Name:        collection.Name,
// 		Description: collection.Description,
// 		TileColor:   collection.TileColor,
// 		Products:    productResponses,
// 	}

// 	return response, nil
// }

func (r *commonRepo) GetTotalCollectionProductsCount(collectionID uint) (int64, error) {
	var count int64

	// Fetch total number of products for the given collection ID
	if err := r.db.Model(&domain.CollectionProduct{}).Joins("JOIN products ON products.product_id = collection_products.product_id").Where("collection_id = ? AND products.is_active = ?", collectionID, true).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commonRepo) FetchProductsByLookID(lookID uint, page, pageSize int) (*domain.Look, error) {
	// Pagination calculations
	offset := (page - 1) * pageSize

	// Fetch the look with the specified look ID with pagination
	var look domain.Look
	result := r.db.Preload("Products", "is_active = ?", true).Preload("Products.Variants.VariantOptions").Where("id = ?", lookID).Offset(offset).Limit(pageSize).First(&look)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &domain.Look{}, fmt.Errorf("look not found: %w", result.Error)
		}
		return &domain.Look{}, fmt.Errorf("failed to fetch look: %w", result.Error)
	}

	// Return the fetched look
	return &look, nil
}

func (r *commonRepo) GetTotalLookProductsCount(lookID uint) (int64, error) {
	var count int64
	// Fetch total number of products for the given look ID
	if err := r.db.Model(&domain.LookProduct{}).Joins("JOIN products ON products.product_id = look_products.product_product_id").Where("look_id = ? AND products.is_active = ?", lookID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *commonRepo) SearchFeatureProductsByName(curatorID uint, searchQuery *dto.SearchProductParams) ([]domain.Product, *dto.Paging, error) {
	var products []domain.Product
	var totalCount int64

	query := r.db.
		Table("products").
		Select("products.*, curator_products.is_feature").
		Joins("INNER JOIN curator_products ON products.product_id = curator_products.product_id").
		Preload("Variants.VariantOptions"). // Preload the variant options for each variant
		Where("products.is_active = ?", true).
		Where("curator_products.curator_id = ? AND LOWER(products.name) LIKE LOWER(?) AND curator_products.is_feature = ?", curatorID, "%"+searchQuery.Product+"%", true)

	if err := query.Model(&domain.Product{}).Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	result := pagination.NewPaginate(query, searchQuery.PageNum, searchQuery.PageSize).Find(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    searchQuery.PageSize,
		CurrentPage: searchQuery.PageNum,
	}
	return products, &paging, nil
}

func (r *commonRepo) SearchLooksByName(queryParams *dto.SearchLookParams) ([]domain.Look, *dto.Paging, error) {
	var looks []domain.Look
	var totalCount int64

	query := r.db.Preload("Products", "is_active = ?", true).Preload("Products.Variants.VariantOptions").Where("LOWER(name) LIKE LOWER(?)", "%"+queryParams.Look+"%")

	// query := r.db.Preload("Products.Variants.VariantOptions").
	// 	Joins("JOIN look_products ON look_products.look_id = looks.id").
	// 	Joins("JOIN products ON products.product_id = look_products.product_product_id").
	// 	Where("products.is_active = ?", true).
	// 	Where("LOWER(looks.name) LIKE LOWER(?)", "%"+queryParams.Look+"%")

	// query := r.db.
	// 	Preload("Products.Variants.VariantOptions").
	// 	Joins("JOIN look_products ON look_products.look_id = looks.id").
	// 	Joins("JOIN products ON products.product_id = look_products.product_product_id").
	// 	Where("products.is_active = ?", true).
	// 	Order("looks.created_at DESC")

	if err := query.Model(&domain.Look{}).Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	result := pagination.NewPaginate(query, queryParams.PageNum, queryParams.PageSize).Find(&looks)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    queryParams.PageSize,
		CurrentPage: queryParams.PageNum,
	}

	return looks, &paging, nil
}

func (r *commonRepo) SearchProductsWithinCuratorLooks(curatorID uint, searchQuery *dto.SearchProductParams) ([]domain.Product, *dto.Paging, error) {
	var products []domain.Product
	var totalCount int64

	query := r.db.
		Table("looks").
		Joins("JOIN look_products ON look_products.look_id = looks.id").
		Joins("JOIN products ON products.product_id = look_products.product_product_id").
		Where("products.is_active = ?", true).
		Where("looks.curator_id = ? AND LOWER(products.name) LIKE LOWER(?)", curatorID, "%"+searchQuery.Product+"%").
		Preload("Variants.VariantOptions")

	if err := query.Model(&domain.Product{}).Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	query.Select("products.*")
	result := pagination.NewPaginate(query, searchQuery.PageNum, searchQuery.PageSize).Find(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    searchQuery.PageSize,
		CurrentPage: searchQuery.PageNum,
	}
	return products, &paging, nil
}

func (r *commonRepo) SearchCollectionsByName(searchQuery *dto.SearchCollectionParams, curatorId uint) ([]dto.CollectionWithProducts, *dto.Paging, error) {
	var results []struct {
		domain.Collection
		domain.Product
	}
	var totalCount int64

	query := r.db.
		Table("collections").
		Select("collections.*, products.*, collection_products.collection_id, collection_products.product_id").
		Joins("LEFT JOIN collection_products ON collection_products.collection_id = collections.id").
		Joins("LEFT JOIN products ON products.product_id = collection_products.product_id").
		Where("collections.curator_id = ? AND LOWER(collections.name) LIKE LOWER(?)", curatorId, "%"+searchQuery.Collection+"%")

	if err := r.db.Table("collections").
		Where("collections.curator_id = ? AND LOWER(collections.name) LIKE LOWER(?)", curatorId, "%"+searchQuery.Collection+"%").
		Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	result := pagination.NewPaginate(query, searchQuery.PageNum, searchQuery.PageSize).Find(&results)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	// Organize the results to avoid duplication
	collectionsMap := make(map[uint]dto.CollectionWithProducts)
	for _, result := range results {
		collectionID := result.Collection.ID
		product := result.Product

		if collection, exists := collectionsMap[collectionID]; exists {
			if len(product.ProductID) > 0 {
				if product.IsActive {
					collection.Products = append(collection.Products, product)
				}
			}
			collectionsMap[collectionID] = collection
		} else {
			products := []domain.Product{}
			if len(product.ProductID) > 0 {
				if product.IsActive {
					products = append(products, product)
				}
			}
			collectionsMap[collectionID] = dto.CollectionWithProducts{
				Collection: result.Collection,
				Products:   products,
			}
		}
	}

	collections := make([]dto.CollectionWithProducts, 0, len(collectionsMap))
	for _, collection := range collectionsMap {
		collections = append(collections, collection)
	}

	for _, collection := range collections {
		for i := range collection.Products {
			if err := r.db.Preload("Variants.VariantOptions").Find(&collection.Products[i]).Error; err != nil {
				return nil, nil, err
			}
		}
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    searchQuery.PageSize,
		CurrentPage: searchQuery.PageNum,
	}
	return collections, &paging, nil
}

func (r *commonRepo) SearchCollectionProductByName(collectionId uint64, searchQuery *dto.SearchProductParams) (*domain.Collection, *dto.Paging, error) {
	var collection domain.Collection
	var totalCount int64

	countQuery := r.db.
		Table("products").
		Joins("JOIN collection_products ON collection_products.product_id = products.product_id").
		Where("products.is_active = ?", true).
		Where("collection_products.collection_id = ? AND LOWER(products.name) LIKE LOWER(?)", collectionId, "%"+searchQuery.Product+"%")

	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	collectionQuery := r.db.
		Where("id = ?", collectionId).
		First(&collection)

	if collectionQuery.Error != nil {
		return nil, nil, collectionQuery.Error
	}

	productQuery := countQuery.
		Preload("Variants.VariantOptions")

	var products []domain.Product
	result := pagination.NewPaginate(productQuery, searchQuery.PageNum, searchQuery.PageSize).Find(&products)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	if productQuery.Error != nil {
		return nil, nil, productQuery.Error
	}

	collection.Products = products

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    searchQuery.PageSize,
		CurrentPage: searchQuery.PageNum,
	}

	return &collection, &paging, nil
}

func (r *commonRepo) SearchCollectionSectionsByName(collectionId uint64, queryParams *dto.SearchSectionParams) ([]domain.CollectionSection, *dto.Paging, error) {
	var sections []domain.CollectionSection
	var totalCount int64

	query := r.db.
		Table("collection_sections AS cs").
		Joins("LEFT JOIN collection_section_products AS csp ON cs.id = csp.collection_section_id").
		Joins("LEFT JOIN products AS p ON csp.product_id = p.product_id").
		Where("p.is_active = ?", true).
		Where("cs.collection_id = ? AND cs.name ILIKE ?", collectionId, "%"+queryParams.Section+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	query.Select("cs.*")
	result := pagination.NewPaginate(query, queryParams.PageNum, queryParams.PageSize).Find(&sections)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	paging := dto.Paging{
		TotalCount:  totalCount,
		PageSize:    queryParams.PageSize,
		CurrentPage: queryParams.PageNum,
	}

	return sections, &paging, nil
}

func (r *commonRepo) GlobalSearch(searchQuery *dto.SearchProductParams) ([]domain.Product, []string, *dto.Paging, error) {
	var totalCount int64
	var productSuggestions, categorySuggestions, brandSuggestions, subCategorySuggestions []string
	paging := &dto.Paging{PageSize: searchQuery.PageSize, CurrentPage: searchQuery.PageNum}
	productNameQuery := searchQuery.Product

	productSuggestions, err := r.getProductSuggestions(productNameQuery)
	if err != nil {
		return nil, nil, nil, err
	}

	brandSuggestions, err = r.getBrandSuggestions(searchQuery.Product)
	if err != nil {
		return nil, nil, nil, err
	}

	categorySuggestions, err = r.getCategorySuggestions(searchQuery.Product)
	if err != nil {
		return nil, nil, nil, err
	}

	subCategorySuggestions, err = r.getSubCategorySuggestions(searchQuery.Product)
	if err != nil {
		return nil, nil, nil, err
	}

	suggestions := append(productSuggestions, brandSuggestions...)
	suggestions = append(suggestions, categorySuggestions...)
	suggestions = append(suggestions, subCategorySuggestions...)

	suggestions = removeDuplicates(suggestions)

	products, totalCount, err := r.getMatchingProducts(productNameQuery, searchQuery, categorySuggestions, brandSuggestions, subCategorySuggestions)
	if err != nil {
		return nil, nil, nil, err
	}

	paging.TotalCount = totalCount
	return products, suggestions, paging, nil
}

func (r *commonRepo) getProductSuggestions(query string) ([]string, error) {
	var suggestions []string

	rows, err := r.db.
		Table("products").
		Select("DISTINCT products.name AS product_name").
		Where("LOWER(products.name) LIKE LOWER(?)", "%"+query+"%").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productName sql.NullString
		if err := rows.Scan(&productName); err != nil {
			return nil, err
		}
		if productName.Valid && !contains(suggestions, productName.String) {
			suggestions = append(suggestions, productName.String)
		}
	}
	return suggestions, nil
}

func (r *commonRepo) getCategorySuggestions(query string) ([]string, error) {
	var suggestions []string

	rows, err := r.db.
		Table("categories").
		Select("categories.name AS category_name, sub_categories.name AS sub_category_name").
		Joins("LEFT JOIN sub_categories ON sub_categories.category_id = categories.id").
		Where("LOWER(categories.name) LIKE LOWER(?) OR LOWER(sub_categories.name) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		var subCategoryName sql.NullString
		if err := rows.Scan(&categoryName, &subCategoryName); err != nil {
			return nil, err
		}
		if categoryName != "" && !contains(suggestions, categoryName) {
			suggestions = append(suggestions, categoryName)
		}
		if subCategoryName.Valid {
			combined := categoryName + " " + subCategoryName.String
			if !contains(suggestions, combined) {
				suggestions = append(suggestions, combined)
			}
		}
	}
	return suggestions, nil
}

func (r *commonRepo) getBrandSuggestions(query string) ([]string, error) {
	var suggestions []string

	rows, err := r.db.
		Table("products").
		Select("DISTINCT products.brand_name AS brand_name").
		Where("LOWER(products.brand_name) LIKE LOWER(?)", "%"+query+"%").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var brandName sql.NullString
		if err := rows.Scan(&brandName); err != nil {
			return nil, err
		}
		if brandName.Valid && !contains(suggestions, brandName.String) {
			suggestions = append(suggestions, brandName.String)
		}
	}
	return suggestions, nil
}

func (r *commonRepo) getSubCategorySuggestions(query string) ([]string, error) {
	var suggestions []string

	rows, err := r.db.
		Table("sub_categories").
		Select("DISTINCT sub_categories.name AS sub_category_name").
		Where("LOWER(sub_categories.name) LIKE LOWER(?)", "%"+query+"%").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subCategoryName sql.NullString
		if err := rows.Scan(&subCategoryName); err != nil {
			return nil, err
		}
		if subCategoryName.Valid && !contains(suggestions, subCategoryName.String) {
			suggestions = append(suggestions, subCategoryName.String)
		}
	}
	return suggestions, nil
}

func (r *commonRepo) getMatchingProducts(query string, searchQuery *dto.SearchProductParams, categorySuggestions, brandSuggestions, subCategorySuggestions []string) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalCount int64

	countQuery := r.db.
		Table("products").
		Joins("JOIN sub_categories ON sub_categories.id = products.sub_category_id").
		Joins("JOIN categories ON categories.id = sub_categories.category_id")

	for _, suggestion := range categorySuggestions {
		if strings.Contains(suggestion, " ") {
			parts := strings.SplitN(suggestion, " ", 2)
			countQuery = countQuery.Or("LOWER(categories.name) LIKE LOWER(?) AND LOWER(sub_categories.name) LIKE LOWER(?)", "%"+parts[0]+"%", "%"+parts[1]+"%")
		} else {
			countQuery = countQuery.Or("LOWER(categories.name) LIKE LOWER(?)", "%"+suggestion+"%")
		}
	}

	for _, suggestion := range brandSuggestions {
		countQuery = countQuery.Or("LOWER(products.brand_name) LIKE LOWER(?)", "%"+suggestion+"%")
	}

	for _, suggestion := range subCategorySuggestions {
		countQuery = countQuery.Or("LOWER(sub_categories.name) LIKE LOWER(?)", "%"+suggestion+"%")
	}

	countQuery = countQuery.Or("LOWER(products.name) LIKE LOWER(?)", "%"+query+"%")

	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	productQuery := countQuery.
		Select("products.*").
		Preload("Variants.VariantOptions")

	result := pagination.NewPaginate(productQuery, searchQuery.PageNum, searchQuery.PageSize).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, totalCount, nil
}

func removeDuplicates(slice []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if _, ok := unique[item]; !ok {
			unique[item] = true
			result = append(result, item)
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (r *commonRepo) GetCuratorByEmail(email string) (*domain.Curator, error) {
	var curator *domain.Curator
	err := r.db.Joins("JOIN users ON users.id = curators.user_id").
		Where("LOWER(users.email) = LOWER(?)", email).
		Preload("User").
		First(&curator).Error
	if err != nil {
		return nil, err
	}
	return curator, nil
}
