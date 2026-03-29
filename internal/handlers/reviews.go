package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/dto"
	model "github.com/harishash/dotshop-be/internal/models"
	service "github.com/harishash/dotshop-be/internal/services"
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService}
}

// CreateReview creates a Product Review of User
//
//	@Summary		Creates a Product Review of User
//	@Description	Creates a Product Review of User
//	@Tags			User Reviews
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateReviewRequest	true	"Review Details"
//	@Success		200		{object}	dto.ReviewResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		500		{object}	fiber.Error
//	@Router			/user/reviews [post]
func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	var review *dto.CreateReviewRequest
	if err := c.BodyParser(&review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := review.Validate(); err != nil {
		return errors.Wrap(err).WithMessage("Invalid request body")
	}
	response, err := h.reviewService.CreateReview(review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetReviewsByProductID gets all reviews of a product
//
//	@Summary		Get all reviews of a product
//	@Description	Get all reviews of a product
//	@Tags			User Reviews
//	@Accept			application/json
//	@Param			product_id	path		string	true	"Product ID"
//	@Success		200			{object}	dto.ReviewResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/user/products/{product_id}/reviews [get]
func (h *ReviewHandler) GetReviewsByProductID(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	reviews, avgRating, err := h.reviewService.GetReviewsByProductID(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch reviews"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"reviews":    reviews,
		"avg_rating": avgRating,
	})
}

// GetReviewsByProductIDAndCuratorID gets all reviews of a product by curator
//
//	@Summary		Get all reviews of a product by curator
//	@Description	Get all reviews of a product by curator
//	@Tags			User Reviews
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			product_id	path		string	true	"Product ID"
//	@Param			curator_id	path		int		true	"Curator ID"
//	@Success		200			{object}	dto.ReviewResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/user/products/{product_id}/reviews/curator/{curator_id} [get]
func (h *ReviewHandler) GetReviewsByProductIDAndCuratorID(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	curatorID, err := strconv.Atoi(c.Params("curator_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	reviews, avgRating, err := h.reviewService.GetReviewsByProductIDAndCuratorID(productID, uint(curatorID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch reviews"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"reviews":    reviews,
		"avg_rating": avgRating,
	})
}

func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	reviewID, err := strconv.Atoi(c.Params("review_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	userID, err := strconv.Atoi(c.Locals("userID").(string)) // Assuming userID is stored in context
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	updatedReview := new(model.Review) /////////////////////// todo
	if err := c.BodyParser(updatedReview); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	response, err := h.reviewService.UpdateReview(uint(reviewID), uint(userID), updatedReview)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Review updated successfully", "review": response})
}

func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	reviewID, err := strconv.Atoi(c.Params("review_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	userID, err := strconv.Atoi(c.Locals("userID").(string)) // Assuming userID is stored in context
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.reviewService.DeleteReview(uint(reviewID), uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Review deleted successfully"})
}
