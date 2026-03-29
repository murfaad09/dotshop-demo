package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"
)

type AdminDashboardHandler struct {
	service service.AdminDashboardService
}

func NewAdminDashboardHandler(service service.AdminDashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{service}
}

// GetTopSellingBrands Get Top Selling Brands
//
//	@Summary		Get top selling brands
//	@Description	Get top selling brands
//	@Tags			Admin Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/dashboard/top-selling-brands [get]

func (h *AdminDashboardHandler) GetTopSellingBrands(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetTopSellingBrands(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetTopSellingProducts Get Top Selling Products
//
//	@Summary		Get top selling products
//	@Description	Get top selling products
//	@Tags			Admin Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/dashboard/top-selling-products [get]

func (h *AdminDashboardHandler) GetTopSellingProducts(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetTopSellingProducts(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetTopCurators Get Top Curators
//
//	@Summary		Get top curators
//	@Description	Get top curators
//	@Tags			Admin Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/dashboard/top-curators [get]

func (h *AdminDashboardHandler) GetTopCurators(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetTopCurators(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetSaleByCategory Get Sale By Category
//
//	@Summary		Get sale by category
//	@Description	Get sale by category
//	@Tags			Admin Dashboard
//	@Security		BearerAuth
//	@Param			from	query		dto.TimeFilter	false	"Time"
//
//	@Success		200		{object}	[]dto.SaleByCategoryResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/admin/dashboard/sale-by-category [get]

func (h *AdminDashboardHandler) GetSaleByCategory(c *fiber.Ctx) error {
	query, err := parseQuery[dto.SaleRequest](c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetSaleByCategory(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetTopWishlistProducts Get Top Wishlist Products
//
//	@Summary		Get top wishlist products
//	@Description	Get top wishlist products
//	@Tags			Admin Dashboard
//	@Security		BearerAuth
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/dashboard/top-wishlist [get]

func (h *AdminDashboardHandler) GetTopWishlistProducts(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CommonProductRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetTopWishlistProducts(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetOrderSales Get Order Sales
//
//	@Summary		Get order sales
//	@Description	Get order sales
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			sort		query		string			false	"Sort by customer_name_asc, customer_name_desc, date_asc, date_desc, items_low_to_high, items_high_to_low, amount_low_to_high, amount_high_to_low"
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/sales/orders [get]
func (h *AdminDashboardHandler) GetOrderSales(c *fiber.Ctx) error {
	query, err := parseQuery[dto.OrderSalesRequest](c)
	if err != nil {
		return err
	}

	if len(query.SortBy) > 0 {
		if err := query.ValidateSortParam(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetOrderSales(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetOrderReturns Get Order Sales

//	@Summary		Get order sales
//	@Description	Get order sales
//	@Tags			AdminBO
//	@Security		BearerAuth
//
// //	@Param			sort		query		string			false	"Sort by customer_name_asc, customer_name_desc, date_asc, date_desc, items_low_to_high, items_high_to_low, amount_low_to_high, amount_high_to_low"
//
//	@Param			pageNum		query		int				false	"Page number"
//	@Param			pageSize	query		int				false	"Page size"
//	@Param			from		query		dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/sales/returns [get]
func (h *AdminDashboardHandler) GetOrderReturns(c *fiber.Ctx) error {
	query, err := parseQuery[dto.OrderSalesRequest](c)
	if err != nil {
		return err
	}

	if len(query.SortBy) > 0 {
		if err := query.ValidateSortParam(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetOrderReturns(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetAllCustomers Get All Customers

//	@Summary		Get All Customers
//	@Description	Get all customers
//	@Tags			AdminBO
//	@Security		BearerAuth
//
// //	@Param			sort		query	string			false	"Sort by customer_name_asc, customer_name_desc, date_asc, date_desc, items_low_to_high, items_high_to_low, amount_low_to_high, amount_high_to_low"
//
//	@Param			pageNum		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//
// //	@Param			from		query	dto.TimeFilter	false	"Time"
//
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/customers/all [get]
func (h *AdminDashboardHandler) GetAllCustomers(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CustomersRequest](c)
	if err != nil {
		return err
	}

	// if len(query.SortBy) > 0 {
	// 	if err := query.ValidateSortParam(); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// }

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.GetAllCustomers(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// CustomerOrdersList Get Consumer Order List
//
//	@Summary		Get order list by user id
//	@Description	This endpoint is used to get order list by user id
//	@Tags			AdminBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			user_id		path		int	true	"User id"
//	@Param			pageNum		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/customers/{user_id}/order-list [get]
func (h *AdminDashboardHandler) ConsumerOrdersListByUserId(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CustomersOrderListRequest](c)
	if err != nil {
		return err
	}

	userId, err := c.ParamsInt("user_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id parameter",
		})
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.ConsumerOrdersListByUserId(uint64(userId), query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// GetAllCurators Get All Curators
//
//	@Summary		Get all curators
//	@Description	Get all curators
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			pageNum		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/curators/all [get]
func (h *AdminDashboardHandler) GetAllCurators(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CuratorRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	curators, err := h.service.GetAllCurators(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch curators",
		})
	}
	return c.Status(fiber.StatusOK).JSON(curators)
}

// CuratorsOrdersListById Get Curators Order List
//
//	@Summary		Get curators order list
//	@Description	This endpoint is used to get curator order listWW
//	@Tags			AdminBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			curator_id	path		int	true	"Curator id"
//	@Param			pageNum		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	dto.Response
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/curators/{curator_id}/order-list [get]
func (h *AdminDashboardHandler) CuratorsOrdersListById(c *fiber.Ctx) error {
	query, err := parseQuery[dto.CuratorsOrderListRequest](c)
	if err != nil {
		return err
	}

	userId, err := c.ParamsInt("curator_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid curator id parameter",
		})
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	resp, err := h.service.CuratorOrdersListById(uint(userId), query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// UpdateReturnStatus Update Return Order Status
//
//	@Summary		Update return order status
//	@Description	Update return order status
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			return_id	path		int								true	"return ID"
//	@Param			body		body		dto.ReturnOrderStatusRequest	true	"status"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/admin/sales/return/{return_id}/status [patch]
func (h *AdminDashboardHandler) UpdateReturnStatus(c *fiber.Ctx) error {
	// Get the ID from the URL params
	id, err := c.ParamsInt("return_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid ID parameter",
		})
	}

	var request dto.ReturnOrderStatusRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to parse request body",
		})
	}

	if err := h.service.UpdateReturnStatus(uint(id), request.Status); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "status updated successfully",
	})
}

// DeleteUser Delete User
//
//	@Summary		Delete User
//	@Description	Deletes a user
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"user ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/customers/{id} [delete]
func (h *AdminDashboardHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.service.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete user, error : " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user deleted successfully"})
}

// BlockCustomer Block Customer
//
//	@Summary		Block Customer
//	@Description	Blocks a customer
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id		path		int							true	"customer ID"
//	@Param			body	body		dto.BlockCustomerRequest	true	"isBlock"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/customers/{id}/block [patch]
func (h *AdminDashboardHandler) BlockCustomer(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	var req dto.BlockCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if err := h.service.BlockCustomer(uint(userID), req.IsBlock); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "customer status updated successfull",
	})
}

// BlockCurator Block Curator
//
//	@Summary		Block Curator
//	@Description	Blocks a Curator
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Curator ID"
//	@Param			body	body		dto.BlockCuratorRequest	true	"isBlock"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/curators/{id}/block [patch]
func (h *AdminDashboardHandler) BlockCurator(c *fiber.Ctx) error {
	curatorIDParam := c.Params("id")
	curatorID, err := strconv.ParseUint(curatorIDParam, 10, 32)
	fmt.Println("curator id", curatorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid curator ID",
		})
	}
	var req dto.BlockCuratorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if err := h.service.BlockCurator(uint(curatorID), req.IsBlock); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "curator status updated successfull",
	})
}

// DeleteUser Delete Curator
//
//	@Summary		Delete Curator
//	@Description	Deletes a Curator
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"curator ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/curators/{id} [delete]
func (h *AdminDashboardHandler) DeleteCurator(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	err = h.service.DeleteCurator(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "curator not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete curator, error : " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "curator deleted successfully"})
}

// GetListedProducts Get Listed Products
//
//	@Summary		Get All Listed Products
//	@Description	Get all products listed by curator
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id	path		int	true	"curator ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/curators/{id}/listed-products [get]
func (p *AdminDashboardHandler) GetListedProducts(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	query, err := parseQuery[dto.PagingParams](c)
	if err != nil {
		return err
	}

	*query = dto.QueryToPagingParams(query)
	products, err := p.service.GetListedProducts(query, uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch products, error: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// DeleteUserReview Delete Review
//
//	@Summary		Delete Review
//	@Description	Delete the review of user
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/reviews/{id} [delete]
func (h *AdminDashboardHandler) DeleteUserReview(c *fiber.Ctx) error {
	reviewID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	if err := h.service.DeleteUserReview(uint(reviewID)); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "review not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete review, error : " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "review deleted successfully"})
}

// GetPaymentDistribution Get Payment Distribution
//
//	@Summary		Get Payment Distribution
//	@Description	Get Payment Distribution
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			pageNum			query		int		false	"Page number"
//	@Param			pageSize		query		int		false	"Page size"
//	@Param			commission_type	query		string	false	"Commission type"
//	@Success		200				{object}	dto.Response
//	@Failure		400				{object}	fiber.Error
//	@Failure		500				{object}	fiber.Error
//	@Router			/admin/financials/payment-distribution [get]
func (h *AdminDashboardHandler) GetPaymentDistribution(c *fiber.Ctx) error {
	query, err := parseQuery[dto.PaymentDistributionRequest](c)
	if err != nil {
		return err
	}

	query.PagingParams = dto.QueryToPagingParams(&query.PagingParams)
	response, err := h.service.GetPaymentDistribution(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch payment distribution, error: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
