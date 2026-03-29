// internal/handlers/notification_handler.go
package handlers

import (
	"strconv"

	"github.com/harishash/dotshop-be/internal/dto"
	services "github.com/harishash/dotshop-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	Service services.NotificationServiceInterface
}

func NewNotificationHandler(service services.NotificationServiceInterface) *NotificationHandler {
	return &NotificationHandler{Service: service}
}

// CreateNotification creates a new notification
//
//	@Summary		Creates a new notification
//	@Description	creates a new notification
//	@Tags			Notification
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.NotificationRequest	true	"Notification"
//	@Success		200		{object}	dto.NotificationResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/notifications [post]
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint64)
	var payload dto.NotificationRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	notification, err := h.Service.CreateNotification(userID, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(notification)
}

// GetNotifications returns a list of notifications
//
//	@Summary		Returns a list of notifications
//	@Description	Returns a list of notifications
//	@Tags			Notification
//	@Security		BearerAuth
//	@Accept			application/json
//	@Success		200	{object}	dto.NotificationResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/notifications [get]
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint64)

	notifications, err := h.Service.GetNotifications(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(notifications)
}

// MarkNotificationAsRead marks a notification as read
//
//	@Summary		Marks a notification as read
//	@Description	Marks a notification as read
//	@Tags			Notification
//	@Security		BearerAuth
//	@Accept			application/json
//
// @Param			id	path		int	true	"Notification ID"
//
//	@Success		200	{object}	dto.NotificationResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/notifications/{id}/read [put]
func (h *NotificationHandler) MarkNotificationAsRead(c *fiber.Ctx) error {
	notificationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid notification ID"})
	}

	notification, err := h.Service.MarkNotificationAsRead(uint(notificationID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(notification)
}
