package repository

import (
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"

	"gorm.io/gorm"
)

type BrandRepository interface {
	CreateBrand(brand models.Brand) (*models.Brand, error)
	UpdateBrand(brand *models.Brand) error
	GetBrands(query *dto.BrandsRequest) ([]*dto.BrandsResponse, *dto.Paging, error)
	GetBrandByID(id uint) (*models.Brand, error)
	GetBrandIDByName(name string) uint
}

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db}
}

func (r *brandRepository) CreateBrand(brand models.Brand) (*models.Brand, error) {
	err := r.db.Create(&brand).Error
	if err != nil {
		return nil, err
	}

	return &brand, nil
}

func (r *brandRepository) UpdateBrand(brand *models.Brand) error {
	return r.db.Save(brand).Error
}

func (r *brandRepository) GetBrands(query *dto.BrandsRequest) ([]*dto.BrandsResponse, *dto.Paging, error) {
	var totalCount int64
	var brands []*dto.BrandsResponse

	err := r.db.Model(&models.Brand{}).Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	offset := (query.PageNum - 1) * query.PageSize

	err = r.db.Table("brands").
		Select("brands.id, brands.name, brands.description, brands.is_active, COUNT(products.product_id) AS number_of_products").
		Joins("LEFT JOIN products ON products.brand_id = brands.id").
		Group("brands.id").
		Limit(query.PageSize).
		Offset(offset).
		Find(&brands).Error
	if err != nil {
		return nil, nil, err
	}

	paging := &dto.Paging{
		TotalCount:  totalCount,
		CurrentPage: query.PageNum,
		PageSize:    query.PageSize,
	}

	return brands, paging, nil
}

func (r *brandRepository) GetBrandByID(id uint) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.First(&brand, id).Error
	if err != nil {
		return nil, err
	}

	return &brand, nil
}

func (r *brandRepository) GetBrandIDByName(name string) uint {
	var brand models.Brand
	err := r.db.Where("name = ?", name).First(&brand).Error
	if err != nil {
		return 0
	}

	return brand.ID
}
