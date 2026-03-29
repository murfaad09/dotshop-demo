package dto

import "time"

type CreateCommentRequest struct {
	ReviewID uint   `json:"reviewId"`
	Content  string `json:"content"`
}

type CommentResponse struct {
	ID       uint   `json:"id"`
	ReviewID uint   `json:"reviewId"`
	Content  string `json:"content"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type Comment struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
