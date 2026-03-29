package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	dto "github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
)

type CuratorBORepo interface {
	InsertProduct(products *dto.CreateFeatureProductRequest) ([]*dto.CreateFeatureProductResponse, error)
	InsertCollection(collection *dto.CreateCollectionRequest) (*dto.CreateCollectionResponse, error)
	InsertCollectionSection(collectionSection *dto.CreateCollectionSectionRequest) (*dto.CreateCollectionSectionResponse, error)
	InsertProductToSection(sectionID uint, products *dto.AddProductToSectionRequest) (*dto.AddProductToSectionResponse, error)
	DeleteProductFromSection(sectionID uint, productID string) error
	UpdateCollectionSection(body *dto.UpdateCollectionSectionRequest, sectionID uint) (*dto.UpdateCollectionSectionResponse, error)
	InsertLook(look *dto.CreateLookRequest) (*dto.CreateLookResponse, error)

	DeleteFromFeatureProduct(curatorID uint, productID string) error
	DeleteCollectionByID(collectionID uint) error
	DeleteProductFromCollection(collectionID uint, productID string) error
	DeleteLookByID(lookID uint) error
	DeleteSectionByID(sectionID uint) error
	DeleteProductFromLook(look uint, productID string) error
	SearchLooksByName(query string) ([]domain.Look, error)
	SearchProductsWithinCuratorLooks(curatorID uint, searchQuery string) ([]domain.Product, error)
	FetchProductsByLookID(curatorID, lookID uint) (domain.Look, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetCuratorByID(id uint64) (*domain.Curator, error)
	GetAllCurators() ([]*dto.GetCuratorResponse, error)
	GetCuratorByCuratorID(id uint64) (*domain.Curator, error)

	ProcessWithdraw(data []byte) error
	UpdateProfile(curator *domain.Curator, user *domain.User) error
	CreateSocialMediaLink(link *domain.SocialMediaLink) (*dto.CreateSocialMediaLinkResponse, error)
	UpdateSocialMediaLink(link *domain.SocialMediaLink, linkID, curatorID uint) (*dto.CreateSocialMediaLinkResponse, error)
	GetSocialMediaLinkByTypeAndCuratorID(linkType string, curatorID uint) (*domain.SocialMediaLink, error)
	DeleteSocialMediaLink(curatorID, linkID uint64) (*dto.DeleteSocialMediaLinkResponse, error)
	UpdatePasswordByEmail(email string, newPassword string) (error, *dto.UpdatePasswordResponse)
	StartTx() (*gorm.DB, error)
	GetLookWithId(id uint, curatorId uint) (*domain.Look, error)
	GetFeatureWithId(id uint, curatorId uint) (*domain.CuratorProduct, error)

	AddProduct(tx *gorm.DB, product *domain.Product) (*gorm.DB, *domain.Product, error)
	GetVariantWithId(id string) (*domain.Variant, error)
	AddProductVariant(tx *gorm.DB, variant *domain.Variant) (*gorm.DB, *domain.Variant, error)
	CheckProductExistsInLook(lookId uint, productId string) (bool, error)
	CheckProductExistsInFeature(featureId uint, productId string) (bool, error)
	AddProductToFeature(tx *gorm.DB, featureProduct *domain.CuratorProduct) (*gorm.DB, error)
	AddProductToLook(tx *gorm.DB, lookProduct *domain.LookProduct) (*gorm.DB, error)
	GetCollectionWithId(id uint, curatorId uint) (*domain.Collection, error)
	GetProductWithId(id string) (*domain.Product, error)
	AddProductToCollection(tx *gorm.DB, collectionProduct *domain.CollectionProduct) (*gorm.DB, error)
	CheckProductExistsInCollection(collectionId uint, productId string) (bool, error)
	GetCuratorWithUserId(id uint64) (*domain.Curator, error)
	UpdateCollection(collection *domain.Collection) error
	UpdateLook(look *domain.Look) error
	GetUserByID(userId uint64) (*domain.User, error)
	GetCuratorAccountDetail(curatorId uint) (*struct {
		domain.BankInformation
		domain.AccountDetails
	}, error)
	UpdateCuratorAccountDetail(curatorId uint, request *dto.UpdateAccountDetailRequest) (*struct {
		domain.BankInformation
		domain.AccountDetails
	}, error)
}

type curatorBORepo struct {
	db *gorm.DB
}

func NewCuratorBORepo() CuratorBORepo {
	instance := GetDatabaseConnection()

	return &curatorBORepo{db: instance.Connection}
}

func (r *curatorBORepo) InsertProduct(products *dto.CreateFeatureProductRequest) ([]*dto.CreateFeatureProductResponse, error) {
	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer tx.Rollback()

	// Check if the curator exists
	var curator *domain.Curator
	if err := tx.Where("id = ?", products.CuratorID).First(&curator).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("curator with ID %d does not exist", products.CuratorID)
		}
		return nil, fmt.Errorf("failed to check curator existence: %v", err)
	}
	var insertedProducts []*dto.CreateFeatureProductResponse
	// Iterate over the list of products
	for _, product := range products.Products {
		// Check if the product is already added by the curator
		var existingCuratorProduct *domain.CuratorProduct
		if err := tx.Where("curator_id = ? AND product_id = ?", products.CuratorID, product.ProductID).First(&existingCuratorProduct).Error; err == nil {
			return nil, fmt.Errorf("product with ID %s is already added by curator with ID %d", product.ProductID, products.CuratorID)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, fmt.Errorf("failed to check product existence in curator: %v", err)
		}

		var existingProduct *domain.Product
		// Check if the product already exists
		err := tx.Where("product_id = ?", product.ProductID).First(&existingProduct).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, fmt.Errorf("failed to check product existence: %v", err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not exist: %v", err)
		}

		// Record in curator_products table
		curatorProduct := &domain.CuratorProduct{
			IsFeature: true,
			CuratorID: uint(products.CuratorID),
			ProductID: product.ProductID,
		}
		if err := tx.Create(&curatorProduct).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert curator product: %v", err)
		}
		insertedProducts = append(insertedProducts, &dto.CreateFeatureProductResponse{
			ProductID:    product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			CreatedAt:    time.Now(),
			IsFeature:    true,
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return insertedProducts, nil
}

func (r *curatorBORepo) InsertCollection(collection *dto.CreateCollectionRequest) (*dto.CreateCollectionResponse, error) {
	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the curator exists
	var curator domain.Curator
	if err := tx.Where("id = ?", collection.CuratorID).First(&curator).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("curator with ID %d does not exist", collection.CuratorID)
		}
		return nil, fmt.Errorf("failed to check curator existence: %v", err)
	}

	// Check if a collection with the same name already exists
	var existingCollection domain.Collection
	if err := tx.Where("name = ? AND curator_id = ?", collection.Name, collection.CuratorID).First(&existingCollection).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("collection with name '%s' already exists", collection.Name)
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fmt.Errorf("failed to check collection existence: %v", err)
	}

	// Create a new collection entry
	newCollection := domain.Collection{
		CuratorID:   collection.CuratorID,
		Name:        collection.Name,
		Description: collection.Description,
		TileColor:   collection.TileColor,
	}

	if err := tx.Create(&newCollection).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}
	var productsResponse []dto.CreateCollectionProductResponse
	for _, product := range collection.Products {
		// Check if the product already exists
		var existingProduct domain.Product
		err := tx.Where("product_id = ?", product.ProductID).First(&existingProduct).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, err
		}

		var productResponse dto.CreateCollectionProductResponse
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not exist: %v", err)
		} else {
			productResponse = dto.CreateCollectionProductResponse{
				ID:           existingProduct.ProductID,
				BrandName:    existingProduct.BrandName,
				SupplierName: existingProduct.SupplierName,
				Name:         existingProduct.Name,
				Description:  existingProduct.Description,
				CreatedAt:    existingProduct.CreatedAt,
			}
		}

		productsResponse = append(productsResponse, productResponse)

		// Update the collection_product table with the relationship between the collection and its products
		if err := tx.Create(&domain.CollectionProduct{
			CollectionID: newCollection.ID,
			ProductID:    product.ProductID,
		}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create collection product relationship: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Construct the response
	response := &dto.CreateCollectionResponse{
		ID:          newCollection.ID,
		Name:        newCollection.Name,
		Description: newCollection.Description,
		TileColor:   newCollection.TileColor,
		Products:    productsResponse,
	}

	return response, nil
}

func (r *curatorBORepo) InsertCollectionSection(collectionSection *dto.CreateCollectionSectionRequest) (*dto.CreateCollectionSectionResponse, error) {
	// Use a transaction for atomic operations
	tx := r.db.Begin()

	// Defer rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the collection exists
	var collection *domain.Collection
	if err := tx.Where("id = ?", collectionSection.CollectionID).First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("collection with ID %d does not exist", collectionSection.CollectionID)
		}
		return nil, fmt.Errorf("failed to check collection existence: %v", err)
	}

	// Check if a section with the same name already exists
	var existingSection *domain.CollectionSection
	if err := tx.Where("name = ? AND collection_id = ?", collectionSection.Name, collectionSection.CollectionID).First(&existingSection).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("collectionSection with name '%s' already exists", *collectionSection.Name)
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fmt.Errorf("failed to check collectionSection existence: %v", err)
	}

	productIDs := make([]string, len(collectionSection.Products))
	for i, product := range collectionSection.Products {
		productIDs[i] = product.ProductID
	}

	// Use the new function to check existing products and get their names
	existingProductNames, err := r.getExistingCollectionProductNames(tx, productIDs, collectionSection.CollectionID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(existingProductNames) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("one or more products already exist in the collection: %s", strings.Join(existingProductNames, ", "))
	}

	// Create the section entity and set the collection ID
	newSection := &domain.CollectionSection{
		Name:         collectionSection.Name,
		ImageURL:     collectionSection.ImageURL,
		Description:  collectionSection.Description,
		CollectionID: collectionSection.CollectionID,
	}
	if err := tx.Create(&newSection).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert into collection_section_products
	var collectionSectionProducts []domain.CollectionSectionProduct
	for _, product := range collectionSection.Products {
		collectionSectionProducts = append(collectionSectionProducts, domain.CollectionSectionProduct{
			CollectionSectionID: newSection.ID,
			ProductID:           product.ProductID,
		})
	}

	if len(collectionSectionProducts) > 0 {
		if err := tx.Create(&collectionSectionProducts).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Create the SectionResponse object
	response := &dto.CreateCollectionSectionResponse{
		CollectionSectionID: int(newSection.ID),
		Name:                newSection.Name,
		ImageURL:            newSection.ImageURL,
		Description:         newSection.Description,
		CollectionID:        int(newSection.CollectionID),
	}

	// Add the products to the response
	for _, product := range collectionSection.Products {
		response.Product = append(response.Product, &dto.CreateCollectionProductResponse{
			ID:           product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			CreatedAt:    time.Now(),
		})
	}

	return response, nil
}

func (r *curatorBORepo) getExistingCollectionProductNames(tx *gorm.DB, productIDs []string, collectionID uint) ([]string, error) {
	var existingProductsDetails []struct {
		ProductID   string
		ProductName string
	}

	if err := tx.Table("collection_products").
		Select("products.product_id as product_id, products.name as product_name").
		Joins("JOIN products ON products.product_id = collection_products.product_id").
		Where("collection_products.product_id IN (?) AND collection_products.collection_id = ?", productIDs, collectionID).
		Scan(&existingProductsDetails).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing product details: %v", err)
	}

	var existingProductNames []string
	for _, item := range existingProductsDetails {
		existingProductNames = append(existingProductNames, item.ProductName)
	}

	return existingProductNames, nil
}

func (r *curatorBORepo) UpdateCollectionSection(
	body *dto.UpdateCollectionSectionRequest,
	sectionID uint) (
	*dto.UpdateCollectionSectionResponse,
	error) {
	section := &domain.CollectionSection{}
	if err := r.db.Where("id = ?", sectionID).First(section).Error; err != nil {
		return nil, err
	}

	section.Name = body.Name
	section.ImageURL = body.ImageURL
	section.Description = body.Description

	if err := r.db.Save(section).Error; err != nil {
		return nil, err
	}

	return &dto.UpdateCollectionSectionResponse{
		Name:        section.Name,
		ImageURL:    section.ImageURL,
		Description: section.Description,
	}, nil
}

func (r *curatorBORepo) InsertProductToSection(
	sectionID uint,
	products *dto.AddProductToSectionRequest) (
	*dto.AddProductToSectionResponse,
	error) {

	// Start a transaction
	tx := r.db.Begin()

	// Defer rollback in case of panic or error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	section := &domain.CollectionSection{}
	if err := r.db.Where("id = ?", sectionID).First(section).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Extract product IDs
	productIDs := make([]string, len(products.Products))
	for i, product := range products.Products {
		productIDs[i] = product.ProductID
	}

	// Check if any of the products already exist in the section
	existingProductNames, err := r.getExistingCollectionProductNames(tx, productIDs, section.CollectionID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to check existing products: %v", err)
	}

	if len(existingProductNames) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("one or more products already exist in the section: %s", strings.Join(existingProductNames, ", "))
	}

	// Insert new products into the section
	var insertedProductResponses []*dto.CreateProductResponse
	for _, product := range products.Products {
		if err := tx.Create(&domain.CollectionSectionProduct{
			CollectionSectionID: sectionID,
			ProductID:           product.ProductID,
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		insertedProductResponses = append(insertedProductResponses, &dto.CreateProductResponse{
			ProductID:    product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			CreatedAt:    time.Now(),
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.AddProductToSectionResponse{
		SectionId: sectionID,
		Products:  insertedProductResponses,
	}, nil
}

func (r *curatorBORepo) DeleteProductFromSection(
	sectionID uint,
	productID string) error {
	return r.db.
		Where("collection_section_id = ? AND product_id = ?", sectionID, productID).
		Delete(&domain.CollectionSectionProduct{}).Error
}

func (r *curatorBORepo) InsertLook(look *dto.CreateLookRequest) (*dto.CreateLookResponse, error) {
	// Use a transaction for atomic operations
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the curator exists
	curator := &domain.Curator{}
	if err := tx.Where("id = ?", look.CuratorID).First(curator).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("curator with ID %d does not exist", look.CuratorID)
		}
		return nil, fmt.Errorf("failed to check curator existence: %v", err)
	}

	// Create the look entity and set the curator ID
	newLook := domain.Look{
		Name:             look.Name,
		ImageURL:         look.ImageURL,
		CuratorID:        look.CuratorID,
		SocialID:         look.SocialID,
		SocialType:       look.SocialType,
		SocialTitle:      look.SocialTitle,
		EmbedLink:        look.EmbedLink,
		VideoDescription: look.VideoDescription,
	}
	if err := tx.Create(&newLook).Error; err != nil {
		return nil, err
	}

	// Insert look products using raw SQL
	sql := `INSERT INTO public.look_products (look_id, product_product_id) VALUES (?, ?)`
	for _, product := range look.Products {
		if err := tx.Exec(sql, newLook.ID, product.ProductID).Error; err != nil {
			return nil, fmt.Errorf("failed to create look product: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Create the LookResponse object
	response := &dto.CreateLookResponse{
		LookID:           int(newLook.ID),
		Name:             newLook.Name,
		ImageURL:         newLook.ImageURL,
		SocialID:         newLook.SocialID,
		SocialType:       newLook.SocialType,
		SocialTitle:      newLook.SocialTitle,
		EmbedLink:        newLook.EmbedLink,
		VideoDescription: newLook.VideoDescription,
		CuratorID:        int(newLook.CuratorID),
		Product:          make([]dto.CreateLookProductResponse, 0, len(look.Products)),
	}

	for _, product := range look.Products {
		response.Product = append(response.Product, dto.CreateLookProductResponse{
			ID:           product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			Category:     product.Category,
			CreatedAt:    time.Now(),
		})
	}

	return response, nil
}
func (r *curatorBORepo) DeleteCollectionByID(collectionID uint) error {
	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var collection domain.Collection
	if err := tx.Where("id = ?", collectionID).First(&collection).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("collection not found")
		}
		return fmt.Errorf("failed to check collection existence: %v", err)
	}

	var sections []domain.CollectionSection
	if err := tx.Where("collection_id = ?", collectionID).Find(&sections).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch collection sections: %v", err)
	}

	for _, section := range sections {
		if err := tx.Where("collection_section_id = ?", section.ID).
			Delete(&domain.CollectionSectionProduct{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete section product entries for section %d: %v", section.ID, err)
		}
	}

	if err := tx.Where("collection_id = ?", collectionID).
		Delete(&domain.CollectionSection{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete collection sections: %v", err)
	}

	if err := tx.Where("collection_id = ?", collectionID).
		Delete(&domain.CollectionProduct{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete collection product entries: %v", err)
	}

	if err := tx.Delete(&collection).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete collection: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *curatorBORepo) DeleteFromFeatureProduct(curatorID uint, productID string) error {
	// Check if the record exists before deletion
	var count int64
	if err := r.db.Model(&domain.CuratorProduct{}).Where("curator_id = ? AND product_id = ? AND is_feature = ?", curatorID, productID, true).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if record exists: %w", err)
	}
	if count == 0 {
		return errors.New("record not found")
	}

	// Define the raw SQL query
	query := `DELETE FROM "curator_products" WHERE "curator_id" = ? AND "product_id" = ? AND "is_feature" = true`

	// Log the query for debugging purposes
	fmt.Printf("Executing query: %s with curatorID=%d, productID=%s\n", query, curatorID, productID)

	// Execute the raw SQL query
	if err := r.db.Exec(query, curatorID, productID).Error; err != nil {
		return fmt.Errorf("failed to delete product from curator products: %w", err)
	}

	return nil
}

func (r *curatorBORepo) DeleteProductFromCollection(collectionID uint, productID string) error {
	// Check if the record exists before deletion
	var count int64
	if err := r.db.Model(&domain.CollectionProduct{}).Where("collection_id = ? AND product_id = ?", collectionID, productID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if record exists: %w", err)
	}
	if count == 0 {
		return errors.New("record not found")
	}

	// Define the raw SQL query
	query := `DELETE FROM "collection_products" WHERE "collection_id" = ? AND "product_id" = ?`

	// Log the query for debugging purposes
	fmt.Printf("Executing query: %s with collectionID=%d, productID=%s\n", query, collectionID, productID)

	// Execute the raw SQL query
	if err := r.db.Exec(query, collectionID, productID).Error; err != nil {
		return fmt.Errorf("failed to delete product from collection: %w", err)
	}

	return nil
}

func (r *curatorBORepo) DeleteLookByID(lookID uint) error {
	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the look exists
	var look domain.Look
	if err := tx.Where("id = ?", lookID).First(&look).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("look not found")
		}
		return fmt.Errorf("failed to check look existence: %v", err)
	}

	// Delete associated lookproduct entries
	if err := tx.Where("look_id = ?", lookID).Delete(&domain.LookProduct{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete look product entries: %v", err)
	}

	// Delete the look
	if err := tx.Delete(&look).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete look: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *curatorBORepo) DeleteSectionByID(sectionID uint) error {
	// Start a transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the section exists
	var section domain.CollectionSection
	if err := tx.Where("id = ?", sectionID).First(&section).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("section with ID %d not found", sectionID)
		}
		return err
	}

	// Delete section products associated with the section
	if err := tx.Where("collection_section_id = ?", sectionID).Delete(&domain.CollectionSectionProduct{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the section itself
	if err := tx.Where("id = ?", sectionID).Delete(&domain.CollectionSection{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *curatorBORepo) DeleteProductFromLook(lookID uint, productID string) error {
	// Check if the record exists before deletion
	var count int64
	if err := r.db.Model(&domain.LookProduct{}).Where("look_id = ? AND product_product_id = ?", lookID, productID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if record exists: %w", err)
	}
	if count == 0 {
		return errors.New("record not found")
	}

	// Define the raw SQL query
	query := `DELETE FROM "look_products" WHERE ("look_products"."look_id", "look_products"."product_product_id") = (?, ?)`

	// Log the query for debugging purposes
	fmt.Printf("Executing query: %s with lookID=%d, productID=%s\n", query, lookID, productID)

	// Execute the raw SQL query
	if err := r.db.Exec(query, lookID, productID).Error; err != nil {
		return fmt.Errorf("failed to delete look product entry: %w", err) // Use %w for error wrapping
	}

	return nil
}

func (r *curatorBORepo) SearchLooksByName(query string) ([]domain.Look, error) {
	var looks []domain.Look
	result := r.db.Preload("Products.Variants.VariantOptions").Where("name LIKE ?", "%"+query+"%").Find(&looks)
	if result.Error != nil {
		return nil, result.Error
	}
	return looks, nil
}

func (r *curatorBORepo) SearchProductsWithinCuratorLooks(curatorID uint, searchQuery string) ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.
		Joins("JOIN look_products ON look_products.look_id = looks.id").
		Joins("JOIN products ON products.product_id = look_products.product_product_id").
		Where("looks.curator_id = ? AND products.name LIKE ?", curatorID, "%"+searchQuery+"%").
		Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *curatorBORepo) FetchProductsByLookID(curatorID uint, lookID uint) (domain.Look, error) {
	// Fetch the look with the specified curator and look IDs
	var look domain.Look
	result := r.db.Preload("Products.Variants").Where("curator_id = ? AND id = ?", curatorID, lookID).First(&look)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Look{}, fmt.Errorf("look not found: %w", result.Error)
		}
		return domain.Look{}, fmt.Errorf("failed to fetch look: %w", result.Error)
	}
	return look, nil
}

func (r *curatorBORepo) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *curatorBORepo) GetCuratorByID(id uint64) (*domain.Curator, error) {
	curator := domain.Curator{}
	result := r.db.Where("id = ?", id).First(&curator)
	if result.Error != nil {
		return nil, result.Error
	}
	return &curator, nil
}
func (r *curatorBORepo) ProcessWithdraw(data []byte) error {
	// Your logic to process a withdrawal
	return nil
}

func (r *curatorBORepo) UpdateProfile(curator *domain.Curator, user *domain.User) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	if err := tx.Model(&curator).Where("id = ?", curator.ID).Updates(curator).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update curator profile: %v", err)
	}

	if err := tx.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Create implements the SocialMediaLinkRepository interface.
func (r *curatorBORepo) CreateSocialMediaLink(link *domain.SocialMediaLink) (*dto.CreateSocialMediaLinkResponse, error) {
	result := r.db.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &dto.CreateSocialMediaLinkResponse{
		ID:           link.ID,
		Type:         link.Type,
		AccessToken:  link.AccessToken,
		OpenID:       link.OpenID,
		RefreshToken: link.RefreshToken,
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
	}

	return response, nil
}

func (r *curatorBORepo) GetSocialMediaLinkByTypeAndCuratorID(linkType string, curatorID uint) (*domain.SocialMediaLink, error) {
	var link domain.SocialMediaLink
	result := r.db.Where("type = ? AND curator_id = ?", linkType, curatorID).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (r *curatorBORepo) UpdateSocialMediaLink(
	link *domain.SocialMediaLink,
	linkID, curatorID uint) (
	*dto.CreateSocialMediaLinkResponse,
	error) {

	result := r.db.Model(&domain.SocialMediaLink{}).
		Where("id = ? AND curator_id = ?", linkID, curatorID).
		Updates(link)

	if result.RowsAffected == 0 {
		return nil, errors.New("no social media link found for the given curatorID and linkID")
	}

	if result.Error != nil {
		return nil, result.Error
	}
	response := &dto.CreateSocialMediaLinkResponse{
		ID:           linkID,
		Type:         link.Type,
		AccessToken:  link.AccessToken,
		OpenID:       link.OpenID,
		RefreshToken: link.RefreshToken,
	}
	return response, nil
}

func (r *curatorBORepo) DeleteSocialMediaLink(curatorID, linkID uint64) (*dto.DeleteSocialMediaLinkResponse, error) {
	result := r.db.Where("curator_id = ? AND id = ?", curatorID, linkID).Delete(&domain.SocialMediaLink{})
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return &dto.DeleteSocialMediaLinkResponse{
			Message: "No social media link found for the given curatorID and linkID",
		}, nil
	}

	response := &dto.DeleteSocialMediaLinkResponse{
		Message: "Social media link deleted successfully",
	}

	return response, nil
}

func (r *curatorBORepo) UpdatePasswordByEmail(email string, newPassword string) (error, *dto.UpdatePasswordResponse) {
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Update("password_hash", newPassword).Error
	if err != nil {
		return fmt.Errorf("failed to update password for email %s: %s", email, err), nil
	}

	res := dto.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}

	return nil, &res
}

func (r *curatorBORepo) GetAllCurators() ([]*dto.GetCuratorResponse, error) {
	var curators []domain.Curator

	// Retrieve curators data with preloaded User and CuratorProduct data
	result := r.db.Preload("User").Preload("CuratorProducts").Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	// Collect all product IDs for batch retrieval of variants
	productIDMap := make(map[string]struct{})
	for _, curator := range curators {
		for _, cp := range curator.CuratorProducts {
			productIDMap[cp.ProductID] = struct{}{}
		}
	}

	// Convert map keys to slice
	productIDs := make([]string, 0, len(productIDMap))
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	// Retrieve all variants for the collected product IDs in one query
	var variants []domain.Variant
	result = r.db.Where("product_id IN ?", productIDs).Find(&variants)
	if result.Error != nil {
		return nil, result.Error
	}

	// Create a map to store the first image of each product
	productImageMap := make(map[string]string)
	for _, variant := range variants {
		if _, exists := productImageMap[variant.ProductID]; !exists && variant.Image != "" {
			productImageMap[variant.ProductID] = variant.Image
		}
	}

	// Build the response
	var responses []*dto.GetCuratorResponse
	for _, curator := range curators {
		var variantImages []string
		for _, cp := range curator.CuratorProducts {
			if image, exists := productImageMap[cp.ProductID]; exists {
				variantImages = append(variantImages, image)
			}
		}
		response := &dto.GetCuratorResponse{
			ID:                curator.ID,
			UserID:            curator.User.ID,
			ShopName:          curator.ShopName,
			Bio:               curator.Bio,
			FirstName:         curator.User.FirstName,
			LastName:          curator.User.LastName,
			Email:             curator.User.Email,
			NumberofFollowers: curator.NumberofFollowers,
			ProfileImageURL:   curator.ProfileImageURL,
			CoverImageURL:     curator.CoverImageURL,
			CreatedAt:         curator.CreatedAt,
			UpdatedAt:         curator.UpdatedAt,
			VariantImages:     variantImages,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (r *curatorBORepo) GetCuratorByCuratorID(id uint64) (*domain.Curator, error) {
	curator := domain.Curator{}
	result := r.db.Where("id = ?", id).
		Preload("User").
		Preload("SocialMediaLinks").
		Find(&curator)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curator, nil
}

func (r *curatorBORepo) GetLookWithId(id uint, curatorId uint) (*domain.Look, error) {
	look := domain.Look{}

	if err := r.db.Where("id = ? AND curator_id = ?", id, curatorId).First(&look).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve look: %v", err)
	}

	return &look, nil
}

func (r *curatorBORepo) GetFeatureWithId(id uint, curatorId uint) (*domain.CuratorProduct, error) {
	feature := domain.CuratorProduct{}

	if err := r.db.Where("feature_id = ? AND curator_id = ?", id, curatorId).First(&feature).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve feature: %v", err)
	}
	return &feature, nil
}

func (r *curatorBORepo) GetCollectionWithId(id uint, curatorId uint) (*domain.Collection, error) {
	collection := domain.Collection{}

	if err := r.db.Where("id = ? AND curator_id = ?", id, curatorId).First(&collection).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve collection: %v", err)
	}

	return &collection, nil
}

func (r *curatorBORepo) GetProductWithId(id string) (*domain.Product, error) {
	product := domain.Product{}

	if err := r.db.Where("product_id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *curatorBORepo) StartTx() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	return tx, nil
}

func (r *curatorBORepo) AddProduct(tx *gorm.DB, product *domain.Product) (*gorm.DB, *domain.Product, error) {
	if err := tx.Create(product).Error; err != nil {
		tx.Rollback()
		return tx, nil, fmt.Errorf("failed to create product: %v", err)
	}

	return tx, product, nil
}

func (r *curatorBORepo) AddProductVariant(tx *gorm.DB, variant *domain.Variant) (*gorm.DB, *domain.Variant, error) {
	if err := tx.Create(variant).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create variant: %v", err)
	}

	return tx, variant, nil
}

func (r *curatorBORepo) AddProductToCollection(tx *gorm.DB, collectionProduct *domain.CollectionProduct) (*gorm.DB, error) {
	if err := tx.Create(collectionProduct).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add product to collection: %v", err)
	}

	return tx, nil
}

func (r *curatorBORepo) CheckProductExistsInCollection(collectionId uint, productId string) (bool, error) {
	var count int64

	// Use count to check existence
	if err := r.db.Model(&domain.CollectionProduct{}).Where("collection_id = ? AND product_id = ?", collectionId, productId).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence: %v", err)
	}

	return count > 0, nil
}

func (r *curatorBORepo) GetVariantWithId(id string) (*domain.Variant, error) {
	variant := domain.Variant{}

	if err := r.db.Where("id = ?", id).First(&variant).Error; err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *curatorBORepo) CheckProductExistsInFeature(featureId uint, productId string) (bool, error) {
	var count int64

	// Use count to check existence
	if err := r.db.Model(&domain.CuratorProduct{}).Where("feature_id = ? AND product_product_id = ?", featureId, productId).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence: %v", err)
	}

	return count > 0, nil
}

func (r *curatorBORepo) CheckProductExistsInLook(lookId uint, productId string) (bool, error) {
	var count int64

	// Use count to check existence
	if err := r.db.Model(&domain.LookProduct{}).Where("look_id = ? AND product_product_id = ?", lookId, productId).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence: %v", err)
	}

	return count > 0, nil
}

func (r *curatorBORepo) AddProductToFeature(tx *gorm.DB, featureProduct *domain.CuratorProduct) (*gorm.DB, error) {
	query := `INSERT INTO "curator_products" ("feature_id", "product_id", "is_feature", "created_at") VALUES (?, ?, true, CURRENT_TIMESTAMP)`
	if err := tx.Exec(query, featureProduct.FeatureID, featureProduct.ProductID).Error; err != nil {
		tx.Rollback()
		return tx, fmt.Errorf("failed to add product to feature: %v", err)
	}

	return tx, nil
}

func (r *curatorBORepo) AddProductToLook(tx *gorm.DB, lookProduct *domain.LookProduct) (*gorm.DB, error) {
	query := `INSERT INTO "look_products" ("look_id","product_product_id") VALUES (?,?)`
	if err := tx.Exec(query, lookProduct.LookID, lookProduct.ProductID).Error; err != nil {
		tx.Rollback()
		return tx, fmt.Errorf("failed to add product to look: %v", err)
	}

	return tx, nil
}

func (r *curatorBORepo) GetCuratorWithUserId(id uint64) (*domain.Curator, error) {
	curators := domain.Curator{}
	result := r.db.Where("user_id = ?", id).Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curators, nil
}

func (r *curatorBORepo) UpdateCollection(collection *domain.Collection) error {
	if err := r.db.Model(&domain.Collection{}).Where("id = ?", collection.ID).Updates(collection).Error; err != nil {
		return fmt.Errorf("failed to update curator collection: %v", err)
	}
	return nil
}

func (r *curatorBORepo) UpdateLook(look *domain.Look) error {
	if err := r.db.Model(&domain.Look{}).Where("id = ?", look.ID).Updates(look).Error; err != nil {
		return fmt.Errorf("failed to update curator look: %v", err)
	}
	return nil
}

func (r *curatorBORepo) GetUserByID(userId uint64) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *curatorBORepo) GetCuratorAccountDetail(curatorId uint) (*struct {
	domain.BankInformation
	domain.AccountDetails
}, error) {
	var results struct {
		domain.BankInformation
		domain.AccountDetails
	}

	query := r.db.Table("bank_informations").
		Select("bank_informations.*, account_details.*").
		Joins("left join account_details on account_details.bank_id = bank_informations.id").
		Where("bank_informations.curator_id = ?", curatorId).
		Scan(&results)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected == 0 {
		return nil, fmt.Errorf("curator account not found")
	}

	return &results, nil
}

func (r *curatorBORepo) UpdateCuratorAccountDetail(curatorId uint, request *dto.UpdateAccountDetailRequest) (*struct {
	domain.BankInformation
	domain.AccountDetails
}, error) {
	var details struct {
		domain.BankInformation
		domain.AccountDetails
	}

	query := r.db.Table("bank_informations").
		Select("bank_informations.*, account_details.*").
		Joins("left join account_details on account_details.bank_id = bank_informations.id").
		Where("bank_informations.curator_id = ?", curatorId).
		First(&details)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected == 0 {
		return nil, fmt.Errorf("curator account not found")
	}

	r.db.Model(&details.BankInformation).Updates(domain.BankInformation{
		Location:    request.Location,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		DateOfBirth: request.DateOfBirth,
	})

	r.db.Model(&details.AccountDetails).Updates(domain.AccountDetails{
		BankAddress:    request.BankAddress,
		BankName:       request.BankName,
		BranchCode:     request.BranchCode,
		AccountNumber:  request.AccountNumber,
		AccountName:    request.AccountName,
		AccountAddress: request.AccountAddress,
		IBAN:           request.IBAN,
	})

	return &details, nil
}
