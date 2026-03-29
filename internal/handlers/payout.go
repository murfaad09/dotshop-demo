package handlers

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/harishash/dotshop-be/internal/dto"

	service "github.com/harishash/dotshop-be/internal/services"
)

type PayoutHandler struct {
	payoutService service.IPayoutService
}

func NewPayoutHandler(payoutService service.IPayoutService) *PayoutHandler {
	return &PayoutHandler{payoutService}
}

// GetPayoutHistory Get Payout History
//
//	@Summary		Get payout history
//	@Description	This endpoint is used to get payout history
//	@Tags			Payout
//	@Accept			application/json
//	@Security		BearerAuth
//
//	@Success		200	{object}	[]dto.PayoutHistoryResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/curator/payout/history [get]
func (p *PayoutHandler) GetPayoutHistory(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)

	payoutHistory, err := p.payoutService.GetPayoutHistory(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(payoutHistory)
}

// GetPayoutDetails Get Payout Details
//
//	@Summary		Get payout details
//	@Description	This endpoint is used to get payout details
//	@Tags			Payout
//	@Accept			application/json
//	@Security		BearerAuth
//
//	@Success		200	{object}	dto.PayoutResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/curator/payout/details [get]
func (p *PayoutHandler) GetPayoutDetails(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)

	payoutDetails, err := p.payoutService.GetPayoutDetails(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(payoutDetails)
}
