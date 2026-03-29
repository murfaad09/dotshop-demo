package handlers

import (
	"github.com/gofiber/fiber/v2"
	curatorOnboarding_service "github.com/harishash/dotshop-be/internal/services"

	dto "github.com/harishash/dotshop-be/internal/dto"
	"github.com/harishash/dotshop-be/internal/utils/errors"
	// "github.com/harishash/dotshop-be/internal/utils/logger"
)

type CuratorOnboardingHandlers struct {
	service curatorOnboarding_service.CuratorOnboardingService
}

func NewCuratorOnboardingHandlers(service curatorOnboarding_service.CuratorOnboardingService) *CuratorOnboardingHandlers {
	return &CuratorOnboardingHandlers{service: service}
}

// CuratorOnboarding creates a new curator
//
//	@Summary		Create a new curator
//	@Description	Create a new curator
//	@Tags			Curator
//	@Accept			application/json
//	@Param			body	body		dto.CuratorOnBoardingRequest	true	"Curator"	dto.CuratorOnBoardingRequest
//	@Success		200		{object}	dto.CuratorOnBoardingResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/onboarding [post]
func (s *CuratorOnboardingHandlers) CreateCurator(c *fiber.Ctx) error {
	var reqBody dto.CuratorOnBoardingRequest

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := reqBody.Validate(); err != nil {
		return errors.Wrap(err).WithMessage("Invalid request body")
	}

	curator, err := s.service.CreateCurator(reqBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(curator)
}

// CheckShopName Check Shop Name is Exists
//
//	@Summary		Check shop name is exists
//	@Description	Check shop name is exists
//	@Tags			Curator
//	@Accept			application/json
//	@Param			shop_name	path		string	true	"Shop Name"
//	@Success		200			{object}	dto.ShopNameResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		409			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/shopname/{shop_name} [get]
func (s *CuratorOnboardingHandlers) CheckShopName(c *fiber.Ctx) error {
	shopName := c.Params("shop_name")
	if len(shopName) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "shop name is required"})
	}

	curator, err := s.service.CheckShopName(shopName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if curator != nil {
		return c.Status(fiber.StatusConflict).JSON(dto.ShopNameResponse{
			Success:  false,
			Property: shopName,
			Message:  "Shop name already exists",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.ShopNameResponse{
		Success:  true,
		Property: shopName,
		Message:  "Shop name is available",
	})
}

// GetCuratorByStoreName Get curator by store name
//
//	@Summary		Get curator by store name
//	@Description	Get curator by store name
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			store_name	path		string	true	"Store Name"
//	@Success		200			{object}	dto.GetStoreCuratorResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/store/{store_name} [get]
func (h *CuratorOnboardingHandlers) GetCuratorByShopName(c *fiber.Ctx) error {
	storeName := c.Params("store_name")
	if len(storeName) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "storeName not provided"})
	}

	curator, err := h.service.GetCuratorByStoreName(storeName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if curator.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "curator not found"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewGetStoreCuratorResponse(*curator))
}
