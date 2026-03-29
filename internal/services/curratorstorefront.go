package service

import (
	"github.com/harishash/dotshop-be/internal/config"
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	curratostore_repo "github.com/harishash/dotshop-be/internal/repositories"
)

// CuratorStoreFrontService interface
type CuratorStoreFrontService interface {
	// GetAllProducts(curatorID uint, isFeature bool, page, limit int) (*dto.GetFeatureProductResponse, error)
	// GetTotalProducts(curatorID uint, isFeature bool) (int, error)
	GetAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error)
	//GetAllLooks() ([]models.Look, error)
	GetAllProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error)
	// GetAllProductsByLookID(lookID string) ([]models.Product, error)
	SearchLookByName(query *dto.SearchLookParams) ([]*dto.LooksResponse, *dto.Paging, error)
	SearchProductsWithinCuratorLooks(curatorID uint, searchQuery *dto.SearchProductParams) ([]*dto.ProductResponse, *dto.Paging, error)
	SearchFeatureProductsByName(curatorID uint, searchQuery *dto.SearchProductParams) ([]*dto.ProductResponse, *dto.Paging, error)
	SearchCollectionByName(searchQuery *dto.SearchCollectionParams, curatorId uint) ([]dto.CollectionWithProducts, *dto.Paging, error)
	SearchCollectionProductByName(collectionId uint64, searchQuery *dto.SearchProductParams) (*dto.CollectionResponse, *dto.Paging, error)
	SearchSectionByName(collectionId uint64, searchQuery *dto.SearchSectionParams) ([]dto.SectionResponse, *dto.Paging, error)
	GlobalSearch(searchQuery *dto.SearchProductParams) (*dto.GlobalSearchResponse, error)
	GetDotShopCuratorId() (*domain.Curator, error)
}

// curatorStoreFrontService struct
type curatorStoreFrontService struct {
	repo curratostore_repo.CommonRepo
}

// NewCuratorStoreFrontService creates a new CuratorStoreFrontService
func NewCuratorStoreFrontService(repo curratostore_repo.CommonRepo) CuratorStoreFrontService {
	return &curatorStoreFrontService{repo: repo}
}

// func (s *curatorStoreFrontService) GetAllProducts(curatorID uint, isFeature bool, page, limit int) (*dto.GetFeatureProductResponse, error) {
// 	offset := (page - 1) * limit
// 	return s.repo.FetchAllProducts(curatorID, isFeature, offset, limit)
// }

// func (s *curatorStoreFrontService) GetTotalProducts(curatorID uint, isFeature bool) (int, error) {
// 	return s.repo.GetTotalProducts(curatorID, "", isFeature)
// }

func (s *curatorStoreFrontService) GetAllCollections(curatorID uint, page, pageSize int) ([]*dto.CollectionResponse, error) {
	return s.repo.FetchAllCollections(curatorID, page, pageSize)
}

// func (s *curatorStoreFrontService) GetAllLooks() ([]models.Look, error) {
// 	return s.repo.FetchAllLooks()
// }

func (s *curatorStoreFrontService) GetAllProductsByCollectionID(collectionID uint, page, pageSize int) (*dto.CollectionResponse, error) {
	return s.repo.FetchProductsByCollectionID(collectionID, page, pageSize)
}

// func (s *curatorStoreFrontService) GetAllProductsByLookID(lookID string) ([]models.Product, error) {
// 	return s.repo.FetchProductsByLookID(lookID)
// }

func (s *curatorStoreFrontService) SearchLookByName(query *dto.SearchLookParams) ([]*dto.LooksResponse, *dto.Paging, error) {
	looks, paging, err := s.repo.SearchLooksByName(query)
	if err != nil {
		return nil, nil, err
	}

	var response []*dto.LooksResponse
	for _, v := range looks {
		var products []*dto.ProductResponse
		for _, product := range v.Products {
			var variants []*dto.VariantResponse
			for _, variant := range product.Variants {

				var variantsOption []*dto.VariantOptionResponse
				for _, variantOption := range variant.VariantOptions {
					variantsOption = append(variantsOption, &dto.VariantOptionResponse{
						Name:      variantOption.Name,
						Value:     variantOption.Value,
						VariantID: variantOption.VariantID,
					})
				}

				variants = append(variants, &dto.VariantResponse{
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
					Units:           variant.Units,
					Attributes:      variant.Attributes,
					VariantOptions:  variantsOption,
				})
			}
			products = append(products, &dto.ProductResponse{
				ID:           product.ProductID,
				BrandName:    product.BrandName,
				SupplierName: product.SupplierName,
				Name:         product.Name,
				Description:  product.Description,
				CreatedAt:    product.CreatedAt,
				UpdatedAt:    product.UpdatedAt,
				Variants:     variants,
			})
		}
		response = append(response, &dto.LooksResponse{
			ID:               v.ID,
			Name:             v.Name,
			ImageURL:         v.ImageURL,
			CuratorID:        v.CuratorID,
			SocialID:         v.SocialID,
			SocialType:       v.SocialType,
			SocialTitle:      v.SocialTitle,
			EmbedLink:        v.EmbedLink,
			VideoDescription: v.VideoDescription,
			Products:         products,
		})
	}

	return response, paging, nil
}

func (s *curatorStoreFrontService) SearchProductsWithinCuratorLooks(curatorID uint, searchQuery *dto.SearchProductParams) ([]*dto.ProductResponse, *dto.Paging, error) {
	products, paging, err := s.repo.SearchProductsWithinCuratorLooks(curatorID, searchQuery)
	if err != nil {
		return nil, nil, err
	}

	return dto.ProductToResponse(products), paging, nil
}

func (s *curatorStoreFrontService) SearchFeatureProductsByName(curatorID uint, searchQuery *dto.SearchProductParams) ([]*dto.ProductResponse, *dto.Paging, error) {
	products, paging, err := s.repo.SearchFeatureProductsByName(curatorID, searchQuery)
	if err != nil {
		return nil, nil, err
	}
	return dto.ProductToResponse(products), paging, nil
}

func (s *curatorStoreFrontService) SearchCollectionByName(searchQuery *dto.SearchCollectionParams, curatorId uint) ([]dto.CollectionWithProducts, *dto.Paging, error) {
	return s.repo.SearchCollectionsByName(searchQuery, curatorId)
}

func (s *curatorStoreFrontService) SearchCollectionProductByName(collectionId uint64, searchQuery *dto.SearchProductParams) (*dto.CollectionResponse, *dto.Paging, error) {
	collection, paging, err := s.repo.SearchCollectionProductByName(collectionId, searchQuery)
	if err != nil {
		return nil, nil, err
	}
	return dto.CollectionToResponse(collection), paging, nil
}

func (s *curatorStoreFrontService) SearchSectionByName(collectionId uint64, searchQuery *dto.SearchSectionParams) ([]dto.SectionResponse, *dto.Paging, error) {
	collection, paging, err := s.repo.SearchCollectionSectionsByName(collectionId, searchQuery)
	if err != nil {
		return nil, nil, err
	}
	return dto.CollectionSectionsToResponse(collection), paging, nil
}

func (s *curatorStoreFrontService) GlobalSearch(searchQuery *dto.SearchProductParams) (*dto.GlobalSearchResponse, error) {
	product, suggestions, paging, err := s.repo.GlobalSearch(searchQuery)
	if err != nil {
		return nil, err
	}

	resp := &dto.GlobalSearchResponse{
		Data:        dto.ProductToResponse(product),
		Suggestions: suggestions,
		Paging:      *paging,
	}

	return resp, nil
}

func (s *curatorStoreFrontService) GetDotShopCuratorId() (*domain.Curator, error) {
	curator, err := s.repo.GetCuratorByEmail(config.GetConfig().DotShopStoreEmail)
	if err != nil {
		return nil, err
	}

	return curator, nil
}
