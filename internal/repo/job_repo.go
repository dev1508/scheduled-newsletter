package repo

import (
	"context"
	"fmt"
	"time"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type jobRepo struct {
	db *db.DB
}

func NewJobRepository(database *db.DB) JobRepository {
	return &jobRepo{
		db: database,
	}
}

func (r *jobRepo) Create(ctx context.Context, req *request.CreateJobRequest) (*models.JobScheduler, error) {
	query := `
		INSERT INTO job_scheduler (content_id, job_type, scheduled_at)
		VALUES ($1, $2, $3)
		RETURNING id, content_id, job_type, scheduled_at, status, attempts, max_attempts, error_message, created_at, updated_at
	`

	var job models.JobScheduler
	err := r.db.Pool.QueryRow(ctx, query, req.ContentID, req.JobType, req.ScheduledAt).Scan(
		&job.ID,
		&job.ContentID,
		&job.JobType,
		&job.ScheduledAt,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &job, nil
}

func (r *jobRepo) CreateTx(ctx context.Context, tx pgx.Tx, contentID uuid.UUID, jobType string, scheduledAt time.Time) (*models.JobScheduler, error) {
	query := `
		INSERT INTO job_scheduler (content_id, job_type, scheduled_at)
		VALUES ($1, $2, $3)
		RETURNING id, content_id, job_type, scheduled_at, status, attempts, max_attempts, error_message, created_at, updated_at
	`

	var job models.JobScheduler
	err := tx.QueryRow(ctx, query, contentID, jobType, scheduledAt).Scan(
		&job.ID,
		&job.ContentID,
		&job.JobType,
		&job.ScheduledAt,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create job in transaction: %w", err)
	}

	return &job, nil
}

func (r *jobRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.JobScheduler, error) {
	query := `
		SELECT id, content_id, job_type, scheduled_at, status, attempts, max_attempts, error_message, created_at, updated_at
		FROM job_scheduler
		WHERE id = $1
	`

	var job models.JobScheduler
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.ContentID,
		&job.JobType,
		&job.ScheduledAt,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}

func (r *jobRepo) GetPendingJobs(ctx context.Context, limit int) ([]*models.JobScheduler, error) {
	query := `
		SELECT id, content_id, job_type, scheduled_at, status, attempts, max_attempts, error_message, created_at, updated_at
		FROM job_scheduler
		WHERE status = $1 AND scheduled_at <= NOW()
		ORDER BY scheduled_at ASC
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.db.Pool.Query(ctx, query, constants.JobStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.JobScheduler
	for rows.Next() {
		var job models.JobScheduler
		err := rows.Scan(
			&job.ID,
			&job.ContentID,
			&job.JobType,
			&job.ScheduledAt,
			&job.Status,
			&job.Attempts,
			&job.MaxAttempts,
			&job.ErrorMessage,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jobs: %w", err)
	}

	return jobs, nil
}

func (r *jobRepo) List(ctx context.Context, limit, offset int) ([]*models.JobScheduler, error) {
	query := `
		SELECT id, content_id, job_type, scheduled_at, status, attempts, max_attempts, error_message, created_at, updated_at
		FROM job_scheduler
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.JobScheduler
	for rows.Next() {
		var job models.JobScheduler
		err := rows.Scan(
			&job.ID,
			&job.ContentID,
			&job.JobType,
			&job.ScheduledAt,
			&job.Status,
			&job.Attempts,
			&job.MaxAttempts,
			&job.ErrorMessage,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jobs: %w", err)
	}

	return jobs, nil
}

func (r *jobRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE job_scheduler
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func (r *jobRepo) UpdateStatusWithError(ctx context.Context, id uuid.UUID, status string, attempts int, errorMessage *string) error {
	query := `
		UPDATE job_scheduler
		SET status = $2, attempts = $3, error_message = $4, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, status, attempts, errorMessage)
	if err != nil {
		return fmt.Errorf("failed to update job status with error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func (r *jobRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM job_scheduler WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}
