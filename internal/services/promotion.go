package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	repository "github.com/harishash/dotshop-be/internal/repositories"
)

type PromotionService interface {
	CreatePromotion(promotion *dto.PromotionRequest) (*dto.PromotionResponse, error)
	UpdatePromotion(promotion *dto.PromotionRequest, id uint) error
	GetPromotions(query *dto.ListPromotionsRequest) (*dto.Response, error)
	GetPromotionByID(id uint) (*dto.PromotionResponse, error)
	DeletePromotion(id uint) error
	ApplyDiscountToProducts(productIDs []string, promotionID uint) error
}

type promotionService struct {
	repo repository.PromotionRepository
}

func NewPromotionService(repo repository.PromotionRepository) PromotionService {
	return &promotionService{repo}
}

func (s *promotionService) CreatePromotion(promotion *dto.PromotionRequest) (*dto.PromotionResponse, error) {
	promo := &models.Promotion{
		Name:          promotion.Name,
		DiscountCode:  promotion.DiscountCode,
		DiscountValue: promotion.DiscountValue,
		ExpiryDate:    promotion.ExpiryDate,
		Status:        promotion.Status,
		Rule:          promotion.Rule,
	}

	createdPromo, err := s.repo.CreatePromotion(promo)
	if err != nil {
		return nil, err
	}

	response := &dto.PromotionResponse{
		ID:            createdPromo.ID,
		Name:          createdPromo.Name,
		DiscountCode:  createdPromo.DiscountCode,
		DiscountValue: createdPromo.DiscountValue,
		ExpiryDate:    createdPromo.ExpiryDate.String(),
		Status:        createdPromo.Status,
		Rule:          createdPromo.Rule,
	}

	return response, nil
}

func (s *promotionService) UpdatePromotion(promotion *dto.PromotionRequest, id uint) error {
	promo := &models.Promotion{
		ID:            id,
		Name:          promotion.Name,
		DiscountCode:  promotion.DiscountCode,
		DiscountValue: promotion.DiscountValue,
		ExpiryDate:    promotion.ExpiryDate,
		Status:        promotion.Status,
		Rule:          promotion.Rule,
	}

	return s.repo.UpdatePromotion(promo)
}

func (s *promotionService) GetPromotions(query *dto.ListPromotionsRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetPromotions(query)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *promotionService) GetPromotionByID(id uint) (*dto.PromotionResponse, error) {
	promotion, err := s.repo.GetPromotionByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.PromotionResponse{
		ID:            promotion.ID,
		Name:          promotion.Name,
		DiscountCode:  promotion.DiscountCode,
		DiscountValue: promotion.DiscountValue,
		ExpiryDate:    promotion.ExpiryDate.String(),
		Status:        promotion.Status,
		Rule:          promotion.Rule,
	}

	return response, nil
}

func (s *promotionService) DeletePromotion(id uint) error {
	return s.repo.DeletePromotion(id)
}

func (s *promotionService) ApplyDiscountToProducts(productIDs []string, promotionID uint) error {
	return s.repo.ApplyDiscountToProducts(productIDs, promotionID)
}
