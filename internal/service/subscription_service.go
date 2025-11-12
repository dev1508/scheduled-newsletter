package service

import (
	"context"
	"fmt"

	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/repo"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type subscriptionService struct {
	subscriptionRepo repo.SubscriptionRepository
	subscriberRepo   repo.SubscriberRepository
	topicRepo        repo.TopicRepository
	logger           *zap.Logger
}

func NewSubscriptionService(
	subscriptionRepo repo.SubscriptionRepository,
	subscriberRepo repo.SubscriberRepository,
	topicRepo repo.TopicRepository,
	logger *zap.Logger,
) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		subscriberRepo:   subscriberRepo,
		topicRepo:        topicRepo,
		logger:           logger,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, req *request.CreateSubscriptionRequest) (*models.Subscription, error) {
	// Validate that subscriber exists
	subscriber, err := s.subscriberRepo.GetByID(ctx, req.SubscriberID)
	if err != nil {
		s.logger.Error("Subscriber not found for subscription", zap.Error(err), zap.String("subscriber_id", req.SubscriberID.String()))
		return nil, fmt.Errorf("subscriber not found")
	}

	if !subscriber.IsActive {
		return nil, fmt.Errorf("subscriber is not active")
	}

	// Validate that topic exists
	topic, err := s.topicRepo.GetByID(ctx, req.TopicID)
	if err != nil {
		s.logger.Error("Topic not found for subscription", zap.Error(err), zap.String("topic_id", req.TopicID.String()))
		return nil, fmt.Errorf("topic not found")
	}

	// Check if subscription already exists
	existingSubscription, err := s.subscriptionRepo.GetBySubscriberAndTopic(ctx, req.SubscriberID, req.TopicID)
	if err == nil {
		// Subscription exists
		if existingSubscription.IsActive {
			return nil, fmt.Errorf("subscriber is already subscribed to this topic")
		}
		// Reactivate existing subscription
		s.logger.Info("Reactivating existing subscription",
			zap.String("subscriber_id", req.SubscriberID.String()),
			zap.String("topic_id", req.TopicID.String()),
		)
		return s.subscriptionRepo.Update(ctx, existingSubscription.ID, true)
	}

	s.logger.Info("Creating new subscription",
		zap.String("subscriber_email", subscriber.Email),
		zap.String("topic_name", topic.Name),
	)

	subscription, err := s.subscriptionRepo.Create(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create subscription", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Subscription created successfully",
		zap.String("id", subscription.ID.String()),
		zap.String("subscriber_email", subscriber.Email),
		zap.String("topic_name", topic.Name),
	)

	return subscription, nil
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subscriberID, topicID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetBySubscriberAndTopic(ctx, subscriberID, topicID)
	if err != nil {
		s.logger.Error("Subscription not found for unsubscribe", zap.Error(err))
		return fmt.Errorf("subscription not found")
	}

	if !subscription.IsActive {
		return fmt.Errorf("subscription is already inactive")
	}

	s.logger.Info("Unsubscribing",
		zap.String("subscription_id", subscription.ID.String()),
		zap.String("subscriber_id", subscriberID.String()),
		zap.String("topic_id", topicID.String()),
	)

	// Deactivate subscription instead of deleting
	_, err = s.subscriptionRepo.Update(ctx, subscription.ID, false)
	if err != nil {
		s.logger.Error("Failed to unsubscribe", zap.Error(err))
		return err
	}

	s.logger.Info("Unsubscribed successfully", zap.String("subscription_id", subscription.ID.String()))
	return nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get subscription", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) ListSubscriberTopics(ctx context.Context, subscriberID uuid.UUID) ([]*models.Subscription, error) {
	// Validate that subscriber exists
	_, err := s.subscriberRepo.GetByID(ctx, subscriberID)
	if err != nil {
		s.logger.Error("Subscriber not found", zap.Error(err), zap.String("subscriber_id", subscriberID.String()))
		return nil, fmt.Errorf("subscriber not found")
	}

	subscriptions, err := s.subscriptionRepo.ListBySubscriber(ctx, subscriberID)
	if err != nil {
		s.logger.Error("Failed to list subscriber topics", zap.Error(err))
		return nil, err
	}

	return subscriptions, nil
}

func (s *subscriptionService) ListTopicSubscribers(ctx context.Context, topicID uuid.UUID) ([]*models.Subscription, error) {
	// Validate that topic exists
	_, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		s.logger.Error("Topic not found", zap.Error(err), zap.String("topic_id", topicID.String()))
		return nil, fmt.Errorf("topic not found")
	}

	subscriptions, err := s.subscriptionRepo.ListByTopic(ctx, topicID)
	if err != nil {
		s.logger.Error("Failed to list topic subscribers", zap.Error(err))
		return nil, err
	}

	return subscriptions, nil
}
