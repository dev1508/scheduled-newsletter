package service

import (
	"context"
	"fmt"
	"strings"

	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TopicService interface {
	CreateTopic(ctx context.Context, req *models.CreateTopicRequest) (*models.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*models.Topic, error)
	GetTopicByName(ctx context.Context, name string) (*models.Topic, error)
	ListTopics(ctx context.Context, limit, offset int) ([]*models.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, req *models.CreateTopicRequest) (*models.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
}

type topicService struct {
	topicRepo repo.TopicRepository
	logger    *zap.Logger
}

func NewTopicService(topicRepo repo.TopicRepository, logger *zap.Logger) TopicService {
	return &topicService{
		topicRepo: topicRepo,
		logger:    logger,
	}
}

func (s *topicService) CreateTopic(ctx context.Context, req *models.CreateTopicRequest) (*models.Topic, error) {
	// Validate and sanitize input
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("topic name cannot be empty")
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		req.Description = &desc
		if *req.Description == "" {
			req.Description = nil
		}
	}

	s.logger.Info("Creating topic", zap.String("name", req.Name))

	topic, err := s.topicRepo.Create(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create topic", zap.Error(err), zap.String("name", req.Name))
		return nil, err
	}

	s.logger.Info("Topic created successfully", 
		zap.String("id", topic.ID.String()),
		zap.String("name", topic.Name),
	)

	return topic, nil
}

func (s *topicService) GetTopic(ctx context.Context, id uuid.UUID) (*models.Topic, error) {
	topic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get topic", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return topic, nil
}

func (s *topicService) GetTopicByName(ctx context.Context, name string) (*models.Topic, error) {
	topic, err := s.topicRepo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("Failed to get topic by name", zap.Error(err), zap.String("name", name))
		return nil, err
	}

	return topic, nil
}

func (s *topicService) ListTopics(ctx context.Context, limit, offset int) ([]*models.Topic, error) {
	// Set default and max limits
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	topics, err := s.topicRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list topics", zap.Error(err))
		return nil, err
	}

	return topics, nil
}

func (s *topicService) UpdateTopic(ctx context.Context, id uuid.UUID, req *models.CreateTopicRequest) (*models.Topic, error) {
	// Validate and sanitize input
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("topic name cannot be empty")
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		req.Description = &desc
		if *req.Description == "" {
			req.Description = nil
		}
	}

	s.logger.Info("Updating topic", zap.String("id", id.String()), zap.String("name", req.Name))

	topic, err := s.topicRepo.Update(ctx, id, req)
	if err != nil {
		s.logger.Error("Failed to update topic", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	s.logger.Info("Topic updated successfully", 
		zap.String("id", topic.ID.String()),
		zap.String("name", topic.Name),
	)

	return topic, nil
}

func (s *topicService) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting topic", zap.String("id", id.String()))

	err := s.topicRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete topic", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	s.logger.Info("Topic deleted successfully", zap.String("id", id.String()))
	return nil
}
