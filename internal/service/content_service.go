package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/models"
	"newsletter-assignment/internal/repo"
	"newsletter-assignment/internal/request"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contentService struct {
	contentRepo repo.ContentRepository
	topicRepo   repo.TopicRepository
	logger      *zap.Logger
}

func NewContentService(
	contentRepo repo.ContentRepository,
	topicRepo repo.TopicRepository,
	logger *zap.Logger,
) ContentService {
	return &contentService{
		contentRepo: contentRepo,
		topicRepo:   topicRepo,
		logger:      logger,
	}
}

func (s *contentService) CreateContent(ctx context.Context, req *request.CreateContentRequest) (*models.Content, error) {
	// Validate and sanitize input
	req.Subject = strings.TrimSpace(req.Subject)
	if req.Subject == "" {
		return nil, fmt.Errorf("subject cannot be empty")
	}

	req.Body = strings.TrimSpace(req.Body)
	if req.Body == "" {
		return nil, fmt.Errorf("body cannot be empty")
	}

	// Validate send_at is in the future
	if req.SendAt.Before(time.Now()) {
		return nil, fmt.Errorf("send_at must be in the future")
	}

	// Validate that topic exists
	topic, err := s.topicRepo.GetByID(ctx, req.TopicID)
	if err != nil {
		s.logger.Error("Topic not found for content", zap.Error(err), zap.String("topic_id", req.TopicID.String()))
		return nil, fmt.Errorf("topic not found")
	}

	s.logger.Info("Creating content",
		zap.String("topic_name", topic.Name),
		zap.String("subject", req.Subject),
		zap.Time("send_at", req.SendAt),
	)

	content, err := s.contentRepo.Create(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create content", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Content created successfully",
		zap.String("id", content.ID.String()),
		zap.String("topic_name", topic.Name),
		zap.String("subject", content.Subject),
		zap.String("status", content.Status),
	)

	return content, nil
}

func (s *contentService) GetContent(ctx context.Context, id uuid.UUID) (*models.Content, error) {
	content, err := s.contentRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get content", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return content, nil
}

func (s *contentService) ListContent(ctx context.Context, limit, offset int) ([]*models.Content, error) {
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

	contents, err := s.contentRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list content", zap.Error(err))
		return nil, err
	}

	return contents, nil
}

func (s *contentService) ListContentByTopic(ctx context.Context, topicID uuid.UUID, limit, offset int) ([]*models.Content, error) {
	// Validate that topic exists
	_, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		s.logger.Error("Topic not found", zap.Error(err), zap.String("topic_id", topicID.String()))
		return nil, fmt.Errorf("topic not found")
	}

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

	contents, err := s.contentRepo.ListByTopic(ctx, topicID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list content by topic", zap.Error(err))
		return nil, err
	}

	return contents, nil
}

func (s *contentService) UpdateContent(ctx context.Context, id uuid.UUID, req *request.UpdateContentRequest) (*models.Content, error) {
	// Validate and sanitize input
	req.Subject = strings.TrimSpace(req.Subject)
	if req.Subject == "" {
		return nil, fmt.Errorf("subject cannot be empty")
	}

	req.Body = strings.TrimSpace(req.Body)
	if req.Body == "" {
		return nil, fmt.Errorf("body cannot be empty")
	}

	// Validate send_at is in the future
	if req.SendAt.Before(time.Now()) {
		return nil, fmt.Errorf("send_at must be in the future")
	}

	s.logger.Info("Updating content",
		zap.String("id", id.String()),
		zap.String("subject", req.Subject),
		zap.Time("send_at", req.SendAt),
	)

	content, err := s.contentRepo.Update(ctx, id, req)
	if err != nil {
		s.logger.Error("Failed to update content", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	s.logger.Info("Content updated successfully",
		zap.String("id", content.ID.String()),
		zap.String("subject", content.Subject),
	)

	return content, nil
}

func (s *contentService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting content", zap.String("id", id.String()))

	err := s.contentRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete content", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	s.logger.Info("Content deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *contentService) ScheduleContent(ctx context.Context, contentID uuid.UUID) error {
	// Get the content to validate it exists and is in correct state
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return err
	}

	if content.Status != constants.ContentStatusScheduled {
		return fmt.Errorf("content is not in scheduled status")
	}

	if content.SendAt.After(time.Now()) {
		return fmt.Errorf("content send time has not arrived yet")
	}

	s.logger.Info("Scheduling content for processing",
		zap.String("id", contentID.String()),
		zap.String("subject", content.Subject),
		zap.Time("send_at", content.SendAt),
	)

	// This is where we would typically enqueue the content to a job queue
	// For now, we'll just log that it's ready for scheduling
	// In future commits, this will integrate with Asynq

	s.logger.Info("Content ready for job scheduling", zap.String("id", contentID.String()))
	return nil
}
