package request

import (
	"time"

	"github.com/google/uuid"
)

// CreateContentRequest represents the request payload for creating content
type CreateContentRequest struct {
	TopicID uuid.UUID `json:"topic_id" binding:"required"`
	Subject string    `json:"subject" binding:"required,min=1,max=500"`
	Body    string    `json:"body" binding:"required,min=1"`
	SendAt  time.Time `json:"send_at" binding:"required"`
}

// UpdateContentRequest represents the request payload for updating content
type UpdateContentRequest struct {
	Subject string    `json:"subject" binding:"required,min=1,max=500"`
	Body    string    `json:"body" binding:"required,min=1"`
	SendAt  time.Time `json:"send_at" binding:"required"`
}
