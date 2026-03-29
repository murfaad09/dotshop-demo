package service

import (
	"errors"

	"github.com/harishash/dotshop-be/internal/dto"
	review_repo "github.com/harishash/dotshop-be/internal/repositories"

	domain "github.com/harishash/dotshop-be/internal/models"
)

type ReviewService interface {
	CreateReview(review *dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	GetReviewsByProductID(productID string) ([]dto.ReviewResponse, float64, error)
	GetReviewsByProductIDAndCuratorID(productID string, curatorID uint) ([]dto.ReviewResponse, float64, error)
	UpdateReview(reviewID uint, userID uint, updatedReview *domain.Review) (*dto.ReviewResponse, error)
	DeleteReview(reviewID uint, userID uint) error
}

type reviewService struct {
	repo review_repo.ReviewRepository
}

func NewReviewService(repo review_repo.ReviewRepository) ReviewService {
	return &reviewService{repo}
}

func (s *reviewService) CreateReview(review *dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	// Check if the user has already reviewed the product for the curator
	existingReviews, err := s.repo.GetReviewsByProductIDAndCuratorID(review.ProductID, review.CuratorID)
	if err != nil {
		return nil, err
	}
	for _, existingReview := range existingReviews {
		if existingReview.UserID == review.UserID {
			return nil, errors.New("user has already reviewed this product for the curator")
		}
	}
	// If no existing review found, create the new review
	Review := domain.Review{
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CuratorID: review.CuratorID,
	}

	reviewreponse, err := s.repo.CreateReview(&Review)
	if err != nil {
		return nil, err
	}

	return &dto.ReviewResponse{
		ID:        reviewreponse.ID,
		ProductID: reviewreponse.ProductID,
		UserID:    reviewreponse.UserID,
		Rating:    reviewreponse.Rating,
		Comment:   reviewreponse.Comment,
		CuratorID: reviewreponse.CuratorID,
		CreatedAt: reviewreponse.CreatedAt,
		UpdatedAt: reviewreponse.UpdatedAt,
	}, nil
}
func (s *reviewService) GetReviewsByProductID(productID string) ([]dto.ReviewResponse, float64, error) {
	reviews, err := s.repo.GetReviewsByProductID(productID)
	if err != nil {
		return nil, 0, err
	}

	var reviewResponses []dto.ReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, dto.ReviewResponse{
			ID:        review.ID,
			ProductID: review.ProductID,
			UserID:    review.UserID,
			FirstName: review.User.FirstName,
			LastName:  review.User.LastName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CuratorID: review.CuratorID,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
			Comments:  processComments(review.Comments),
		})
	}
	averageRating := calculateAverageRating(reviewResponses)
	return reviewResponses, averageRating, nil
}

func (s *reviewService) GetReviewsByProductIDAndCuratorID(productID string, curatorID uint) ([]dto.ReviewResponse, float64, error) {
	reviews, err := s.repo.GetReviewsByProductIDAndCuratorID(productID, curatorID)
	if err != nil {
		return nil, 0, err
	}
	var reviewResponses []dto.ReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, dto.ReviewResponse{
			ID:        review.ID,
			ProductID: review.ProductID,
			UserID:    review.UserID,
			Rating:    review.Rating,
			CuratorID: review.CuratorID,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
			Comments:  processComments(review.Comments),
		})
	}
	averageRating := calculateAverageRating(reviewResponses)
	return reviewResponses, averageRating, nil
}

func (s *reviewService) UpdateReview(reviewID uint, userID uint, updatedReview *domain.Review) (*dto.ReviewResponse, error) {
	review, err := s.repo.GetReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	if review.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return s.repo.UpdateReview(reviewID, updatedReview)
}

func (s *reviewService) DeleteReview(reviewID uint, userID uint) error {
	review, err := s.repo.GetReviewByID(reviewID)
	if err != nil {
		return err
	}
	if review.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.DeleteReview(reviewID)
}

func calculateAverageRating(reviews []dto.ReviewResponse) float64 {
	if len(reviews) == 0 {
		return 0
	}
	sum := 0
	for _, review := range reviews {
		sum += review.Rating
	}
	return float64(sum) / float64(len(reviews))
}

func processComments(comments []*domain.Comment) []*dto.Comment {
	var commentResponses []*dto.Comment
	for _, comment := range comments {
		commentResponses = append(commentResponses, &dto.Comment{
			ID:        comment.ID,
			User:      processUser(comment.User),
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	return commentResponses
}

func processUser(user *domain.User) dto.User {
	var username, firstName, lastName string

	if user.Username != nil {
		username = *user.Username
	}

	if user.FirstName != nil {
		firstName = *user.FirstName
	}

	if user.LastName != nil {
		lastName = *user.LastName
	}

	return dto.User{
		ID:        user.ID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     user.Email,
	}
}
