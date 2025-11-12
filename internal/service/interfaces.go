package service

import (
	"context"

	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
)

// TopicService defines the interface for topic business logic
type TopicService interface {
	CreateTopic(ctx context.Context, req *request.CreateTopicRequest) (*models.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*models.Topic, error)
	GetTopicByName(ctx context.Context, name string) (*models.Topic, error)
	ListTopics(ctx context.Context, limit, offset int) ([]*models.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, req *request.UpdateTopicRequest) (*models.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
}

// SubscriberService defines the interface for subscriber business logic
type SubscriberService interface {
	CreateSubscriber(ctx context.Context, req *request.CreateSubscriberRequest) (*models.Subscriber, error)
	GetSubscriber(ctx context.Context, id uuid.UUID) (*models.Subscriber, error)
	GetSubscriberByEmail(ctx context.Context, email string) (*models.Subscriber, error)
	ListSubscribers(ctx context.Context, limit, offset int) ([]*models.Subscriber, error)
	UpdateSubscriber(ctx context.Context, id uuid.UUID, req *request.UpdateSubscriberRequest) (*models.Subscriber, error)
	DeleteSubscriber(ctx context.Context, id uuid.UUID) error
}

// SubscriptionService defines the interface for subscription business logic
type SubscriptionService interface {
	Subscribe(ctx context.Context, req *request.CreateSubscriptionRequest) (*models.Subscription, error)
	Unsubscribe(ctx context.Context, subscriberID, topicID uuid.UUID) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	ListSubscriberTopics(ctx context.Context, subscriberID uuid.UUID) ([]*models.Subscription, error)
	ListTopicSubscribers(ctx context.Context, topicID uuid.UUID) ([]*models.Subscription, error)
}

// ContentService defines the interface for content business logic
type ContentService interface {
	CreateContent(ctx context.Context, req *request.CreateContentRequest) (*models.Content, error)
	GetContent(ctx context.Context, id uuid.UUID) (*models.Content, error)
	ListContent(ctx context.Context, limit, offset int) ([]*models.Content, error)
	ListContentByTopic(ctx context.Context, topicID uuid.UUID, limit, offset int) ([]*models.Content, error)
	UpdateContent(ctx context.Context, id uuid.UUID, req *request.UpdateContentRequest) (*models.Content, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error
	ScheduleContent(ctx context.Context, contentID uuid.UUID) error
}

// DeliveryService defines the interface for delivery business logic
type DeliveryService interface {
	CreateDelivery(ctx context.Context, req *request.CreateDeliveryRequest) (*models.Delivery, error)
	GetDelivery(ctx context.Context, id uuid.UUID) (*models.Delivery, error)
	ListDeliveriesByContent(ctx context.Context, contentID uuid.UUID) ([]*models.Delivery, error)
	ListDeliveriesBySubscriber(ctx context.Context, subscriberID uuid.UUID, limit, offset int) ([]*models.Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status string) error
}

// JobService defines the interface for job scheduling business logic
type JobService interface {
	CreateJob(ctx context.Context, req *request.CreateJobRequest) (*models.JobScheduler, error)
	GetJob(ctx context.Context, id uuid.UUID) (*models.JobScheduler, error)
	ProcessPendingJobs(ctx context.Context) error
	UpdateJobStatus(ctx context.Context, id uuid.UUID, status string, attempts int, errorMessage *string) error
}
