package dto

import "time"

type UpdatePasswordRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword" validate:"required,password" binding:"required,password"`
	NewPassword     string `json:"newPassword" validate:"required,password" binding:"required,password"`
}

type UpdatePasswordResponse struct {
	Success bool
	Message string
}

type UpdateProfileRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Bio             string `json:"bio"`
	CoverImageURL   string `json:"coverImageURL"`
	ProfileImageURL string `json:"profileImageURL"`
}

type UpdateProfileResponse struct {
	UserId          uint64 `json:"user_id"`
	CuratorId       uint64 `json:"curator_id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Bio             string `json:"bio"`
	CoverImageURL   string `json:"coverImageURL"`
	ProfileImageURL string `json:"profileImageURL"`
}

type UpdateConsumerProfileRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DOB         string `json:"dob"`
	PhoneNumber string `json:"phoneNumber"`
}

type UpdateConsumerProfileResponse struct {
	AccessToken string `json:"access_token"`
	UserId      uint64 `json:"userId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DOB         string `json:"dob"`
	PhoneNumber string `json:"phoneNumber"`
}

type AddConsumerAddressRequest struct {
	FirstName              string `json:"firstName"`
	LastName               string `json:"lastName"`
	Address                string `json:"address"`
	City                   string `json:"city"`
	State                  string `json:"state"`
	Country                string `json:"country"`
	PostCode               string `json:"postCode"`
	PhoneNumber            string `json:"phoneNumber"`
	DefaultShippingAddress bool   `json:"defaultShippingAddress"`
	DefaultBillingAddress  bool   `json:"defaultBillingAddress"`
}

type AddConsumerAddressResponse struct {
	Id                     uint      `json:"id"`
	UserId                 uint      `json:"userId"`
	FirstName              *string   `json:"firstName"`
	LastName               *string   `json:"lastName"`
	Address                string    `json:"address"`
	City                   string    `json:"city"`
	State                  string    `json:"state"`
	Country                string    `json:"country"`
	PostCode               string    `json:"postCode"`
	PhoneNumber            string    `json:"phoneNumber"`
	DefaultShippingAddress bool      `json:"defaultShippingAddress"`
	DefaultBillingAddress  bool      `json:"defaultBillingAddress"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type UpdateConsumerAddressRequest struct {
	FirstName              string `json:"firstName"`
	LastName               string `json:"lastName"`
	Address                string `json:"address"`
	City                   string `json:"city"`
	State                  string `json:"state"`
	Country                string `json:"country"`
	PostCode               string `json:"postCode"`
	PhoneNumber            string `json:"phoneNumber"`
	DefaultShippingAddress *bool  `json:"defaultShippingAddress"`
	DefaultBillingAddress  *bool  `json:"defaultBillingAddress"`
}

type ForgotPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,password" binding:"required,password"`
}
