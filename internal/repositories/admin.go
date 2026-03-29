package repository

import (
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

// AdminRepo interface
type AdminRepo interface {
	UpdateCuratorStatus(curatorID uint64, status string) error
	GetCuratorByCuratorID(curatorID uint64) (*domain.Curator, error)
	GetAdminUserIDsByRole(roleID uint) ([]uint, error)
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo() AdminRepo {
	instance := GetDatabaseConnection()
	return &adminRepo{db: instance.Connection}
}

func (r *adminRepo) UpdateCuratorStatus(curatorID uint64, status string) error {
	result := r.db.Model(&domain.Curator{}).Where("id = ?", curatorID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *adminRepo) GetCuratorByCuratorID(curatorID uint64) (*domain.Curator, error) {
	curators := domain.Curator{}
	result := r.db.Where("id = ?", curatorID).Preload("User").Find(&curators)
	if result.Error != nil {
		return nil, result.Error
	}

	return &curators, nil
}

func (r *adminRepo) GetAdminUserIDsByRole(roleID uint) ([]uint, error) {
	var adminIDs []uint

	query := r.db.Table("users").
		Where("role_id = ?", roleID).
		Pluck("id", &adminIDs)

	if query.Error != nil {
		return nil, query.Error
	}

	return adminIDs, nil
}
