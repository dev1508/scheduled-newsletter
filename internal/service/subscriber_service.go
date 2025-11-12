package service

import (
	"context"
	"fmt"
	"strings"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/repo"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type subscriberService struct {
	subscriberRepo repo.SubscriberRepository
	logger         *zap.Logger
}

func NewSubscriberService(subscriberRepo repo.SubscriberRepository, logger *zap.Logger) SubscriberService {
	return &subscriberService{
		subscriberRepo: subscriberRepo,
		logger:         logger,
	}
}

func (s *subscriberService) CreateSubscriber(ctx context.Context, req *request.CreateSubscriberRequest) (*models.Subscriber, error) {
	// Validate and sanitize input
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		req.Name = &name
		if *req.Name == "" {
			req.Name = nil
		}
	}

	s.logger.Info("Creating subscriber", zap.String("email", req.Email))

	subscriber, err := s.subscriberRepo.Create(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create subscriber", zap.Error(err), zap.String("email", req.Email))
		return nil, err
	}

	s.logger.Info("Subscriber created successfully",
		zap.String("id", subscriber.ID.String()),
		zap.String("email", subscriber.Email),
	)

	return subscriber, nil
}

func (s *subscriberService) GetSubscriber(ctx context.Context, id uuid.UUID) (*models.Subscriber, error) {
	subscriber, err := s.subscriberRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get subscriber", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return subscriber, nil
}

func (s *subscriberService) GetSubscriberByEmail(ctx context.Context, email string) (*models.Subscriber, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	subscriber, err := s.subscriberRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get subscriber by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	return subscriber, nil
}

func (s *subscriberService) ListSubscribers(ctx context.Context, limit, offset int) ([]*models.Subscriber, error) {
	// Set default and max limits
	if limit <= 0 {
		limit = constants.DefaultLimit
	}
	if limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}
	if offset < 0 {
		offset = constants.DefaultOffset
	}

	subscribers, err := s.subscriberRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list subscribers", zap.Error(err))
		return nil, err
	}

	return subscribers, nil
}

func (s *subscriberService) UpdateSubscriber(ctx context.Context, id uuid.UUID, req *request.UpdateSubscriberRequest) (*models.Subscriber, error) {
	// Validate and sanitize input
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		req.Name = &name
		if *req.Name == "" {
			req.Name = nil
		}
	}

	s.logger.Info("Updating subscriber", zap.String("id", id.String()), zap.String("email", req.Email))

	subscriber, err := s.subscriberRepo.Update(ctx, id, req)
	if err != nil {
		s.logger.Error("Failed to update subscriber", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	s.logger.Info("Subscriber updated successfully",
		zap.String("id", subscriber.ID.String()),
		zap.String("email", subscriber.Email),
	)

	return subscriber, nil
}

func (s *subscriberService) DeleteSubscriber(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting subscriber", zap.String("id", id.String()))

	err := s.subscriberRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete subscriber", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	s.logger.Info("Subscriber deleted successfully", zap.String("id", id.String()))
	return nil
}
