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

type subscriptionRepo struct {
	db *db.DB
}

func NewSubscriptionRepository(database *db.DB) SubscriptionRepository {
	return &subscriptionRepo{
		db: database,
	}
}

func (r *subscriptionRepo) Create(ctx context.Context, req *request.CreateSubscriptionRequest) (*models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (subscriber_id, topic_id)
		VALUES ($1, $2)
		RETURNING id, subscriber_id, topic_id, subscribed_at, is_active
	`

	var subscription models.Subscription
	err := r.db.Pool.QueryRow(ctx, query, req.SubscriberID, req.TopicID).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.TopicID,
		&subscription.SubscribedAt,
		&subscription.IsActive,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "subscriptions_subscriber_id_topic_id_key" (SQLSTATE 23505)` {
			return nil, fmt.Errorf("subscription already exists for this subscriber and topic")
		}
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, subscriber_id, topic_id, subscribed_at, is_active
		FROM subscriptions
		WHERE id = $1
	`

	var subscription models.Subscription
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.TopicID,
		&subscription.SubscribedAt,
		&subscription.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) GetBySubscriberAndTopic(ctx context.Context, subscriberID, topicID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, subscriber_id, topic_id, subscribed_at, is_active
		FROM subscriptions
		WHERE subscriber_id = $1 AND topic_id = $2
	`

	var subscription models.Subscription
	err := r.db.Pool.QueryRow(ctx, query, subscriberID, topicID).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.TopicID,
		&subscription.SubscribedAt,
		&subscription.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) ListBySubscriber(ctx context.Context, subscriberID uuid.UUID) ([]*models.Subscription, error) {
	query := `
		SELECT id, subscriber_id, topic_id, subscribed_at, is_active
		FROM subscriptions
		WHERE subscriber_id = $1 AND is_active = true
		ORDER BY subscribed_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, subscriberID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions by subscriber: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.SubscriberID,
			&subscription.TopicID,
			&subscription.SubscribedAt,
			&subscription.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *subscriptionRepo) ListByTopic(ctx context.Context, topicID uuid.UUID) ([]*models.Subscription, error) {
	query := `
		SELECT id, subscriber_id, topic_id, subscribed_at, is_active
		FROM subscriptions
		WHERE topic_id = $1 AND is_active = true
		ORDER BY subscribed_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions by topic: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.SubscriberID,
			&subscription.TopicID,
			&subscription.SubscribedAt,
			&subscription.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, isActive bool) (*models.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET is_active = $2
		WHERE id = $1
		RETURNING id, subscriber_id, topic_id, subscribed_at, is_active
	`

	var subscription models.Subscription
	err := r.db.Pool.QueryRow(ctx, query, id, isActive).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.TopicID,
		&subscription.SubscribedAt,
		&subscription.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}
