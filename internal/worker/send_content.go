package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/email"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/repo"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// subscriberData holds subscriber information for email sending
type subscriberData struct {
	ID    uuid.UUID
	Email string
}

// SendContentWorker handles sending newsletter content to subscribers
type SendContentWorker struct {
	contentRepo      repo.ContentRepository
	subscriptionRepo repo.SubscriptionRepository
	subscriberRepo   repo.SubscriberRepository
	jobRepo          repo.JobRepository
	deliveryRepo     repo.DeliveryRepository
	emailSender      *email.SMTPSender
	logger           *zap.Logger
}

// NewSendContentWorker creates a new send content worker
func NewSendContentWorker(
	contentRepo repo.ContentRepository,
	subscriptionRepo repo.SubscriptionRepository,
	subscriberRepo repo.SubscriberRepository,
	jobRepo repo.JobRepository,
	deliveryRepo repo.DeliveryRepository,
	emailSender *email.SMTPSender,
	logger *zap.Logger,
) *SendContentWorker {
	return &SendContentWorker{
		contentRepo:      contentRepo,
		subscriptionRepo: subscriptionRepo,
		subscriberRepo:   subscriberRepo,
		jobRepo:          jobRepo,
		deliveryRepo:     deliveryRepo,
		emailSender:      emailSender,
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

	// Filter active subscriptions and collect subscriber data
	activeSubscribers := 0
	subscribersData := make([]subscriberData, 0)

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
			subscribersData = append(subscribersData, subscriberData{
				ID:    subscriber.ID,
				Email: subscriber.Email,
			})
		}
	}

	// Send emails in parallel with actual SMTP sending
	w.sendEmailsInParallel(ctx, content, subscribersData)

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

// sendEmailsInParallel sends emails to multiple subscribers concurrently
func (w *SendContentWorker) sendEmailsInParallel(ctx context.Context, content *models.Content, subscribers []subscriberData) {
	const maxConcurrency = 20 // Increased concurrency for better performance

	// Create a semaphore to limit concurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	w.logger.Info("Starting parallel email sending",
		zap.Int("total_emails", len(subscribers)),
		zap.Int("max_concurrency", maxConcurrency),
	)

	for i, subscriber := range subscribers {
		wg.Add(1)

		go func(sub subscriberData, index int) {
			defer wg.Done()

			// Acquire semaphore (limit concurrency)
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Send email with actual SMTP
			w.sendSingleEmail(ctx, content, sub.Email, sub.ID, index)
		}(subscriber, i)
	}

	// Wait for all emails to be sent
	wg.Wait()

	w.logger.Info("Parallel email sending completed",
		zap.Int("total_emails", len(subscribers)),
	)
}

// sendSingleEmail sends an email to a single subscriber and tracks delivery
func (w *SendContentWorker) sendSingleEmail(ctx context.Context, content *models.Content, subscriberEmail string, subscriberID uuid.UUID, index int) {
	start := time.Now()

	// Create delivery record
	delivery, err := w.deliveryRepo.CreateDelivery(ctx, content.ID, subscriberID, subscriberEmail, constants.DeliveryStatusPending)
	if err != nil {
		w.logger.Error("Failed to create delivery record",
			zap.String("content_id", content.ID.String()),
			zap.String("subscriber_email", subscriberEmail),
			zap.Error(err),
		)
		return
	}

	// Prepare email request
	emailReq := &email.EmailRequest{
		To:       subscriberEmail,
		Subject:  content.Subject,
		HTMLBody: content.Body, // Assuming content.Body contains HTML
		TextBody: content.Body, // For now, use same content for text
	}

	// Send email via SMTP
	err = w.emailSender.Send(emailReq)

	// Update delivery status based on result
	now := time.Now()
	if err != nil {
		// Email failed
		errorMsg := err.Error()
		updateErr := w.deliveryRepo.UpdateDeliveryStatus(ctx, delivery.ID, constants.DeliveryStatusFailed, nil, &errorMsg)
		if updateErr != nil {
			w.logger.Error("Failed to update delivery status to failed", zap.Error(updateErr))
		}

		w.logger.Error("Failed to send email",
			zap.String("content_id", content.ID.String()),
			zap.String("recipient", subscriberEmail),
			zap.String("delivery_id", delivery.ID.String()),
			zap.Int("email_index", index),
			zap.Duration("send_duration", time.Since(start)),
			zap.Error(err),
		)
	} else {
		// Email sent successfully
		updateErr := w.deliveryRepo.UpdateDeliveryStatus(ctx, delivery.ID, constants.DeliveryStatusSent, &now, nil)
		if updateErr != nil {
			w.logger.Error("Failed to update delivery status to sent", zap.Error(updateErr))
		}

		w.logger.Info("Email sent successfully",
			zap.String("content_id", content.ID.String()),
			zap.String("recipient", subscriberEmail),
			zap.String("delivery_id", delivery.ID.String()),
			zap.Int("email_index", index),
			zap.Duration("send_duration", time.Since(start)),
		)
	}
}
