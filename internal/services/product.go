package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	product_repo "github.com/harishash/dotshop-be/internal/repositories"
)

type IProductService interface {
	GetProductsWithFilter(query *dto.Filter) ([]dto.FilteredProductResponse, *dto.Paging, error)
	GetProductById(id string) (*domain.Product, error)
	CreateProduct(product domain.Product) (*domain.Product, error)
	UpdateProduct(product domain.Product) (*dto.CreateProductResponse, error)
	DeleteProducts(ids []*string) ([]*domain.Product, error)
	GetBrands(searchStr string) ([]*dto.Brands, error)
	GetAllProductsWithStats(params *dto.ListProductReviewsRequest) (*dto.Response, error)
}

type ProductService struct {
	productRepository product_repo.IProductRepository
}

func NewProductService(productRepository product_repo.IProductRepository) IProductService {
	return &ProductService{
		productRepository: productRepository}
}

func (s *ProductService) GetProductsWithFilter(query *dto.Filter) ([]dto.FilteredProductResponse, *dto.Paging, error) {
	return s.productRepository.GetProductsWithFilter(query)
}

func (s *ProductService) GetProductById(id string) (*domain.Product, error) {
	return s.productRepository.GetProductById(id)
}

func (s *ProductService) CreateProduct(product domain.Product) (*domain.Product, error) {
	return s.productRepository.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(product domain.Product) (*dto.CreateProductResponse, error) {
	updatedProduct, err := s.productRepository.UpdateProduct(&product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProducts(ids []*string) ([]*domain.Product, error) {

	var products []*domain.Product
	for _, id := range ids {
		product, err := s.productRepository.DeleteProduct(*id)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil
}

func (s *ProductService) GetBrands(searchStr string) ([]*dto.Brands, error) {
	return s.productRepository.GetBrands(searchStr)
}

func (s *ProductService) GetAllProductsWithStats(params *dto.ListProductReviewsRequest) (*dto.Response, error) {
	data, paging, err := s.productRepository.GetAllProductsWithStats(params)
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
