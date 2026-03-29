package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/harishash/dotshop-be/internal/dto"
	_ "github.com/harishash/dotshop-be/internal/dto"
	_ "github.com/harishash/dotshop-be/internal/models"

	curatorstore_services "github.com/harishash/dotshop-be/internal/services"
)

type CuratorStoreFrontHandlers struct {
	commonService    curatorstore_services.CommonService
	curatorSFService curatorstore_services.CuratorStoreFrontService
}

func NewCuratorStoreFrontHandlers(commonService curatorstore_services.CommonService, curatorSFService curatorstore_services.CuratorStoreFrontService) *CuratorStoreFrontHandlers {
	return &CuratorStoreFrontHandlers{commonService: commonService, curatorSFService: curatorSFService}
}

// GetAllProducts get curator all feature products
//
//	@Summary		Get curator all feature products
//	@Description	Get curator all feature products
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			curator_id		path		int		true	"curator_id"
//	@Param			subCategoryIds	query		string	false	"sub category ids"
//	@Success		200				{object}	fiber.Map
//	@Failure		400				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/{curator_id}/allproducts [get]
func (h *CuratorStoreFrontHandlers) GetAllProducts(c *fiber.Ctx) error {
	return NewCommonHandlers(h.commonService).GetAllProducts(c)
}

// GetAllCollections get curator all collection
//
//	@Summary		Get curator all collection
//	@Description	Get curator all collection
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Param			curator_id	path		int	true	"curator_id"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/allcollections [get]
func (h *CuratorStoreFrontHandlers) GetAllCollections(c *fiber.Ctx) error {

	return NewCommonHandlers(h.commonService).GetAllCollections(c)
}

// GetCuratorAllLooks get curator all looks
//
//	@Summary		Get curator all looks
//	@Description	Get curator all looks
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Param			curator_id	path		int	true	"curator_id"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/alllooks [get]
func (h *CuratorStoreFrontHandlers) GetCuratorAllLooks(c *fiber.Ctx) error {

	return NewCommonHandlers(h.commonService).GetCuratorAllLooks(c)
}

// GetAllLooks get curator all looks
//
//	@Summary		Get curator all looks
//	@Description	Get curator all looks
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/looks [get]
func (h *CuratorStoreFrontHandlers) GetAllLooks(c *fiber.Ctx) error {

	return NewCommonHandlers(h.commonService).GetAllLooks(c)
}

// GetAllProductsByCollectionID  get curator all products by collection
//
//	@Summary		Get curator all products by collection
//	@Description	Get curator all products by collection
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			collection_id	path		int	true	"collection_id"
//	@Success		200				{object}	fiber.Map
//	@Failure		400				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collections/{collection_id}/products [get]
func (h *CuratorStoreFrontHandlers) GetAllProductsByCollectionID(c *fiber.Ctx) error {
	return NewCommonHandlers(h.commonService).GetAllProductsByCollectionID(c)
}

// GetAllProductsByLookID  get curator all products by look
//
//	@Summary		Get curator all products by look
//	@Description	Get curator all products by look
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			look_id	path		int	true	"look_id"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/looks/{look_id}/products [get]
func (h *CuratorStoreFrontHandlers) GetAllProductsByLookID(c *fiber.Ctx) error {
	return NewCommonHandlers(h.commonService).GetAllProductsByLookID(c)
}

// FetchSectionByID  get  section by id
//
//	@Summary		Get section by id
//	@Description	Get section by id
//	@Tags			CuratorSF
//	@Produce		application/json
//	@Param			section_id	path		int	true	"section_id"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/section/{section_id} [get]
func (h *CuratorStoreFrontHandlers) FetchSectionByID(c *fiber.Ctx) error {

	return NewCommonHandlers(h.commonService).FetchSectionByID(c)
}

//  func (h *CuratorStoreFrontHandlers) GetAllLooksByCuratorID(c *fiber.Ctx) error {

//  	return NewCommonHandlers(h.service).GetAllLooksByCuratorID(c)
//  }

// SearchLookByName Search Looks
//
//	@Summary		Search Looks by name
//	@Description	Search Looks by name
//	@Tags			Search
//	@Accept			application/json
//
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			look		query		string	true	"look name"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/look/search [get]
func (h *CuratorStoreFrontHandlers) SearchLookByName(c *fiber.Ctx) error {
	query, err := parseQuery[dto.SearchLookParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	looks, paging, err := h.curatorSFService.SearchLookByName(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   looks,
		Paging: *paging,
	})
}

// SearchProductsWithinCuratorLooks Search Products within Curator Looks
//
//	@Summary		Search product within curator looks
//	@Description	Search product within curator looks
//	@Tags			Search
//	@Accept			application/json
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			curator_id	path		string	true	"curator id"
//	@Param			product		query		string	true	"product name"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/look/search/products [get]
func (h *CuratorStoreFrontHandlers) SearchProductsWithinCuratorLooks(c *fiber.Ctx) error {
	curatorIDStr := c.Params("curator_id")
	if curatorIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "curator id is required")
	}
	curatorID, err := strconv.ParseUint(curatorIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.SearchProductParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	products, paging, err := h.curatorSFService.SearchProductsWithinCuratorLooks(uint(curatorID), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   products,
		Paging: *paging,
	})
}

// SearchFeatureProductsByName Search Feature Products By Name
//
//	@Summary		Search Feature Products By Name
//	@Description	Search Feature Products By Name
//	@Tags			Search
//	@Accept			application/json
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			curator_id	path		string	true	"curator id"
//	@Param			product		query		string	true	"product name"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/products/search [get]
func (h *CuratorStoreFrontHandlers) SearchFeatureProductsByName(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "curator id is required")
	}
	curatorid, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.SearchProductParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	products, paging, err := h.curatorSFService.SearchFeatureProductsByName(uint(curatorid), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   products,
		Paging: *paging,
	})
}

// SearchCollectionByName Search Collection
//
//	@Summary		Search collection by collection name
//	@Description	Search collection by collection name
//	@Tags			Search
//	@Accept			application/json
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			curator_id	path		string	true	"curator id"
//	@Param			collection	query		string	true	"collection name"
//
//	@Success		200			{object}	[]dto.CollectionWithProducts
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/collection/search [get]
func (h *CuratorStoreFrontHandlers) SearchCollectionByName(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "curator id is required")
	}
	curatorId, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.SearchCollectionParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	collections, paging, err := h.curatorSFService.SearchCollectionByName(query, uint(curatorId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   collections,
		Paging: *paging,
	})
}

// SearchCollectionProductByName Search Collection
//
//	@Summary		Search collection product by product name
//	@Description	Search collection product by product name
//	@Tags			Search
//	@Accept			application/json
//	@Param			pageNum			query		int		false	"Page number"
//	@Param			pageSize		query		int		false	"Page size"
//	@Param			collection_id	path		string	true	"collection id"
//	@Param			product			query		string	true	"product name"
//
//	@Success		200				{object}	[]dto.ProductResponse
//	@Failure		400				{object}	fiber.Error
//	@Failure		401				{object}	fiber.Error
//	@Failure		403				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id}/products/search [get]
func (h *CuratorStoreFrontHandlers) SearchCollectionProductByName(c *fiber.Ctx) error {
	collectionIdStr := c.Params("collection_id")
	if collectionIdStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "collection id is required")
	}
	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting collectionID to uint: %v", err))
	}

	query, err := parseQuery[dto.SearchProductParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	collections, paging, err := h.curatorSFService.SearchCollectionProductByName(collectionId, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   collections,
		Paging: *paging,
	})
}

// SearchSectionByName Search Section By Name
//
//	@Summary		Search section by name
//	@Description	Search section by name
//	@Tags			Search
//	@Accept			application/json
//
//	@Param			pageNum			query		int		false	"Page number"
//	@Param			pageSize		query		int		false	"Page size"
//	@Param			collection_id	path		string	true	"collection id"
//	@Param			section			query		string	true	"section name"
//
//	@Success		200				{object}	[]dto.Response
//	@Failure		400				{object}	fiber.Error
//	@Failure		401				{object}	fiber.Error
//	@Failure		403				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id}/section/search [get]
func (h *CuratorStoreFrontHandlers) SearchSectionByName(c *fiber.Ctx) error {
	query, err := parseQuery[dto.SearchSectionParams](c)
	if err != nil {
		return err
	}

	collectionIdStr := c.Params("collection_id")
	if collectionIdStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "collection id is required")
	}
	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting collectionID to uint: %v", err))
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	collectionSections, paging, err := h.curatorSFService.SearchSectionByName(collectionId, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Data:   collectionSections,
		Paging: *paging,
	})
}

// GlobalSearch Global Search
//
//	@Summary		Search product, brands, categories, and subcategories
//	@Description	This endpoint is used to search products, brands, categories, and subcategories
//	@Tags			Search
//	@Accept			application/json
//	@Param			pageNum		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			product		query		string	false	"product name"
//
//	@Success		200			{object}	dto.GlobalSearchResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/global-search [get]
func (h *CuratorStoreFrontHandlers) GlobalSearch(c *fiber.Ctx) error {
	query, err := parseQuery[dto.SearchProductParams](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.curatorSFService.GlobalSearch(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetDotShopCuratorID DotShop curator id
//
//	@Summary		This endpoint is used to get dotShop curator id
//	@Description	This endpoint is used to get dotShop curator id
//	@Tags			CuratorSF
//	@Accept			application/json
//	@Success		200	{object}	dto.GetStoreCuratorResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/dotshop/curator [get]
func (h *CuratorStoreFrontHandlers) GetDotShopCuratorID(c *fiber.Ctx) error {
	curator, err := h.curatorSFService.GetDotShopCuratorId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewGetStoreCuratorResponse(*curator))
}
