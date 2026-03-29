package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	common_repo "github.com/harishash/dotshop-be/internal/repositories"
)

// CommonService interface
type CommonService interface {
	GetAllProducts(curatorID uint, subCategories string, isFeature bool, page, pageSize int) (*dto.GetFeatureProductResponse, error)
	GetTotalProducts(curatorID uint, subCategories string, isFeature bool) (int, error)
	GetAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error)
	FetchSectionByID(sectionID uint) (*dto.SectionResponse, error)
	GetTotalProductsCount(curatorID uint) (int64, error)
	GetCuratorAllLooks(curatorid uint, page, pageSize int) ([]*dto.LooksResponse, error)
	GetTotalCuratorLooksCount(curatorID uint) (int64, error)
	GetAllProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error)
	GetTotalCollectionProductsCount(collectionID uint) (int64, error)
	GetAllProductsByLookID(lookID uint, page, pageSize int) (*dto.LooksResponse, error)
	GetTotalLookProductsCount(lookID uint) (int64, error)
	GetAllLooks(page, pageSize int) ([]*dto.LooksResponse, error)
	GetTotalLooksCount() (int64, error)
}

// commonService struct
type commonService struct {
	repo common_repo.CommonRepo
}

// NewCommonService creates a new CommonService
func NewCommonService(repo common_repo.CommonRepo) CommonService {
	return &commonService{repo: repo}
}

func (s *commonService) GetAllProducts(curatorID uint, subCategories string, isFeature bool, page, pageSize int) (*dto.GetFeatureProductResponse, error) {
	return s.repo.FetchAllProducts(curatorID, subCategories, isFeature, page, pageSize)
}

func (s *commonService) GetTotalProducts(curatorID uint, subCategories string, isFeature bool) (int, error) {
	return s.repo.GetTotalProducts(curatorID, subCategories, isFeature)
}

func (s *commonService) GetAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error) {
	return s.repo.FetchAllCollections(curatorID, page, pageSize)
}

func (s *commonService) FetchSectionByID(sectionID uint) (*dto.SectionResponse, error) {
	return s.repo.FetchSectionByID(sectionID)
}

func (s *commonService) GetTotalProductsCount(curatorID uint) (int64, error) {
	return s.repo.GetTotalProductsCount(curatorID)
}

func (s *commonService) GetCuratorAllLooks(curatorid uint, page, pageSize int) ([]*dto.LooksResponse, error) {
	looks, err := s.repo.FetchCuratorAllLooks(curatorid, page, pageSize)
	if err != nil {
		return nil, err
	}

	var looksResponse []*dto.LooksResponse
	for _, look := range looks {
		lookResp := &dto.LooksResponse{
			ID:               look.ID,
			Name:             look.Name,
			ImageURL:         look.ImageURL,
			CuratorID:        look.CuratorID,
			SocialID:         look.SocialID,
			SocialType:       look.SocialType,
			SocialTitle:      look.SocialTitle,
			EmbedLink:        look.EmbedLink,
			VideoDescription: look.VideoDescription,
			CreatedAt:        look.CreatedAt,
			UpdatedAt:        look.UpdatedAt,
			Products:         mapProductsResponse(look.Products),
		}
		looksResponse = append(looksResponse, lookResp)
	}

	return looksResponse, nil
}

func (s *commonService) GetAllLooks(page, pageSize int) ([]*dto.LooksResponse, error) {
	looks, err := s.repo.FetchAllLooks(page, pageSize)
	if err != nil {
		return nil, err
	}

	var looksResponse []*dto.LooksResponse
	for _, look := range looks {
		lookResp := &dto.LooksResponse{
			ID:               look.ID,
			Name:             look.Name,
			ImageURL:         look.ImageURL,
			CuratorID:        look.CuratorID,
			SocialID:         look.SocialID,
			SocialType:       look.SocialType,
			SocialTitle:      look.SocialTitle,
			EmbedLink:        look.EmbedLink,
			VideoDescription: look.VideoDescription,
			CreatedAt:        look.CreatedAt,
			UpdatedAt:        look.UpdatedAt,
			Products:         mapProductsResponse(look.Products),
		}
		looksResponse = append(looksResponse, lookResp)
	}

	return looksResponse, nil
}

func mapProductsResponse(products []domain.Product) []*dto.ProductResponse {
	var productResponses []*dto.ProductResponse
	for _, product := range products {
		productResp := &dto.ProductResponse{
			ID:           product.ProductID,
			BrandName:    product.BrandName,
			SupplierName: product.SupplierName,
			Name:         product.Name,
			Description:  product.Description,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
			Variants:     mapVariantsResponse(product.Variants),
		}
		productResponses = append(productResponses, productResp)
	}
	return productResponses
}

func mapVariantsResponse(variants []domain.Variant) []*dto.VariantResponse {
	var variantResponses []*dto.VariantResponse
	for _, variant := range variants {
		variantResp := &dto.VariantResponse{
			ID:              variant.ID,
			ProductID:       variant.ProductID,
			SKU:             variant.SKU,
			Title:           variant.Title,
			InventoryAmount: variant.InventoryAmount,
			Image:           variant.Image,
			RetailPrice:     variant.RetailPrice,
			RetailCurrency:  variant.RetailCurrency,
			BasePrice:       variant.BasePrice,
			BaseCurrency:    variant.BaseCurrency,
			VariantOptions:  mapVariantOptionsResponse(variant.VariantOptions),
			Units:           variant.Units,
			Attributes:      variant.Attributes,
		}
		variantResponses = append(variantResponses, variantResp)
	}
	return variantResponses
}

func mapVariantOptionsResponse(variantOptions []domain.VariantOption) []*dto.VariantOptionResponse {
	var variantOptionResponses []*dto.VariantOptionResponse
	for _, variantOption := range variantOptions {
		variantOptionResp := &dto.VariantOptionResponse{
			Name:      variantOption.Name,
			Value:     variantOption.Value,
			VariantID: variantOption.VariantID,
		}
		variantOptionResponses = append(variantOptionResponses, variantOptionResp)
	}
	return variantOptionResponses
}

func (s *commonService) GetTotalCuratorLooksCount(curatorID uint) (int64, error) {
	return s.repo.GetTotalCuratorLooksCount(curatorID)
}

func (s *commonService) GetTotalLooksCount() (int64, error) {
	return s.repo.GetTotalLooksCount()
}

func (s *commonService) GetAllProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error) {
	return s.repo.FetchProductsByCollectionID(collectionID, page, pageSize)
}

func (s *commonService) GetTotalCollectionProductsCount(collectionID uint) (int64, error) {
	return s.repo.GetTotalCollectionProductsCount(collectionID)
}

func (s *commonService) GetAllProductsByLookID(lookID uint, page, pageSize int) (*dto.LooksResponse, error) {
	look, err := s.repo.FetchProductsByLookID(lookID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &dto.LooksResponse{
		ID:               look.ID,
		Name:             look.Name,
		ImageURL:         look.ImageURL,
		CuratorID:        look.CuratorID,
		SocialID:         look.SocialID,
		SocialType:       look.SocialType,
		SocialTitle:      look.SocialTitle,
		EmbedLink:        look.EmbedLink,
		VideoDescription: look.VideoDescription,
		Products:         mapProductsResponse(look.Products),
	}, nil
}

func mapLookResponse(look *domain.Look) *dto.LooksResponse {
	return &dto.LooksResponse{
		ID:               look.ID,
		Name:             look.Name,
		ImageURL:         look.ImageURL,
		CuratorID:        look.CuratorID,
		SocialID:         look.SocialID,
		SocialType:       look.SocialType,
		SocialTitle:      look.SocialTitle,
		EmbedLink:        look.EmbedLink,
		VideoDescription: look.VideoDescription,
		Products:         mapProductsResponse(look.Products),
	}
}
func (s *commonService) GetTotalLookProductsCount(lookID uint) (int64, error) {
	return s.repo.GetTotalLookProductsCount(lookID)
}
