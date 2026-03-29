package repository

import (
	"errors"
	"time"

	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

// CartsRepo interface
type CartsRepo interface {
	BuyNow(userID string) error
	CreateCart(carts *dto.CartRequest) (*dto.CartResponse, error)
	UpdateCartItemQuantity(cartID uint64, variantID string, newQuantity uint) (*dto.CartItemResponse, error)
	AddCartItems(cartID uint, items []*models.CartItem) error
	GetCartByID(cartID uint) (*models.Cart, error)
	UpdateCartItems(cartID uint, items []*models.CartItem) error
	DeleteCart(cartID uint) error
	DeleteCartItem(cartID uint, variantID string) error
	GetCartByUserID(userID uint) (*models.Cart, error)
	GetInventoryByVariantID(variantID string) (uint, error)
}

// cartsRepo struct
type cartsRepo struct {
	db *gorm.DB
}

// NewCartsRepo creates a new CartsRepo
func NewCartsRepo() CartsRepo {
	instance := GetDatabaseConnection()
	return &cartsRepo{db: instance.Connection}
}

func (r *cartsRepo) BuyNow(userID string) error {
	// Your logic to handle the buy now action
	return nil
}

func (r *cartsRepo) CreateCart(carts *dto.CartRequest) (*dto.CartResponse, error) {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the user has an active cart or not
	var existingCart models.Cart
	err := tx.Where("user_id = ? AND deleted_at IS NULL", carts.UserID).First(&existingCart).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// If the user has an active cart, do not allow creating a new one
	if existingCart.ID != 0 {
		return nil, errors.New("user already has an active cart")
	}

	// Create Cart object
	cart := models.Cart{
		UserID: carts.UserID,
		Items:  []*models.CartItem{},
	}

	if err := tx.Create(&cart).Error; err != nil {
		return nil, err
	}

	// Populate CartItems from DTO
	for _, item := range carts.CartItem {
		cart.Items = append(cart.Items, &models.CartItem{
			CartID:      cart.ID, // Use the cart ID here
			CuratorID:   int(item.CuratorID),
			ProductID:   item.ProductID,
			Price:       item.Price,
			ProductName: item.ProductName,
			BrandName:   item.BrandName,
			VariantID:   item.VariantID,
			ImageURL:    item.ImageURL,
			Color:       item.Color,
			Size:        item.Size,
			Quantity:    item.Quantity,
		})
	}

	// Create the new cart items
	for _, item := range cart.Items {
		if err := tx.Create(item).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Prepare the response
	responseItems := make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		responseItems[i] = dto.CartItemResponse{
			CartID:      item.CartID,
			ProductID:   item.ProductID,
			CuratorID:   uint(item.CuratorID),
			Price:       item.Price,
			ProductName: item.ProductName,
			BrandName:   item.BrandName,
			VariantID:   item.VariantID,
			ImageURL:    item.ImageURL,
			Color:       item.Color,
			Size:        item.Size,
			Quantity:    item.Quantity,
		}
	}

	addToCartResponse := &dto.CartResponse{
		CartID:   cart.ID,
		UserID:   cart.UserID,
		CartItem: responseItems,
	}

	return addToCartResponse, nil
}

func (r *cartsRepo) UpdateCartItemQuantity(cartID uint64, variantID string, newQuantity uint) (*dto.CartItemResponse, error) {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find the cart item
	var cartItem models.CartItem
	err := tx.Where("cart_id = ? AND variant_id = ?", cartID, variantID).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, errors.New("cart item not found")
		}
		tx.Rollback()
		return nil, err
	}

	// Update the quantity
	cartItem.Quantity = newQuantity
	if err := tx.Save(&cartItem).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Prepare the response
	updatedCartItemResponse := &dto.CartItemResponse{
		CartID:      cartItem.CartID,
		CuratorID:   uint(cartItem.CuratorID),
		ProductID:   cartItem.ProductID,
		Price:       cartItem.Price,
		ProductName: cartItem.ProductName,
		BrandName:   cartItem.BrandName,
		VariantID:   cartItem.VariantID,
		ImageURL:    cartItem.ImageURL,
		Color:       cartItem.Color,
		Size:        cartItem.Size,
		Quantity:    cartItem.Quantity,
	}

	return updatedCartItemResponse, nil
}

func (r *cartsRepo) GetCartByID(cartID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("id = ? AND deleted_at IS NULL", cartID).Preload("Items").First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartsRepo) AddCartItems(cartID uint, items []*models.CartItem) error {
	for _, item := range items {
		item.CartID = cartID
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *cartsRepo) UpdateCartItems(cartID uint, items []*models.CartItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// If item ID is 0, it's a new item and should be created
			if item.ID == 0 {
				if err := tx.Create(item).Error; err != nil {
					return err
				}
			} else {
				// Otherwise, update the existing item
				if err := tx.Model(item).Where("id = ?", item.ID).Updates(map[string]interface{}{
					"quantity": item.Quantity,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// func (r *cartsRepo) UpdateCartItems(cartID uint, items []*models.CartItem) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		for _, item := range items {
// 			item.CartID = cartID  // Ensure the CartID is set for each item
// 			if err := tx.Create(item).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

func (r *cartsRepo) DeleteCart(cartID uint) error {
	return r.db.Model(&models.Cart{}).Where("id = ?", cartID).Update("deleted_at", time.Now()).Error
}

func (r *cartsRepo) DeleteCartItem(cartID uint, variantID string) error {
	return r.db.Where("cart_id = ? AND variant_id = ?", cartID, variantID).Delete(&models.CartItem{}).Error
}

func (r *cartsRepo) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN products ON products.product_id = cart_items.product_id").
			Where("products.is_active = ?", true)
	}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartsRepo) GetInventoryByVariantID(variantID string) (uint, error) {
	var inventory uint
	err := r.db.Table("variants").Where("id = ?", variantID).Select("inventory_amount").Scan(&inventory).Error
	if err != nil {
		return 0, err
	}
	return inventory, nil
}
