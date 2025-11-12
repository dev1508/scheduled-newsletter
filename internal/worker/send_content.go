package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/repo"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// SendContentWorker handles sending newsletter content to subscribers
type SendContentWorker struct {
	contentRepo      repo.ContentRepository
	subscriptionRepo repo.SubscriptionRepository
	subscriberRepo   repo.SubscriberRepository
	jobRepo          repo.JobRepository
	logger           *zap.Logger
}

// NewSendContentWorker creates a new send content worker
func NewSendContentWorker(
	contentRepo repo.ContentRepository,
	subscriptionRepo repo.SubscriptionRepository,
	subscriberRepo repo.SubscriberRepository,
	jobRepo repo.JobRepository,
	logger *zap.Logger,
) *SendContentWorker {
	return &SendContentWorker{
		contentRepo:      contentRepo,
		subscriptionRepo: subscriptionRepo,
		subscriberRepo:   subscriberRepo,
		jobRepo:          jobRepo,
		logger:           logger,
	}
}

// HandleSendContent processes the send content task
func (w *SendContentWorker) HandleSendContent(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		ContentID string `json:"content_id"`
		JobID     string `json:"job_id"`
	}

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		w.logger.Error("Failed to unmarshal task payload", zap.Error(err))
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	contentID, err := uuid.Parse(payload.ContentID)
	if err != nil {
		w.logger.Error("Invalid content ID", zap.String("content_id", payload.ContentID), zap.Error(err))
		return fmt.Errorf("invalid content ID: %w", err)
	}

	jobID, err := uuid.Parse(payload.JobID)
	if err != nil {
		w.logger.Error("Invalid job ID", zap.String("job_id", payload.JobID), zap.Error(err))
		return fmt.Errorf("invalid job ID: %w", err)
	}

	w.logger.Info("Processing send content task",
		zap.String("content_id", contentID.String()),
		zap.String("job_id", jobID.String()),
	)

	// Fetch content
	content, err := w.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		w.logger.Error("Failed to fetch content", zap.String("content_id", contentID.String()), zap.Error(err))
		
		// Update job status to failed
		errorMsg := fmt.Sprintf("Failed to fetch content: %v", err)
		w.jobRepo.UpdateStatusWithError(ctx, jobID, constants.JobStatusFailed, 1, &errorMsg)
		
		return fmt.Errorf("failed to fetch content: %w", err)
	}

	// Fetch active subscriptions for this topic
	subscriptions, err := w.subscriptionRepo.ListByTopic(ctx, content.TopicID)
	if err != nil {
		w.logger.Error("Failed to fetch subscriptions", 
			zap.String("topic_id", content.TopicID.String()), 
			zap.Error(err),
		)
		
		// Update job status to failed
		errorMsg := fmt.Sprintf("Failed to fetch subscriptions: %v", err)
		w.jobRepo.UpdateStatusWithError(ctx, jobID, constants.JobStatusFailed, 1, &errorMsg)
		
		return fmt.Errorf("failed to fetch subscriptions: %w", err)
	}

	// Filter active subscriptions and get subscriber count
	activeSubscribers := 0
	subscriberEmails := make([]string, 0)

	for _, subscription := range subscriptions {
		if subscription.IsActive {
			subscriber, err := w.subscriberRepo.GetByID(ctx, subscription.SubscriberID)
			if err != nil {
				w.logger.Warn("Failed to fetch subscriber", 
					zap.String("subscriber_id", subscription.SubscriberID.String()),
					zap.Error(err),
				)
				continue
			}
			activeSubscribers++
			subscriberEmails = append(subscriberEmails, subscriber.Email)
		}
	}

	// Log the dry-run execution (no actual email sending yet)
	w.logger.Info("Would send content to subscribers",
		zap.String("content_id", content.ID.String()),
		zap.String("subject", content.Subject),
		zap.String("topic_id", content.TopicID.String()),
		zap.Int("subscriber_count", activeSubscribers),
		zap.Strings("subscriber_emails", subscriberEmails),
		zap.Time("scheduled_at", content.SendAt),
	)

	// Simulate processing time and success
	w.logger.Info("Content processing completed successfully",
		zap.String("content_id", content.ID.String()),
		zap.Int("emails_sent", activeSubscribers),
	)

	// Update content status to sent
	if err := w.contentRepo.UpdateStatus(ctx, contentID, constants.ContentStatusSent); err != nil {
		w.logger.Error("Failed to update content status", zap.Error(err))
		// Don't return error here as the main processing was successful
	}

	// Update job status to completed
	if err := w.jobRepo.UpdateStatus(ctx, jobID, constants.JobStatusCompleted); err != nil {
		w.logger.Error("Failed to update job status", zap.Error(err))
		// Don't return error here as the main processing was successful
	}

	w.logger.Info("Send content task completed successfully",
		zap.String("content_id", contentID.String()),
		zap.String("job_id", jobID.String()),
		zap.Int("subscribers_notified", activeSubscribers),
	)

	return nil
}
