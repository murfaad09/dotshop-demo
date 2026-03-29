package repository

import (
	"gorm.io/gorm"

	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	"github.com/harishash/dotshop-be/internal/utils/pagination"
)

type PromotionRepository interface {
	CreatePromotion(promotion *models.Promotion) (*models.Promotion, error)
	UpdatePromotion(promotion *models.Promotion) error
	GetPromotions(params *dto.ListPromotionsRequest) ([]*dto.ListPromotionResponse, *dto.Paging, error)
	GetPromotionByID(id uint) (*models.Promotion, error)
	DeletePromotion(id uint) error
	ApplyDiscountToProducts(productIDs []string, promotionID uint) error
}

type promotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) PromotionRepository {
	return &promotionRepository{db}
}

func (r *promotionRepository) CreatePromotion(promotion *models.Promotion) (*models.Promotion, error) {
	err := r.db.Create(promotion).Error
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *promotionRepository) UpdatePromotion(promotion *models.Promotion) error {
	return r.db.Model(&models.Promotion{}).Where("id = ?", promotion.ID).Updates(promotion).Error
}

func (r *promotionRepository) GetPromotions(params *dto.ListPromotionsRequest) ([]*dto.ListPromotionResponse, *dto.Paging, error) {
	var totalCount int64
	var promotions []*dto.ListPromotionResponse

	query := r.db.Model(&models.Promotion{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.StartValue > 0 {
		query = query.Where("discount_value >= ?", params.StartValue)
	}
	if params.EndValue > 0 {
		query = query.Where("discount_value <= ?", params.EndValue)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	pagination.NewPaginate(query, params.PageNum, params.PageSize)
	if err := query.Find(&promotions).Error; err != nil {
		return nil, nil, err
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		CurrentPage: params.PageNum,
		PageSize:    params.PageSize,
	}

	return promotions, paging, nil
}

func (r *promotionRepository) GetPromotionByID(id uint) (*models.Promotion, error) {
	var promotion models.Promotion
	err := r.db.First(&promotion, id).Error
	if err != nil {
		return nil, err
	}

	return &promotion, nil
}

func (r *promotionRepository) DeletePromotion(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.Promotion{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *promotionRepository) ApplyDiscountToProducts(productIDs []string, promotionID uint) error {
	return r.db.Model(&models.Product{}).Where("product_id IN ?", productIDs).Update("promotion_id", promotionID).Error
}
