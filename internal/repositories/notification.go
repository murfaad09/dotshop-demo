package repository

import (
	"fmt"
	"time"

	constants "github.com/harishash/dotshop-be/internal/constants"
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(notification *domain.Notification) error
	GetByUserID(userID uint64) ([]domain.Notification, error)
	MarkAsRead(notificationID uint) error
	GetByID(notificationID uint) (*domain.Notification, error)
	NotifyCurators(orderVariants []*domain.OrderVariants, message string) error
	NotifyAdmins(message string) error
}

type NotificationRepository struct {
	DB        *gorm.DB
	AdminRepo AdminRepo
}

func NewNotificationRepository(db *gorm.DB, adminRepo AdminRepo) *NotificationRepository {
	return &NotificationRepository{DB: db, AdminRepo: adminRepo}
}
func (repo *NotificationRepository) Create(notification *domain.Notification) error {
	return repo.DB.Create(notification).Error
}

func (repo *NotificationRepository) GetByUserID(userID uint64) ([]domain.Notification, error) {
	var notifications []domain.Notification
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err := repo.DB.Order("created_at DESC").Where("user_id = ? AND (is_read = false OR created_at >= ?)", userID, sevenDaysAgo).Find(&notifications).Error
	return notifications, err
}

func (repo *NotificationRepository) MarkAsRead(notificationID uint) error {
	return repo.DB.Model(&domain.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (repo *NotificationRepository) GetByID(notificationID uint) (*domain.Notification, error) {
	var notification domain.Notification
	err := repo.DB.Where("id = ?", notificationID).First(&notification).Error
	return &notification, err
}

func (repo *NotificationRepository) NotifyCurators(orderVariants []*domain.OrderVariants, message string) error {
	curatorIDs := make(map[uint]struct{})
	for _, variant := range orderVariants {
		if _, exists := curatorIDs[uint(variant.CuratorID)]; !exists {
			curatorIDs[uint(variant.CuratorID)] = struct{}{}
			user, err := repo.AdminRepo.GetCuratorByCuratorID(variant.CuratorID)
			if err != nil {
				return fmt.Errorf("failed to retrieve curator: %v", err)
			}
			notification := &domain.Notification{
				UserID:  uint64(user.ID),
				Title:   message,
				Message: message,
				IsRead:  false,
			}

			if err := repo.Create(notification); err != nil {
				return fmt.Errorf("failed to create notification for curator: %v", err)
			}
		}
	}
	return nil
}

func (repo *NotificationRepository) NotifyAdmins(message string) error {
	userIDs, err := repo.AdminRepo.GetAdminUserIDsByRole(constants.ADMIN_ROLE_ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve admins: %v", err)
	}
	for _, userID := range userIDs {
		notification := &domain.Notification{
			UserID:  uint64(userID),
			Title:   message,
			Message: message,
			IsRead:  false,
		}
		if err := repo.Create(notification); err != nil {
			return fmt.Errorf("failed to create notification for admin: %v", err)
		}
	}
	return nil
}
