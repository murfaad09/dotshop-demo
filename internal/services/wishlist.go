package service

import (
	repository "github.com/harishash/dotshop-be/internal/repositories"

	"github.com/harishash/dotshop-be/internal/dto"
)

type WishlistService interface {
	AddToWishlist(req *dto.AddWishlistItemRequest, userId uint) (*dto.AddWishlistItemResponse, error)
	RemoveFromWishlist(userID, curatorID uint, productID, variantID string) error
	GetWishlist(userID uint) (*dto.GetWishlistResponse, error)
}

type wishlistService struct {
	repo repository.WishlistRepository
}

func NewWishlistService(repo repository.WishlistRepository) WishlistService {
	return &wishlistService{repo}
}

func (s *wishlistService) AddToWishlist(req *dto.AddWishlistItemRequest, userId uint) (*dto.AddWishlistItemResponse, error) {
	wishlistItem, err := s.repo.AddToWishlist(req, userId)
	if err != nil {

		return nil, err
	}

	return wishlistItem, nil
}

func (s *wishlistService) RemoveFromWishlist(userID, curatorID uint, productID, variantID string) error {
	return s.repo.RemoveFromWishlist(userID, curatorID, productID, variantID)
}

func (s *wishlistService) GetWishlist(userID uint) (*dto.GetWishlistResponse, error) {
	data, err := s.repo.GetWishlist(userID)
	if err != nil {
		return nil, err
	}

	// Convert []*model.WishlistItem to []*dto.AddWishlistItemRequest
	var convertedData []*dto.AddWishlistItemRequest
	for _, item := range data {
		convertedItem := &dto.AddWishlistItemRequest{
			CuratorID:          item.CuratorID,
			ProductID:          item.ProductID,
			ProductName:        item.ProductName,
			ProductImage:       item.ProductImage,
			ProductBrand:       item.ProductBrand,
			ProductPrice:       item.ProductPrice,
			VariantID:          item.VariantID,
			VariantOptionName:  item.VariantOptionName,
			VariantOptionValue: item.VariantOptionValue,
		}
		convertedData = append(convertedData, convertedItem)
	}

	res := &dto.GetWishlistResponse{
		Success: true,
		Data:    convertedData,
	}
	return res, nil
}
