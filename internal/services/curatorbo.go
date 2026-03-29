package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/harishash/dotshop-be/integration/aws"
	"gorm.io/gorm"

	dto "github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	curatorbo_repo "github.com/harishash/dotshop-be/internal/repositories"
	repo "github.com/harishash/dotshop-be/internal/repositories"
	"github.com/harishash/dotshop-be/internal/utils/logger"
)

// CuratorBOService interface
type CuratorBOService interface {
	AddProduct(products *dto.CreateFeatureProductRequest) ([]*dto.CreateFeatureProductResponse, error)
	AddCollection(collection *dto.CreateCollectionRequest) (*dto.CreateCollectionResponse, error)
	AddCollectionSection(collectionSection *dto.CreateCollectionSectionRequest) (*dto.CreateCollectionSectionResponse, error)
	AddProductToSection(sectionID uint, products *dto.AddProductToSectionRequest) (*dto.AddProductToSectionResponse, error)
	UpdateCollectionSection(body *dto.UpdateCollectionSectionRequest, sectionID uint) (*dto.UpdateCollectionSectionResponse, error)
	DeleteProductFromSectionByID(sectionID uint, productID string) error

	AddLook(look *dto.CreateLookRequest) (*dto.CreateLookResponse, error)
	DeleteFromFeatureProduct(curatorID uint, productID string) error
	DeleteCollectionByID(collectionID uint) error
	DeleteProductFromCollectionByID(collectionID uint, productID string) error
	DeleteLookByID(lookID uint) error
	DeleteSectionByID(sectionID uint) error
	DeleteProductFromLookByID(lookID uint, productID string) error
	// GetOrdersStatus() ([]models.OrderStatus, error)
	// GetPayoutDetails() ([]models.PayoutDetail, error)
	Withdraw(data []byte) error
	// GetProfile() (models.Profile, error)
	UpdateProfile(userId uint64, body *dto.UpdateProfileRequest) (*domain.Curator, *domain.User, error)
	AddSocialMediaLink(curatorID uint, link *dto.SocialMediaLinksRequest) (*dto.CreateSocialMediaLinkResponse, error)
	CheckSocialMediaLinkExists(linkType string, curatorID uint) (bool, error)
	// not used
	UpdateSocialMediaLink(curatorID, linkID uint, link *dto.SocialMediaLinksRequest) (*dto.CreateSocialMediaLinkResponse, error)
	UpdateSocialMediaLinks(request *dto.SocialMediaLinksRequest, curatorID uint) (*dto.CreateSocialMediaLinkResponse, error)
	RemoveSocialMediaLink(curatorID, linkID uint64) (*dto.DeleteSocialMediaLinkResponse, error)
	ChangePassword(body *dto.UpdatePasswordRequest) (error, *dto.UpdatePasswordResponse)
	GetAllCurators() ([]*dto.GetCuratorResponse, error)
	GetCuratorByCuratorID(id uint64) (*domain.Curator, error)
	AddProductToFeature(body *dto.AddProductToFeatureRequest, featureId uint, userId uint64) ([]*dto.AddFeatureProductResponse, error)
	AddProductToLook(body *dto.AddProductToLookRequest, lookId uint, userId uint64) ([]*dto.AddLookProductResponse, error)
	AddProductToCollection(body *dto.AddProductToCollectionRequest, collectionId uint, userId uint64) ([]*dto.AddCollectionProductResponse, error)
	UpdateCollection(body *dto.UpdateCollectionRequest, collectionId uint, userId uint64) (*dto.UpdateCollectionResponse, error)
	UpdateLook(body *dto.UpdateLookRequest, collectionId uint, userId uint64) (*dto.UpdateLookResponse, error)
	GetCuratorAccountDetail(curatorId uint) (*dto.AccountDetailResponse, error)
	UpdateCuratorAccountDetail(curatorId uint, request *dto.UpdateAccountDetailRequest) (*dto.AccountDetailResponse, error)
}

// curatorBOService struct
type curatorBOService struct {
	repo             curatorbo_repo.CuratorBORepo
	awsService       aws.AWSService
	notificationRepo repo.NotificationRepositoryInterface
}

// NewCuratorBOService creates a new CuratorBOService
func NewCuratorBOService(repo curatorbo_repo.CuratorBORepo, awsService aws.AWSService, notificationRepo repo.NotificationRepositoryInterface) CuratorBOService {
	return &curatorBOService{repo: repo, awsService: awsService, notificationRepo: notificationRepo}
}

func (s *curatorBOService) AddProduct(products *dto.CreateFeatureProductRequest) ([]*dto.CreateFeatureProductResponse, error) {
	resp, err := s.repo.InsertProduct(products)
	if err != nil {
		return nil, err
	}
	curator, err := s.repo.GetCuratorByCuratorID(uint64(products.CuratorID))
	if err != nil {
		return nil, err
	}

	for _, product := range resp {
		message := fmt.Sprintf("A new Feature Product #%s has been Added by %s", product.ProductID, curator.Name)
		if err := s.notificationRepo.NotifyAdmins(message); err != nil {
			logger.Warnf("warning: failed to notify admin: %v", err)
		}
	}

	return resp, nil
}

func (s *curatorBOService) AddCollection(collection *dto.CreateCollectionRequest) (*dto.CreateCollectionResponse, error) {
	resp, err := s.repo.InsertCollection(collection)
	if err != nil {
		return nil, err
	}

	curator, err := s.repo.GetCuratorByCuratorID(uint64(collection.CuratorID))
	if err == nil {
		if err := s.awsService.CreateCollectionMail(curator.User.Email, *curator.User.FirstName, resp.Name); err != nil {
			log.Warnf("failed to send collection mail: %v", err)
		}
	}

	return resp, nil
}

func (s *curatorBOService) AddCollectionSection(collectionSection *dto.CreateCollectionSectionRequest) (*dto.CreateCollectionSectionResponse, error) {
	return s.repo.InsertCollectionSection(collectionSection)
}

func (s *curatorBOService) AddProductToSection(
	sectionID uint,
	products *dto.AddProductToSectionRequest) (
	*dto.AddProductToSectionResponse,
	error) {
	return s.repo.InsertProductToSection(sectionID, products)
}

func (s *curatorBOService) UpdateCollectionSection(
	body *dto.UpdateCollectionSectionRequest,
	sectionID uint) (
	*dto.UpdateCollectionSectionResponse,
	error) {
	return s.repo.UpdateCollectionSection(body, sectionID)
}

func (s *curatorBOService) DeleteProductFromSectionByID(sectionID uint, productID string) error {
	return s.repo.DeleteProductFromSection(sectionID, productID)
}

func (s *curatorBOService) AddLook(look *dto.CreateLookRequest) (*dto.CreateLookResponse, error) {

	resp, err := s.repo.InsertLook(look)
	if err != nil {
		return nil, err
	}

	curator, err := s.repo.GetCuratorByCuratorID(uint64(look.CuratorID))
	if err == nil {
		if err := s.awsService.NewLookCreatedMail(curator.User.Email, *curator.User.FirstName, resp.Name); err != nil {
			log.Warnf("failed to send look mail: %v", err)
		}
	}
	return resp, nil
}

func (s *curatorBOService) DeleteFromFeatureProduct(curatorID uint, productID string) error {
	if err := s.repo.DeleteFromFeatureProduct(curatorID, productID); err != nil {
		return err
	}
	return nil
}
func (s *curatorBOService) DeleteCollectionByID(collectionID uint) error {
	err := s.repo.DeleteCollectionByID(collectionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorBOService) DeleteProductFromCollectionByID(collectionID uint, productID string) error {
	if err := s.repo.DeleteProductFromCollection(collectionID, productID); err != nil {
		return err
	}
	return nil
}

func (s *curatorBOService) DeleteLookByID(lookID uint) error {
	err := s.repo.DeleteLookByID(lookID)
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorBOService) DeleteSectionByID(sectionID uint) error {
	err := s.repo.DeleteSectionByID(sectionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorBOService) DeleteProductFromLookByID(lookID uint, productID string) error {

	if err := s.repo.DeleteProductFromLook(lookID, productID); err != nil {
		return err
	}
	return nil
}

// func (s *curatorBOService) GetOrdersStatus() ([]models.OrderStatus, error) {
//     return s.repo.FetchOrdersStatus()
// }

// func (s *curatorBOService) GetPayoutDetails() ([]models.PayoutDetail, error) {
//     return s.repo.FetchPayoutDetails()
// }

func (s *curatorBOService) Withdraw(data []byte) error {
	return s.repo.ProcessWithdraw(data)
}

// func (s *curatorBOService) GetProfile() (models.Profile, error) {
//     return s.repo.FetchProfile()
// }

func (s *curatorBOService) UpdateProfile(userId uint64, body *dto.UpdateProfileRequest) (*domain.Curator, *domain.User, error) {
	curator, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		return nil, nil, err
	}

	curatorDetails, err := s.repo.GetCuratorByID(uint64(curator.ID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch curator from database: %v", err)
	}

	if curatorDetails == nil {
		return nil, nil, errors.New("curator not found")
	}

	if len(body.Bio) > 0 {
		curatorDetails.Bio = body.Bio
	}

	if len(body.CoverImageURL) > 0 {
		curatorDetails.CoverImageURL = body.CoverImageURL
	}

	if len(body.ProfileImageURL) > 0 {
		curatorDetails.ProfileImageURL = body.ProfileImageURL
	}

	if len(body.FirstName) > 0 {
		user.FirstName = &body.FirstName
	}

	if len(body.LastName) > 0 {
		user.LastName = &body.LastName
	}

	curatorDetails.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateProfile(curatorDetails, user); err != nil {
		return nil, nil, err
	}

	return curatorDetails, user, nil
}

// AddSocialMediaLink implements the SocialMediaLinkService interface.
func (s *curatorBOService) AddSocialMediaLink(curatorID uint, link *dto.SocialMediaLinksRequest) (*dto.CreateSocialMediaLinkResponse, error) {

	linkRequest := &domain.SocialMediaLink{
		Type:         domain.SocialMediaType(link.Platform),
		AccessToken:  link.AccessToken,
		OpenID:       link.OpenID,
		CuratorID:    curatorID,
		RefreshToken: link.RefreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateSocialMediaLink(linkRequest)
}

// not used
func (s *curatorBOService) UpdateSocialMediaLink(
	curatorID, linkID uint,
	link *dto.SocialMediaLinksRequest) (
	*dto.CreateSocialMediaLinkResponse,
	error) {

	linkRequest := &domain.SocialMediaLink{
		Type:         domain.SocialMediaType(link.Platform),
		AccessToken:  link.AccessToken,
		OpenID:       link.OpenID,
		CuratorID:    curatorID,
		RefreshToken: link.RefreshToken,
		UpdatedAt:    time.Now(),
	}

	return s.repo.UpdateSocialMediaLink(linkRequest, linkID, curatorID)
}

func (s *curatorBOService) CheckSocialMediaLinkExists(linkType string, curatorID uint) (bool, error) {
	_, err := s.repo.GetSocialMediaLinkByTypeAndCuratorID(linkType, curatorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *curatorBOService) UpdateSocialMediaLinks(request *dto.SocialMediaLinksRequest, curatorID uint) (*dto.CreateSocialMediaLinkResponse, error) {
	link, err := s.repo.GetSocialMediaLinkByTypeAndCuratorID(request.Platform, curatorID)
	if err != nil {
		return nil, err
	}

	link.AccessToken = request.AccessToken
	link.OpenID = request.OpenID
	link.RefreshToken = request.RefreshToken

	return s.repo.UpdateSocialMediaLink(link, link.ID, curatorID)
}

// RemoveSocialMediaLink implements the SocialMediaLinkService interface.
func (s *curatorBOService) RemoveSocialMediaLink(curatorID, linkID uint64) (*dto.DeleteSocialMediaLinkResponse, error) {
	return s.repo.DeleteSocialMediaLink(curatorID, linkID)
}
func (s *curatorBOService) ChangePassword(body *dto.UpdatePasswordRequest) (error, *dto.UpdatePasswordResponse) {
	user, err := s.repo.GetUserByEmail(body.Email)
	if err != nil {
		return fmt.Errorf("failed to fetch user from database: %v", err), nil
	}

	if user == nil {
		return errors.New("email not found"), nil
	}

	currentPassword := hashPassword(body.CurrentPassword)
	if *user.PasswordHash != currentPassword {
		return fmt.Errorf("current password does not match the user's actual password"), nil
	}

	newPassword := hashPassword(body.NewPassword)
	body.NewPassword = newPassword

	if newPassword == currentPassword {
		return errors.New("the new password must be different from the current password"), nil
	}

	err, res := s.repo.UpdatePasswordByEmail(body.Email, body.NewPassword)
	return err, res
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *curatorBOService) GetAllCurators() ([]*dto.GetCuratorResponse, error) {
	curators, err := s.repo.GetAllCurators()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curators from database: %v", err)
	}

	if curators == nil {
		return nil, errors.New("curators not found")
	}

	return curators, nil
}

func (s *curatorBOService) GetCuratorByCuratorID(id uint64) (*domain.Curator, error) {
	curators, err := s.repo.GetCuratorByCuratorID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curator from database: %v", err)
	}

	return curators, nil
}

func (s *curatorBOService) AddProductToFeature(body *dto.AddProductToFeatureRequest, featureId uint, userId uint64) ([]*dto.AddFeatureProductResponse, error) {
	curatorId, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	// Check if the feature exists
	feature, err := s.repo.GetFeatureWithId(featureId, curatorId.ID)
	if err != nil {
		return nil, err
	}
	if feature == nil {
		return nil, errors.New("feature not found")
	}

	var productsResponse []*dto.AddFeatureProductResponse

	// Use a single transaction for all operations
	tx, err := s.repo.StartTx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, product := range body.Products {
		exists, err := s.repo.CheckProductExistsInFeature(featureId, product.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if !exists {
			if _, err = s.repo.AddProductToFeature(tx, &domain.CuratorProduct{
				FeatureID: &featureId,
				ProductID: product.ProductID,
			}); err != nil {
				tx.Rollback()
				return nil, err
			}

			productsResponse = append(productsResponse, &dto.AddFeatureProductResponse{
				ProductID: product.ProductID,
			})
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("product already exists in feature")
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return productsResponse, nil
}

func (s *curatorBOService) AddProductToLook(body *dto.AddProductToLookRequest, lookId uint, userId uint64) ([]*dto.AddLookProductResponse, error) {
	curatorId, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	// Check if the look exists
	look, err := s.repo.GetLookWithId(lookId, curatorId.ID)
	if err != nil {
		return nil, err
	}
	if look == nil {
		return nil, errors.New("look not found")
	}

	var productsResponse []*dto.AddLookProductResponse

	// Use a single transaction for all operations
	tx, err := s.repo.StartTx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, product := range body.Products {
		exists, err := s.repo.CheckProductExistsInLook(lookId, product.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if !exists {
			if _, err = s.repo.AddProductToLook(tx, &domain.LookProduct{
				LookID:    lookId,
				ProductID: product.ProductID,
			}); err != nil {
				tx.Rollback()
				return nil, err
			}

			productsResponse = append(productsResponse, &dto.AddLookProductResponse{
				ProductID: product.ProductID,
			})
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("product already exists in look")
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return productsResponse, nil
}

func (s *curatorBOService) AddProductToCollection(body *dto.AddProductToCollectionRequest, collectionId uint, userId uint64) ([]*dto.AddCollectionProductResponse, error) {
	curatorId, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	// Check if the collection exists
	collection, err := s.repo.GetCollectionWithId(collectionId, curatorId.ID)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, errors.New("collection not found")
	}

	var productsResponse []*dto.AddCollectionProductResponse

	// Use a single transaction for all operations
	tx, err := s.repo.StartTx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, product := range body.Products {
		exists, err := s.repo.CheckProductExistsInCollection(collectionId, product.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if !exists {
			if _, err = s.repo.AddProductToCollection(tx, &domain.CollectionProduct{
				CollectionID: collectionId,
				ProductID:    product.ProductID,
			}); err != nil {
				tx.Rollback()
				return nil, err
			}

			productsResponse = append(productsResponse, &dto.AddCollectionProductResponse{
				ProductID: product.ProductID,
			})
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return productsResponse, nil
}

func (s *curatorBOService) UpdateCollection(body *dto.UpdateCollectionRequest, collectionId uint, userId uint64) (*dto.UpdateCollectionResponse, error) {
	curator, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	collection, err := s.repo.GetCollectionWithId(collectionId, curator.ID)
	if err != nil {
		return nil, err
	}

	if len(body.Name) > 0 {
		collection.Name = body.Name
	}

	if len(body.Description) > 0 {
		collection.Description = body.Description
	}

	if len(body.TileColor) > 0 {
		collection.TileColor = body.TileColor
	}

	if err := s.repo.UpdateCollection(collection); err != nil {
		return nil, err
	}

	resp := &dto.UpdateCollectionResponse{
		Id:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		TileColor:   collection.TileColor,
	}

	return resp, nil
}

func (s *curatorBOService) UpdateLook(body *dto.UpdateLookRequest, collectionId uint, userId uint64) (*dto.UpdateLookResponse, error) {
	curator, err := s.repo.GetCuratorWithUserId(userId)
	if err != nil {
		return nil, err
	}

	look, err := s.repo.GetLookWithId(collectionId, curator.ID)
	if err != nil {
		return nil, err
	}

	if len(body.Name) > 0 {
		look.Name = body.Name
	}

	if len(body.ImageURL) > 0 {
		look.ImageURL = body.ImageURL
	}

	if err := s.repo.UpdateLook(look); err != nil {
		return nil, err
	}

	resp := &dto.UpdateLookResponse{
		Id:               look.ID,
		Name:             look.Name,
		ImageURL:         look.ImageURL,
		SocialID:         look.SocialID,
		SocialType:       look.SocialType,
		SocialTitle:      look.SocialTitle,
		EmbedLink:        look.EmbedLink,
		VideoDescription: look.VideoDescription,
	}

	return resp, nil
}

func (s *curatorBOService) GetCuratorAccountDetail(curatorId uint) (*dto.AccountDetailResponse, error) {
	details, err := s.repo.GetCuratorAccountDetail(curatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch curators from database: %v", err)
	}

	if details == nil {
		return nil, errors.New("curators account detials not found")
	}

	resp := &dto.AccountDetailResponse{
		ID:             details.BankInformation.ID,
		CuratorID:      details.BankInformation.CuratorID,
		Location:       details.BankInformation.Location,
		FirstName:      details.BankInformation.FirstName,
		LastName:       details.BankInformation.LastName,
		DateOfBirth:    details.BankInformation.DateOfBirth,
		BankAddress:    details.AccountDetails.BankAddress,
		BankName:       details.AccountDetails.BankName,
		BranchCode:     details.AccountDetails.BranchCode,
		AccountNumber:  details.AccountDetails.AccountNumber,
		AccountName:    details.AccountDetails.AccountName,
		AccountAddress: details.AccountDetails.AccountAddress,
		IBAN:           details.AccountDetails.IBAN,
		CreatedAt:      details.BankInformation.CreatedAt,
		UpdatedAt:      details.BankInformation.UpdatedAt,
	}

	return resp, nil
}

func (s *curatorBOService) UpdateCuratorAccountDetail(curatorId uint, request *dto.UpdateAccountDetailRequest) (*dto.AccountDetailResponse, error) {
	details, err := s.repo.UpdateCuratorAccountDetail(curatorId, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update curator account details: %v", err)
	}

	resp := &dto.AccountDetailResponse{
		ID:             details.BankInformation.ID,
		CuratorID:      details.BankInformation.CuratorID,
		Location:       details.BankInformation.Location,
		FirstName:      details.BankInformation.FirstName,
		LastName:       details.BankInformation.LastName,
		DateOfBirth:    details.BankInformation.DateOfBirth,
		BankAddress:    details.AccountDetails.BankAddress,
		BankName:       details.AccountDetails.BankName,
		BranchCode:     details.AccountDetails.BranchCode,
		AccountNumber:  details.AccountDetails.AccountNumber,
		AccountName:    details.AccountDetails.AccountName,
		AccountAddress: details.AccountDetails.AccountAddress,
		IBAN:           details.AccountDetails.IBAN,
		CreatedAt:      details.BankInformation.CreatedAt,
		UpdatedAt:      details.BankInformation.UpdatedAt,
	}

	return resp, nil
}
