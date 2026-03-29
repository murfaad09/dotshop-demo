package dto

import (
	"time"

	domain "github.com/harishash/dotshop-be/internal/models"
)

type SigninUserRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required,email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
}

type SingupUserRequest struct {
	FirstName string  `json:"firstName" validate:"required,min=3,max=100"`
	LastName  string  `json:"lastName" validate:"required,min=3,max=100"`
	Email     string  `json:"email" validate:"required,email"`
	Password  *string `json:"password"`
	RoleID    uint    `json:"roleID" validate:"required"`
	Username  *string `json:"username"`
}

type ProfileResponse struct {
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            string    `gorm:"uniqueIndex" json:"email"`
	AuthProviderType string    `json:"auth_provider"` // e.g., "google", "apple", "facebook"
	RoleID           uint      `gorm:"not null" json:"role_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewAddressDS(address *domain.ShippingInfo) *AddConsumerAddressResponse {
	return &AddConsumerAddressResponse{
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
}

type SendEmailForgotPasswordRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_ame"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
