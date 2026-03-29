package handlers

import (
	"github.com/gofiber/fiber/v2"

	dto "github.com/harishash/dotshop-be/internal/dto"
	service "github.com/harishash/dotshop-be/internal/services"
)

type CommentHandler struct {
	commentService service.CommentService
	// productService service.IProductService
}

func NewCommentHandler(comment service.CommentService, product service.IProductService) *CommentHandler {
	return &CommentHandler{
		commentService: comment,
		// productService: product,
	}
}

// CreateComment Create Comment
//
//	@Summary		Create Comment
//	@Description	Create Comment
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			comment	body		dto.CreateCommentRequest	true	"Comment data"
//	@Success		201		{object}	dto.CommentResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/reviews/comments [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId, ok := c.Locals("user_id").(uint64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	resp, err := h.commentService.CreateComment(&req, userId)
	if err != nil {
		if err.Error() == "ERROR: insert or update on table \"comments\" violates foreign key constraint \"fk_reviews_comments\" (SQLSTATE 23503)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid review ID"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create comment, error: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateComment Update Comment
//
//	@Summary		Update Comment
//	@Description	Update Comment
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Comment ID"
//	@Param			comment	body		dto.UpdateCommentRequest	true	"Comment data"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/admin/reviews/comments/{id} [patch]
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	var req dto.UpdateCommentRequest
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment ID"})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.commentService.UpdateComment(&req, uint(id)); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "comment not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment updated successfully"})
}

// DeleteComment Delete Comment
//
//	@Summary		Delete Comment
//	@Description	Delete Comment
//	@Tags			AdminBO
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Comment ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/admin/reviews/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment ID"})
	}

	if err := h.commentService.DeleteComment(uint(id)); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "comment not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete comment, error: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "comment deleted successfully"})
}
