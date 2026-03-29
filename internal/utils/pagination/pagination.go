package pagination

import (
	"gorm.io/gorm"
)

func NewPaginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}

// The Paginate function is not recommended for use anymore. It is only used in older functions.
// Please use NewPaginate() instead of this.
func Paginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}

func SetProductSort(query *gorm.DB, sortBy string) *gorm.DB {
	switch sortBy {
	case "new_in":
		return query.Order("products.created_at DESC")
	case "price_low_to_high":
		return query.
			Select("products.*").
			Joins("LEFT JOIN variants ON variants.product_id = products.product_id").
			Group("products.product_id").
			Order("MIN(variants.retail_price) ASC NULLS LAST")
	case "price_high_to_low":
		return query.
			Select("products.*").
			Joins("LEFT JOIN variants ON variants.product_id = products.product_id").
			Group("products.product_id").
			Order("MIN(variants.retail_price) DESC NULLS LAST")
	default:
		return query.Order("products.created_at DESC")
	}
}

func SetOrderSalesSort(query *gorm.DB, sortBy string) *gorm.DB {
	switch sortBy {
	case "customer_name_asc":
		query = query.Joins("JOIN users ON users.id = orders.user_id").Order("CONCAT(users.first_name, ' ', users.last_name) ASC")
	case "customer_name_desc":
		query = query.Joins("JOIN users ON users.id = orders.user_id").Order("CONCAT(users.first_name, ' ', users.last_name) DESC")
	case "date_asc":
		query = query.Order("orders.created_at ASC")
	case "date_desc":
		query = query.Order("orders.created_at DESC")
	case "items_low_to_high":
		query = query.Order("orders.total_quantity ASC")
	case "items_high_to_low":
		query = query.Order("orders.total_quantity DESC")
	case "amount_low_to_high":
		query = query.Order("orders.total_amount ASC")
	case "amount_high_to_low":
		query = query.Order("orders.total_amount DESC")
	}

	return query
}

func SetCatalogProductsSort(query *gorm.DB, sortBy string) *gorm.DB {
	switch sortBy {
	case "product_name_asc":
		query = query.Order("products.name ASC")
	case "product_name_desc":
		query = query.Order("products.name DESC")
	case "brand_name_asc":
		query = query.Order("products.brand_name ASC")
	case "brand_name_desc":
		query = query.Order("products.brand_name DESC")
	case "discount_low_to_high":
		query = query.Order("created_at ASC") // TODO
	case "discount_high_to_low":
		query = query.Order("created_at DESC") // TODO
	case "price_low_to_high":
		query = query.Order("created_at ASC") // TODO
	case "price_high_to_low":
		query = query.Order("created_at DESC") // TODO
	}

	return query
}
