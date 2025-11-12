package handler

import (
	"net/http"

	"newsletter-assignment/internal/request"
	"newsletter-assignment/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	logger              *zap.Logger
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req request.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	subscription, err := h.subscriptionService.Subscribe(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "subscriber not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		case "topic not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		case "subscriber is not active":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Subscriber is not active",
			})
			return
		case "subscriber is already subscribed to this topic":
			c.JSON(http.StatusConflict, gin.H{
				"error": "Subscriber is already subscribed to this topic",
			})
			return
		}

		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create subscription",
		})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	subscriberIDStr := c.Param("subscriber_id")
	subscriberID, err := uuid.Parse(subscriberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscriber ID format",
		})
		return
	}

	topicIDStr := c.Param("id")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

	err = h.subscriptionService.Unsubscribe(c.Request.Context(), subscriberID, topicID)
	if err != nil {
		switch err.Error() {
		case "subscription not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscription not found",
			})
			return
		case "subscription is already inactive":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Subscription is already inactive",
			})
			return
		}

		h.logger.Error("Failed to unsubscribe", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unsubscribe",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unsubscribed",
	})
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscription ID format",
		})
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscription not found",
			})
			return
		}

		h.logger.Error("Failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscription",
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) ListSubscriberTopics(c *gin.Context) {
	subscriberIDStr := c.Param("id")
	subscriberID, err := uuid.Parse(subscriberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid subscriber ID format",
		})
		return
	}

	subscriptions, err := h.subscriptionService.ListSubscriberTopics(c.Request.Context(), subscriberID)
	if err != nil {
		if err.Error() == "subscriber not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscriber not found",
			})
			return
		}

		h.logger.Error("Failed to list subscriber topics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list subscriber topics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
	})
}

func (h *SubscriptionHandler) ListTopicSubscribers(c *gin.Context) {
	topicIDStr := c.Param("id")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID format",
		})
		return
	}

	subscriptions, err := h.subscriptionService.ListTopicSubscribers(c.Request.Context(), topicID)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}

		h.logger.Error("Failed to list topic subscribers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list topic subscribers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
	})
}
