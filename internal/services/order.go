package service

import (
	"fmt"
	"time"

	"github.com/harishash/dotshop-be/integration/aws"
	"github.com/harishash/dotshop-be/internal/dto"
	user_repo "github.com/harishash/dotshop-be/internal/repositories"
	"github.com/harishash/dotshop-be/internal/utils/errors"
	"github.com/harishash/dotshop-be/internal/utils/logger"

	api "github.com/harishash/dotshop-be/integration/vndr/convictional/client"
	domain "github.com/harishash/dotshop-be/internal/models"
	repo "github.com/harishash/dotshop-be/internal/repositories"
)

type IOrderService interface {
	CreateOrder(body *dto.OrderRequest) (*dto.OrderResponse, error)
	OrdersList(userId uint64) ([]dto.OrdersListResponse, error)
	CancelOrder(body *dto.CancelOrderRequest, orderId string) error
	CreateReturn(body *dto.ReturnRequest, userId uint64) ([]*dto.ReturnResponse, error)
}

type OrderService struct {
	client               *api.APIClient
	orderRepository      repo.IOrderRepository
	userRepository       repo.IUserRepository
	noticationRepository repo.NotificationRepositoryInterface
	adminRepository      repo.AdminRepo
	awsService           aws.AWSService
	productRepository    user_repo.IProductRepository
}

func NewOrderService(
	repo repo.IOrderRepository,
	userRepo repo.IUserRepository,
	notificationRepo repo.NotificationRepositoryInterface,
	adminRepository repo.AdminRepo,
	awsService aws.AWSService,
	productRepository user_repo.IProductRepository,

) *OrderService {
	return &OrderService{
		client:               api.Client,
		orderRepository:      repo,
		userRepository:       userRepo,
		noticationRepository: notificationRepo,
		adminRepository:      adminRepository,
		awsService:           awsService,
		productRepository:    productRepository,
	}
}

func (h *OrderService) CreateOrder(body *dto.OrderRequest) (*dto.OrderResponse, error) {
	// check all ordered variants exists in our store
	userAddress, err := h.userRepository.GetUserAddressByID(body.AddressID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.GetUserByID(uint(body.UserId))
	if err != nil {
		return nil, err
	}

	var variants []domain.Variant
	for _, v := range body.Items {
		variant, err := h.orderRepository.GetVariant(v.VariantID)
		if err != nil {
			return nil, err
		}

		variants = append(variants, *variant)
	}

	request := dto.ConvicationalOrderRequest(body, userAddress)
	request.CustomerEmail = user.Email

	if userAddress.FirstName == nil || userAddress.LastName == nil {
		return nil, errors.New("first and last name is required in address")
	}

	request.Address.Name = *userAddress.FirstName + " " + *userAddress.LastName

	// send request for create order to convicational api
	cnResp, err := h.client.CreateOrder(request)
	if err != nil {
		return nil, err
	}

	var totalPrice float64
	var totalQuantity uint
	var orderVarients []*domain.OrderVariants
	var fulfillments []domain.Fulfillments
	for i, v := range variants {

		brandName := h.productRepository.GetBrandName(v.ProductID)
		totalPrice += float64(body.Items[i].Quantity) * v.RetailPrice
		totalQuantity += body.Items[i].Quantity
		orderVarients = append(orderVarients, &domain.OrderVariants{
			OrderID:           cnResp.ID,
			CuratorID:         body.Items[i].CuratorID,
			ProductID:         v.ProductID,
			BrandName:         brandName,
			SellerOrderId:     cnResp.Items[i].SellerOrderID,
			SellerOrderItemId: cnResp.Items[i].SellerOrderItemID,
			BuyerReference:    cnResp.Items[i].BuyerReference,
			Price:             float64(body.Items[i].Quantity) * v.RetailPrice,
			VariantID:         cnResp.Items[i].VariantID,
			Quantity:          cnResp.Items[i].Quantity,
			VariantOptionName: body.Items[i].VariantOptionName,
			VariantSize:       body.Items[i].VariantSize,
		})
		for i := range cnResp.Fulfillments {
			fulfillments = append(fulfillments, domain.Fulfillments{
				ID:           cnResp.Fulfillments[i].ID,
				OrderID:      cnResp.ID,
				Posted:       cnResp.Fulfillments[i].Posted,
				Carrier:      cnResp.Fulfillments[i].Carrier,
				TrackingCode: cnResp.Fulfillments[i].TrackingCode,
			})
		}
	}

	orderDomain := &domain.Order{
		ID:              cnResp.ID,
		UserID:          body.UserId,
		BuyerReference:  body.BuyerReference,
		ShippingMethod:  body.ShippingMethod,
		PaymentID:       body.PaymentID,
		Note:            body.Note,
		TotalAmount:     totalPrice,
		TotalQuantity:   totalQuantity,
		IsTest:          body.IsTest,
		ShippingAddress: &userAddress.AddressOne.String,
		ShippingCity:    &userAddress.City.String,
		ShippingState:   &userAddress.State.String,
		ShippingZip:     &userAddress.Zip.String,
		ShippingCountry: &userAddress.Country.String,
	}

	orderResp, err := h.orderRepository.CreateOrder(orderDomain, orderVarients, fulfillments)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("A new order #%s has been placed.", cnResp.ID)
	if err := h.noticationRepository.NotifyCurators(orderVarients, message); err != nil {
		logger.Warnf("warning: failed to notify curators: %v", err)
	}

	if err := h.noticationRepository.NotifyAdmins(message); err != nil {
		logger.Warnf("warning: failed to notify admin: %v", err)
	}

	resp := &dto.OrderResponse{
		UserId:         body.UserId,
		OrderId:        cnResp.ID,
		Address:        dto.NewAddressDS(userAddress),
		Items:          body.Items,
		BuyerReference: orderResp.BuyerReference,
		Note:           orderResp.Note,
		IsTest:         body.IsTest,
		FulFillments:   cnResp.Fulfillments,
		CreatedAt:      orderResp.CreatedAt,
	}

	return resp, nil
}

func (h *OrderService) OrdersList(userId uint64) ([]dto.OrdersListResponse, error) {
	curator, err := h.orderRepository.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	var response []dto.OrdersListResponse
	orderVariants, err := h.orderRepository.GetOrderIDsByCuratorID(uint64(curator.ID))
	if err != nil {
		return nil, err
	}

	for _, orderId := range orderVariants {
		order, err := h.orderRepository.GetOrderById(orderId)
		if err != nil {
			return nil, err
		}

		orderVariants, err := h.orderRepository.GetOrderVariantsByOrderId(orderId, &curator.ID)
		if err != nil {
			return nil, err
		}

		var variantResponses []dto.OrderListVariantResponse
		var orderAmount float64
		var orderQuantity uint
		for _, variant := range orderVariants {
			product, variantDB, err := h.productRepository.GetProductByVariantId(variant.VariantID)
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

		response = append(response, dto.OrdersListResponse{
			ID:             order.ID,
			UserID:         order.UserID,
			TotalAmount:    orderAmount,
			TotalQuantity:  orderQuantity,
			Status:         order.Status,
			BuyerReference: order.BuyerReference,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
			Note:           order.Note,
			Variant:        variantResponses,
		})

	}

	return response, nil
}

func (h *OrderService) CancelOrder(body *dto.CancelOrderRequest, orderId string) error {
	orderdb, err := h.orderRepository.GetOrderById(orderId)
	if err != nil {
		return err
	}

	if orderdb.Cancelled {
		return errors.New("order already cancelled")
	}

	orderVariant, err := h.orderRepository.GetOrderVariantsByOrderId(orderId, nil)
	if err != nil {
		return err
	}

	tx, err := h.orderRepository.StartTx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now().UTC()
	order := &domain.Order{
		ID:              orderId,
		Status:          "cancelled",
		Cancelled:       true,
		CancelledReason: body.Reason,
		CancelledData:   &now,
	}
	if err := h.orderRepository.UpdateCancelOrder(tx, order); err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range orderVariant {
		variant := &domain.OrderVariants{
			ID:              v.ID,
			Cancelled:       true,
			CancelledReason: body.Reason,
			CancelledData:   &now,
		}

		if err := h.orderRepository.UpdateCancelOrderVariant(tx, variant); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := h.client.CancelOrder(body, orderId); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err.Error())
	}

	message := fmt.Sprintf("Order #%s has been returned.", orderId)
	if err := h.noticationRepository.NotifyCurators(orderVariant, message); err != nil {
		logger.Warnf("warning: failed to notify curators: %v", err)
	}

	if err := h.noticationRepository.NotifyAdmins(message); err != nil {
		logger.Warnf("warning: failed to notify curators: %v", err)
	}

	return nil
}

func (h *OrderService) CreateReturn(body *dto.ReturnRequest, userId uint64) ([]*dto.ReturnResponse, error) {
	var returnResponses []*dto.ReturnResponse
	var orderVariants []*domain.OrderVariants
	var returnRequests []*domain.ReturnOrder
	var returnTotalAmount float64
	var returnTotalQuantity uint

	orderdb, err := h.orderRepository.GetOrderById(body.OrderId)
	if err != nil {
		return nil, err
	}

	for _, req := range body.ReturnVariants {
		orderVariant, err := h.orderRepository.GetOrderVariant(body.OrderId, req.VariantId, userId)
		if err != nil {
			return nil, err
		}

		if err := validateReturnVariant(req, orderVariant); err != nil {
			return nil, err
		}

		returnRequests = append(returnRequests, &domain.ReturnOrder{
			UserId:         uint(userId),
			CuratorId:      uint(orderVariant.CuratorID),
			OrderId:        body.OrderId,
			OrderVariantId: orderVariant.ID,
			Status:         "pending",
			Reason:         body.Reason,
			Quantity:       req.Quantity,
			Amount:         orderVariant.Variant.RetailPrice * float64(req.Quantity),
		})

		returnTotalAmount += float64(req.Quantity) * orderVariant.Variant.RetailPrice
		returnTotalQuantity += req.Quantity

		orderVariants = append(orderVariants, &domain.OrderVariants{
			ID:                orderVariant.ID,
			OrderID:           orderVariant.OrderID,
			VariantID:         orderVariant.VariantID,
			SellerOrderId:     orderVariant.SellerOrderId,
			SellerOrderItemId: orderVariant.SellerOrderItemId,
			BuyerReference:    orderVariant.BuyerReference,
			Quantity:          orderVariant.Quantity - req.Quantity,
			Price:             orderVariant.Price - orderVariant.Variant.RetailPrice,
		})
	}

	order := &domain.Order{
		ID:            body.OrderId,
		TotalQuantity: orderdb.TotalQuantity - returnTotalQuantity,
		TotalAmount:   orderdb.TotalAmount - returnTotalAmount,
	}

	returnResp, err := h.orderRepository.ReturnOrder(order, orderVariants, returnRequests)
	if err != nil {
		return nil, err
	}

	for i, v := range returnResp {
		resp := &dto.ReturnResponse{
			Id:                v.ID,
			UserId:            uint(userId),
			OrderId:           body.OrderId,
			VariantId:         orderVariants[i].VariantID,
			OrderVariantId:    v.OrderVariantId,
			SellerOrderId:     orderVariants[i].SellerOrderId,
			SellerOrderItemId: orderVariants[i].SellerOrderItemId,
			BuyerCode:         orderVariants[i].BuyerReference,
			Quantity:          v.Quantity,
			Status:            v.Status,
			Reason:            v.Reason,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		}

		returnResponses = append(returnResponses, resp)

	}

	message := fmt.Sprintf("Order #%s has been returned.", body.OrderId)
	if err := h.noticationRepository.NotifyCurators(orderVariants, message); err != nil {
		logger.Warnf("warning: failed to notify curators: %v", err)
	}

	if err := h.noticationRepository.NotifyAdmins(message); err != nil {
		logger.Warnf("warning: failed to notify admins: %v", err)
	}

	return returnResponses, nil
}

func validateReturnVariant(req dto.ReturnOrderVariantRequest, variant *domain.OrderVariants) error {
	if variant.ID <= 0 {
		return fmt.Errorf("order variant not found for variant: %v", variant.ID)
	}
	if time.Since(variant.CreatedAt) > 21*24*time.Hour {
		return fmt.Errorf("returns are only accepted within 21 days of purchase for variant: %v", variant.ID)
	}
	if variant.SellerOrderId != req.SellerOrderId || variant.SellerOrderItemId != req.SellerOrderItemId || variant.BuyerReference != req.BuyerCode {
		return fmt.Errorf("order details do not match for variant: %v", variant.ID)
	}
	if variant.Quantity < req.Quantity {
		return fmt.Errorf("returned quantity must be less than or equal to ordered quantity for variant: %v", variant.ID)
	}
	return nil
}
