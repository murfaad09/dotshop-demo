package service

import (
	// "github.com/harishash/dotshop-be/internal/dto"
	admin_repo "github.com/harishash/dotshop-be/internal/repositories"

	// "github.com/harishash/dotshop-be/internal/utils/email"
	"github.com/harishash/dotshop-be/internal/utils/errors"
	// "github.com/harishash/dotshop-be/internal/utils/email"
	domain "github.com/harishash/dotshop-be/internal/models"
)

// AdminService interface
type AdminService interface {
	ChangeCuratorStatus(curatorID uint64, status string) *errors.Error
}

type adminService struct {
	repo admin_repo.AdminRepo
}

func NewAdminService(repo admin_repo.AdminRepo) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) ChangeCuratorStatus(curatorID uint64, status string) *errors.Error {

	// check if curator exists or not
	curator, err := s.repo.GetCuratorByCuratorID(curatorID)
	if err != nil {
		return errors.Wrap(err).WithMessage("failed to get curator")
	}

	if curator == nil {
		return errors.New("curator not found")
	}

	previousCuratorStatus := curator.Status

	newCuratorStatus := domain.Status(status)

	if previousCuratorStatus == newCuratorStatus {
		return errors.New("curator status is already the requested status")
	}

	err = s.repo.UpdateCuratorStatus(curatorID, status)
	if err != nil {
		return errors.Wrap(err).WithMessage("failed to change curator status")
	}

	return nil
}
