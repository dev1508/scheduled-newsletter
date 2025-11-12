package request

import "github.com/google/uuid"

// CreateSubscriptionRequest represents the request payload for creating a subscription
type CreateSubscriptionRequest struct {
	SubscriberID uuid.UUID `json:"subscriber_id" binding:"required"`
	TopicID      uuid.UUID `json:"topic_id" binding:"required"`
}
