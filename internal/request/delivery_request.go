package request

import "github.com/google/uuid"

// CreateDeliveryRequest represents the request payload for creating a delivery
type CreateDeliveryRequest struct {
	ContentID    uuid.UUID `json:"content_id" binding:"required"`
	SubscriberID uuid.UUID `json:"subscriber_id" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
}
