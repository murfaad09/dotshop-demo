package dto

import "time"

type NotificationRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
