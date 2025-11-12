package repo

import (
	"context"
	"fmt"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type contentRepo struct {
	db *db.DB
}

func NewContentRepository(database *db.DB) ContentRepository {
	return &contentRepo{
		db: database,
	}
}

func (r *contentRepo) Create(ctx context.Context, req *request.CreateContentRequest) (*models.Content, error) {
	query := `
		INSERT INTO content (topic_id, subject, body, send_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, topic_id, subject, body, send_at, status, created_at, updated_at
	`

	var content models.Content
	err := r.db.Pool.QueryRow(ctx, query, req.TopicID, req.Subject, req.Body, req.SendAt).Scan(
		&content.ID,
		&content.TopicID,
		&content.Subject,
		&content.Body,
		&content.SendAt,
		&content.Status,
		&content.CreatedAt,
		&content.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	return &content, nil
}

func (r *contentRepo) CreateTx(ctx context.Context, tx pgx.Tx, req *request.CreateContentRequest) (*models.Content, error) {
	query := `
		INSERT INTO content (topic_id, subject, body, send_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, topic_id, subject, body, send_at, status, created_at, updated_at
	`

	var content models.Content
	err := tx.QueryRow(ctx, query, req.TopicID, req.Subject, req.Body, req.SendAt).Scan(
		&content.ID,
		&content.TopicID,
		&content.Subject,
		&content.Body,
		&content.SendAt,
		&content.Status,
		&content.CreatedAt,
		&content.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create content in transaction: %w", err)
	}

	return &content, nil
}

func (r *contentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Content, error) {
	query := `
		SELECT id, topic_id, subject, body, send_at, status, created_at, updated_at
		FROM content
		WHERE id = $1
	`

	var content models.Content
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&content.ID,
		&content.TopicID,
		&content.Subject,
		&content.Body,
		&content.SendAt,
		&content.Status,
		&content.CreatedAt,
		&content.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("content not found")
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	return &content, nil
}

func (r *contentRepo) List(ctx context.Context, limit, offset int) ([]*models.Content, error) {
	query := `
		SELECT id, topic_id, subject, body, send_at, status, created_at, updated_at
		FROM content
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list content: %w", err)
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID,
			&content.TopicID,
			&content.Subject,
			&content.Body,
			&content.SendAt,
			&content.Status,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating content: %w", err)
	}

	return contents, nil
}

func (r *contentRepo) ListByTopic(ctx context.Context, topicID uuid.UUID, limit, offset int) ([]*models.Content, error) {
	query := `
		SELECT id, topic_id, subject, body, send_at, status, created_at, updated_at
		FROM content
		WHERE topic_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, topicID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list content by topic: %w", err)
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID,
			&content.TopicID,
			&content.Subject,
			&content.Body,
			&content.SendAt,
			&content.Status,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating content: %w", err)
	}

	return contents, nil
}

func (r *contentRepo) ListScheduled(ctx context.Context, limit int) ([]*models.Content, error) {
	query := `
		SELECT id, topic_id, subject, body, send_at, status, created_at, updated_at
		FROM content
		WHERE status = $1 AND send_at <= NOW()
		ORDER BY send_at ASC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, constants.ContentStatusScheduled, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduled content: %w", err)
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID,
			&content.TopicID,
			&content.Subject,
			&content.Body,
			&content.SendAt,
			&content.Status,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating content: %w", err)
	}

	return contents, nil
}

func (r *contentRepo) Update(ctx context.Context, id uuid.UUID, req *request.UpdateContentRequest) (*models.Content, error) {
	query := `
		UPDATE content
		SET subject = $2, body = $3, send_at = $4, updated_at = NOW()
		WHERE id = $1 AND status = $5
		RETURNING id, topic_id, subject, body, send_at, status, created_at, updated_at
	`

	var content models.Content
	err := r.db.Pool.QueryRow(ctx, query, id, req.Subject, req.Body, req.SendAt, constants.ContentStatusScheduled).Scan(
		&content.ID,
		&content.TopicID,
		&content.Subject,
		&content.Body,
		&content.SendAt,
		&content.Status,
		&content.CreatedAt,
		&content.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("content not found or cannot be updated (already sent)")
		}
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return &content, nil
}

func (r *contentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE content
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update content status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("content not found")
	}

	return nil
}

func (r *contentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM content 
		WHERE id = $1 AND status = $2
	`

	result, err := r.db.Pool.Exec(ctx, query, id, constants.ContentStatusScheduled)
	if err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("content not found or cannot be deleted (already sent)")
	}

	return nil
}
