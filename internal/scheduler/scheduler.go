package scheduler

import (
	"context"
	"fmt"
	"time"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/queue"
	"newsletter-assignment/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewsletterJobPayload represents the payload for newsletter sending jobs
type NewsletterJobPayload struct {
	ContentID uuid.UUID `json:"content_id"`
	JobID     uuid.UUID `json:"job_id"`
}

// Scheduler handles the periodic processing of scheduled jobs
type Scheduler struct {
	jobRepo   repo.JobRepository
	queue     queue.Queue
	logger    *zap.Logger
	interval  time.Duration
	batchSize int
	stopCh    chan struct{}
}

// NewScheduler creates a new scheduler instance
func NewScheduler(
	jobRepo repo.JobRepository,
	queue queue.Queue,
	logger *zap.Logger,
	interval time.Duration,
	batchSize int,
) *Scheduler {
	return &Scheduler{
		jobRepo:   jobRepo,
		queue:     queue,
		logger:    logger,
		interval:  interval,
		batchSize: batchSize,
		stopCh:    make(chan struct{}),
	}
}

// Start begins the scheduler's job processing loop
func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("Starting job scheduler",
		zap.Duration("interval", s.interval),
		zap.Int("batch_size", s.batchSize),
	)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Process jobs immediately on start
	s.processJobs(ctx)

	for {
		select {
		case <-ticker.C:
			s.processJobs(ctx)
		case <-s.stopCh:
			s.logger.Info("Scheduler stopped")
			return
		case <-ctx.Done():
			s.logger.Info("Scheduler context cancelled")
			return
		}
	}
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stopCh)
}

// processJobs fetches pending jobs and enqueues them to Asynq
func (s *Scheduler) processJobs(ctx context.Context) {
	jobs, err := s.jobRepo.GetPendingJobs(ctx, s.batchSize)
	if err != nil {
		s.logger.Error("Failed to get pending jobs", zap.Error(err))
		return
	}

	if len(jobs) == 0 {
		s.logger.Debug("No pending jobs found")
		return
	}

	s.logger.Info("Processing pending jobs", zap.Int("count", len(jobs)))

	for _, job := range jobs {
		if err := s.processJob(ctx, job); err != nil {
			s.logger.Error("Failed to process job",
				zap.String("job_id", job.ID.String()),
				zap.Error(err),
			)
			
			// Update job status to failed with error message
			errorMsg := err.Error()
			updateErr := s.jobRepo.UpdateStatusWithError(
				ctx,
				job.ID,
				constants.JobStatusFailed,
				job.Attempts+1,
				&errorMsg,
			)
			if updateErr != nil {
				s.logger.Error("Failed to update job status", zap.Error(updateErr))
			}
		}
	}
}

// processJob processes a single job by enqueuing it to Asynq
func (s *Scheduler) processJob(ctx context.Context, job *models.JobScheduler) error {
	switch job.JobType {
	case constants.JobTypeSendNewsletter:
		return s.enqueueNewsletterJob(ctx, job)
	default:
		return fmt.Errorf("unknown job type: %s", job.JobType)
	}
}

// enqueueNewsletterJob enqueues a newsletter sending job to Asynq
func (s *Scheduler) enqueueNewsletterJob(ctx context.Context, job *models.JobScheduler) error {
	// Enqueue the task using queue
	info, err := s.queue.EnqueueSendContent(job.ContentID.String(), job.ID.String())
	if err != nil {
		return fmt.Errorf("failed to enqueue task to Asynq: %w", err)
	}

	s.logger.Info("Job enqueued to Asynq",
		zap.String("job_id", job.ID.String()),
		zap.String("content_id", job.ContentID.String()),
		zap.String("asynq_id", info.ID),
		zap.String("queue", info.Queue),
	)

	// Update job status to enqueued
	err = s.jobRepo.UpdateStatus(ctx, job.ID, constants.JobStatusEnqueued)
	if err != nil {
		s.logger.Error("Failed to update job status to enqueued",
			zap.String("job_id", job.ID.String()),
			zap.Error(err),
		)
		// Don't return error here as the job was successfully enqueued to Asynq
		// The status update failure is not critical
	}

	return nil
}

// GetStats returns scheduler statistics
func (s *Scheduler) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// This could be enhanced to return more detailed statistics
	return map[string]interface{}{
		"interval":    s.interval.String(),
		"batch_size":  s.batchSize,
		"status":      "running",
		"last_run":    time.Now(), // In a real implementation, you'd track this
	}, nil
}
