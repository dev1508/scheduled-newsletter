package request

import (
	"time"

	"github.com/google/uuid"
)

// CreateJobRequest represents the request payload for creating a scheduled job
type CreateJobRequest struct {
	ContentID   uuid.UUID `json:"content_id" binding:"required"`
	JobType     string    `json:"job_type" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}
