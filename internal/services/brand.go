package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	"github.com/harishash/dotshop-be/internal/utils/logger"

	repository "github.com/harishash/dotshop-be/internal/repositories"
)

type BrandService interface {
	GetBrands(query *dto.BrandsRequest) (*dto.Response, error)
	UpdateBrandStatus(id int, isActive bool) error
}

type brandService struct {
	repo        repository.BrandRepository
	productRepo repository.IProductRepository
}

func NewBrandService(repo repository.BrandRepository, productRepo repository.IProductRepository) BrandService {
	return &brandService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *brandService) GetBrands(query *dto.BrandsRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetBrands(query)
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

func (s *brandService) UpdateBrandStatus(id int, isActive bool) error {
	brand, err := s.repo.GetBrandByID(uint(id))
	if err != nil {
		return err
	}

	brand.IsActive = isActive
	if err = s.repo.UpdateBrand(brand); err != nil {
		return err
	}

	if err = s.productRepo.UpdateProductStatusByBrandID(brand.ID, isActive); err != nil {
		logger.Errorf("failed to update product statuses for brand: %v, error : %v", id, err)
		return err
	}

	return nil
}
