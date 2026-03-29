package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/dto"

	service "github.com/harishash/dotshop-be/internal/services"
)

// GraphHandler handles the HTTP requests for graph data
type CuratorDashboardHandler struct {
	service service.CuratorDashboardService
}

func NewCuratorDashboardHandler(service service.CuratorDashboardService) *CuratorDashboardHandler {
	return &CuratorDashboardHandler{service}
}

// GetGraphDataForRevenue Get Graph Data for Revenue
//
//	@Summary		Get Graph Data for Revenue
//	@Description	Get Graph Data for Revenue
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/revenue [get]
func (h *CuratorDashboardHandler) GetGraphDataForRevenue(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForRevenue(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)

}

// GetGraphDataForSales Get Graph Data for Sales
//
//	@Summary		Get Graph Data for Sales
//	@Description	Get Graph Data for Sales
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/sales [get]
func (h *CuratorDashboardHandler) GetGraphDataForSales(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForSales(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)

}

// OrderGraphData Get Graph Data for Orders
//
//	@Summary		Get Graph Data for Orders
//	@Description	Get Graph Data for Orders
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	dto.GraphResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/orders [get]
func (h *CuratorDashboardHandler) OrderGraphData(c *fiber.Ctx) error {
	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForOrder(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)
}

// GetGraphDataForAOV Get Graph Data for AOV
//
//	@Summary		Get Graph Data for AOV
//	@Description	Get Graph Data for AOV
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	dto.GraphResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/average-order-value [get]
func (h *CuratorDashboardHandler) GetGraphDataForAOV(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForAvgOrderValue(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)
}

// GetGraphDataForAUPOrder Get Graph Data for AUPOrder
//
//	@Summary		Get Graph Data for AUPOrder
//	@Description	Get Graph Data for AUPOrder
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	dto.GraphResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/average-units-per-order [get]
func (h *CuratorDashboardHandler) GetGraphDataForAUPOrder(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForAUPOrder(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)
}

// GetGraphDataForUnits Get Graph Data for Units
//
//	@Summary		Get Graph Data for Units
//	@Description	Get Graph Data for Units
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Param			from		query		string	false	"From"
//	@Param			to			query		string	false	"To"
//	@Success		200			{object}	dto.GraphResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/units [get]
func (h *CuratorDashboardHandler) GetGraphDataForUnits(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForUnitsSold(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)
}

// GetCuratorTopWishlistProducts Get Curator Top Wishlist Products
//
//	@Summary		Get curator top wishlist products
//	@Description	Get curator top wishlist products
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/top-wishlist [get]
func (h *CuratorDashboardHandler) GetCuratorTopWishlistProducts(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}

	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetCuratorTopWishlist(uint(curatorUint), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetCuratorTopWishlistProducts Get Curator Top Selling Products
//
//	@Summary		Get curator top selling products
//	@Description	Get curator top selling products
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/top-selling-products [get]
func (h *CuratorDashboardHandler) GetCuratorTopSellingProducts(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}

	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetCuratorTopSellingProducts(uint(curatorUint), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetCuratorTopSellingBrands Get Curator Top Selling Brands
//
//	@Summary		Get curator top selling brands
//	@Description	Get curator top selling brands
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/top-selling-brands [get]
func (h *CuratorDashboardHandler) GetCuratorTopSellingBrands(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}

	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetCuratorTopSellingBrands(uint(curatorUint), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetCuratorTopPurchasers Get Curator Top Pruchasers
//
//	@Summary		Get curator top purchasers
//	@Description	Get curator top purchasers
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/top-purchasers [get]
func (h *CuratorDashboardHandler) GetCuratorTopPurchasers(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}

	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetCuratorTopPurchasers(uint(curatorUint), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetCuratorSaleByCategory Get Curator Sale By Category
//
//	@Summary		Get curator sale by category
//	@Description	Get curator sale by category
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//
//	@Success		200			{object}	[]dto.SaleByCategoryResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/sale-by-category [get]
func (h *CuratorDashboardHandler) GetCuratorSaleByCategory(c *fiber.Ctx) error {
	curatorID := c.Params("curator_id")
	if curatorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "CuratorID is required")
	}

	curatorUint, err := strconv.ParseUint(curatorID, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error converting curatorID to uint: %v", err))
	}

	query, err := parseQuery[dto.SaleRequest](c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetCuratorSaleByCategory(uint(curatorUint), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GraphDataForReturns Graph Data For Returns
//
//	@Summary		Graph data for returns
//	@Description	Graph data for returns
//	@Tags			Curator Dashboard
//	@Security		BearerAuth
//	@Param			from		query		dto.TimeFilter	false	"Time"
//	@Param			curator_id	path		string			true	"Curator Id"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/dashboard/returns [get]
func (h *CuratorDashboardHandler) GraphDataForReturns(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	graphData, err := h.service.GraphDataForReturns(
		uint(curatorID),
		query.From, query.TimeFilter.To)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"error": "Error fetching graph data",
			})
	}

	return c.JSON(graphData)
}
