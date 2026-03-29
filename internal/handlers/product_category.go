package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	category_services "github.com/harishash/dotshop-be/internal/services"
)

type CategoryHandler struct {
	Service category_services.CategoryService
}

func NewCategoryHandler(Service category_services.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service}
}

// AddNewCategory adds a new category
//
//	@Summary		adds a new category
//	@Description	adds a new category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateCategoryRequest	true	"Product Category"
//	@Success		200		{object}	dto.CategoryResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/admin/categories [post]
func (h *CategoryHandler) AddNewCategory(c *fiber.Ctx) error {

	var categoryReq dto.CreateCategoryRequest
	if err := c.BodyParser(&categoryReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	category, err := h.Service.CreateCategory(&categoryReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateExistingCategory updates an existing category
//
//	@Summary		updates an existing category
//	@Description	updates an existing category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id		path		int							true	"Category ID"
//	@Param			body	body		dto.UpdateCategoryRequest	true	"Product Category"
//	@Success		200		{object}	dto.CategoryResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/admin/categories/{id} [patch]
func (h *CategoryHandler) UpdateExistingCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	req := new(dto.UpdateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	category, err := h.Service.UpdateCategory(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// DeleteCategory deletes a category
//
//	@Summary		deletes a category
//	@Description	deletes a category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	category, err := h.Service.DeleteCategory(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Category deleted with products successfully",
		"category": category,
	})
}

// GetAllCategories returns all categories
//
//	@Summary
//	@Description
//	@Tags		AdminBO
//	@Accept		application/json
//
//	@Param		pageNum		query		int	false	"Page number"
//	@Param		pageSize	query		int	false	"Page size"
//
//	@Success	200			{object}	fiber.Map
//	@Failure	400			{object}	fiber.Error
//	@Failure	401			{object}	fiber.Error
//	@Failure	403			{object}	fiber.Error
//	@Failure	404			{object}	fiber.Error
//	@Router		/admin/categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit number"})
	}

	categories, total, err := h.Service.GetAllCategories(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// GetCategoryByID gets a category by ID
//
//	@Summary		Get category by ID
//	@Description	Get category by ID
//	@Tags			AdminBO
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	dto.CategoryResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}
	category, err := h.Service.GetCategory(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(category)
}

// CreateSubCategory creates a new sub-category
//
//	@Summary		Create sub-category
//	@Description	Create sub-category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int								true	"Category ID"
//	@Param			body		body		dto.CreateSubCategoryRequest	true	"Sub-category details"
//	@Success		201			{object}	dto.CreateSubCategoryResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/admin/categories/{category_id}/subcategories [post]
func (h *CategoryHandler) CreateSubCategory(c *fiber.Ctx) error {

	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	var request dto.CreateSubCategoryRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	subCategory, err := h.Service.AddSubCategory(uint(categoryID), request.Name, request.Products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(subCategory)
}

// AddProductInSubCategory adds products to a sub-category
//
//	@Summary		Add products to sub-category
//	@Description	Add products to sub-category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int							true	"Category ID"
//	@Param			id			path		int							true	"Sub-category ID"
//	@Param			body		body		[]dto.CreateProductRequest	true	"Products"
//	@Success		200			{object}	dto.CreateSubCategoryResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/admin/categories/{category_id}/subcategories/{id} [post]
func (h *CategoryHandler) AddProductInSubCategory(c *fiber.Ctx) error {

	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	subCategoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sub-category ID"})
	}

	var request []*dto.CreateProductRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	products, err := h.Service.UpdateSubCategory(uint(categoryID), uint(subCategoryID), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(products)
}

// UpdateProductInCategory updates the products in a category
//
//	@Summary		Updates the products in a category
//	@Description	Updates the products in a category
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ChangeCategoryRequest	true	"Product IDs"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/categories/products/{product_id} [put]
func (h *CategoryHandler) UpdateProductInCategory(c *fiber.Ctx) error {

	productIdString := c.Params("product_id")

	if productIdString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product IDs"})
	}

	productIDs := strings.Split(productIdString, ",")

	for i, id := range productIDs {
		productIDs[i] = strings.TrimSpace(id)
	}

	if productIDs[0] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product IDs"})
	}

	var request dto.ChangeCategoryRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.Service.UpdateProductCategory(productIDs, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// UpdateSingleProduct updates single product
//
//	@Summary		updates single product
//	@Description	updates single product
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		string							true	"Product ID"
//	@Param			body		body		dto.UpdateSingleProductRequest	true	"Product IDs"
//	@Success		200			{object}	dto.UpdateSingleProductResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/admin/categories/product/{product_id} [put]
func (h *CategoryHandler) UpdateSingleProduct(c *fiber.Ctx) error {

	productId := c.Params("product_id")
	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product IDs"})
	}

	var request dto.UpdateSingleProductRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.Service.UpdateSingleProduct(productId, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// UpdateExistingSubCategory Updates An Existing Subcategory
//
//	@Summary		updates an existing subcategory
//	@Description	updates an existing subcategory
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			category_id	path		int								true	"Category ID"
//	@Param			id			path		int								true	"Subcategory ID"
//	@Param			body		body		dto.UpdateSubcategoryRequest	true	"Update Subcategory Request"
//	@Success		200			{object}	dto.CategoryResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/categories/{category_id}/subcategories/{id} [patch]
func (h *CategoryHandler) UpdateExistingSubcategory(c *fiber.Ctx) error {
	categoryId, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	subcategoryId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subcategory ID"})
	}

	req := new(dto.UpdateSubcategoryRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	category, err := h.Service.UpdateSubcategoryName(uint(categoryId), uint(subcategoryId), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// DeleteSubcategory deletes a subcategory
//
//	@Summary		deletes a subcategory
//	@Description	deletes a subcategory
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id	path		int	true	"Subcategory ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/categories/{category_id}/subcategories/{id} [delete]
func (h *CategoryHandler) DeleteSubcategory(c *fiber.Ctx) error {

	subCategoryID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subcategory ID"})
	}

	subCategory, err := h.Service.DeleteSubCategory(uint(subCategoryID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Subcategory deleted successfully with associated Products",
		"subCategory": subCategory,
	})
}

// GetCatalogProducts get catalog products
//
//	@Summary		Get catalog products
//	@Description	Get catalog products
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			product				query	string	false	"Product name"
//	@Param			sort				query	string	false	"Sort by product_name_asc, product_name_desc, brand_name_asc, brand_name_desc, discount_low_to_high, discount_high_to_low, price_low_to_high, price_high_to_low"
//	@Param			subCategoryIds		query	string	false	"sub category ids"
//	@Param			searchByBrandName	query	string	false	"brand name"
//	@Param			pageNum				query	int		false	"Page number"
//	@Param			pageSize			query	int		false	"Page size"
//	@Produce		json
//	@Success		200	{object}	dto.Response
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/admin/catalog/products [get]
func (h *CategoryHandler) GetCatalogProducts(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CatalogProducts](c)
	if err != nil {
		return err
	}

	if len(query.SortBy) > 0 {
		if err := query.ValidateSortParam(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	category, err := h.Service.GetCatalogProducts(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}
