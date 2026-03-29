package handler

import (
	"github.com/gofiber/fiber/v2"

	order_service "github.com/harishash/dotshop-be/integration/vndr/convictional/service"
)

type OrderHandler struct {
	orderService order_service.IOrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: order_service.NewOrderService(),
	}
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	// Call service to get orders
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get orders",
		})
	}
	// Return orders as JSON response
	return c.JSON(orders)
}
