package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	common_services "github.com/harishash/dotshop-be/internal/services"
)

type CommonHandlers struct {
	service common_services.CommonService
}

func NewCommonHandlers(service common_services.CommonService) *CommonHandlers {
	return &CommonHandlers{service: service}
}

func (h *CommonHandlers) GetAllProducts(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}
	curatorid, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	subCategories := c.Query("subCategoryIds")

	// Get the pagination parameters from the request, use defaults if not provided
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit")
	}

	products, err := h.service.GetAllProducts(uint(curatorid), subCategories, true, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Get the total count of products
	totalProducts, err := h.service.GetTotalProducts(uint(curatorid), subCategories, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Add pagination metadata to the response
	return c.JSON(fiber.Map{
		// "featureID":       products.FeatureID,
		"products":        products.Products,
		"currentPage":     page,
		"totalPages":      (totalProducts + (limit) - 1) / (limit),
		"totalProducts":   totalProducts,
		"productsPerPage": limit,
	})
}

func (h *CommonHandlers) GetAllCollections(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}
	curatorIDUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	page := 1      // Default page
	pageSize := 10 // Default page size
	pageQuery := c.Query("page")
	pageSizeQuery := c.Query("pageSize")

	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
		}
	}

	if pageSizeQuery != "" {
		pageSize, err = strconv.Atoi(pageSizeQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page size")
		}
	}

	// Fetch collections with pagination
	collections, err := h.service.GetAllCollections(uint(curatorIDUint), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalProducts, err := h.service.GetTotalProductsCount(uint(curatorIDUint))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate total pages
	totalPages := (int(totalProducts) + pageSize - 1) / pageSize

	// Return response with pagination metadata
	return c.JSON(fiber.Map{
		"currentPage":     page,
		"totalPages":      totalPages,
		"totalProducts":   totalProducts,
		"productsPerPage": pageSize,
		"collections":     collections,
	})
}

func (h *CommonHandlers) FetchSectionByID(c *fiber.Ctx) error {
	sectionID := c.Params("section_id")
	if sectionID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "sectionID is required")
	}
	sectionIDUint, err := strconv.ParseUint(sectionID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting sectionID to uint: %v", err))
	}

	// Fetch collections with pagination
	sectionProducts, err := h.service.FetchSectionByID(uint(sectionIDUint))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return response with pagination metadata
	return c.JSON(fiber.Map{
		"section": sectionProducts,
	})
}

func (h *CommonHandlers) GetCuratorAllLooks(c *fiber.Ctx) error {
	// Parse curator ID from request params
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}
	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}
	curatorIDUint := uint(curatorUint)

	page := 1   // Default page
	limit := 10 // Default page size
	pageQ := c.Query("page")
	pageSizeQ := c.Query("pageSize")

	// Convert page and limit to integers
	if pageQ != "" {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
		}
	}
	if pageSizeQ != "" {
		limit, err = strconv.Atoi(pageSizeQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
		}
	}

	// Fetch looks with pagination
	looks, err := h.service.GetCuratorAllLooks(uint(curatorIDUint), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch total number of looks
	totalLooks, err := h.service.GetTotalCuratorLooksCount(curatorIDUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate total pages
	totalPages := (int(totalLooks) + limit - 1) / limit

	// Return response with pagination metadata
	return c.JSON(fiber.Map{
		"currentPage":  page,
		"totalPages":   totalPages,
		"totalLooks":   totalLooks,
		"looksPerPage": limit,
		"looks":        looks,
	})
}

func (h *CommonHandlers) GetAllLooks(c *fiber.Ctx) error {

	page := 1   // Default page
	limit := 10 // Default page size
	pageQ := c.Query("page")
	pageSizeQ := c.Query("pageSize")

	var err error
	// Convert page and limit to integers
	if pageQ != "" {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
		}
	}
	if pageSizeQ != "" {
		limit, err = strconv.Atoi(pageSizeQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
		}
	}

	// Fetch looks with pagination
	looks, err := h.service.GetAllLooks(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch total number of looks
	totalLooks, err := h.service.GetTotalLooksCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate total pages
	totalPages := (int(totalLooks) + limit - 1) / limit

	// Return response with pagination metadata
	return c.JSON(fiber.Map{
		"currentPage":  page,
		"totalPages":   totalPages,
		"totalLooks":   totalLooks,
		"looksPerPage": limit,
		"looks":        looks,
	})
}

func (h *CommonHandlers) GetAllProductsByCollectionID(c *fiber.Ctx) error {
	collectionID := c.Params("collection_id")
	if collectionID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "collectionID is required")
	}
	collectID, err := strconv.ParseUint(collectionID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting collectionID to uint: %v", err))
	}

	page := 1      // Default page
	pageSize := 10 // Default page size
	pageQuery := c.Query("page")
	pageSizeQuery := c.Query("pageSize")

	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
		}
	}

	if pageSizeQuery != "" {
		pageSize, err = strconv.Atoi(pageSizeQuery)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page size")
		}
	}

	// Fetch products with pagination
	products, err := h.service.GetAllProductsByCollectionID(uint(collectID), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch total number of products (assuming it's available in the service layer)
	totalProducts, err := h.service.GetTotalCollectionProductsCount(uint(collectID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	paging := dto.Paging{
		TotalCount:  totalProducts,
		PageSize:    pageSize,
		CurrentPage: page,
	}

	// Return response with pagination metadata
	return c.JSON(dto.Response{
		Data:   products,
		Paging: paging,
	})
}

func (h *CommonHandlers) GetAllProductsByLookID(c *fiber.Ctx) error {
	lookID := c.Params("look_id")

	if lookID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "LookID is required")
	}

	lookIDUint, err := strconv.ParseUint(lookID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting lookID to uint: %v", err))
	}

	// Parse pagination parameters from query params
	page := 1      // Default page
	pageSize := 10 // Default page size
	pageQ := c.Query("page")
	pageSizeQ := c.Query("pageSize")

	// Convert page and pageSize to integers if provided by the user
	if pageQ != "" {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
		}
	}
	if pageSizeQ != "" {
		pageSize, err = strconv.Atoi(pageSizeQ)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
		}
	}

	// Fetch products with pagination
	product, err := h.service.GetAllProductsByLookID(uint(lookIDUint), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch total number of products
	totalProducts, err := h.service.GetTotalLookProductsCount(uint(lookIDUint))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate total pages
	totalPages := (int(totalProducts) + pageSize - 1) / pageSize

	// Return response with pagination metadata
	return c.JSON(fiber.Map{
		"currentPage":     page,
		"totalPages":      totalPages,
		"totalProducts":   totalProducts,
		"productsPerPage": pageSize,
		"look":            product,
	})
}
