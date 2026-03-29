package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/harishash/dotshop-be/integration/stripe"
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	cart_repo "github.com/harishash/dotshop-be/internal/repositories"
)

// CartsService interface
type CartsService interface {
	BuyNow(userID string) error
	CreateCart(carts *dto.CartRequest) (*dto.CartResponse, error)
	UpdateCartItemQuantity(cartID uint64, variantID string, newQuantity uint) (*dto.CartItemResponse, error)
	AddItemsToCart(cartID uint, cartItems []*dto.AddCartItemsRequest) (*dto.CartResponse, error)
	DeleteCart(cartID uint) error
	DeleteCartItem(cartID uint, variantID string) error
	GetCartByUserID(userID uint, currency, postalCode, country string) (*dto.CartResponse, error)
}

// cartsService struct
type cartsService struct {
	repo cart_repo.CartsRepo
}

// NewCartsService creates a new CartsService
func NewCartsService(repo cart_repo.CartsRepo) CartsService {
	return &cartsService{repo: repo}
}
func (s *cartsService) BuyNow(userID string) error {
	return s.repo.BuyNow(userID)
}

func (s *cartsService) CreateCart(carts *dto.CartRequest) (*dto.CartResponse, error) {
	return s.repo.CreateCart(carts)
}

func (s *cartsService) UpdateCartItemQuantity(cartID uint64, variantID string, newQuantity uint) (*dto.CartItemResponse, error) {
	return s.repo.UpdateCartItemQuantity(cartID, variantID, newQuantity)
}

func (s *cartsService) AddItemsToCart(cartID uint, cartItems []*dto.AddCartItemsRequest) (*dto.CartResponse, error) {
	// Check if the cart exists
	cart, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("cart does not exist")
	}

	// Create a map of existing cart items for quick lookup
	existingItemsMap := make(map[string]*models.CartItem)
	for _, item := range cart.Items {
		existingItemsMap[item.VariantID] = item
	}

	// Process incoming cart items
	for _, newItem := range cartItems {
		// Fetch available inventory for the product variant
		availableInventory, err := s.repo.GetInventoryByVariantID(newItem.VariantID)
		if err != nil {
			return nil, err
		}

		// Check if the requested quantity is available
		if newItem.Quantity > availableInventory {
			return nil, fmt.Errorf("requested quantity for product %s exceeds available inventory", newItem.ProductName)
		}

		if existingItem, found := existingItemsMap[newItem.VariantID]; found {
			// Check if the updated quantity exceeds the available inventory
			if existingItem.Quantity+newItem.Quantity > availableInventory {
				return nil, fmt.Errorf("total quantity for product %s exceeds available inventory", newItem.ProductName)
			}
			// Update the quantity of existing item
			existingItem.Quantity += newItem.Quantity
		} else {
			// Add new item to the cart
			cart.Items = append(cart.Items, &models.CartItem{
				CuratorID:   newItem.CuratorID,
				CartID:      cartID,
				ProductID:   newItem.ProductID,
				Price:       newItem.Price,
				ProductName: newItem.ProductName,
				BrandName:   newItem.BrandName,
				VariantID:   newItem.VariantID,
				ImageURL:    newItem.ImageURL,
				Color:       newItem.Color,
				Size:        newItem.Size,
				Quantity:    newItem.Quantity,
			})
		}
	}

	// Persist changes to the cart
	err = s.repo.UpdateCartItems(cartID, cart.Items)
	if err != nil {
		return nil, err
	}

	// Prepare the response
	responseItems := make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		responseItems[i] = dto.CartItemResponse{
			CartID:      item.CartID,
			CuratorID:   uint(item.CuratorID),
			ProductID:   item.ProductID,
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

func (s *cartsService) DeleteCart(cartID uint) error {
	// Check if the cart exists
	cart, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart does not exist")
	}

	// Mark the cart as deleted by setting the deleted_at timestamp
	err = s.repo.DeleteCart(cartID)
	if err != nil {
		return err
	}

	return nil
}

func (s *cartsService) GetCartByUserID(userID uint, currency, postalCode, country string) (*dto.CartResponse, error) {

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("active cart does not exist for this user")
	}

	var totalAmount float64

	responseItems := make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		responseItems[i] = dto.CartItemResponse{
			CuratorID:   uint(item.CuratorID),
			CartID:      item.CartID,
			ProductID:   item.ProductID,
			Price:       item.Price,
			ProductName: item.ProductName,
			BrandName:   item.BrandName,
			VariantID:   item.VariantID,
			ImageURL:    item.ImageURL,
			Color:       item.Color,
			Size:        item.Size,
			Quantity:    item.Quantity,
		}

		totalAmount += item.Price * float64(item.Quantity)
	}

	paramsForTax := stripe.TaxReq{
		Amount:     int64(totalAmount),
		Currency:   currency,
		Country:    country,
		PostalCode: postalCode,
	}
	taxPtr, err := stripe.PerformTaxCalculation(paramsForTax)
	if err != nil {
		log.Printf("Error performing tax calculation: %v", err)
	}

	var tax float64

	if taxPtr == nil {
		tax = float64(0)
	} else {
		tax = float64(*taxPtr)
	}

	getCartResponse := &dto.CartResponse{
		CartID:      cart.ID,
		UserID:      cart.UserID,
		CartItem:    responseItems,
		SubTotal:    totalAmount,
		Tax:         tax,
		TotalAmount: totalAmount + tax,
	}

	return getCartResponse, nil
}

func (s *cartsService) DeleteCartItem(cartID uint, variantID string) error {
	// Check if the cart exists and contains the item
	cart, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart does not exist")
	}

	var itemExists bool
	for _, item := range cart.Items {
		if item.VariantID == variantID {
			itemExists = true
			break
		}
	}
	if !itemExists {
		return errors.New("cart item does not exist in the cart")
	}

	// Perform the deletion
	if err := s.repo.DeleteCartItem(cartID, variantID); err != nil {
		return err
	}

	// Refresh the cart to check if it's empty
	cart, err = s.repo.GetCartByID(cartID)
	if err != nil {
		return err
	}
	if cart == nil || len(cart.Items) == 0 {
		if err := s.repo.DeleteCart(cartID); err != nil {
			return err
		}
	}

	return nil
}
