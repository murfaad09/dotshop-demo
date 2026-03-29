package handlers

import (
	"strconv"

	"github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

type WishlistHandler struct {
	service service.WishlistService
}

func NewWishlistHandler(service service.WishlistService) *WishlistHandler {
	return &WishlistHandler{service}
}

// AddToWishlist
//
//	@Summary		Add item to wishlist
//	@Description	Add item to wishlist
//	@Tags			Wishlist
//	@Security		BearerAuth
//	@Param			user_id	path		string						true	"user_id"
//	@Param			body	body		dto.AddWishlistItemRequest	true	"Wishlist Item"
//	@Success		200		{object}	dto.AddWishlistItemResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/wishlist/{user_id}/add [post]
func (h *WishlistHandler) AddToWishlist(c *fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var request dto.AddWishlistItemRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := dto.ValidateWishlistItemRequest(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	wishlistItem, err := h.service.AddToWishlist(&request, uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(wishlistItem)
}

// RemoveFromWishlist
//
//	@Summary		Remove item from wishlist
//	@Description	Remove item from wishlist
//	@Tags			Wishlist
//	@Security		BearerAuth
//	@Param			user_id		path		string	true	"user_id"
//	@Param			curator_id	query		string	true	"curator_id"
//	@Param			product_id	query		string	true	"product_id"
//	@Param			variant_id	query		string	false	"variant_id"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/wishlist/{user_id}/remove [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	productID := c.Query("product_id")
	variantID := c.Query("variant_id")
	curatorID, err := strconv.Atoi(c.Query("curator_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	if err := h.service.RemoveFromWishlist(uint(userID), uint(curatorID), productID, variantID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item removed from wishlist"})
}

// GetWishlist
//
//	@Summary		Get wishlist
//	@Description	Get wishlist
//	@Tags			Wishlist
//	@Security		BearerAuth
//	@Param			user_id	path		string	true	"user_id"
//	@Success		200		{object}	dto.GetWishlistResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/wishlist/{user_id}/products [get]
func (h *WishlistHandler) GetWishlist(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	wishlistItems, err := h.service.GetWishlist(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(wishlistItems)
}
