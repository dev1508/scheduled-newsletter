package repo

import (
	"context"
	"fmt"

	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type topicRepo struct {
	db *db.DB
}

func NewTopicRepository(database *db.DB) TopicRepository {
	return &topicRepo{
		db: database,
	}
}

func (r *topicRepo) Create(ctx context.Context, req *request.CreateTopicRequest) (*models.Topic, error) {
	query := `
		INSERT INTO topics (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`

	var topic models.Topic
	err := r.db.Pool.QueryRow(ctx, query, req.Name, req.Description).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Description,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "topics_name_key" (SQLSTATE 23505)` {
			return nil, fmt.Errorf("topic with name '%s' already exists", req.Name)
		}
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return &topic, nil
}

func (r *topicRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM topics
		WHERE id = $1
	`

	var topic models.Topic
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Description,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	return &topic, nil
}

func (r *topicRepo) GetByName(ctx context.Context, name string) (*models.Topic, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM topics
		WHERE name = $1
	`

	var topic models.Topic
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Description,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	return &topic, nil
}

func (r *topicRepo) List(ctx context.Context, limit, offset int) ([]*models.Topic, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM topics
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Description,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating topics: %w", err)
	}

	return topics, nil
}

func (r *topicRepo) Update(ctx context.Context, id uuid.UUID, req *request.UpdateTopicRequest) (*models.Topic, error) {
	query := `
		UPDATE topics
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`

	var topic models.Topic
	err := r.db.Pool.QueryRow(ctx, query, id, req.Name, req.Description).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Description,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return &topic, nil
}

func (r *topicRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM topics WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("topic not found")
	}

	return nil
}
