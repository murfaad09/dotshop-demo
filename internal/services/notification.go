package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	repositories "github.com/harishash/dotshop-be/internal/repositories"
)

type NotificationServiceInterface interface {
	CreateNotification(userID uint64, body dto.NotificationRequest) (*dto.NotificationResponse, error)
	GetNotifications(userID uint64) ([]*dto.NotificationResponse, error)
	MarkNotificationAsRead(notificationID uint) (*dto.NotificationResponse, error)
}

type NotificationService struct {
	Repo repositories.NotificationRepositoryInterface
}

func NewNotificationService(
	repo repositories.NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{Repo: repo}
}

func (service *NotificationService) CreateNotification(userID uint64, body dto.NotificationRequest) (*dto.NotificationResponse, error) {
	notification := &domain.Notification{
		UserID:  userID,
		Title:   body.Title,
		Message: body.Message,
		IsRead:  false,
	}

	if err := service.Repo.Create(notification); err != nil {
		return nil, err
	}

	responseDTO := &dto.NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		Read:      notification.IsRead,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}

	return responseDTO, nil
}

func (service *NotificationService) GetNotifications(userID uint64) ([]*dto.NotificationResponse, error) {
	notifications, err := service.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responseDTOs []*dto.NotificationResponse
	for _, notification := range notifications {
		responseDTOs = append(responseDTOs, &dto.NotificationResponse{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Message:   notification.Message,
			Read:      notification.IsRead,
			CreatedAt: notification.CreatedAt,
			UpdatedAt: notification.UpdatedAt,
		})
	}

	return responseDTOs, nil
}

func (service *NotificationService) MarkNotificationAsRead(notificationID uint) (*dto.NotificationResponse, error) {
	if err := service.Repo.MarkAsRead(notificationID); err != nil {
		return nil, err
	}

	notification, err := service.Repo.GetByID(notificationID)
	if err != nil {
		return nil, err
	}

	responseDTO := &dto.NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		Read:      notification.IsRead,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}

	return responseDTO, nil
}
