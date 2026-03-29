package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/harishash/dotshop-be/internal/config"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	model "github.com/harishash/dotshop-be/internal/models"
	"github.com/harishash/dotshop-be/internal/utils/errors"
	"github.com/harishash/dotshop-be/internal/utils/logger"

	repository "github.com/harishash/dotshop-be/internal/repositories"
)

type AdminDashboardService interface {
	GetTopSellingBrands(query *dto.CommonProductRequest) (*dto.Response, error)
	GetTopSellingProducts(query *dto.CommonProductRequest) (*dto.Response, error)
	GetTopCurators(query *dto.CommonProductRequest) (*dto.Response, error)
	GetSaleByCategory(query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error)
	GetTopWishlistProducts(query *dto.CommonProductRequest) (*dto.Response, error)
	GetOrderSales(query *dto.OrderSalesRequest) (*dto.Response, error)
	GetOrderReturns(query *dto.OrderSalesRequest) (*dto.Response, error)
	GetAllCustomers(query *dto.CustomersRequest) (*dto.Response, error)
	GetAllCurators(query *dto.CuratorRequest) (*dto.Response, error)
	GetPaymentDistribution(query *dto.PaymentDistributionRequest) (*dto.Response, error)
	UpdateReturnStatus(id uint, status bool) error
	DeleteUser(id uint) error
	DeleteCurator(id uint) error
	BlockCustomer(userID uint, isBlock bool) error
	BlockCurator(curatorID uint, isBlock bool) error
	ConsumerOrdersListByUserId(userId uint64, query *dto.CustomersOrderListRequest) (*dto.Response, error)
	CuratorOrdersListById(curatorId uint, query *dto.CuratorsOrderListRequest) (*dto.Response, error)
	GetListedProducts(params *dto.PagingParams, curatorID uint) (*dto.Response, error)
	DeleteUserReview(id uint) error
}

type adminDashboardService struct {
	adminRepo   repository.AdminDashboardRepository
	orderRepo   repository.IOrderRepository
	productRepo repository.IProductRepository
	userRepo    repository.IUserRepository
	reviewRepo  repository.ReviewRepository
	commonRepo  repository.CommonRepo
}

func NewAdminDashboardService(adminRepo repository.AdminDashboardRepository, orderRepo repository.IOrderRepository,
	productRepo repository.IProductRepository, userRepo repository.IUserRepository,
	reviewRepo repository.ReviewRepository, commonRepo repository.CommonRepo) AdminDashboardService {
	return &adminDashboardService{
		adminRepo:   adminRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		reviewRepo:  reviewRepo,
		commonRepo:  commonRepo,
	}
}

func (s *adminDashboardService) GetTopSellingBrands(query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetTopSellingBrands(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetTopSellingProducts(query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetTopSellingProducts(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetTopCurators(query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetTopCurators(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetSaleByCategory(query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error) {
	return s.adminRepo.GetSalesByCategory(query)
}

func (s *adminDashboardService) GetTopWishlistProducts(query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetTopWishlistProducts(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetAllCurators(query *dto.CuratorRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetAllCurators(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetPaymentDistribution(query *dto.PaymentDistributionRequest) (*dto.Response, error) {
	curator, err := s.commonRepo.GetCuratorByEmail(config.GetConfig().DotShopStoreEmail)
	if err != nil {
		logger.Errorf("failed to get curator: %v, error : %v", config.GetConfig().DotShopStoreEmail, err)
		return nil, err
	}

	data, paging, err := s.adminRepo.GetPaymentDistribution(query, curator.ID)
	if err != nil {
		return nil, err
	}

	for i := range data {
		data[i].DotShopProfitPercentage = constants.DotShopProfitPercentage
		data[i].CuratorCommissionPercentage = constants.CuratorCommissionPercentage
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetOrderSales(query *dto.OrderSalesRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetOrderSales(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetOrderReturns(query *dto.OrderSalesRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetOrderReturns(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) GetAllCustomers(query *dto.CustomersRequest) (*dto.Response, error) {
	data, paging, err := s.adminRepo.GetAllCustomers(query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}

	return res, nil
}

func (s *adminDashboardService) UpdateReturnStatus(id uint, status bool) error {
	// Add any business logic or validation here if needed
	if id == 0 {
		return errors.New("invalid ID")
	}

	statusStr := "Declined"
	if status {
		statusStr = "Approved"
	}

	return s.adminRepo.UpdateReturnStatus(id, statusStr)
}

func (s *adminDashboardService) DeleteUser(id uint) error {
	return s.adminRepo.DeleteUser(id)
}

func (s *adminDashboardService) DeleteCurator(id uint) error {
	return s.adminRepo.DeleteCurator(id)
}
func (s *adminDashboardService) BlockCustomer(userID uint, isBlock bool) error {
	user, err := s.adminRepo.FindCustomerByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("customer not found")
	}

	user.IsBlock = isBlock
	return s.adminRepo.UpdateCustomerStatus(user)
}

func (s *adminDashboardService) GetListedProducts(params *dto.PagingParams, curatorID uint) (*dto.Response, error) {
	var response []*dto.ListedProductResponse

	products, paging, err := s.adminRepo.GetListedProducts(params, curatorID)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return &dto.Response{
			Data:   response,
			Paging: *paging,
		}, nil
	}

	for _, product := range products {
		variants, err := s.productRepo.GetVariantsByProductId(product.ProductID)
		if err != nil {
			return nil, err
		}

		productResponse := &dto.ListedProductResponse{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			BrandName:   product.BrandName,
			Variants:    processProductVariants(variants),
		}

		response = append(response, productResponse)
	}

	res := &dto.Response{
		Data:   response,
		Paging: *paging,
	}

	return res, nil
}

func processProductVariants(variants []*model.Variant) []*dto.ListedProductVariant {
	var productVariants []*dto.ListedProductVariant
	for _, variant := range variants {
		productVariants = append(productVariants, &dto.ListedProductVariant{
			ID:             variant.ID,
			Title:          variant.Title,
			Image:          variant.Image,
			RetailPrice:    variant.RetailPrice,
			RetailCurrency: variant.RetailCurrency,
			BasePrice:      variant.BasePrice,
			BaseCurrency:   variant.BaseCurrency,
			Units:          variant.Units,
		})
	}

	return productVariants
}

func (s *adminDashboardService) BlockCurator(curatorID uint, isBlock bool) error {
	curator, err := s.adminRepo.FindCuratorByID(curatorID)
	if err != nil {
		return err
	}
	fmt.Println("curator", curator)
	if curator == nil {
		return errors.New("curator not found")
	}

	curator.IsBlock = isBlock
	fmt.Println("curator", curator)

	if err := s.adminRepo.UpdateCuratorStatus(curator); err != nil {
		return err
	}
	fmt.Println("curator.UserID", curator.UserID)

	if err1 := s.BlockCustomer(curator.UserID, isBlock); err1 != nil {
		return err
	}
	return nil
}

func (h *adminDashboardService) ConsumerOrdersListByUserId(userId uint64, query *dto.CustomersOrderListRequest) (*dto.Response, error) {
	var response []dto.CustomerOrdersListResponse

	orders, paging, err := h.orderRepo.GetOrdersByUserIdForCustomer(userId, query)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		orderVariants, err := h.orderRepo.GetOrderVariantByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		var variantResponses []dto.OrderListVariantResponse
		var orderAmount float64
		var orderQuantity uint
		for _, variant := range orderVariants {
			product, variantDB, err := h.productRepo.GetProductByVariantId(variant.VariantID)
			if err != nil {
				return nil, err
			}

			variantResponses = append(variantResponses, dto.OrderListVariantResponse{
				ID:                variant.VariantID,
				ProductID:         variant.ProductID,
				ProductName:       product.Name,
				BrandName:         product.BrandName,
				CuratorID:         variant.CuratorID,
				Quantity:          variant.Quantity,
				Description:       product.Description,
				SKU:               variantDB.SKU,
				Title:             variantDB.Title,
				Image:             variantDB.Image,
				Price:             variant.Price,
				RetailPrice:       variant.Price,
				RetailCurrency:    variantDB.RetailCurrency,
				VariantOptionName: variant.VariantOptionName,
				VariantSize:       variant.VariantSize,
			})

			orderAmount += variant.Price
			orderQuantity += variant.Quantity
		}

		resp := dto.CustomerOrdersListResponse{
			ID:              order.ID,
			UserID:          order.UserID,
			TotalAmount:     orderAmount,
			TotalQuantity:   orderQuantity,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			ShippingCity:    order.ShippingCity,
			ShippingState:   order.ShippingState,
			ShippingCountry: order.ShippingCountry,
			ShippingZip:     order.ShippingZip,
			BuyerReference:  order.BuyerReference,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			Note:            order.Note,
			Variant:         variantResponses,
		}
		response = append(response, resp)
	}

	res := &dto.Response{
		Data:   response,
		Paging: *paging,
	}
	return res, nil
}

func (h *adminDashboardService) CuratorOrdersListById(curatorId uint, query *dto.CuratorsOrderListRequest) (*dto.Response, error) {

	var response []dto.CuratorOrdersListResponse
	orderVariants, paging, err := h.orderRepo.GetCuratorOrderIDsByID(curatorId, query)
	if err != nil {
		return nil, err
	}

	for _, orderId := range orderVariants {
		order, err := h.orderRepo.GetOrderById(orderId)
		if err != nil {
			return nil, err
		}

		orderVariants, err := h.orderRepo.GetOrderVariantsByOrderId(orderId, &curatorId)
		if err != nil {
			return nil, err
		}

		var variantResponses []dto.OrderListVariantResponse
		var orderAmount float64
		var orderQuantity uint
		for _, variant := range orderVariants {
			product, variantDB, err := h.productRepo.GetProductByVariantId(variant.VariantID)
			if err != nil {
				return nil, err
			}

			variantResponses = append(variantResponses, dto.OrderListVariantResponse{
				ID:                variant.VariantID,
				ProductID:         variant.ProductID,
				ProductName:       product.Name,
				BrandName:         product.BrandName,
				CuratorID:         variant.CuratorID,
				Quantity:          variant.Quantity,
				Description:       product.Description,
				SKU:               variantDB.SKU,
				Title:             variantDB.Title,
				Image:             variantDB.Image,
				Price:             variant.Price,
				RetailPrice:       variant.Price,
				RetailCurrency:    variantDB.RetailCurrency,
				VariantOptionName: variant.VariantOptionName,
				VariantSize:       variant.VariantSize,
			})

			orderAmount += variant.Price
			orderQuantity += variant.Quantity
		}

		user, err := h.userRepo.GetUserByID(uint(order.UserID))
		if err != nil {
			log.Errorf("error getting user: %v", err)
		}

		response = append(response, dto.CuratorOrdersListResponse{
			ID:                order.ID,
			UserID:            order.UserID,
			CustomerFirstName: user.FirstName,
			CustomerLastName:  user.LastName,
			TotalAmount:       orderAmount,
			TotalQuantity:     orderQuantity,
			Status:            order.Status,
			BuyerReference:    order.BuyerReference,
			CreatedAt:         order.CreatedAt,
			UpdatedAt:         order.UpdatedAt,
			Note:              order.Note,
			Variant:           variantResponses,
		})
	}

	res := &dto.Response{
		Data:   response,
		Paging: *paging,
	}

	return res, nil
}

func (h *adminDashboardService) DeleteUserReview(id uint) error {
	return h.reviewRepo.DeleteReview(id)
}
