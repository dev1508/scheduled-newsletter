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

type SubscriberHandler struct {
	subscriberService service.SubscriberService
	logger            *zap.Logger
}

func NewSubscriberHandler(subscriberService service.SubscriberService, logger *zap.Logger) *SubscriberHandler {
	return &SubscriberHandler{
		subscriberService: subscriberService,
		logger:            logger,
	}
}

func (h *SubscriberHandler) CreateSubscriber(c *gin.Context) {
	var req request.CreateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	subscriber, err := h.subscriberService.CreateSubscriber(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "subscriber with email '"+req.Email+"' already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to create subscriber", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create subscriber",
		})
		return
	}

	c.JSON(http.StatusCreated, subscriber)
}

func (h *SubscriberHandler) GetSubscriber(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscriber ID format",
		})
		return
	}

	subscriber, err := h.subscriberService.GetSubscriber(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "subscriber not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		}

		h.logger.Error("Failed to get subscriber", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscriber",
		})
		return
	}

	c.JSON(http.StatusOK, subscriber)
}

func (h *SubscriberHandler) GetSubscriberByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email parameter is required",
		})
		return
	}

	subscriber, err := h.subscriberService.GetSubscriberByEmail(c.Request.Context(), email)
	if err != nil {
		if err.Error() == "subscriber not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		}

		h.logger.Error("Failed to get subscriber by email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscriber",
		})
		return
	}

	c.JSON(http.StatusOK, subscriber)
}

func (h *SubscriberHandler) ListSubscribers(c *gin.Context) {
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

	subscribers, err := h.subscriberService.ListSubscribers(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list subscribers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list subscribers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribers": subscribers,
		"limit":       limit,
		"offset":      offset,
	})
}

func (h *SubscriberHandler) UpdateSubscriber(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscriber ID format",
		})
		return
	}

	var req request.UpdateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	subscriber, err := h.subscriberService.UpdateSubscriber(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "subscriber not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		}

		h.logger.Error("Failed to update subscriber", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update subscriber",
		})
		return
	}

	c.JSON(http.StatusOK, subscriber)
}

func (h *SubscriberHandler) DeleteSubscriber(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscriber ID format",
		})
		return
	}

	err = h.subscriberService.DeleteSubscriber(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "subscriber not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		}

		h.logger.Error("Failed to delete subscriber", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete subscriber",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
