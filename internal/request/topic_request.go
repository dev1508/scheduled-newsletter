package request

// CreateTopicRequest represents the request payload for creating a topic
type CreateTopicRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}

// UpdateTopicRequest represents the request payload for updating a topic
type UpdateTopicRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}
