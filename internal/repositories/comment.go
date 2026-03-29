package repository

import (
	models "github.com/harishash/dotshop-be/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment *models.Comment) (*models.Comment, error)
	GetCommentByID(id uint) (*models.Comment, error)
	UpdateComment(comment *models.Comment) error
	DeleteComment(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) CreateComment(comment *models.Comment) (*models.Comment, error) {
	err := r.db.Create(&comment).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *commentRepository) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) UpdateComment(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) DeleteComment(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
