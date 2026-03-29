package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/dto"
	"github.com/harishash/dotshop-be/internal/utils/errors"

	cart_services "github.com/harishash/dotshop-be/internal/services"
)

type CartsHandlers struct {
	service cart_services.CartsService
}

func NewCartsHandlers(service cart_services.CartsService) *CartsHandlers {
	return &CartsHandlers{service: service}
}

// func (h *CartsHandlers) BuyNow(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(string)
// 	err := h.service.BuyNow(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return c.SendStatus(fiber.StatusOK)
// }

// CreateCart creates a User Cart
//
//	@Summary		Create  a User Cart
//	@Description	Create  a User Cart
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CartRequest	true	"Cart Items"
//	@Success		200		{object}	dto.CartResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/cart/add [post]
func (h *CartsHandlers) CreateCart(c *fiber.Ctx) error {
	var cart *dto.CartRequest
	if err := c.BodyParser(&cart); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdCart, err := h.service.CreateCart(cart)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdCart)
}

// UpdateCartItemQuantity 	updates the quantity of an item in the cart
//
//	@Summary		Update the quantity of an item in the cart
//	@Description	Update the quantity of an item in the cart
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			cart_id		path		int							true	"Cart ID"
//	@Param			variant_id	path		string						true	"Variant ID"
//	@Param			body		body		dto.CartItemQuantityRequest	true	"Quantity"
//	@Success		200			{object}	dto.CartItemResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/cart/{cart_id}/product/{variant_id} [put]
func (h *CartsHandlers) UpdateCartItemQuantity(c *fiber.Ctx) error {
	cartIDStr := c.Params("cart_id")
	if cartIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CartID is required")
	}
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	variantID := c.Params("variant_id")
	if variantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "variant ID is required"})
	}

	var request dto.CartItemQuantityRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	if err := request.Validate(); err != nil {
		return errors.Wrap(err).WithMessage("Invalid request body")
	}

	response, err := h.service.UpdateCartItemQuantity(cartID, variantID, request.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

// AddItemsToCart adds items to the cart
//
//	@Summary		Add items to the cart
//	@Description	Add items to the cart
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			cart_id	path		int							true	"Cart ID"
//	@Param			body	body		[]dto.AddCartItemsRequest	true	"Cart Items"
//	@Success		200		{object}	dto.CartResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/cart/{cart_id}/items [post]
func (h *CartsHandlers) AddItemsToCart(c *fiber.Ctx) error {

	// Extract the cart ID from the URL
	cartIDParam := c.Params("cart_id")
	cartID, err := strconv.Atoi(cartIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart ID",
		})
	}

	var request []*dto.AddCartItemsRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	response, err := h.service.AddItemsToCart(uint(cartID), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

// DeleteCart deletes the cart
//
//	@Summary		Delete cart
//	@Description	Delete cart
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			cart_id	path		int	true	"Cart ID"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/cart/{cart_id} [delete]
func (h *CartsHandlers) DeleteCart(c *fiber.Ctx) error {
	cartIDStr := c.Params("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	err = h.service.DeleteCart(uint(cartID))
	if err != nil {
		if err.Error() == "cart does not exist" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Send success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart with ID " + strconv.Itoa(int(cartID)) + " deleted successfully",
	})
}

// DeleteCart deletes the cart item
//
//	@Summary		Delete cart item
//	@Description	Delete cart item
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			cart_id		path		int	true	"Cart ID"
//	@Param			variant_id	path		int	true	"Variant ID"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/cart/{cart_id}/items/{variant_id} [delete]
func (h *CartsHandlers) DeleteCartItem(c *fiber.Ctx) error {
	cartIDStr := c.Params("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	variantIDStr := c.Params("variant_id")
	if variantIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	err = h.service.DeleteCartItem(uint(cartID), variantIDStr)
	if err != nil {
		if err.Error() == "cart does not exist" || err.Error() == "cart item does not exist in the cart" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Send success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart item with ID " + variantIDStr + " deleted successfully",
	})
}

// GetCartByUserID gets the cart by user ID
//
//	@Summary		Get cart by user ID
//	@Description	Get cart by user ID
//	@Tags			Checkout
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	dto.CartResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/cart/user/{user_id} [get]
func (h *CartsHandlers) GetCartByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	currency := c.Query("currency")
	postalCode := c.Query("postalCode")
	country := c.Query("country")

	if currency == "" {
		currency = "usd"
	}
	if postalCode == "" {
		postalCode = "98104"
	}
	if country == "" {
		country = "US"
	}

	cartResponse, err := h.service.GetCartByUserID(uint(userID), currency, postalCode, country)
	if err != nil {
		if err.Error() == "active cart does not exist for this user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(cartResponse)
}
