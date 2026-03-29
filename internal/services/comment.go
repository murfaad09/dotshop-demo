package service

import (
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	repository "github.com/harishash/dotshop-be/internal/repositories"
)

type CommentService interface {
	CreateComment(comment *dto.CreateCommentRequest, userID uint64) (*dto.CommentResponse, error)
	UpdateComment(comment *dto.UpdateCommentRequest, id uint) error
	DeleteComment(id uint) error
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo}
}

func (s *commentService) CreateComment(req *dto.CreateCommentRequest, userID uint64) (*dto.CommentResponse, error) {
	comment := &models.Comment{
		UserID:   uint(userID),
		ReviewID: req.ReviewID,
		Content:  req.Content,
	}

	createdComment, err := s.repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:       createdComment.ID,
		ReviewID: createdComment.ReviewID,
		Content:  createdComment.Content,
	}

	return response, nil
}

func (s *commentService) UpdateComment(req *dto.UpdateCommentRequest, id uint) error {
	comment, err := s.repo.GetCommentByID(id)
	if err != nil {
		return err
	}

	comment.Content = req.Content
	return s.repo.UpdateComment(comment)
}

func (s *commentService) DeleteComment(id uint) error {
	return s.repo.DeleteComment(id)
}
