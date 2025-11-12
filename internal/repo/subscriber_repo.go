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

type subscriberRepo struct {
	db *db.DB
}

func NewSubscriberRepository(database *db.DB) SubscriberRepository {
	return &subscriberRepo{
		db: database,
	}
}

func (r *subscriberRepo) Create(ctx context.Context, req *request.CreateSubscriberRequest) (*models.Subscriber, error) {
	query := `
		INSERT INTO subscribers (email, name)
		VALUES ($1, $2)
		RETURNING id, email, name, is_active, created_at, updated_at
	`

	var subscriber models.Subscriber
	err := r.db.Pool.QueryRow(ctx, query, req.Email, req.Name).Scan(
		&subscriber.ID,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.IsActive,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "subscribers_email_key" (SQLSTATE 23505)` {
			return nil, fmt.Errorf("subscriber with email '%s' already exists", req.Email)
		}
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	return &subscriber, nil
}

func (r *subscriberRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscriber, error) {
	query := `
		SELECT id, email, name, is_active, created_at, updated_at
		FROM subscribers
		WHERE id = $1
	`

	var subscriber models.Subscriber
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&subscriber.ID,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.IsActive,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscriber not found")
		}
		return nil, fmt.Errorf("failed to get subscriber: %w", err)
	}

	return &subscriber, nil
}

func (r *subscriberRepo) GetByEmail(ctx context.Context, email string) (*models.Subscriber, error) {
	query := `
		SELECT id, email, name, is_active, created_at, updated_at
		FROM subscribers
		WHERE email = $1
	`

	var subscriber models.Subscriber
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&subscriber.ID,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.IsActive,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscriber not found")
		}
		return nil, fmt.Errorf("failed to get subscriber: %w", err)
	}

	return &subscriber, nil
}

func (r *subscriberRepo) List(ctx context.Context, limit, offset int) ([]*models.Subscriber, error) {
	query := `
		SELECT id, email, name, is_active, created_at, updated_at
		FROM subscribers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscribers: %w", err)
	}
	defer rows.Close()

	var subscribers []*models.Subscriber
	for rows.Next() {
		var subscriber models.Subscriber
		err := rows.Scan(
			&subscriber.ID,
			&subscriber.Email,
			&subscriber.Name,
			&subscriber.IsActive,
			&subscriber.CreatedAt,
			&subscriber.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscriber: %w", err)
		}
		subscribers = append(subscribers, &subscriber)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscribers: %w", err)
	}

	return subscribers, nil
}

func (r *subscriberRepo) Update(ctx context.Context, id uuid.UUID, req *request.UpdateSubscriberRequest) (*models.Subscriber, error) {
	query := `
		UPDATE subscribers
		SET email = $2, name = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, name, is_active, created_at, updated_at
	`

	var subscriber models.Subscriber
	err := r.db.Pool.QueryRow(ctx, query, id, req.Email, req.Name, req.IsActive).Scan(
		&subscriber.ID,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.IsActive,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscriber not found")
		}
		return nil, fmt.Errorf("failed to update subscriber: %w", err)
	}

	return &subscriber, nil
}

func (r *subscriberRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscribers WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscriber: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscriber not found")
	}

	return nil
}
