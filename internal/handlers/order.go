package handlers

import (
	"github.com/gofiber/fiber/v2"

	dto "github.com/harishash/dotshop-be/internal/dto"
	order_service "github.com/harishash/dotshop-be/internal/services"
)

type OrderHandler struct {
	orderService order_service.IOrderService
}

func NewOrderHandler(service order_service.IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: service,
	}
}

// CreateOrder Create Order
//
//	@Summary		Create Order
//	@Description	Create new order
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.OrderRequest	true	"Create Order Request"
//
//	@Success		200		{object}	dto.OrderResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/order/create [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	body, err := parseBody[dto.OrderRequest](c)
	if err != nil {
		return err
	}

	resp, err := h.orderService.CreateOrder(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// OrdersList Get Order List
//
//	@Summary		Get Order List
//	@Description	This endpoint is used to get our order list
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//
//	@Success		200	{object}	[]dto.OrdersListResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/curator/orders [get]
func (h *OrderHandler) OrdersList(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)
	resp, err := h.orderService.OrdersList(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// CancelOrder cancel order
//
//	@Summary		This endpoint is used to cancel an order
//	@Description	This endpoint is used to cancel an order
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			order_id	path		string					true	"Order ID"
//	@Param			body		body		dto.CancelOrderRequest	true	"Order Cancel Request"
//	@Success		200			{object}	dto.CartResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/user/order/{order_id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	orderId := c.Params("order_id")
	if orderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id is required",
		})
	}

	body, err := parseBody[dto.CancelOrderRequest](c)
	if err != nil {
		return err
	}

	err = h.orderService.CancelOrder(body, orderId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order cancelled successfully"})
}

// CreateReturn Create Return
//
//	@Summary		Create Return
//	@Description	Create return
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.ReturnRequest	true	"Create Return Request"
//
//	@Success		200		{object}	[]dto.ReturnResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/order/return [post]
func (h *OrderHandler) CreateReturn(c *fiber.Ctx) error {
	body, err := parseBody[dto.ReturnRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)
	for _, req := range body.ReturnVariants {
		if req.Quantity <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity must be greater than 0 for variant: " + req.VariantId,
			})
		}
	}

	resp, err := h.orderService.CreateReturn(body, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
