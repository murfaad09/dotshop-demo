package handlers

import (
	"github.com/gofiber/fiber/v2"

	dto "github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"
)

type PromotionHandler struct {
	service service.PromotionService
}

func NewPromotionHandler(service service.PromotionService) *PromotionHandler {
	return &PromotionHandler{service}
}

// CreatePromotion Create Promotion
//
//	@Summary		Create Promotion
//	@Description	Create Promotion
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			promotion	body		dto.PromotionRequest	true	"Promotion data"
//	@Success		201			{object}	dto.PromotionResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/admin/promotions [post]
func (h *PromotionHandler) CreatePromotion(c *fiber.Ctx) error {
	var req dto.PromotionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.service.CreatePromotion(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdatePromotion Update Promotion
//
//	@Summary		Update Promotion
//	@Description	Update Promotion
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int						true	"Promotion ID"
//	@Param			promotion	body		dto.PromotionRequest	true	"Promotion data"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Router			/admin/promotions/{id} [put]
func (h *PromotionHandler) UpdatePromotion(c *fiber.Ctx) error {
	var req dto.PromotionRequest

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid promotion ID"})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdatePromotion(&req, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Promotion updated successfully"})
}

// GetPromotions Get Promotions
//
//	@Summary		Get Promotions
//	@Description	Get Promotions
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			status		query		string	false	"Promotion status"
//	@Param			startValue	query		float64	false	"Start value"
//	@Param			endValue	query		float64	false	"End value"
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/promotions [get]
func (h *PromotionHandler) GetPromotions(c *fiber.Ctx) error {
	query, err := parseQuery[dto.ListPromotionsRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetPromotions(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetPromotionByID Get Promotion by ID
//
//	@Summary		Get Promotion by ID
//	@Description	Get Promotion by ID
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Promotion ID"
//	@Success		200	{object}	dto.PromotionResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/promotions/{id} [get]
func (h *PromotionHandler) GetPromotionByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid promotion ID"})
	}

	resp, err := h.service.GetPromotionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeletePromotion Delete Promotion
//
//	@Summary		Delete Promotion
//	@Description	Delete Promotion
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Promotion ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/promotions/{id} [delete]
func (h *PromotionHandler) DeletePromotion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid promotion ID"})
	}

	err = h.service.DeletePromotion(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Promotion deleted successfully"})
}

// ApplyBulkDiscount Apply Bulk Discount
//
//	@Summary		Apply Bulk Discount
//	@Description	Apply a discount to multiple products
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ApplyBulkDiscountRequest	true	"Apply Bulk Discount Request"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/admin/promotions/discount [post]
func (h *PromotionHandler) ApplyBulkDiscount(c *fiber.Ctx) error {
	var req dto.ApplyBulkDiscountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	if len(req.ProductIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product IDs cannot be empty"})
	}
	if req.PromotionID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid promotion ID"})
	}

	err := h.service.ApplyDiscountToProducts(req.ProductIDs, req.PromotionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "discount applied successfully"})
}
