package handlers

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"
)

type ProductsHandler struct {
	productService service.IProductService
}

func NewProductsHandler(Service service.IProductService) *ProductsHandler {
	return &ProductsHandler{Service}
}

// GetAllProducts Get Products With Filter
//
//	@Summary		Search products with filter by name
//	@Description	Search products with filter by name
//	@Tags			Products
//	@Accept			application/json
//	@Param			pageNum				query		int		false	"Page number"
//	@Param			pageSize			query		int		false	"Page size"
//	@Param			subCategoryIds		query		string	false	"sub category ids"
//	@Param			searchByBrandName	query		string	false	"brand name"
//	@Param			productIds			query		string	false	"product ids"
//	@Param			searchBy			query		string	false	"product name"
//	@Param			sort				query		string	false	"Sort by new_in, price_low_to_high, price_high_to_low"
//
//	@Success		200					{object}	[]dto.Response
//	@Failure		400					{object}	fiber.Error
//	@Failure		401					{object}	fiber.Error
//	@Failure		403					{object}	fiber.Error
//	@Failure		404					{object}	fiber.Error
//	@Router			/products [get]
func (p *ProductsHandler) GetAllProducts(c *fiber.Ctx) error {
	query, err := parseQuery[dto.Filter](c)
	if err != nil {
		return err
	}

	// validate sort param
	if len(query.SortBy) > 0 {
		if err := query.ValidateSortParam(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// set default paging param if not provided
	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	products, paging, err := p.productService.GetProductsWithFilter(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   products,
		Paging: *paging,
	})
}

// DeleteProducts Delete Products
//
//	@Summary		Delete Products
//	@Description	Delete Products
//	@Tags			Products
//	@Accept			application/json
//	@Param			ids	path		string	true	"Product IDs"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/product/{ids} [delete]
func (p *ProductsHandler) DeleteProducts(c *fiber.Ctx) error {
	ids := c.Params("ids")

	decodedIDs, err := url.QueryUnescape(ids)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to decode product IDs",
		})
	}

	elements := strings.Split(decodedIDs, ",")

	var ptrIds []*string
	for _, id := range elements {
		trimmedID := strings.TrimSpace(id)
		idCopy := trimmedID
		ptrIds = append(ptrIds, &idCopy)
	}

	deletedProducts, err := p.productService.DeleteProducts(ptrIds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products deleted successfully",
		"deleted": deletedProducts,
	})
}

// GetBrands Get Brands
//
//	@Summary		Get Brands
//	@Description	Get Brands
//	@Tags			Products
//	@Accept			application/json
//	@Param			searchBy	query		string	false	"brand name"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/brands [get]
func (p *ProductsHandler) GetBrands(c *fiber.Ctx) error {

	searchStr := c.Query("search")

	brands, err := p.productService.GetBrands(searchStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"brands": brands,
	})
}

// GetAllProductsWithStats Get Brands
//
//	@Summary		Get All Products With Stats
//	@Description	Get all products with reviews and average rating
//	@Tags			AdminBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			stars		query		string	false	"product rating"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/admin/reviews/products [get]
func (p *ProductsHandler) GetAllProductsWithStats(c *fiber.Ctx) error {
	query, err := parseQuery[dto.ListProductReviewsRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query parameters"})
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	products, err := p.productService.GetAllProductsWithStats(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
