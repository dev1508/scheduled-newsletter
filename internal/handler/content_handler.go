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

type ContentHandler struct {
	contentService service.ContentService
	logger         *zap.Logger
}

func NewContentHandler(contentService service.ContentService, logger *zap.Logger) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		logger:         logger,
	}
}

func (h *ContentHandler) CreateContent(c *gin.Context) {
	var req request.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	content, err := h.contentService.CreateContent(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "topic not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		case "subject cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Subject cannot be empty",
			})
			return
		case "body cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Body cannot be empty",
			})
			return
		case "send_at must be in the future":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Send time must be in the future",
			})
			return
		}

		h.logger.Error("Failed to create content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create content",
		})
		return
	}

	c.JSON(http.StatusCreated, content)
}

func (h *ContentHandler) GetContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content ID format",
		})
		return
	}

	content, err := h.contentService.GetContent(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "content not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Content not found",
			})
			return
		}

		h.logger.Error("Failed to get content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get content",
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *ContentHandler) ListContent(c *gin.Context) {
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

	contents, err := h.contentService.ListContent(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": contents,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *ContentHandler) ListContentByTopic(c *gin.Context) {
	topicIDStr := c.Param("id")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

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

	contents, err := h.contentService.ListContentByTopic(c.Request.Context(), topicID, limit, offset)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}

		h.logger.Error("Failed to list content by topic", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": contents,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *ContentHandler) UpdateContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content ID format",
		})
		return
	}

	var req request.UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	content, err := h.contentService.UpdateContent(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "content not found or cannot be updated (already sent)":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Content not found or cannot be updated (already sent)",
			})
			return
		case "subject cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Subject cannot be empty",
			})
			return
		case "body cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Body cannot be empty",
			})
			return
		case "send_at must be in the future":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Send time must be in the future",
			})
			return
		}

		h.logger.Error("Failed to update content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update content",
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *ContentHandler) DeleteContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content ID format",
		})
		return
	}

	err = h.contentService.DeleteContent(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "content not found or cannot be deleted (already sent)" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Content not found or cannot be deleted (already sent)",
			})
			return
		}

		h.logger.Error("Failed to delete content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete content",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ContentHandler) ScheduleContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content ID format",
		})
		return
	}

	err = h.contentService.ScheduleContent(c.Request.Context(), id)
	if err != nil {
		switch err.Error() {
		case "content not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Content not found",
			})
			return
		case "content is not in scheduled status":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Content is not in scheduled status",
			})
			return
		case "content send time has not arrived yet":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Content send time has not arrived yet",
			})
			return
		}

		h.logger.Error("Failed to schedule content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to schedule content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Content scheduled successfully",
	})
}
