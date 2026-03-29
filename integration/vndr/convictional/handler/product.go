package handler

import (
	"github.com/gofiber/fiber/v2"
	product_service "github.com/harishash/dotshop-be/integration/vndr/convictional/service"
)

type ProductHandler struct {
	productService product_service.IProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: product_service.NewProductService(),
	}
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get products",
		})
	}

	return c.JSON(products)
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get products",
		})
	}
	return c.JSON(products)
}

// Get Product by product id
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}
	return c.JSON(product)
}

// Delete Product image by product and image id
func (h *ProductHandler) DeleteProductImage(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	imageID := c.Params("image_id")

	if productID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Product ID is required")
	}
	if imageID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Image ID is required")
	}
	_, err := h.productService.DeleteProductImage(productID, imageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product image",
		})
	}
	return c.JSON(fiber.StatusOK)
}
