package repo

import (
	"context"
	"time"

	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TopicRepository defines the interface for topic data operations
type TopicRepository interface {
	Create(ctx context.Context, req *request.CreateTopicRequest) (*models.Topic, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error)
	GetByName(ctx context.Context, name string) (*models.Topic, error)
	List(ctx context.Context, limit, offset int) ([]*models.Topic, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateTopicRequest) (*models.Topic, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SubscriberRepository defines the interface for subscriber data operations
type SubscriberRepository interface {
	Create(ctx context.Context, req *request.CreateSubscriberRequest) (*models.Subscriber, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscriber, error)
	GetByEmail(ctx context.Context, email string) (*models.Subscriber, error)
	List(ctx context.Context, limit, offset int) ([]*models.Subscriber, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateSubscriberRequest) (*models.Subscriber, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	Create(ctx context.Context, req *request.CreateSubscriptionRequest) (*models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetBySubscriberAndTopic(ctx context.Context, subscriberID, topicID uuid.UUID) (*models.Subscription, error)
	ListBySubscriber(ctx context.Context, subscriberID uuid.UUID) ([]*models.Subscription, error)
	ListByTopic(ctx context.Context, topicID uuid.UUID) ([]*models.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, isActive bool) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ContentRepository defines the interface for content data operations
type ContentRepository interface {
	Create(ctx context.Context, req *request.CreateContentRequest) (*models.Content, error)
	CreateTx(ctx context.Context, tx pgx.Tx, req *request.CreateContentRequest) (*models.Content, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Content, error)
	List(ctx context.Context, limit, offset int) ([]*models.Content, error)
	ListByTopic(ctx context.Context, topicID uuid.UUID, limit, offset int) ([]*models.Content, error)
	ListScheduled(ctx context.Context, limit int) ([]*models.Content, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateContentRequest) (*models.Content, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}


// JobRepository defines the interface for job scheduler data operations
type JobRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, contentID uuid.UUID, jobType string, scheduledAt time.Time) (*models.JobScheduler, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.JobScheduler, error)
	GetPendingJobs(ctx context.Context, limit int) ([]*models.JobScheduler, error)
	List(ctx context.Context, limit, offset int) ([]*models.JobScheduler, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateStatusWithError(ctx context.Context, id uuid.UUID, status string, attempts int, errorMessage *string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// DeliveryRepository defines the interface for delivery data operations
type DeliveryRepository interface {
	CreateDelivery(ctx context.Context, contentID, subscriberID uuid.UUID, email, status string) (*models.Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status string, sentAt *time.Time, errorMessage *string) error
	GetDeliveryByContentAndSubscriber(ctx context.Context, contentID, subscriberID uuid.UUID) (*models.Delivery, error)
	ListDeliveriesByContent(ctx context.Context, contentID uuid.UUID) ([]*models.Delivery, error)
}
