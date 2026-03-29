package repository

import (
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

type PayoutRepository interface {
	GetCuratorWithUserId(id uint64) (*domain.Curator, error)
	GetPayoutHistoryWithCuratorId(id uint64) ([]domain.PayoutHistory, error)
	GetPayoutDetailsWithCuratorId(id uint64) (*domain.Payout, error)
}

type payoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository(db *gorm.DB) PayoutRepository {
	return &payoutRepository{db}
}

func (r *payoutRepository) GetCuratorWithUserId(id uint64) (*domain.Curator, error) {
	curators := domain.Curator{}
	result := r.db.Where("user_id = ?", id).Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curators, nil
}

func (r *payoutRepository) GetPayoutHistoryWithCuratorId(id uint64) ([]domain.PayoutHistory, error) {
	var payoutHistory []domain.PayoutHistory
	result := r.db.Where("curator_id = ?", id).Find(&payoutHistory)
	if result.Error != nil {
		return nil, result.Error
	}

	return payoutHistory, nil
}

func (r *payoutRepository) GetPayoutDetailsWithCuratorId(id uint64) (*domain.Payout, error) {
	payout := domain.Payout{}
	if err := r.db.Where("curator_id = ?", id).First(&payout).Error; err != nil {
		return nil, err
	}

	return &payout, nil
}
