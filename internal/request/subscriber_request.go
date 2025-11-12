package request

// CreateSubscriberRequest represents the request payload for creating a subscriber
type CreateSubscriberRequest struct {
	Email string  `json:"email" binding:"required,email,max=255"`
	Name  *string `json:"name" binding:"omitempty,max=255"`
}

// UpdateSubscriberRequest represents the request payload for updating a subscriber
type UpdateSubscriberRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Name     *string `json:"name" binding:"omitempty,max=255"`
	IsActive bool   `json:"is_active"`
}
