package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"

	domain "github.com/harishash/dotshop-be/internal/models"
	pagination "github.com/harishash/dotshop-be/internal/utils/pagination"
)

const (
	STATUS_ACTIVE  = "Active"
	STATUS_BLOCKED = "Blocked"
)

type AdminDashboardRepository interface {
	GetTopSellingBrands(query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingBrandsResponse, *dto.Paging, error)
	GetTopSellingProducts(query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingProductResponse, *dto.Paging, error)
	GetTopCurators(query *dto.CommonProductRequest) ([]dto.GetTopCuratorsResponse, *dto.Paging, error)
	GetSalesByCategory(query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error)
	GetTopWishlistProducts(query *dto.CommonProductRequest) ([]dto.GetCuratorTopWishlistResponse, *dto.Paging, error)
	GetOrderSales(query *dto.OrderSalesRequest) ([]dto.GetOrderSalesResponse, *dto.Paging, error)
	GetOrderReturns(query *dto.OrderSalesRequest) ([]dto.GetOrderReturnsResponse, *dto.Paging, error)
	GetAllCustomers(query *dto.CustomersRequest) ([]dto.GetAllCustomerResponse, *dto.Paging, error)
	GetAllCurators(query *dto.CuratorRequest) ([]dto.GetAllCuratorResponse, *dto.Paging, error)
	GetPaymentDistribution(query *dto.PaymentDistributionRequest, curatorID uint) ([]*dto.ProductData, *dto.Paging, error)
	UpdateReturnStatus(id uint, status string) error
	DeleteUser(id uint) error
	DeleteCurator(id uint) error
	UpdateCustomerStatus(user *models.User) error
	UpdateCuratorStatus(curator *models.Curator) error
	FindCustomerByID(id uint) (*models.User, error)
	FindCuratorByID(id uint) (*models.Curator, error)
	GetListedProducts(params *dto.PagingParams, curatorID uint) ([]*dto.ListedProducts, *dto.Paging, error)
}

type adminDashboardRepository struct {
	db          *gorm.DB
	productRepo IProductRepository
}

func NewAdminDashboardRepository(db *gorm.DB, productRepo IProductRepository) AdminDashboardRepository {
	return &adminDashboardRepository{db: db, productRepo: productRepo}
}

func (r *adminDashboardRepository) GetTopSellingBrands(query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingBrandsResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		return db
	}

	return getTopSellingBrands(r.db, query, adminFilter)
}

func (r *adminDashboardRepository) GetTopSellingProducts(query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingProductResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		return db
	}

	return getTopSellingProducts(r.db, query, adminFilter)
}

func (r *adminDashboardRepository) GetTopCurators(query *dto.CommonProductRequest) ([]dto.GetTopCuratorsResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		return db
	}

	return getTopCurators(r.db, query, adminFilter)
}

type topCurators struct {
	UserID        uint64
	TotalPrice    float64
	TotalQuantity uint
}

func getTopCurators(db *gorm.DB, query *dto.CommonProductRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetTopCuratorsResponse, *dto.Paging, error) {
	var topCurators []topCurators
	var totalCount int64

	dbQuery := db.Table("orders").
		Select("user_id, SUM(total_amount) as total_price, SUM(total_quantity) as total_quantity").
		Where("total_quantity != 0 AND cancelled = false").
		Group("user_id")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToQuery(dbQuery)

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	err = dbQuery.
		Order("total_quantity DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&topCurators).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetTopCuratorsResponse
	for _, v := range topCurators {
		var curator domain.Curator
		if err := db.Table("curators").Where("user_id = ?", v.UserID).Preload("User").Find(&curator).Error; err != nil {
			log.Warnf("unable to get curator detail with user id %v: %v", v.UserID, err)
			continue
		}

		resp = append(resp, dto.GetTopCuratorsResponse{
			CuratorID:    curator.ID,
			FirstName:    curator.User.FirstName,
			LastName:     curator.User.LastName,
			ProfileImage: curator.ProfileImageURL,
			Earnings:     fmt.Sprintf("%.2f", v.TotalPrice),
			Sales:        v.TotalQuantity,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *adminDashboardRepository) GetSalesByCategory(query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		return db
	}

	return getSalesByCategory(r.db, query, adminFilter)
}

func (r *adminDashboardRepository) GetTopWishlistProducts(query *dto.CommonProductRequest) ([]dto.GetCuratorTopWishlistResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB { return db }
	return getTopWishlistProducts(r.db, query, adminFilter)
}

func (r *adminDashboardRepository) GetAllCurators(query *dto.CuratorRequest) ([]dto.GetAllCuratorResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		return db.Order("curators.created_at DESC")
	}

	return getAllCurators(r.db, query, adminFilter)
}

func (r *adminDashboardRepository) GetPaymentDistribution(query *dto.PaymentDistributionRequest, curatorID uint) (
	[]*dto.ProductData, *dto.Paging, error) {
	var totalCount int64
	var result []*dto.ProductData
	var commission float64 = constants.CuratorCommissionPercentage / 100

	subQuery := r.db.Table("variants v").
		Select("v.image").
		Where("v.product_id = p.product_id").
		Order("v.id").
		Limit(1)

	mainQuery := r.db.Table("products p").
		Select(`p.name as "product_name",
			(?) as "product_image",
			b.margin as "total_margin",
			COALESCE(SUM(ov.price * ov.quantity), 0) as "total_sales",
			COALESCE(SUM(CASE WHEN ov.curator_id = ? THEN (ov.price * ov.quantity) * 0.25 ELSE 0 END), 0) as "dotshop_profit",
			COALESCE(SUM(CASE WHEN ov.curator_id != ? THEN (ov.price * ov.quantity) * ? ELSE 0 END), 0) as "curator_commission"`,
			subQuery, curatorID, curatorID, commission).
		Joins("LEFT JOIN brands b on p.brand_id = b.id").
		Joins("LEFT JOIN order_variants ov on p.product_id = ov.product_id").
		Group("p.product_id, p.name, b.margin")

	err := mainQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	pagination.NewPaginate(mainQuery, query.PageNum, query.PageSize)
	if err = mainQuery.Find(&result).Error; err != nil {
		return nil, nil, err
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return result, paging, nil
}

func getAllCurators(db *gorm.DB, query *dto.CuratorRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetAllCuratorResponse, *dto.Paging, error) {
	var totalCount int64
	var curatorData []struct {
		CuratorID    uint    `json:"curator_id"`
		FirstName    string  `json:"first_name"`
		LastName     string  `json:"last_name"`
		Email        string  `json:"email"`
		ShopName     string  `json:"shop_name"`
		Image        string  `json:"image"`
		TotalRevenue float64 `json:"total_revenue"`
		NoOfOrders   int64   `json:"no_of_orders"`
		ItemsSold    int64   `json:"items_sold"`
		IsBlock      bool    `json:"is_block"`
	}

	dbQuery := db.Model(&domain.Curator{}).Select(`
		curators.id as curator_id,
		users.first_name as first_name,
		users.last_name as last_name,
		users.email as email,
		curators.name as name,
		curators.is_block as is_block,
		curators.profile_image_url as image,
		curators.shop_name as shop_name,
		COALESCE(SUM(order_variants.price)* 0.15, 0) AS total_revenue,
		COUNT(DISTINCT orders.id) as no_of_orders,
		COALESCE(SUM(order_variants.quantity), 0) as items_sold
	`).
		Joins("LEFT JOIN users ON users.id = curators.user_id").
		Joins("LEFT JOIN order_variants ON order_variants.curator_id = curators.id").
		Joins("LEFT JOIN orders ON orders.id = order_variants.order_id").
		Group("curators.id, users.first_name, users.last_name, users.email, curators.shop_name")

	dbQuery = filterFunc(dbQuery)

	err := db.Model(&domain.Curator{}).
		Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	err = dbQuery.
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&curatorData).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetAllCuratorResponse
	for _, data := range curatorData {
		status := STATUS_ACTIVE
		if data.IsBlock {
			status = STATUS_BLOCKED
		}
		resp = append(resp, dto.GetAllCuratorResponse{
			CuratorID:    data.CuratorID,
			FirstName:    data.FirstName,
			LastName:     data.LastName,
			Email:        data.Email,
			Image:        data.Image,
			ShopName:     data.ShopName,
			TotalRevenue: data.TotalRevenue,
			NoOfOrders:   int(data.NoOfOrders),
			ItemsSold:    int(data.ItemsSold),
			Status:       status,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *adminDashboardRepository) GetOrderSales(query *dto.OrderSalesRequest) ([]dto.GetOrderSalesResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		if len(query.SortBy) > 0 {
			db = pagination.SetOrderSalesSort(db, query.SortBy)
		} else {
			db = db.Order("updated_at DESC")
		}

		return db
	}

	return getOrderSales(r.db, query, adminFilter)
}

func getOrderSales(db *gorm.DB, query *dto.OrderSalesRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetOrderSalesResponse, *dto.Paging, error) {
	var totalCount int64
	var orders []domain.Order

	dbQuery := db.Model(&domain.Order{}).Preload("User").Preload("OrderVariants")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToQuery(dbQuery)

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	err = dbQuery.
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&orders).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetOrderSalesResponse
	for _, v := range orders {
		var userAddress domain.ShippingInfo
		err := db.Where("user_id = ? AND default_address = ?", v.UserID, true).First(&userAddress).Error
		if err != nil {
			err = db.Where("user_id = ?", v.UserID).First(&userAddress).Error
			if err != nil {
				log.Warnf("Unable to get user address detail with user id %v: %v", v.UserID, err)
			}
		}

		var variantList []dto.OrderVariantResponse

		for _, ov := range v.OrderVariants {
			// Initialize default values
			var productName, variantImage string

			// Fetch the variant and product details in a single query, where possible
			if err := db.Model(&domain.Variant{}).
				Where("id = ?", ov.VariantID).
				Select("image").
				Take(&variantImage).
				Error; err != nil {
				log.Warnf("Unable to get variant detail with variant id %v: %v", ov.VariantID, err)
			}

			if err := db.Model(&domain.Product{}).
				Where("product_id = ?", ov.ProductID).
				Select("name").
				Take(&productName).
				Error; err != nil {
				log.Warnf("Unable to get product detail with product id %v: %v", ov.ProductID, err)
			}

			// Create a variant response object
			variantResponse := dto.OrderVariantResponse{
				ProductName:  productName,
				Brand:        ov.BrandName,
				Size:         ov.VariantSize,
				Price:        ov.Price,
				Quantity:     int(ov.Quantity),
				VariantImage: variantImage,
			}

			// Append to variant list
			variantList = append(variantList, variantResponse)
		}

		resp = append(resp, dto.GetOrderSalesResponse{
			ID:                v.ID,
			UserID:            v.UserID,
			CustomerFirstName: v.User.FirstName,
			CustomerLastName:  v.User.LastName,
			ShippingMethod:    v.ShippingMethod,
			PaymentID:         v.PaymentID,
			CustomerAddress:   &userAddress.AddressOne.String,
			CustomerCity:      &userAddress.City.String,
			CustomerState:     &userAddress.State.String,
			CustomerCountry:   &userAddress.Country.String,
			CustomerZip:       &userAddress.Zip.String,
			ShippingAddress:   v.ShippingAddress,
			ShippingCity:      v.ShippingCity,
			ShippingState:     v.ShippingState,
			ShippingCountry:   v.ShippingCountry,
			ShippingZip:       v.ShippingZip,
			TotalAmount:       v.TotalAmount,
			TotalQuantity:     v.TotalQuantity,
			Status:            v.Status,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
			Variants:          variantList,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *adminDashboardRepository) GetOrderReturns(query *dto.OrderSalesRequest) ([]dto.GetOrderReturnsResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		// if len(query.SortBy) > 0 {
		// 	db = pagination.SetOrderSalesSort(db, query.SortBy)
		// } else {
		// 	db = db.Order("updated_at DESC")
		// }

		return db
	}

	return getOrderReturns(r.db, query, adminFilter)
}

func getOrderReturns(db *gorm.DB, query *dto.OrderSalesRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetOrderReturnsResponse, *dto.Paging, error) {
	var totalCount int64
	var returnOrders []domain.ReturnOrder

	dbQuery := db.Model(&domain.ReturnOrder{}).Preload("OrderVariants.Order")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToQuery(dbQuery)

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	err = dbQuery.
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&returnOrders).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetOrderReturnsResponse
	for _, v := range returnOrders {
		var userAddress domain.ShippingInfo
		err := db.Where("user_id = ? AND default_address = ?", v.UserId, true).First(&userAddress).Error
		if err != nil {
			err = db.Where("user_id = ?", v.UserId).First(&userAddress).Error
			if err != nil {
				log.Warnf("Unable to get user address detail with user id %v: %v", v.UserId, err)
			}
		}
		product := domain.Product{}
		if err := db.Where("product_id = ?", v.OrderVariants.ProductID).First(&product).Error; err != nil {
			log.Warn("Unable to get product detail with product id %v: %v", v.OrderVariants.ProductID, err)
		}

		variant := domain.Variant{}
		if err := db.Where("id = ?", v.OrderVariants.VariantID).First(&variant).Error; err != nil {
			log.Warn("Unable to get variant image with variant id %v: %v", v.OrderVariants.VariantID, err)
		}

		orderReturnDetails := dto.OrderReturnDetails{
			OrderVariantId:    uint(v.OrderVariants.ID),
			ProductId:         v.OrderVariants.ProductID,
			BrandName:         v.OrderVariants.BrandName,
			ProductName:       product.Name,
			VariantOptionName: v.OrderVariants.VariantOptionName,
			VariantSize:       v.OrderVariants.VariantSize,
			VariantImage:      variant.Image,
			Reason:            v.Reason,
			TotalAmount:       v.Amount,
			TotalQuantity:     v.Quantity,
			CreatedAt:         v.CreatedAt,
		}

		resp = append(resp, dto.GetOrderReturnsResponse{
			ID:                uint(v.ID),
			UserID:            uint64(v.UserId),
			ReturnId:          v.ReturnId,
			OrderVariantId:    uint(v.OrderVariantId),
			CuratorId:         uint(v.CuratorId),
			CustomerFirstName: v.User.FirstName,
			CustomerLastName:  v.User.LastName,
			CustomerAddress:   &userAddress.AddressOne.String,
			CustomerCity:      &userAddress.City.String,
			CustomerState:     &userAddress.State.String,
			CustomerCountry:   &userAddress.Country.String,
			CustomerZip:       &userAddress.Zip.String,
			ShippingAddress:   v.OrderVariants.Order.ShippingAddress,
			ShippingCity:      v.OrderVariants.Order.ShippingCity,
			ShippingState:     v.OrderVariants.Order.ShippingState,
			ShippingCountry:   v.OrderVariants.Order.ShippingCountry,
			ShippingZip:       v.OrderVariants.Order.ShippingZip,
			TotalAmount:       v.Amount,
			TotalQuantity:     v.Quantity,
			Status:            v.Status,
			OrderReturnDetail: orderReturnDetails,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *adminDashboardRepository) GetAllCustomers(query *dto.CustomersRequest) ([]dto.GetAllCustomerResponse, *dto.Paging, error) {
	adminFilter := func(db *gorm.DB) *gorm.DB {
		// if len(query.SortBy) > 0 {
		// 	db = pagination.SetOrderSalesSort(db, query.SortBy)
		// } else {
		// 	db = db.Order("updated_at DESC")
		// }

		return db.Order("created_at DESC")
	}

	return getAllCustomers(r.db, query, adminFilter)
}

func getAllCustomers(db *gorm.DB, query *dto.CustomersRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetAllCustomerResponse, *dto.Paging, error) {
	var totalCount int64
	var customerData []struct {
		ID            uint       `json:"id"`
		Email         string     `json:"email"`
		FirstName     *string    `json:"first_name"`
		LastName      *string    `json:"last_name"`
		PhoneNumber   *string    `json:"phone_number"`
		TotalPrice    float64    `json:"total_price"`
		TotalQuantity uint       `json:"total_quantity"`
		LastOrderDate *time.Time `json:"last_order_date"`
		OrderCount    int64      `json:"order_count"`
		CreatedAt     time.Time  `json:"created_at"`
	}

	dbQuery := db.Model(&domain.User{}).Select(`
		users.id,
		users.email,
		users.first_name,
		users.last_name,
		users.phone_number,
		users.created_at,
		COALESCE(SUM(orders.total_amount), 0) as total_price,
		COALESCE(SUM(orders.total_quantity), 0) as total_quantity,
		MAX(orders.created_at) as last_order_date,
		COUNT(orders.id) as order_count
	`).
		Joins("LEFT JOIN curators ON curators.user_id = users.id").
		Joins("LEFT JOIN orders ON orders.user_id = users.id AND orders.total_quantity != 0").
		Where("curators.id IS NULL").
		Group("users.id")

	dbQuery = filterFunc(dbQuery)

	err := db.Model(&domain.User{}).
		Joins("LEFT JOIN curators ON curators.user_id = users.id").
		Where("curators.id IS NULL").
		Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	err = dbQuery.
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&customerData).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetAllCustomerResponse
	for _, data := range customerData {
		resp = append(resp, dto.GetAllCustomerResponse{
			ID:            data.ID,
			Email:         data.Email,
			FirstName:     data.FirstName,
			LastName:      data.LastName,
			PhoneNumber:   data.PhoneNumber,
			LastOrderDate: data.LastOrderDate,
			LifeTimeSpend: data.TotalPrice,
			Orders:        data.OrderCount,
			Items:         data.TotalQuantity,
			CreatedAt:     data.CreatedAt,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *adminDashboardRepository) UpdateReturnStatus(id uint, status string) error {
	result := r.db.Model(&domain.ReturnOrder{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (r *adminDashboardRepository) DeleteUser(id uint) error {
	return r.db.Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *adminDashboardRepository) DeleteCurator(id uint) error {
	// Begin a transaction
	tx := r.db.Begin()

	// Retrieve the user_id from the Curator table
	var curator domain.Curator
	if err := tx.Where("id = ?", id).First(&curator).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the deleted_at value in the Users table to the current time
	if err := tx.Model(&domain.User{}).Where("id = ?", curator.UserID).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the curator
	if err := tx.Where("id = ?", id).Delete(&domain.Curator{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *adminDashboardRepository) FindCustomerByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *adminDashboardRepository) UpdateCustomerStatus(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *adminDashboardRepository) UpdateCuratorStatus(curator *models.Curator) error {
	return r.db.Save(curator).Error

}
func (r *adminDashboardRepository) FindCuratorByID(id uint) (*models.Curator, error) {
	var curator models.Curator
	result := r.db.First(&curator, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &curator, result.Error
}
func (r *adminDashboardRepository) GetListedProducts(params *dto.PagingParams, curatorID uint) ([]*dto.ListedProducts, *dto.Paging, error) {
	var products []*dto.ListedProducts

	query := `
    WITH CollectionProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN collection_products cp ON p.product_id = cp.product_id
        WHERE cp.collection_id IN (
            SELECT id FROM collections WHERE curator_id = ?
        )
    ),
    SectionProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN collection_section_products csp ON p.product_id = csp.product_id
        WHERE csp.collection_section_id IN (
            SELECT id FROM collection_sections WHERE collection_id IN (
                SELECT id FROM collections WHERE curator_id = ?
            )
        )
    ),
    LookProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN look_products lp ON p.product_id = lp.product_product_id
        WHERE lp.look_id IN (
            SELECT id FROM looks WHERE curator_id = ?
        )
    ),
    CombinedProducts AS (
        SELECT product_id FROM CollectionProducts
        UNION
        SELECT product_id FROM SectionProducts
        UNION
        SELECT product_id FROM LookProducts
    ),
    ProductPrices AS (
        SELECT p.product_id, p.name, p.brand_name
        FROM CombinedProducts cp
        JOIN products p ON cp.product_id = p.product_id
        GROUP BY p.product_id, p.name, p.brand_name
    )
    SELECT * FROM ProductPrices
    LIMIT ? OFFSET ?;
    `

	// Execute query
	err := r.db.Raw(query, curatorID, curatorID, curatorID, params.PageSize, (params.PageNum-1)*params.PageSize).Scan(&products).Error
	if err != nil {
		return nil, nil, err
	}

	// Total count for pagination
	var totalCount int64
	countQuery := r.db.Raw(`
    WITH CollectionProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN collection_products cp ON p.product_id = cp.product_id
        WHERE cp.collection_id IN (
            SELECT id FROM collections WHERE curator_id = ?
        )
    ),
    SectionProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN collection_section_products csp ON p.product_id = csp.product_id
        WHERE csp.collection_section_id IN (
            SELECT id FROM collection_sections WHERE collection_id IN (
                SELECT id FROM collections WHERE curator_id = ?
            )
        )
    ),
    LookProducts AS (
        SELECT p.product_id
        FROM products p
        JOIN look_products lp ON p.product_id = lp.product_product_id
        WHERE lp.look_id IN (
            SELECT id FROM looks WHERE curator_id = ?
        )
    ),
    CombinedProducts AS (
        SELECT product_id FROM CollectionProducts
        UNION
        SELECT product_id FROM SectionProducts
        UNION
        SELECT product_id FROM LookProducts
    )
    SELECT COUNT(*) FROM CombinedProducts;
    `, curatorID, curatorID, curatorID).Scan(&totalCount).Error
	if countQuery != nil {
		return nil, nil, countQuery
	}

	// Create paging object
	paging := &dto.Paging{
		TotalCount:  totalCount,
		CurrentPage: params.PageNum,
		PageSize:    params.PageSize,
	}

	return products, paging, nil
}
