package repository

import (
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateReview(review *domain.Review) (*domain.Review, error)
	GetReviewsByProductID(productID string) ([]*domain.Review, error)
	GetReviewsByProductIDAndCuratorID(productID string, curatorID uint) ([]*domain.Review, error)
	GetReviewByID(reviewID uint) (*dto.ReviewResponse, error)
	UpdateReview(reviewID uint, updatedReview *domain.Review) (*dto.ReviewResponse, error)
	DeleteReview(reviewID uint) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db}
}

func (r *reviewRepository) CreateReview(Review *domain.Review) (*domain.Review, error) {

	err := r.db.Create(&Review).Error
	if err != nil {
		return nil, err
	}

	return Review, nil
}

func (r *reviewRepository) GetReviewsByProductID(productID string) ([]*domain.Review, error) {
	var reviews []*domain.Review
	err := r.db.Preload("Comments").Preload("Comments.User").Where("product_id = ?", productID).Preload("User").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetReviewsByProductIDAndCuratorID(productID string, curatorID uint) ([]*domain.Review, error) {
	var reviews []*domain.Review
	err := r.db.Preload("Comments").Preload("Comments.User").Where("product_id = ? AND curator_id = ?", productID, curatorID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetReviewByID(reviewID uint) (*dto.ReviewResponse, error) {
	var review domain.Review
	err := r.db.First(&review, reviewID).Error
	if err != nil {
		return nil, err
	}

	return &dto.ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CuratorID: review.CuratorID,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}, nil
}
func (r *reviewRepository) UpdateReview(reviewID uint, updatedReview *domain.Review) (*dto.ReviewResponse, error) {
	err := r.db.Model(&domain.Review{}).Where("id = ?", reviewID).Updates(updatedReview).Error
	if err != nil {
		return nil, err
	}

	updatedReview.ID = reviewID // Ensure the ID is set in the updatedReview
	return &dto.ReviewResponse{
		ID:        updatedReview.ID,
		ProductID: updatedReview.ProductID,
		UserID:    updatedReview.UserID,
		Rating:    updatedReview.Rating,
		Comment:   updatedReview.Comment,
		CuratorID: updatedReview.CuratorID,
		CreatedAt: updatedReview.CreatedAt,
		UpdatedAt: updatedReview.UpdatedAt,
	}, nil
}

func (r *reviewRepository) DeleteReview(reviewID uint) error {
	return r.db.Delete(&domain.Review{}, reviewID).Error
}
