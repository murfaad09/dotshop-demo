package dto

import (
	"time"

	"github.com/harishash/dotshop-be/internal/utils/errors"

	domain "github.com/harishash/dotshop-be/internal/models"
)

type SocialMedia string

const (
	Instagram SocialMedia = "instagram"
	Youtube   SocialMedia = "youtube"
	TikTok    SocialMedia = "tiktok"
	Other     SocialMedia = "other"
)

type SocialMediaLinks struct {
	Youtube   string `json:"1"`
	TikTok    string `json:"2"`
	Instagram string `json:"3"`
}

type CuratorDetails struct {
	CoverImage        string           `json:"coverImage"`
	ProfileImage      string           `json:"profileImage"`
	CuratorName       string           `json:"curatorName"`
	Description       string           `json:"description"`
	SocialMediaLinks  SocialMediaLinks `json:"socialMediaLinks"`
	NumberOfFollowers string           `json:"numberOfFollowers"`
}

type ShopMyLook struct {
	ProductID    string `json:"productId"`
	ProductImage string `json:"productImage"`
}

type CuratorStore struct {
	StoreCuratorDetails CuratorDetails   `json:"storeCuratorDetails"`
	FeaturedProducts    []ProductRequest `json:"featuredProducts"`
	Collections         []ProductRequest `json:"collections"`
	ShopMyLook          []ShopMyLook     `json:"shopMyLook"`
}

type CuratorOnBoardingRequest struct {
	FirstName    string `json:"firstName" validate:"required,min=3,max=100"`
	LastName     string `json:"lastName" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,password"`
	Bio          string `json:"bio"`
	ShopName     string `json:"shopName" validate:"required" binding:"required"`
	ProfileImage string `json:"profileImage"`
	CoverImage   string `json:"coverImage"`

	Socials         []*SocialMediaLinksRequest
	BankInformation []*BankInformationRequest
}

type CuratorOnBoardingResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName" validate:"required,min=3,max=100"`
	LastName  string `json:"lastName" validate:"required,min=3,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Bio       string `json:"bio"`
	ShopName  string `json:"shopName" validate:"required" binding:"required"`
	Status    string `json:"status"`

	ProfileImage string `json:"profileImage"`
	CoverImage   string `json:"coverImage"`

	Socials         []*SocialMediaLinksRequest
	BankInformation []*BankInformationRequest
}

type SocialMediaLinksRequest struct {
	Platform     string `json:"platform"`
	Link         string `json:"link"`
	AccessToken  string `json:"accessToken"`
	OpenID       string `json:"openID"`
	RefreshToken string `json:"refreshToken"`
}

type BankInformationRequest struct {
	Location    string `json:"location"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`

	AccountDetails AccountDetailsRequest
}

type AccountDetailsRequest struct {
	BankName       string `json:"bankName"`
	BankAddress    string `json:"bankAddress"`
	BranchCode     string `json:"branchCode"`
	IBAN           string `json:"IBAN"`
	AccountNumber  string `json:"accountNumber"`
	AccountName    string `json:"accountName"`
	AccountAddress string `json:"accountAddress"`
}

func (r *CuratorOnBoardingRequest) Validate() *errors.Error {
	if r.FirstName == "" {
		return errors.New("firstName is required")
	}
	if r.LastName == "" {
		return errors.New("lastName is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.ShopName == "" {
		return errors.New("shopName is required")
	}
	return nil
}

type GetAllCuratorsResponse struct {
	Curators []GetCuratorResponse `json:"curators"`
}

type SocialMediaLinkResponse struct {
	ID           uint                   `json:"id"`
	Type         domain.SocialMediaType `json:"type"`
	URL          string                 `json:"url"`
	AccessToken  string                 `json:"accessToken"`
	OpenID       string                 `json:"openID"`
	RefreshToken string                 `json:"refreshToken"`
	CreatedAt    time.Time              `json:"createdAt"`
}

type GetCuratorResponse struct {
	ID                uint                       `json:"id"`
	UserID            uint                       `json:"userID"`
	ShopName          string                     `json:"shopName"`
	Bio               string                     `json:"bio"`
	FirstName         *string                    `json:"firstName"`
	LastName          *string                    `json:"lastName"`
	Email             string                     `json:"email"`
	NumberofFollowers uint                       `json:"numberOfFollowers"`
	ProfileImageURL   string                     `json:"profileImageURL"`
	CoverImageURL     string                     `json:"coverImageURL"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
	VariantImages     []string                   `json:"variantImages"`
	SocialMediaLinks  []*SocialMediaLinkResponse `json:"socialMediaLinks"`
}

type GetStoreCuratorResponse struct {
	ID                uint      `json:"id"`
	UserID            uint      `json:"userID"`
	ShopName          string    `json:"shopName"`
	Bio               string    `json:"bio"`
	FirstName         *string   `json:"firstName"`
	LastName          *string   `json:"lastName"`
	Email             string    `json:"email"`
	NumberofFollowers uint      `json:"numberOfFollowers"`
	ProfileImageURL   string    `json:"profileImageURL"`
	CoverImageURL     string    `json:"coverImageURL"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	VariantImages     []string  `json:"variantImages"`
}

func NewGetCuratorResponse(curator domain.Curator) GetCuratorResponse {
	socialMediaLinks := make([]*SocialMediaLinkResponse, len(curator.SocialMediaLinks))
	for i, link := range curator.SocialMediaLinks {
		modifiedAccessToken := link.AccessToken
		if len(modifiedAccessToken) > 0 {
			modifiedAccessToken = modifiedAccessToken[4:] + modifiedAccessToken[:4]
		}

		socialMediaLinks[i] = &SocialMediaLinkResponse{
			ID:           link.ID,
			Type:         link.Type,
			URL:          link.URL,
			AccessToken:  modifiedAccessToken,
			OpenID:       link.OpenID,
			RefreshToken: link.RefreshToken,
			CreatedAt:    link.CreatedAt,
		}
	}
	if len(socialMediaLinks) == 0 {
		socialMediaLinks = nil
	}
	return GetCuratorResponse{
		ID:                curator.ID,
		UserID:            curator.UserID,
		ShopName:          curator.ShopName,
		Bio:               curator.Bio,
		FirstName:         curator.User.FirstName,
		LastName:          curator.User.LastName,
		Email:             curator.User.Email,
		NumberofFollowers: curator.NumberofFollowers,
		ProfileImageURL:   curator.ProfileImageURL,
		CoverImageURL:     curator.CoverImageURL,
		CreatedAt:         curator.CreatedAt,
		UpdatedAt:         curator.UpdatedAt,
		SocialMediaLinks:  socialMediaLinks,
	}
}

func NewGetStoreCuratorResponse(curator domain.Curator) GetStoreCuratorResponse {

	return GetStoreCuratorResponse{
		ID:                curator.ID,
		UserID:            curator.UserID,
		ShopName:          curator.ShopName,
		Bio:               curator.Bio,
		FirstName:         curator.User.FirstName,
		LastName:          curator.User.LastName,
		Email:             curator.User.Email,
		NumberofFollowers: curator.NumberofFollowers,
		ProfileImageURL:   curator.ProfileImageURL,
		CoverImageURL:     curator.CoverImageURL,
		CreatedAt:         curator.CreatedAt,
		UpdatedAt:         curator.UpdatedAt,
	}
}

type ShopNameResponse struct {
	Success  bool   `json:"success"`
	Property string `json:"property"`
	Message  string `json:"message"`
}

type CreateSocialMediaLinkResponse struct {
	ID           uint                   `json:"id"`
	Type         domain.SocialMediaType `json:"type"`
	URL          string                 `json:"url"`
	AccessToken  string                 `json:"accessToken"`
	OpenID       string                 `json:"openID"`
	RefreshToken string                 `json:"refreshToken"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type DeleteSocialMediaLinkResponse struct {
	Message string `json:"message"`
}

type AccountDetailResponse struct {
	ID             uint      `json:"id"`
	CuratorID      uint      `json:"curatorID"`
	Location       string    `json:"location"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	DateOfBirth    string    `json:"dateOfBirth"`
	BankAddress    string    `json:"bankAddress"`
	BankName       string    `json:"bankName"`
	BranchCode     string    `json:"branchCode"`
	AccountNumber  string    `json:"accountNumber"`
	AccountName    string    `json:"accountName"`
	AccountAddress string    `json:"accountAddress"`
	IBAN           string    `json:"IBAN"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type UpdateAccountDetailRequest struct {
	Location       string `json:"location"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DateOfBirth    string `json:"dateOfBirth"`
	BankAddress    string `json:"bankAddress"`
	BankName       string `json:"bankName"`
	BranchCode     string `json:"branchCode"`
	AccountNumber  string `json:"accountNumber"`
	AccountName    string `json:"accountName"`
	AccountAddress string `json:"accountAddress"`
	IBAN           string `json:"IBAN"`
}
