package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/harishash/dotshop-be/integration/aws"
	constants "github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	user_repo "github.com/harishash/dotshop-be/internal/repositories"
	auth "github.com/harishash/dotshop-be/internal/utils/auth"
	"gorm.io/gorm"
)

type IUserService interface {
	GetUserByEmail(email string) (*domain.User, error)
	ValidateUser(email, password string) (*domain.User, *uint, error)
	GetUsers() ([]*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
	RoleExists(roleID uint) (bool, error)
	GetUserByID(id uint) (*domain.User, error)
	GetProfile(email string) (*dto.ProfileResponse, error)
	UpdateConsumerProfile(userId uint64, body *dto.UpdateConsumerProfileRequest) (*domain.User, error)
	AddConsumerAddress(userId uint64, body *dto.AddConsumerAddressRequest) (*dto.AddConsumerAddressResponse, error)
	GetAllUserAddresses(userId uint) ([]dto.AddConsumerAddressResponse, error)
	GetUserAddressById(userId, addressId uint64) (*dto.AddConsumerAddressResponse, error)
	ConsumerOrdersList(userId uint64) ([]dto.OrdersListResponse, error)
	UpdateUserAddress(userId, addressId uint64, body *dto.UpdateConsumerAddressRequest) (*dto.AddConsumerAddressResponse, error)
	DeleteAddress(userId, addressId uint64) error
	SendForgotPasswordEmail(email string) error
	ForgotPassword(body *dto.ForgotPasswordRequest, userId uint64) (*dto.UpdatePasswordResponse, error)
}

type UserService struct {
	awsService        aws.AWSService
	userRepository    user_repo.IUserRepository
	orderRepository   user_repo.IOrderRepository
	productRepository user_repo.IProductRepository
}

func NewUserService(awsService aws.AWSService, repo user_repo.IUserRepository, OrderRepository user_repo.IOrderRepository, productRepository user_repo.IProductRepository) IUserService {
	return &UserService{awsService: awsService, userRepository: repo, orderRepository: OrderRepository, productRepository: productRepository}
}

func (s *UserService) GetUsers() ([]*domain.User, error) {
	// Call repository to fetch users
	return s.userRepository.GetUsers()
}

func (s *UserService) ValidateUser(email, password string) (*domain.User, *uint, error) {
	return s.userRepository.ValidateUser(email, password)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) CreateUser(user domain.User) (*domain.User, error) {
	existingUser, err := s.userRepository.GetUserByEmailUnscoped(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingUser != nil {
		if existingUser.DeletedAt.Valid {
			user.ID = existingUser.ID
			updatedUser, err := s.userRepository.UpdateUserProfileFields(user)
			if err != nil {
				return nil, err
			}

			return updatedUser, nil
		} else {
			return nil, errors.New(constants.UserAlreadyExists.String())
		}
	}

	return s.userRepository.CreateUser(user)
}

func (s *UserService) RoleExists(roleID uint) (bool, error) {
	return s.userRepository.RoleExists(roleID)
}

func (s *UserService) GetProfile(email string) (*dto.ProfileResponse, error) {
	user, err := s.userRepository.GetProfileWithEmail(email)
	if err != nil {
		return nil, err
	}

	profile := &dto.ProfileResponse{
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Email:            user.Email,
		AuthProviderType: user.AuthProviderType,
		RoleID:           user.RoleID,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	return profile, nil
}

func (s *UserService) UpdateConsumerProfile(userId uint64, body *dto.UpdateConsumerProfileRequest) (*domain.User, error) {
	user, err := s.userRepository.GetUserByID(uint(userId))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	emailExist, err := s.userRepository.CheckEmailExists(body.Email)
	if err != nil {
		return nil, err
	}

	if emailExist != nil && emailExist.ID != uint(userId) {
		return nil, errors.New("email already exists")
	}

	user.FirstName = &body.FirstName
	user.LastName = &body.LastName
	user.Email = body.Email
	user.DOB = &body.DOB
	user.PhoneNumber = &body.PhoneNumber
	user.ID = uint(userId)
	user.UpdatedAt = time.Now().UTC()
	if err := s.userRepository.UpdateConsumerProfile(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) AddConsumerAddress(userId uint64, body *dto.AddConsumerAddressRequest) (*dto.AddConsumerAddressResponse, error) {
	user, err := s.userRepository.GetUserByID(uint(userId))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if body.DefaultBillingAddress || body.DefaultShippingAddress {
		addresses, err := s.userRepository.GetAllUserAddresses(uint(userId))
		if err != nil {
			return nil, err
		}

		updates := make(map[string]interface{})
		for _, address := range addresses {
			if body.DefaultBillingAddress {
				updates["default_billing"] = false
			}

			if body.DefaultShippingAddress {
				updates["default_address"] = false
			}

			if err := s.userRepository.PatchUserAddress(address.ID, updates); err != nil {
				return nil, err
			}
		}
	}

	address := &domain.ShippingInfo{
		UserId:         uint(userId),
		FirstName:      &body.FirstName,
		LastName:       &body.LastName,
		AddressOne:     sql.NullString{String: body.Address, Valid: true},
		City:           sql.NullString{String: body.City, Valid: true},
		State:          sql.NullString{String: body.State, Valid: true},
		Country:        sql.NullString{String: body.Country, Valid: true},
		Zip:            sql.NullString{String: body.PostCode, Valid: true},
		PhoneNumber:    sql.NullString{String: body.PhoneNumber, Valid: true},
		DefaultAddress: body.DefaultShippingAddress,
		DefaultBilling: body.DefaultBillingAddress,
	}

	if err := s.userRepository.AddConsumerAddress(address); err != nil {
		return nil, err
	}

	resp := &dto.AddConsumerAddressResponse{
		Id:                     address.ID,
		UserId:                 address.UserId,
		FirstName:              address.FirstName,
		LastName:               address.LastName,
		Address:                address.AddressOne.String,
		City:                   address.City.String,
		State:                  address.State.String,
		Country:                address.Country.String,
		PostCode:               address.Zip.String,
		PhoneNumber:            address.PhoneNumber.String,
		DefaultShippingAddress: address.DefaultAddress,
		DefaultBillingAddress:  address.DefaultBilling,
		CreatedAt:              address.CreatedAt,
		UpdatedAt:              address.UpdatedAt,
	}

	return resp, nil
}

func (s *UserService) GetAllUserAddresses(userId uint) ([]dto.AddConsumerAddressResponse, error) {
	user, err := s.userRepository.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	addresses, err := s.userRepository.GetAllUserAddresses(userId)
	if err != nil {
		return nil, err
	}

	var resp []dto.AddConsumerAddressResponse
	for _, v := range addresses {
		resp = append(resp, dto.AddConsumerAddressResponse{
			Id:                     v.ID,
			UserId:                 v.UserId,
			FirstName:              v.FirstName,
			LastName:               v.LastName,
			Address:                v.AddressOne.String,
			City:                   v.City.String,
			State:                  v.State.String,
			Country:                v.Country.String,
			PostCode:               v.Zip.String,
			PhoneNumber:            v.PhoneNumber.String,
			DefaultShippingAddress: v.DefaultAddress,
			DefaultBillingAddress:  v.DefaultBilling,
			CreatedAt:              v.CreatedAt,
			UpdatedAt:              v.UpdatedAt,
		})

	}

	return resp, nil
}

func (s *UserService) GetUserAddressById(userId, addressId uint64) (*dto.AddConsumerAddressResponse, error) {
	user, err := s.userRepository.GetUserByID(uint(userId))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	address, err := s.userRepository.GetUserAddressByID(addressId)
	if err != nil {
		return nil, err
	}

	resp := &dto.AddConsumerAddressResponse{
		Id:                     address.ID,
		UserId:                 address.UserId,
		FirstName:              address.FirstName,
		LastName:               address.LastName,
		Address:                address.AddressOne.String,
		City:                   address.City.String,
		State:                  address.State.String,
		Country:                address.Country.String,
		PostCode:               address.Zip.String,
		PhoneNumber:            address.PhoneNumber.String,
		DefaultShippingAddress: address.DefaultAddress,
		DefaultBillingAddress:  address.DefaultBilling,
		CreatedAt:              address.CreatedAt,
		UpdatedAt:              address.UpdatedAt,
	}

	return resp, nil
}

func (h *UserService) ConsumerOrdersList(userId uint64) ([]dto.OrdersListResponse, error) {
	var response []dto.OrdersListResponse

	orders, err := h.orderRepository.GetOrdersByUserId(userId)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		orderVariants, err := h.orderRepository.GetOrderVariantByOrderID(order.ID)
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

		resp := dto.OrdersListResponse{
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
		}
		response = append(response, resp)
	}

	return response, nil
}

func (s *UserService) UpdateUserAddress(userId, addressId uint64, body *dto.UpdateConsumerAddressRequest) (*dto.AddConsumerAddressResponse, error) {
	user, err := s.userRepository.GetUserByID(uint(userId))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	address, err := s.userRepository.GetUserAddressByID(addressId)
	if err != nil {
		return nil, err
	}

	if address.FirstName != nil {
		address.FirstName = &body.FirstName
	}

	if address.LastName != nil {
		address.LastName = &body.LastName
	}

	if len(body.Address) > 0 {
		address.AddressOne = sql.NullString{String: body.Address, Valid: true}
	}

	if len(body.City) > 0 {
		address.City = sql.NullString{String: body.City, Valid: true}
	}

	if len(body.State) > 0 {
		address.State = sql.NullString{String: body.State, Valid: true}
	}
	if len(body.Country) > 0 {
		address.Country = sql.NullString{String: body.Country, Valid: true}
	}
	if len(body.PostCode) > 0 {
		address.Zip = sql.NullString{String: body.PostCode, Valid: true}
	}
	if len(body.PhoneNumber) > 0 {
		address.PhoneNumber = sql.NullString{String: body.PhoneNumber, Valid: true}
	}

	if body.DefaultShippingAddress != nil {
		address.DefaultAddress = *body.DefaultShippingAddress
	}

	if body.DefaultBillingAddress != nil {
		address.DefaultBilling = *body.DefaultBillingAddress
	}

	if body.DefaultBillingAddress != nil || body.DefaultShippingAddress != nil {
		addresses, err := s.userRepository.GetAllUserAddresses(user.ID)
		if err != nil {
			return nil, err
		}

		updates := make(map[string]interface{})
		for _, address := range addresses {
			if body.DefaultBillingAddress != nil {
				if *body.DefaultBillingAddress {
					updates["default_billing"] = false
				}
			}

			if body.DefaultShippingAddress != nil {
				if *body.DefaultShippingAddress {
					updates["default_address"] = false
				}
			}

			if err := s.userRepository.PatchUserAddress(address.ID, updates); err != nil {
				return nil, err
			}
		}
	}

	address.UpdatedAt = time.Now().UTC()
	if err := s.userRepository.UpdateUserAddress(address); err != nil {
		return nil, err
	}

	resp := &dto.AddConsumerAddressResponse{
		Id:                     address.ID,
		UserId:                 address.UserId,
		FirstName:              address.FirstName,
		LastName:               address.LastName,
		Address:                address.AddressOne.String,
		City:                   address.City.String,
		State:                  address.State.String,
		Country:                address.Country.String,
		PostCode:               address.Zip.String,
		PhoneNumber:            address.PhoneNumber.String,
		DefaultShippingAddress: address.DefaultAddress,
		DefaultBillingAddress:  address.DefaultBilling,
		CreatedAt:              address.CreatedAt,
		UpdatedAt:              address.UpdatedAt,
	}

	return resp, nil
}

func (s *UserService) DeleteAddress(userId, addressId uint64) error {
	return s.userRepository.DeleteAddress(userId, addressId)
}

func (s *UserService) SendForgotPasswordEmail(email string) error {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if user != nil {
		tokenString, err := auth.GenerateJWTForForgotPassword(user.Email, uint64(user.RoleID), uint64(user.ID))
		if err != nil {
			return errors.New(constants.FailedTokenCreation.String())
		}

		if err := s.awsService.ForgotPasswordMail(user.Email, *user.FirstName, tokenString); err != nil {
			return fmt.Errorf("failed to send forgotpassword mail: %v", err)
		}
	}

	return nil
}

func (s *UserService) ForgotPassword(body *dto.ForgotPasswordRequest, userId uint64) (*dto.UpdatePasswordResponse, error) {
	return s.userRepository.UpdatePasswordByID(userId, hashPassword(body.NewPassword))
}
