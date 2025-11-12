package handler

import (
	"net/http"
	"strconv"

	"newsletter-assignment/internal/request"
	"newsletter-assignment/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TopicHandler struct {
	topicService service.TopicService
	logger       *zap.Logger
}

func NewTopicHandler(topicService service.TopicService, logger *zap.Logger) *TopicHandler {
	return &TopicHandler{
		topicService: topicService,
		logger:       logger,
	}
}

func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var req request.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	topic, err := h.topicService.CreateTopic(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "topic with name '"+req.Name+"' already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to create topic", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create topic",
		})
		return
	}

	c.JSON(http.StatusCreated, topic)
}

func (h *TopicHandler) GetTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

	topic, err := h.topicService.GetTopic(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}

		h.logger.Error("Failed to get topic", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get topic",
		})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) ListTopics(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset parameter",
		})
		return
	}

	topics, err := h.topicService.ListTopics(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list topics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list topics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

	var req request.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	topic, err := h.topicService.UpdateTopic(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}

		h.logger.Error("Failed to update topic", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update topic",
		})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

	err = h.topicService.DeleteTopic(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}

		h.logger.Error("Failed to delete topic", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete topic",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
