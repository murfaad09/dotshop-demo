package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	admin_service "github.com/harishash/dotshop-be/internal/services"
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type AdminHandlers struct {
	service admin_service.AdminService
}

func NewAdminHandlers(service admin_service.AdminService) *AdminHandlers {
	return &AdminHandlers{service: service}
}

func (h *AdminHandlers) ChangeCuratorStatus(c *fiber.Ctx) error {
	var reqBody dto.ChangeCuratorStatusRequest
	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := reqBody.Validate(); err != nil {
		return errors.Wrap(err).WithMessage("Invalid request body")
	}

	// check if curator exists or not

	err = h.service.ChangeCuratorStatus(curatorID, string(reqBody.Status))

	return c.SendStatus(fiber.StatusOK)
}
