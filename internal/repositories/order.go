package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	"github.com/harishash/dotshop-be/internal/utils/pagination"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	GetVariant(id string) (*domain.Variant, error)
	CreateOrder(orderDomain *domain.Order, orderVarients []*domain.OrderVariants, fulfillments []domain.Fulfillments) (*domain.Order, error)
	GetOrderIDsByCuratorID(id uint64) ([]string, error)
	GetOrderVariantByUserID(id uint64) ([]domain.OrderVariants, error)

	GetOrderById(id string) (*domain.Order, error)
	GetVariantsByVariantID(id string) ([]domain.Variant, error)
	UpdateCancelOrderVariant(tx *gorm.DB, orderVariant *domain.OrderVariants) error
	GetOrderVariantsByOrderId(id string, curatorId *uint) ([]*domain.OrderVariants, error)
	UpdateCancelOrder(tx *gorm.DB, order *domain.Order) error
	GetCuratorWithUserId(id uint64) (*domain.Curator, error)
	GetOrderVariant(orderId string, variantId string, userId uint64) (*domain.OrderVariants, error)
	ReturnOrder(order *domain.Order, orderVariants []*domain.OrderVariants, returnVariants []*domain.ReturnOrder) ([]*domain.ReturnOrder, error)
	StartTx() (*gorm.DB, error)
	GetOrdersByUserId(userId uint64) ([]domain.Order, error)
	GetOrderVariantByOrderID(orderId string) ([]domain.OrderVariants, error)
	GetOrdersByUserIdForCustomer(userId uint64, query *dto.CustomersOrderListRequest) ([]domain.Order, *dto.Paging, error)
	GetCuratorOrderIDsByID(curatorId uint, query *dto.CuratorsOrderListRequest) ([]string, *dto.Paging, error)
}

type orderRepo struct {
	db *gorm.DB
	NotificationRepository
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) StartTx() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	return tx, nil
}

func (r *orderRepo) GetVariant(id string) (*domain.Variant, error) {
	var row domain.Variant

	tx := r.db.Where("id = ?", id).First(&row)

	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("find variant by ID '%v': %v", id, err)
	}

	return &row, nil
}

func (r *orderRepo) CreateOrder(orderDomain *domain.Order, orderVarients []*domain.OrderVariants, fulfillments []domain.Fulfillments) (*domain.Order, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	if err := tx.Create(&orderDomain).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("unable to create order: %v", tx.Error)
	}

	if err := tx.Create(&orderVarients).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("unable to create order variant: %v", tx.Error)
	}

	if len(fulfillments) != 0 {
		if err := tx.Create(&fulfillments).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("unable to create order fulfillment: %v", tx.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %v", tx.Error)
	}
	return orderDomain, nil
}

func (r *orderRepo) GetOrderIDsByCuratorID(id uint64) ([]string, error) {
	var orderIDs []string

	tx := r.db.Model(&domain.OrderVariants{}).
		Distinct("order_id").
		Where("curator_id = ?", id).
		Pluck("order_id", &orderIDs)

	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get curator orders list by ID '%v': %v", id, err)
	}

	return orderIDs, nil
}

func (r *orderRepo) GetOrderVariantByUserID(userID uint64) ([]domain.OrderVariants, error) {
	var orderVariants []domain.OrderVariants

	// Use a join to get order variants based on the user ID from the orders table
	tx := r.db.Joins("JOIN orders ON orders.id = order_variants.order_id").
		Where("orders.user_id = ?", userID).
		Find(&orderVariants)

	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, fmt.Errorf("unable to get order variants for user ID '%v': %v", userID, err)
	}

	return orderVariants, nil
}

func (r *orderRepo) GetOrderById(id string) (*domain.Order, error) {
	var row domain.Order

	tx := r.db.Where("id = ?", id).Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get order by ID '%v': %v", id, err)
	}

	return &row, nil
}

func (r *orderRepo) GetVariantsByVariantID(id string) ([]domain.Variant, error) {
	var row []domain.Variant

	tx := r.db.Where("id = ?", id).
		Preload("VariantOptions").
		Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get order variant by ID '%v': %v", id, err)
	}

	return row, nil
}

func (r *orderRepo) UpdateCancelOrderVariant(tx *gorm.DB, orderVariant *domain.OrderVariants) error {
	if err := tx.Model(&domain.OrderVariants{}).Where("id = ?", orderVariant.ID).Updates(map[string]interface{}{
		"cancelled":        true,
		"cancelled_reason": orderVariant.CancelledReason,
		"cancelled_data":   orderVariant.CancelledData,
	}).Error; err != nil {
		return fmt.Errorf("unable to update cancelledorder variant: %v", err.Error())
	}

	return nil
}

func (r *orderRepo) UpdateCancelOrder(tx *gorm.DB, order *domain.Order) error {
	if err := tx.Model(&domain.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":           order.Status,
		"cancelled":        true,
		"cancelled_reason": order.CancelledReason,
		"cancelled_data":   order.CancelledData,
	}).Error; err != nil {
		return fmt.Errorf("unable to update cancelled order: %v", err.Error())
	}
	return nil
}

func (r *orderRepo) GetOrderVariantsByOrderId(id string, curatorId *uint) ([]*domain.OrderVariants, error) {
	var row []*domain.OrderVariants

	query := r.db.Where("order_id = ?", id)
	if curatorId != nil {
		query = query.Where("curator_id = ?", *curatorId)
	}

	tx := query.Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get order by ID '%v': %v", id, err)
	}

	return row, nil
}

func (r *orderRepo) GetCuratorWithUserId(id uint64) (*domain.Curator, error) {
	curators := domain.Curator{}
	result := r.db.Where("user_id = ?", id).Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curators, nil
}

func GetOrders(
	r *gorm.DB,
	curator_id uint,
	startDate, endDate time.Time) (
	[]domain.Order,
	error) {
	var orders []domain.Order
	err := r.Select("id, total_amount, created_at").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("curator_id = ?", curator_id).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, err
}

func (r *orderRepo) GetOrderVariant(orderId string, variantId string, userId uint64) (*domain.OrderVariants, error) {
	var row domain.OrderVariants
	tx := r.db.Preload("Variant").
		Joins("JOIN orders ON orders.id = order_variants.order_id").
		Where("order_variants.order_id = ? AND order_variants.variant_id = ? AND orders.user_id = ?", orderId, variantId, userId).
		Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get order variant by variant Id '%v': %v", variantId, err)
	}

	return &row, nil
}

func (r *orderRepo) UpdateReturnOrder(tx *gorm.DB, order *domain.Order) error {
	if err := tx.Model(&domain.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"total_amount":   order.TotalAmount,
		"total_quantity": order.TotalQuantity,
	}).Error; err != nil {
		return fmt.Errorf("unable to update order: %v", err.Error())
	}
	return nil
}

func (r *orderRepo) UpdateReturnOrderVariants(tx *gorm.DB, orderID string, orderVariants []*domain.OrderVariants) error {
	for _, variant := range orderVariants {
		if err := tx.Model(&domain.OrderVariants{}).Where("order_id = ? AND variant_id = ?", orderID, variant.VariantID).Updates(map[string]interface{}{
			"quantity": variant.Quantity,
			"price":    variant.Price,
		}).Error; err != nil {
			return fmt.Errorf("unable to update order variant: %v", err.Error())
		}
	}
	return nil
}

func (r *orderRepo) CreateReturnVariants(tx *gorm.DB, returnVariants []*domain.ReturnOrder) error {
	if err := tx.Create(&returnVariants).Error; err != nil {
		return fmt.Errorf("unable to create return order: %v", err.Error())
	}
	return nil
}

func (r *orderRepo) ReturnOrder(order *domain.Order, orderVariants []*domain.OrderVariants, returnVariants []*domain.ReturnOrder) ([]*domain.ReturnOrder, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := r.UpdateReturnOrder(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := r.UpdateReturnOrderVariants(tx, order.ID, orderVariants); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := r.CreateReturnVariants(tx, returnVariants); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err.Error())
	}

	return returnVariants, nil
}

func (r *orderRepo) CreateReturn(tx *gorm.DB, returnOrder []*domain.ReturnOrder) ([]*domain.ReturnOrder, error) {
	if err := r.db.Create(&returnOrder).Error; err != nil {
		return nil, fmt.Errorf("unable to create order: %v", err.Error())
	}

	return returnOrder, nil
}

func GetOrderReturns(
	r *gorm.DB,
	curator_id uint,
	startDate, endDate time.Time) (
	[]domain.ReturnOrder,
	error) {
	var returnorders []domain.ReturnOrder
	err := r.Select("id, total_amount, created_at").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("curator_id = ?", curator_id).
		Find(&returnorders).Error

	if err != nil {
		return nil, err
	}
	return returnorders, err
}

func (r *orderRepo) GetOrdersByUserId(userId uint64) ([]domain.Order, error) {
	var row []domain.Order

	tx := r.db.Order("updated_at DESC").Where("user_id = ?", userId).Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}

		return nil, fmt.Errorf("unable to get order by user id '%v': %v", userId, err)
	}

	return row, nil
}

func (r *orderRepo) GetOrderVariantByOrderID(orderId string) ([]domain.OrderVariants, error) {
	var orderVariants []domain.OrderVariants

	tx := r.db.Where("order_id = ?", orderId).Find(&orderVariants)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, fmt.Errorf("unable to get order variants for order ID '%v': %v", orderId, err)
	}

	return orderVariants, nil
}

func (r *orderRepo) GetOrdersByUserIdForCustomer(userId uint64, query *dto.CustomersOrderListRequest) ([]domain.Order, *dto.Paging, error) {
	var totalCount int64
	var row []domain.Order

	dbQuery := r.db.Model(&domain.Order{}).Where("user_id = ?", userId)

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	dbQuery.Order("updated_at DESC")

	tx := pagination.NewPaginate(dbQuery, query.PageNum, query.PageSize).Find(&row)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, domain.ErrResourceNotFound
		}

		return nil, nil, fmt.Errorf("unable to get order by user id '%v': %v", userId, err)
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return row, paging, nil
}

func (r *orderRepo) GetCuratorOrderIDsByID(curatorId uint, query *dto.CuratorsOrderListRequest) ([]string, *dto.Paging, error) {
	var totalCount int64
	var orderIDs []string

	dbQuery := r.db.Model(&domain.OrderVariants{}).
		Distinct("order_id").
		Where("curator_id = ?", curatorId)

	err := dbQuery.Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	tx := pagination.NewPaginate(dbQuery, query.PageNum, query.PageSize).Pluck("order_id", &orderIDs)
	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, domain.ErrResourceNotFound
		}

		return nil, nil, fmt.Errorf("unable to get order by curator id '%v': %v", curatorId, err)
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		PageSize:    query.PageSize,
		CurrentPage: query.PageNum,
	}

	return orderIDs, paging, nil
}
