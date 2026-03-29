package repository

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/harishash/dotshop-be/integration/klaviyo"
	constants "github.com/harishash/dotshop-be/internal/constants"
	dto "github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
)

type CuratorOnboardingRepo interface {
	AddCurator(curator dto.CuratorOnBoardingRequest) (dto.CuratorOnBoardingResponse, error)
	CheckShopName(shopName string) (*domain.Curator, error)
	GetCuratorByStoreName(storeName string) (*domain.Curator, error)
}

type curatorOnboardingRepo struct {
	db      *gorm.DB
	klaviyo *klaviyo.Klaviyo
}

func NewCuratorOnboardingRepo(db *gorm.DB, klaviyo *klaviyo.Klaviyo) CuratorOnboardingRepo {
	return &curatorOnboardingRepo{db: db, klaviyo: klaviyo}
}

func (r *curatorOnboardingRepo) AddCurator(curator dto.CuratorOnBoardingRequest) (dto.CuratorOnBoardingResponse, error) {
	tx := r.db.Begin()

	var existingUser domain.User
	if err := r.db.Unscoped().Where("email = ?", curator.Email).First(&existingUser).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.CuratorOnBoardingResponse{}, err
	}

	var curatorModel domain.Curator
	if existingUser.ID != 0 && existingUser.DeletedAt.Valid {
		existingUser.FirstName = &curator.FirstName
		existingUser.LastName = &curator.LastName
		existingUser.PasswordHash = &curator.Password
		existingUser.RoleID = constants.CURATOR_ROLE_ID
		existingUser.DeletedAt = gorm.DeletedAt{}

		if err := tx.Save(&existingUser).Error; err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		if err := r.db.Unscoped().Where("user_id = ?", existingUser.ID).First(&curatorModel).Error; err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		curatorModel.Name = curator.FirstName + " " + curator.LastName
		curatorModel.ShopName = curator.ShopName
		curatorModel.Bio = curator.Bio
		curatorModel.ProfileImageURL = curator.ProfileImage
		curatorModel.CoverImageURL = curator.CoverImage
		curatorModel.DeletedAt = gorm.DeletedAt{}
		curatorModel.UpdatedAt = time.Now().UTC()

		if err := tx.Save(&curatorModel).Error; err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		if err := r.updateSocialMediaLinks(tx, curatorModel.ID, curator.Socials); err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		if err := r.updateBankInformation(tx, curatorModel.ID, curator.BankInformation); err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}
	} else {
		user, err := r.createUser(tx, curator)
		if err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		curatorModel, err = r.createCurator(tx, user.ID, curator)
		if err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		if err := r.addSocialMediaLinks(tx, curatorModel.ID, curator.Socials); err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}

		if err := r.addBankInformation(tx, curatorModel.ID, curator.BankInformation); err != nil {
			tx.Rollback()
			return dto.CuratorOnBoardingResponse{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return dto.CuratorOnBoardingResponse{}, err
	}

	return r.prepareResponse(curator, curatorModel), nil
}

func (r *curatorOnboardingRepo) updateSocialMediaLinks(tx *gorm.DB, curatorID uint, socials []*dto.SocialMediaLinksRequest) error {
	if err := tx.Where("curator_id = ?", curatorID).Delete(&domain.SocialMediaLink{}).Error; err != nil {
		return err
	}

	return r.addSocialMediaLinks(tx, curatorID, socials)
}

func (r *curatorOnboardingRepo) updateBankInformation(tx *gorm.DB, curatorID uint, bankInfos []*dto.BankInformationRequest) error {
	if err := tx.Where("bank_id IN (SELECT id FROM bank_informations WHERE curator_id = ?)", curatorID).Delete(&domain.AccountDetails{}).Error; err != nil {
		return err
	}

	if err := tx.Where("curator_id = ?", curatorID).Delete(&domain.BankInformation{}).Error; err != nil {
		return err
	}

	return r.addBankInformation(tx, curatorID, bankInfos)
}

func (r *curatorOnboardingRepo) createUser(tx *gorm.DB, curator dto.CuratorOnBoardingRequest) (domain.User, error) {
	user := domain.User{
		Email:        curator.Email,
		PasswordHash: &curator.Password,
		RoleID:       constants.CURATOR_ROLE_ID,
		FirstName:    &curator.FirstName,
		LastName:     &curator.LastName,
	}
	if err := tx.Create(&user).Error; err != nil {
		return domain.User{}, err
	}

	klaviyoProfile, err := (*r.klaviyo).CreateProfile(curator.Email, curator.FirstName, curator.LastName)
	if err != nil {
		log.Printf("Error creating Klaviyo profile: %v\n", err)
	}

	if klaviyoProfile != nil {
		err := tx.Model(&user).Where("id = ?", user.ID).Update("klaviyo_id", klaviyoProfile.Data.ID).Error
		if err != nil {
			log.Printf("Error updating Klaviyo ID: %v\n", err)
		}

		err = (*r.klaviyo).AddProfileToList(constants.KlaviyoListIdForPendingCurators, klaviyoProfile.Data.ID)
		if err != nil {
			log.Printf("Error adding Klaviyo profile to list: %v\n", err)
		}
	}

	return user, nil
}

func (r *curatorOnboardingRepo) createCurator(tx *gorm.DB, userID uint, curator dto.CuratorOnBoardingRequest) (domain.Curator, error) {
	fullname := curator.FirstName + " " + curator.LastName
	curatorModel := domain.Curator{
		UserID:            userID,
		Name:              fullname,
		ShopName:          curator.ShopName,
		Bio:               curator.Bio,
		NumberofFollowers: 0,
		// Status:            "pending",
		ProfileImageURL: curator.ProfileImage,
		CoverImageURL:   curator.CoverImage,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	if err := tx.Create(&curatorModel).Error; err != nil {
		return domain.Curator{}, err
	}
	return curatorModel, nil
}

func (r *curatorOnboardingRepo) addSocialMediaLinks(tx *gorm.DB, curatorID uint, socials []*dto.SocialMediaLinksRequest) error {
	for _, socialMediaLink := range socials {
		socialMedia := domain.SocialMediaLink{
			CuratorID:    curatorID,
			Type:         domain.SocialMediaType(socialMediaLink.Platform),
			URL:          socialMediaLink.Link,
			AccessToken:  socialMediaLink.AccessToken,
			OpenID:       socialMediaLink.OpenID,
			RefreshToken: socialMediaLink.RefreshToken,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}
		if err := tx.Create(&socialMedia).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *curatorOnboardingRepo) addBankInformation(tx *gorm.DB, curatorID uint, bankInfos []*dto.BankInformationRequest) error {
	for _, bankInfo := range bankInfos {
		bankInformation := domain.BankInformation{
			CuratorID:   curatorID,
			Location:    bankInfo.Location,
			FirstName:   bankInfo.FirstName,
			LastName:    bankInfo.LastName,
			DateOfBirth: bankInfo.DateOfBirth,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		if err := tx.Create(&bankInformation).Error; err != nil {
			return err
		}

		accountDetail := domain.AccountDetails{
			BankID:         bankInformation.ID,
			BankAddress:    bankInfo.AccountDetails.BankAddress,
			BankName:       bankInfo.AccountDetails.BankName,
			BranchCode:     bankInfo.AccountDetails.BranchCode,
			AccountNumber:  bankInfo.AccountDetails.AccountNumber,
			AccountName:    bankInfo.AccountDetails.AccountName,
			AccountAddress: bankInfo.AccountDetails.AccountAddress,
			IBAN:           bankInfo.AccountDetails.IBAN,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}
		if err := tx.Create(&accountDetail).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *curatorOnboardingRepo) prepareResponse(curator dto.CuratorOnBoardingRequest, curatorModel domain.Curator) dto.CuratorOnBoardingResponse {
	response := dto.CuratorOnBoardingResponse{
		ID:        curatorModel.ID,
		FirstName: curator.FirstName,
		LastName:  curator.LastName,
		Email:     curator.Email,
		Bio:       curator.Bio,
		// Status:          curatorModel.Status,
		ShopName:        curator.ShopName,
		ProfileImage:    curator.ProfileImage,
		CoverImage:      curator.CoverImage,
		Socials:         curator.Socials,
		BankInformation: curator.BankInformation,
	}
	return response
}

func (r *curatorOnboardingRepo) CheckShopName(shopName string) (*domain.Curator, error) {
	curator := domain.Curator{}
	result := r.db.Where("shop_name = ?", shopName).First(&curator)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &curator, nil
}

func (r *curatorOnboardingRepo) GetCuratorByStoreName(storeName string) (*domain.Curator, error) {
	curators := domain.Curator{}
	result := r.db.Where("shop_name ILIKE ?", "%"+storeName+"%").
		Preload("User").
		Preload("SocialMediaLinks").Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curators, nil
}
