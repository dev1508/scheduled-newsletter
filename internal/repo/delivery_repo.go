package repo

import (
	"context"
	"fmt"
	"time"

	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/models"

	"github.com/google/uuid"
)

type deliveryRepo struct {
	db *db.DB
}

// NewDeliveryRepository creates a new delivery repository
func NewDeliveryRepository(database *db.DB) DeliveryRepository {
	return &deliveryRepo{
		db: database,
	}
}

// CreateDelivery creates a new delivery record
func (r *deliveryRepo) CreateDelivery(ctx context.Context, contentID, subscriberID uuid.UUID, email, status string) (*models.Delivery, error) {
	query := `
		INSERT INTO deliveries (content_id, subscriber_id, email, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, content_id, subscriber_id, email, status, sent_at, error_message, created_at, updated_at
	`

	var delivery models.Delivery
	err := r.db.Pool.QueryRow(ctx, query, contentID, subscriberID, email, status).Scan(
		&delivery.ID,
		&delivery.ContentID,
		&delivery.SubscriberID,
		&delivery.Email,
		&delivery.Status,
		&delivery.SentAt,
		&delivery.ErrorMessage,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create delivery: %w", err)
	}

	return &delivery, nil
}

// UpdateDeliveryStatus updates the delivery status
func (r *deliveryRepo) UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status string, sentAt *time.Time, errorMessage *string) error {
	query := `
		UPDATE deliveries 
		SET status = $2, sent_at = $3, error_message = $4, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id, status, sentAt, errorMessage)
	if err != nil {
		return fmt.Errorf("failed to update delivery status: %w", err)
	}

	return nil
}

// GetDeliveryByContentAndSubscriber gets delivery by content and subscriber
func (r *deliveryRepo) GetDeliveryByContentAndSubscriber(ctx context.Context, contentID, subscriberID uuid.UUID) (*models.Delivery, error) {
	query := `
		SELECT id, content_id, subscriber_id, email, status, sent_at, error_message, created_at, updated_at
		FROM deliveries
		WHERE content_id = $1 AND subscriber_id = $2
	`

	var delivery models.Delivery
	err := r.db.Pool.QueryRow(ctx, query, contentID, subscriberID).Scan(
		&delivery.ID,
		&delivery.ContentID,
		&delivery.SubscriberID,
		&delivery.Email,
		&delivery.Status,
		&delivery.SentAt,
		&delivery.ErrorMessage,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	return &delivery, nil
}

// ListDeliveriesByContent lists all deliveries for a content
func (r *deliveryRepo) ListDeliveriesByContent(ctx context.Context, contentID uuid.UUID) ([]*models.Delivery, error) {
	query := `
		SELECT id, content_id, subscriber_id, email, status, sent_at, error_message, created_at, updated_at
		FROM deliveries
		WHERE content_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []*models.Delivery
	for rows.Next() {
		var delivery models.Delivery
		err := rows.Scan(
			&delivery.ID,
			&delivery.ContentID,
			&delivery.SubscriberID,
			&delivery.Email,
			&delivery.Status,
			&delivery.SentAt,
			&delivery.ErrorMessage,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, &delivery)
	}

	return deliveries, nil
}
