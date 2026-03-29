package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"
)

type BrandHandler struct {
	service service.BrandService
}

func NewBrandHandler(service service.BrandService) *BrandHandler {
	return &BrandHandler{service}
}

// GetBrands Get Brands
//
//	@Summary		Get Brands
//	@Description	Get Brands
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			pageNum		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/catalog/brands [get]
func (h *BrandHandler) GetBrands(c *fiber.Ctx) error {
	query, err := parseQuery[dto.BrandsRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetBrands(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateBrandStatus Update the status of the Brand (show/hide)
//
//	@Summary		Update Brand Status
//	@Description	Update the status (is_active field) of the Brand (show/hide)
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id		path	int								true	"Brand ID"
//	@Param			request	body	dto.UpdateBrandStatusRequest	true	"Update Brand Status Request"
//	@Success		200		"OK"
//	@Failure		400		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/catalog/brands/{id} [patch]
func (h *BrandHandler) UpdateBrandStatus(c *fiber.Ctx) error {
	brandID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Brand ID"})
	}

	req := new(dto.UpdateBrandStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateBrandStatus(brandID, req.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "brand status updated successfully"})
}
