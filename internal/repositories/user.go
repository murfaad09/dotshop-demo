package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/harishash/dotshop-be/integration/klaviyo"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByEmailUnscoped(email string) (*domain.User, error)
	ValidateUser(email string, password string) (*domain.User, *uint, error)
	GetUsers() ([]*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
	UpdateUserProfileFields(user domain.User) (*domain.User, error)
	GetUserAddressByID(id uint64) (*domain.ShippingInfo, error)
	RoleExists(id uint) (bool, error)
	GetUserByID(id uint) (*domain.User, error)
	GetProfileWithEmail(email string) (*domain.User, error)
	UpdateConsumerProfile(user *domain.User) error
	CheckEmailExists(email string) (*domain.User, error)
	AddConsumerAddress(address *domain.ShippingInfo) error
	GetAllUserAddresses(userId uint) ([]domain.ShippingInfo, error)
	UpdateUserAddress(address *domain.ShippingInfo) error
	DeleteAddress(userId, addressId uint64) error
	UpdatePasswordByID(userId uint64, newPassword string) (*dto.UpdatePasswordResponse, error)
	PatchUserAddress(addressId uint, updates map[string]interface{}) error
}

type userRepo struct {
	db      *gorm.DB
	klaviyo *klaviyo.Klaviyo
}

func NewUserRepository(db *gorm.DB, klaviyo *klaviyo.Klaviyo) IUserRepository {
	return &userRepo{db: db, klaviyo: klaviyo}
}

func (r *userRepo) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(id uint) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) GetUserByEmailUnscoped(email string) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Unscoped().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepo) ValidateUser(email string, password string) (*domain.User, *uint, error) {
	user := domain.User{}
	result := r.db.Where("email = ? AND password_hash = ?", email, password).First(&user)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var curatorId uint
	if user.RoleID == 2 {
		result = r.db.Select("id").Model(&domain.Curator{}).
			Where("user_id = ? AND deleted_at IS NULL", user.ID).First(&curatorId)
		if result.Error != nil {
			// If the curator is found but has been soft-deleted, return an error indicating the account doesn't exist
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil, fmt.Errorf("account does not exist")
			}
			return nil, nil, result.Error
		}
	}

	return &user, &curatorId, nil
}

func (r *userRepo) RoleExists(id uint) (bool, error) {
	var role domain.Role
	if err := r.db.First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userRepo) GetUsers() ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	results := r.db.Find(&users)
	if results.Error != nil {
		return nil, results.Error
	}
	return users, nil
}

func (r *userRepo) CreateUser(user domain.User) (*domain.User, error) {

	results := r.db.Create(&user)
	if results.Error != nil {
		return nil, results.Error
	}

	klaviyoProfile, err := (*r.klaviyo).CreateProfile(user.Email, *user.FirstName, *user.LastName)
	if err != nil {
		log.Printf("Error creating Klaviyo profile: %v\n", err)
	}

	if klaviyoProfile != nil {
		err := r.db.Model(&user).Where("id = ?", user.ID).Update("klaviyo_id", klaviyoProfile.Data.ID).Error
		if err != nil {
			log.Printf("Error updating Klaviyo ID: %v\n", err)
		}

		err = (*r.klaviyo).AddProfileToList(constants.KlaviyoListIdForEmailList, klaviyoProfile.Data.ID)
		if err != nil {
			log.Printf("Error adding Klaviyo profile to list: %v\n", err)
		}
	}

	return &user, nil
}

func (r *userRepo) UpdateUserProfileFields(user domain.User) (*domain.User, error) {
	err := r.db.Exec(`
		UPDATE users 
		SET 
			email = ?, 
			first_name = ?, 
			last_name = ?, 
			password_hash = ?, 
			role_id = ?, 
			updated_at = ?, 
			deleted_at = NULL 
		WHERE id = ?`,
		user.Email, user.FirstName, user.LastName, user.PasswordHash, user.RoleID, time.Now(), user.ID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update user profile fields, error: %v", err)
	}

	var updatedUser domain.User
	if err := r.db.Where("id = ?", user.ID).First(&updatedUser).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated user, error: %v", err)
	}

	return &updatedUser, nil
}

func (r *userRepo) GetUserAddressByID(id uint64) (*domain.ShippingInfo, error) {
	address := domain.ShippingInfo{}
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&address)
	if result.Error != nil {
		return nil, fmt.Errorf("user address not found with id: %v", id)
	}
	return &address, nil
}

func (r *userRepo) GetProfileWithEmail(email string) (*domain.User, error) {
	var user *domain.User
	results := r.db.Where("email = ?", email).Find(&user)
	if results.Error != nil {
		return nil, results.Error
	}
	return user, nil
}

func (r *userRepo) UpdateConsumerProfile(user *domain.User) error {
	if err := r.db.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return fmt.Errorf("failed to update consumer profile: %v", err)
	}

	return nil
}

func (r *userRepo) CheckEmailExists(email string) (*domain.User, error) {
	user := domain.User{}
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) AddConsumerAddress(address *domain.ShippingInfo) error {
	results := r.db.Create(&address)
	if results.Error != nil {
		return results.Error
	}
	return nil
}

func (r *userRepo) GetAllUserAddresses(userId uint) ([]domain.ShippingInfo, error) {
	var addresses []domain.ShippingInfo
	results := r.db.Order("created_at DESC").Where("user_id = ? AND deleted_at IS NULL", userId).Find(&addresses)
	if results.Error != nil {
		return nil, results.Error
	}

	return addresses, nil
}

func (r *userRepo) UpdateUserAddress(address *domain.ShippingInfo) error {
	return r.db.Save(address).Error
}

func (r *userRepo) DeleteAddress(userId, addressId uint64) error {
	return r.db.Where("id = ? AND user_id = ?", addressId, userId).Delete(&domain.ShippingInfo{}).Error
}

func (r *userRepo) UpdatePasswordByID(userId uint64, newPassword string) (*dto.UpdatePasswordResponse, error) {
	err := r.db.Model(&domain.User{}).Where("id = ?", userId).Update("password_hash", newPassword).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update password for id %v: %v", userId, err)
	}

	res := dto.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}

	return &res, nil
}

func (r *userRepo) PatchUserAddress(addressId uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.ShippingInfo{}).Where("id = ?", addressId).Updates(updates).Error
}
