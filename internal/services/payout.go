package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	payout_repo "github.com/harishash/dotshop-be/internal/repositories"
)

type IPayoutService interface {
	GetPayoutHistory(userId uint64) ([]dto.PayoutHistoryResponse, error)
	GetPayoutDetails(userId uint64) (*dto.PayoutResponse, error)
}

type payoutService struct {
	repo payout_repo.PayoutRepository
}

func NewPayoutService(repo payout_repo.PayoutRepository) IPayoutService {
	return &payoutService{repo}
}

func (s *payoutService) GetPayoutHistory(userId uint64) ([]dto.PayoutHistoryResponse, error) {
	curatorId, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	payoutHistory, err := s.repo.GetPayoutHistoryWithCuratorId(uint64(curatorId.ID))
	if err != nil {
		return nil, err
	}

	var resp []dto.PayoutHistoryResponse
	for _, v := range payoutHistory {
		resp = append(resp, dto.PayoutHistoryResponse{
			Id:               v.ID,
			CuratorId:        v.CuratorID,
			ReturnAmount:     v.ReturnAmount,
			CommissionAmount: v.CommissionAmount,
			PayoutAmount:     v.PayoutAmount,
			CreatedAt:        v.CreatedAt,
			UpdatedAt:        v.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *payoutService) GetPayoutDetails(userId uint64) (*dto.PayoutResponse, error) {
	curatorId, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	payoutHistory, err := s.repo.GetPayoutDetailsWithCuratorId(uint64(curatorId.ID))
	if err != nil {
		return nil, err
	}

	resp := &dto.PayoutResponse{
		ID:               payoutHistory.ID,
		CuratorID:        payoutHistory.CuratorID,
		PayoutAmount:     payoutHistory.PayoutAmount,
		CommissionAmount: payoutHistory.CommissionAmount,
		ReturnAmount:     payoutHistory.ReturnAmount,
		CreatedAt:        payoutHistory.CreatedAt,
		UpdatedAt:        payoutHistory.UpdatedAt,
	}

	return resp, nil
}
