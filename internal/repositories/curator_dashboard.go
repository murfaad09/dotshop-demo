package repository

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2/log"
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
)

type CuratorDashboardRepository interface {
	GetOrders(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.OrderCount, error)
	GetUnitsSold(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.UnitsSold, error)
	GetSales(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.SalesIntervalResult, error)
	GetRevenue(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.RevenueIntervalResult, error)
	GetCuratorTopWishlistProduct(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopWishlistResponse, *dto.Paging, error)
	GetCuratorTopSellingBrands(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingBrandsResponse, *dto.Paging, error)
	GetCuratorTopSellingProducts(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingProductResponse, *dto.Paging, error)
	GetOrderVariants(curatorID uint, startDate, endDate *time.Time) ([]domain.OrderVariants, error)
	GetCuratorTopPurchasers(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopPurchasersResponse, *dto.Paging, error)
	GetAverageOrderValue(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.AOVIntervalResultResponse, error)
	GetOrderReturns(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.OrderReturns, error)
	GetAverageUnitsPerOrder(fromTimestamp, toTimestamp, interval string, curatorID uint) ([]dto.UnitsSoldPerOrder, error)
	GetCuratorSalesByCategory(curatorID uint, query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error)
}

type curatorDashboardRepository struct {
	db          *gorm.DB
	productRepo IProductRepository
}

func NewCuratorDashboardRepository(db *gorm.DB, productRepo IProductRepository) CuratorDashboardRepository {
	return &curatorDashboardRepository{db: db, productRepo: productRepo}
}

func (repo *curatorDashboardRepository) GetOrders(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.OrderCount,
	error) {

	var results []dto.OrderCount

	query := fmt.Sprintf(`
		WITH date_intervals AS (
			SELECT generate_series(
				?::timestamp, 
				?::timestamp, 
				'1 %s'::interval
			) AS interval_start
		)
		SELECT
			interval_start,
			COUNT(DISTINCT orders.id) AS order_count
		FROM
			date_intervals
		LEFT JOIN
			order_variants
		ON
			order_variants.created_at >= interval_start
			AND order_variants.created_at < interval_start + '1 %s'::interval
			AND order_variants.curator_id = ?
		LEFT JOIN
			orders
		ON
			orders.id = order_variants.order_id
		GROUP BY
			interval_start
		ORDER BY
			interval_start;
		`, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *curatorDashboardRepository) GetUnitsSold(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.UnitsSold,
	error) {

	var results []dto.UnitsSold

	query := fmt.Sprintf(`
        WITH date_intervals AS (
            SELECT generate_series(
                ?::timestamp, 
                ?::timestamp, 
                '1 %s'::interval
            ) AS interval_start
        )
        SELECT
            interval_start,
            COALESCE(SUM(order_variants.quantity), 0) AS units_sold
        FROM
            date_intervals
        LEFT JOIN
            order_variants
        ON
            order_variants.created_at >= interval_start
            AND order_variants.created_at < interval_start + '1 %s'::interval
			AND order_variants.curator_id = ?
        GROUP BY
            interval_start
        ORDER BY
            interval_start;
    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *curatorDashboardRepository) GetOrderVariants(
	curatorID uint,
	startDate, endDate *time.Time) (
	[]domain.OrderVariants, error) {

	dbQuery := r.db.Table("order_variants").
		Where("created_at BETWEEN ? AND ?", *startDate, *endDate).
		Where("curator_id = ?", curatorID)

	var orderVariants []domain.OrderVariants
	err := dbQuery.Find(&orderVariants).Error
	if err != nil {
		return nil, err
	}

	return orderVariants, nil
}

func (repo *curatorDashboardRepository) GetOrderReturns(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.OrderReturns,
	error) {

	var orderReturns []dto.OrderReturns

	query := fmt.Sprintf(`
        WITH date_intervals AS (
            SELECT generate_series(
                ?::timestamp, 
                ?::timestamp, 
                '1 %s'::interval
            ) AS interval_start
        )
        SELECT
			interval_start,
			COALESCE(SUM(COALESCE(return_orders.quantity, 0)), 0) AS total_returned_quantity
		FROM
			date_intervals
		LEFT JOIN
			return_orders
		ON
			return_orders.created_at >= interval_start
					AND return_orders.created_at < interval_start + '1 %s'::interval
		LEFT JOIN
			order_variants
		ON
			return_orders.order_variant_id = order_variants.id
			AND order_variants.curator_id = ?
		GROUP BY
			interval_start
		ORDER BY
			interval_start;
    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).Scan(&orderReturns).Error; err != nil {
		return nil, err
	}

	// Log the results for debugging
	fmt.Println("Query results:", orderReturns)

	// Check for empty results and handle accordingly
	if len(orderReturns) == 0 {
		return []dto.OrderReturns{}, nil
	}

	return orderReturns, nil
}

type topWishlistProduct struct {
	ProductID    string
	ProductName  string
	ProductImage string
	ProductBrand string
	ProductPrice string
	Count        int64
}

func (r *curatorDashboardRepository) GetCuratorTopWishlistProduct(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopWishlistResponse, *dto.Paging, error) {
	curatorFilter := func(db *gorm.DB) *gorm.DB {
		return db.Where("curator_id = ?", curatorID)
	}

	return getTopWishlistProducts(r.db, query, curatorFilter)
}

func getTopWishlistProducts(db *gorm.DB, query *dto.CommonProductRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetCuratorTopWishlistResponse, *dto.Paging, error) {
	var topWishlistProducts []topWishlistProduct
	var totalCount int64

	dbQuery := db.Table("wishlist_items").
		Select("wishlist_items.product_id, wishlist_items.product_name, wishlist_items.product_image, wishlist_items.product_brand, wishlist_items.product_price, COUNT(*) as count").
		Joins("JOIN products ON products.product_id = wishlist_items.product_id").
		Where("products.is_active = ?", true).
		Group("wishlist_items.product_id, wishlist_items.product_name, wishlist_items.product_image, wishlist_items.product_brand, wishlist_items.product_price")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToQuery(dbQuery)

	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	err := dbQuery.
		Order("count DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&topWishlistProducts)
	if err.Error != nil {
		return nil, nil, err.Error
	}

	var resp []dto.GetCuratorTopWishlistResponse
	for _, v := range topWishlistProducts {
		resp = append(resp, dto.GetCuratorTopWishlistResponse{
			ProductID:    v.ProductID,
			ProductName:  v.ProductName,
			ProductImage: v.ProductImage,
			ProductBrand: v.ProductBrand,
			ProductPrice: v.ProductPrice,
			Count:        v.Count,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (r *curatorDashboardRepository) GetCuratorTopSellingProducts(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingProductResponse, *dto.Paging, error) {
	curatorFilter := func(db *gorm.DB) *gorm.DB {
		return db.Where("curator_id = ?", curatorID)
	}

	return getTopSellingProducts(r.db, query, curatorFilter)
}

type topSellingProduct struct {
	ProductID     string
	CuratorID     uint64
	TotalPrice    float64
	TotalQuantity uint
}

func getTopSellingProducts(db *gorm.DB, query *dto.CommonProductRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetCuratorTopSellingProductResponse, *dto.Paging, error) {
	var topSellingProducts []topSellingProduct
	var totalCount int64

	dbQuery := db.Table("order_variants").
		Select("product_id, curator_id, SUM(price * quantity) as total_price, SUM(quantity) as total_quantity").
		Group("product_id, curator_id")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToQuery(dbQuery)

	err := dbQuery.
		Count(&totalCount)
	if err.Error != nil {
		return nil, nil, err.Error
	}

	err = dbQuery.
		Order("total_quantity DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&topSellingProducts)
	if err.Error != nil {
		return nil, nil, err.Error
	}

	var resp []dto.GetCuratorTopSellingProductResponse
	for _, v := range topSellingProducts {
		var product domain.Product
		var variant domain.Variant

		productResult := db.Where("product_id = ?", v.ProductID).First(&product)
		if productResult.Error != nil {
			totalCount--
			log.Warnf("unable to get product with id %v: %v", v.ProductID, err)
			continue
		}

		variantResult := db.Where("product_id = ?", v.ProductID).Order("id").First(&variant)
		if variantResult.Error != nil {
			totalCount--
			log.Warnf("unable to get product with id %v: %v", v.ProductID, err)
			continue
		}

		resp = append(resp, dto.GetCuratorTopSellingProductResponse{
			ProductID:    v.ProductID,
			CuratorID:    v.CuratorID,
			ProductName:  product.Name,
			ProductImage: variant.Image,
			ProductBrand: product.BrandName,
			ProductPrice: variant.RetailPrice,
			Earnings:     fmt.Sprintf("%.2f", v.TotalPrice),

			Sales: v.TotalQuantity,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

type topSellingBrands struct {
	BrandName     string
	CuratorID     uint64
	TotalPrice    float64
	TotalQuantity uint
}

func (r *curatorDashboardRepository) GetCuratorTopSellingBrands(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopSellingBrandsResponse, *dto.Paging, error) {
	curatorFilter := func(db *gorm.DB) *gorm.DB {
		return db.Where("curator_id = ?", curatorID)
	}

	return getTopSellingBrands(r.db, query, curatorFilter)
}

func getTopSellingBrands(db *gorm.DB, query *dto.CommonProductRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.GetCuratorTopSellingBrandsResponse, *dto.Paging, error) {
	var topSellingBrands []topSellingBrands
	var totalCount int64

	dbQuery := db.Table("order_variants").
		Select("brand_name, curator_id, SUM(price * quantity) as total_price, SUM(quantity) as total_quantity").
		Where("brand_name != ''").
		Group("curator_id, brand_name")

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
		Scan(&topSellingBrands).Error
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.GetCuratorTopSellingBrandsResponse
	for _, v := range topSellingBrands {
		var brandImage string
		brandResult := db.Table("products").Where("LOWER(brand_name) = ?", strings.ToLower(v.BrandName)).Pluck("brand_image", &brandImage)
		if brandResult.Error != nil {
			log.Warnf("unable to get brand image with brand name %v: %v", v.BrandName, brandResult.Error)
			continue
		}

		resp = append(resp, dto.GetCuratorTopSellingBrandsResponse{
			BrandName:  v.BrandName,
			CuratorID:  v.CuratorID,
			BrandImage: brandImage,
			Earnings:   fmt.Sprintf("%.2f", v.TotalPrice),
			Sales:      v.TotalQuantity,
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

type topPurchasers struct {
	FirstName      string
	LastName       string
	UserID         uint
	TotalPurchases uint
	TotalSpent     float64
}

func (r *curatorDashboardRepository) GetCuratorTopPurchasers(curatorID uint, query *dto.CommonProductRequest) ([]dto.GetCuratorTopPurchasersResponse, *dto.Paging, error) {
	var topPurchasers []topPurchasers
	var totalCount int64

	dbQuery := r.db.Table("order_variants").
		Select("users.first_name, users.last_name, users.id as user_id, SUM(order_variants.quantity) as total_purchases, SUM(order_variants.price * order_variants.quantity) as total_spent").
		Joins("JOIN orders ON orders.id = order_variants.order_id").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("order_variants.curator_id = ?", curatorID).
		Group("users.id, users.first_name, users.last_name")

	dbQuery = query.TimeFilter.ApplyTimeFilterToSpecificTableQuery(dbQuery, "order_variants")

	err := dbQuery.Count(&totalCount)
	if err.Error != nil {
		return nil, nil, err.Error
	}

	err = dbQuery.
		Order("total_spent DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&topPurchasers)
	if err.Error != nil {
		return nil, nil, err.Error
	}

	var resp []dto.GetCuratorTopPurchasersResponse
	for _, purchaser := range topPurchasers {
		resp = append(resp, dto.GetCuratorTopPurchasersResponse{
			FirstName:      purchaser.FirstName,
			LastName:       purchaser.LastName,
			UserID:         purchaser.UserID,
			CuratorID:      curatorID,
			TotalPurchases: purchaser.TotalPurchases,
			TotalSpent:     fmt.Sprintf("%.2f", purchaser.TotalSpent),
		})
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return resp, paging, nil
}

func (repo *curatorDashboardRepository) GetSales(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.SalesIntervalResult,
	error) {

	var results []dto.SalesIntervalResult

	query := fmt.Sprintf(`
	WITH date_intervals AS (
		SELECT generate_series(
			?::timestamp, 
			?::timestamp, 
			'1 %s'::interval
		) AS interval_start
	),
	unique_orders AS (
		SELECT
			DISTINCT ON (orders.id) orders.id,
			orders.total_amount,
			orders.created_at,
			order_variants.curator_id
		FROM
			orders
		JOIN
			order_variants
		ON
			orders.id = order_variants.order_id
	)
	SELECT
		interval_start,
		SUM(unique_orders.total_amount) AS total_amount_sum
	FROM
		date_intervals
	LEFT JOIN
		unique_orders
	ON
		unique_orders.created_at >= interval_start
		AND unique_orders.created_at < interval_start + '1 %s'::interval
		AND unique_orders.curator_id = ?
	GROUP BY
		interval_start
	ORDER BY
		interval_start;
	
	
    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *curatorDashboardRepository) GetRevenue(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.RevenueIntervalResult,
	error) {

	var results []dto.RevenueIntervalResult

	query := fmt.Sprintf(`
	WITH date_intervals AS (
		SELECT generate_series(
			?::timestamp, 
			?::timestamp, 
			'1 %s'::interval
		) AS interval_start
	),
	unique_orders AS (
		SELECT
			DISTINCT ON (orders.id) orders.id,
			orders.total_amount,
			orders.created_at,
			order_variants.curator_id
		FROM
			orders
		JOIN
			order_variants
		ON
			orders.id = order_variants.order_id
	)
	SELECT
		interval_start,
		SUM(unique_orders.total_amount) AS total_revenue
	FROM
		date_intervals
	LEFT JOIN
		unique_orders
	ON
		unique_orders.created_at >= interval_start
		AND unique_orders.created_at < interval_start + '1 %s'::interval
		AND unique_orders.curator_id = ?
	GROUP BY
		interval_start
	ORDER BY
		interval_start;
	
	
    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *curatorDashboardRepository) GetAverageOrderValue(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.AOVIntervalResultResponse,
	error) {

	var results []dto.AOVIntervalResultResponse

	query := fmt.Sprintf(`
        WITH date_intervals AS (
            SELECT generate_series(
                ?::timestamp, 
                ?::timestamp, 
                '1 %s'::interval
            ) AS interval_start
        )
        SELECT
    interval_start,
    COALESCE(SUM(price * quantity), 0) AS total_order_value,
    COUNT(DISTINCT order_id) AS total_number_of_orders
FROM
    date_intervals
LEFT JOIN
    order_variants
ON
    order_variants.created_at >= interval_start
    AND order_variants.created_at < interval_start + '1 %s'::interval
    AND order_variants.curator_id = ?
    and order_variants.cancelled = false 
GROUP BY
    interval_start
ORDER BY
    interval_start
    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *curatorDashboardRepository) GetAverageUnitsPerOrder(
	fromTimestamp, toTimestamp, interval string,
	curatorID uint) (
	[]dto.UnitsSoldPerOrder,
	error) {

	var results []dto.UnitsSoldPerOrder

	query := fmt.Sprintf(`
        WITH date_intervals AS (
    SELECT generate_series(
        ?::timestamp, 
        ?::timestamp, 
        '1 %s'::interval
    ) AS interval_start
),
order_totals AS (
    SELECT 
        order_id,
        SUM(quantity) AS total_units,
        MIN(created_at) AS created_at
    FROM
        order_variants
    WHERE
        curator_id = ?
    GROUP BY
        order_id
)
SELECT
    interval_start,
    COUNT(order_totals.order_id) AS total_orders,
    COALESCE(SUM(order_totals.total_units), 0) AS total_units_sold
FROM
    date_intervals
LEFT JOIN
    order_totals
ON
    order_totals.created_at >= interval_start
    AND order_totals.created_at < interval_start + '1 %s'::interval
GROUP BY
    interval_start
ORDER BY
    interval_start;

    `, interval, interval)

	if err := repo.db.Raw(query, fromTimestamp, toTimestamp, curatorID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

type saleByCategory struct {
	CategoryID    uint
	CategoryName  string
	TotalQuantity uint
}

func (r *curatorDashboardRepository) GetCuratorSalesByCategory(curatorID uint, query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error) {
	curatorFilter := func(db *gorm.DB) *gorm.DB {
		return db.Where("order_variants.curator_id = ?", curatorID)
	}

	return getSalesByCategory(r.db, query, curatorFilter)
}

func getSalesByCategory(db *gorm.DB, query *dto.SaleRequest, filterFunc func(*gorm.DB) *gorm.DB) ([]dto.SaleByCategoryResponse, error) {
	var salesByCategory []saleByCategory
	var totalSalesQuantity uint

	dbQuery := db.Table("order_variants").
		Select("categories.id as category_id, categories.name as category_name, SUM(order_variants.quantity) as total_quantity").
		Joins("JOIN products ON products.product_id = order_variants.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Group("categories.id, categories.name")

	dbQuery = filterFunc(dbQuery)
	dbQuery = query.TimeFilter.ApplyTimeFilterToSpecificTableQuery(dbQuery, "order_variants")

	if err := dbQuery.Scan(&salesByCategory).Error; err != nil {
		return nil, err
	}

	for _, sale := range salesByCategory {
		totalSalesQuantity += sale.TotalQuantity
	}

	var response []dto.SaleByCategoryResponse
	for _, sale := range salesByCategory {
		percentage := (float64(sale.TotalQuantity) / float64(totalSalesQuantity)) * 100
		response = append(response, dto.SaleByCategoryResponse{
			CategoryName: sale.CategoryName,
			Percentage:   percentage,
		})
	}

	return response, nil
}
