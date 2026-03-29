package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/middlewares/payments"
)

const PaymentError = "following parameters are required"

type PaymentHandler struct {
	paypalPaymentService payments.IPaymentGateway
	stripePaymentService payments.IPaymentGateway
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paypalPaymentService: payments.NewPaymentGateway("paypal"),
		stripePaymentService: payments.NewPaymentGateway("stripe"),
	}
}

func (h *PaymentHandler) GetConfig(c *fiber.Ctx) error {
	gateway, err := getGateWay(h, c.Params("gateway"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	response, err := gateway.GetPublishableKey()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(constants.InvalidStripeKey)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PaymentHandler) Authorize(c *fiber.Ctx) error {
	input := new(constants.PaymentAuthoriseRequest)
	gateway, err := getGateWay(h, c.Params("gateway"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": PaymentError,
			"params":  constants.PaymentAuthoriseRequest{},
		})
	}
	auth, err := gateway.Authorize(&input.Amount, &input.Currency)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"transaction": auth})
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	input := new(constants.PaymentCreateRequest)
	gateway, err := getGateWay(h, c.Params("gateway"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	fmt.Println(string(c.BodyRaw()))
	fmt.Println("")
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": PaymentError,
			"params":  constants.PaymentCreateRequest{},
		})
	}

	response, err := gateway.Create(&input.Amount, &input.Currency)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PaymentHandler) CapturePayment(c *fiber.Ctx) error {
	input := new(constants.PaymentCaptureRequest)
	gateway, err := getGateWay(h, c.Params("gateway"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	fmt.Println(string(c.BodyRaw()))
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": PaymentError,
			"params":  constants.PaymentCaptureRequest{},
		})
	}

	response, err := gateway.Capture(&input.Amount, &input.Currency, &input.OrderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PaymentHandler) UpdatePayment(c *fiber.Ctx) error {
	return nil
}

func getGateWay(h *PaymentHandler, gateway string) (payments.IPaymentGateway, error) {

	switch gateway {
	case "stripe":
		return h.stripePaymentService, nil
	case "paypal":
		return h.paypalPaymentService, nil
	}
	return nil, errors.New("invalid gateway, only stripe and paypal are supported as URLParams")
}
