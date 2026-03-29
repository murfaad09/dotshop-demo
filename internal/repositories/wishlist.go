package repository

import (
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"

	"github.com/harishash/dotshop-be/internal/dto"
)

type WishlistRepository interface {
	AddToWishlist(reqBody *dto.AddWishlistItemRequest, userId uint) (*dto.AddWishlistItemResponse, error)
	RemoveFromWishlist(userID, curatorID uint, productID, variantID string) error
	GetWishlist(userID uint) ([]*domain.WishlistItem, error)
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db}
}

func (r *wishlistRepository) AddToWishlist(reqBody *dto.AddWishlistItemRequest, userId uint) (*dto.AddWishlistItemResponse, error) {
	wishlistItem := &domain.WishlistItem{
		UserID:             userId,
		CuratorID:          reqBody.CuratorID,
		ProductID:          reqBody.ProductID,
		ProductName:        reqBody.ProductName,
		ProductImage:       reqBody.ProductImage,
		ProductBrand:       reqBody.ProductBrand,
		ProductPrice:       reqBody.ProductPrice,
		VariantID:          reqBody.VariantID,
		VariantOptionName:  reqBody.VariantOptionName,
		VariantOptionValue: reqBody.VariantOptionValue,
	}

	if r.getWishlistItemsFromDB(userId, reqBody.CuratorID, reqBody.ProductID, reqBody.VariantID) {
		res := &dto.AddWishlistItemResponse{
			Success: false,
			Message: "Item already exists in wishlist",
		}
		return res, nil
	}

	if err := r.db.Create(wishlistItem).Error; err != nil {
		return nil, err
	}
	res := &dto.AddWishlistItemResponse{
		Success: true,
		Message: "Item added to wishlist successfully",
	}
	return res, nil
}

func (r *wishlistRepository) RemoveFromWishlist(userID, curatorID uint, productID, variantID string) error {
	query := r.db.Where("user_id = ? AND product_id = ? AND curator_id = ?", userID, productID, curatorID)
	if len(variantID) > 0 {
		query = query.Where("variant_id = ?", variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}
	if err := query.Delete(&domain.WishlistItem{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *wishlistRepository) GetWishlist(userID uint) ([]*domain.WishlistItem, error) {
	var wishlistItems []*domain.WishlistItem
	if err := r.db.Joins("JOIN products ON products.product_id = wishlist_items.product_id").Order("created_at DESC").Where("user_id = ? AND products.is_active = ?", userID, true).Find(&wishlistItems).Error; err != nil {
		return nil, err
	}
	return wishlistItems, nil
}

func (r *wishlistRepository) getWishlistItemsFromDB(userID, curatorID uint, productID string, variantID *string) bool {
	query := r.db.Where("user_id = ? AND product_id = ? AND curator_id = ?", userID, productID, curatorID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}
	if err := query.First(&domain.WishlistItem{}).Error; err != nil {
		return false
	}
	return true
}
