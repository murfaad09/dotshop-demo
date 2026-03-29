package service

import (
	"fmt"

	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"

	curatoronboarding_repo "github.com/harishash/dotshop-be/internal/repositories"
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

// CuratorOnboardingService interface
type CuratorOnboardingService interface {
	// CreateCuratorOnboarding
	CreateCurator(dto.CuratorOnBoardingRequest) (dto.CuratorOnBoardingResponse, *errors.Error)
	CheckShopName(shopName string) (*domain.Curator, error)
	GetCuratorByStoreName(storeName string) (*domain.Curator, error)
}

type curatorOnboardingService struct {
	repo curatoronboarding_repo.CuratorOnboardingRepo
}

func NewCuratorOnboardingService(repo curatoronboarding_repo.CuratorOnboardingRepo) CuratorOnboardingService {
	return &curatorOnboardingService{repo: repo}
}

func (s *curatorOnboardingService) CreateCurator(req dto.CuratorOnBoardingRequest) (dto.CuratorOnBoardingResponse, *errors.Error) {
	req.Password = hashPassword(req.Password)
	curator, err := s.repo.AddCurator(req)
	if err != nil {
		return dto.CuratorOnBoardingResponse{}, errors.Wrap(err).WithMessage("failed to create curator")
	}

	return curator, nil
}

func (s *curatorOnboardingService) CheckShopName(shopName string) (*domain.Curator, error) {
	curator, err := s.repo.CheckShopName(shopName)
	if err != nil {
		return nil, errors.Wrap(err).WithMessage("failed to check shop name")
	}

	return curator, nil
}

func (s *curatorOnboardingService) GetCuratorByStoreName(storeName string) (*domain.Curator, error) {
	curators, err := s.repo.GetCuratorByStoreName(storeName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curator from database: %v", err)
	}

	return curators, nil
}
