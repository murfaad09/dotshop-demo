package service

import (
	"encoding/json"
	"fmt"

	dto "github.com/harishash/dotshop-be/integration/vndr/convictional/buyer/dto"
	api "github.com/harishash/dotshop-be/integration/vndr/convictional/client"
)

type IProductService interface {
	GetAllProducts() (dto.ProductList, error)
	GetProductByID(id string) (dto.Product, error)
	DeleteProductImage(productID, imageID string) (dto.Product, error)
	//UpdateProduct(product dto.Product) (dto.Product, error)
}

type ProductService struct {
	client *api.APIClient
}

func NewProductService() IProductService {
	return &ProductService{
		client: api.Client,
	}
}

func (p *ProductService) GetAllProducts() (dto.ProductList, error) {
	allProductsResponse, err := p.client.GetAll("products")
	if err != nil {
		return dto.ProductList{}, fmt.Errorf("error getting all products: %v", err)
	}
	var allProducts dto.ProductList
	err = json.Unmarshal(allProductsResponse, &allProducts)

	if err != nil {
		return dto.ProductList{}, fmt.Errorf("error parsing all products response: %v", err)
	}

	return allProducts, nil
}

func (p *ProductService) GetProductByID(id string) (dto.Product, error) {
	url := fmt.Sprintf("%s/%s", "products", id)
	product, err := p.client.Get(url)
	if err != nil {
		return dto.Product{}, fmt.Errorf("error getting product by id: %v", err)
	}
	var products dto.Product
	err = json.Unmarshal(product, &products)
	if err != nil {
		return dto.Product{}, fmt.Errorf("error parsing product response: %v", err)
	}
	return products, nil
}

func (p *ProductService) DeleteProductImage(productID, imageID string) (dto.Product, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", "products", productID, "images", imageID)
	err := p.client.Delete(url)
	if err != nil {
		return dto.Product{}, fmt.Errorf("error getting while delete product image: %v", err)
	}

	return dto.Product{}, nil
}
