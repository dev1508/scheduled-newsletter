package http

import (
	"net/http"

	"newsletter-assignment/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	topicHandler        *handler.TopicHandler
	subscriberHandler   *handler.SubscriberHandler
	subscriptionHandler *handler.SubscriptionHandler
}

func NewHandler(
	topicHandler *handler.TopicHandler,
	subscriberHandler *handler.SubscriberHandler,
	subscriptionHandler *handler.SubscriptionHandler,
) *Handler {
	return &Handler{
		topicHandler:        topicHandler,
		subscriberHandler:   subscriberHandler,
		subscriptionHandler: subscriptionHandler,
	}
}

func (h *Handler) SetupRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/healthz", h.healthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Topic routes
		topics := v1.Group("/topics")
		{
			topics.POST("", h.topicHandler.CreateTopic)
			topics.GET("", h.topicHandler.ListTopics)
			topics.GET("/:id", h.topicHandler.GetTopic)
			topics.PUT("/:id", h.topicHandler.UpdateTopic)
			topics.DELETE("/:id", h.topicHandler.DeleteTopic)
		}

		// Subscriber routes
		subscribers := v1.Group("/subscribers")
		{
			subscribers.POST("", h.subscriberHandler.CreateSubscriber)
			subscribers.GET("", h.subscriberHandler.ListSubscribers)
			subscribers.GET("/search", h.subscriberHandler.GetSubscriberByEmail) // ?email=user@example.com
			subscribers.GET("/:id", h.subscriberHandler.GetSubscriber)
			subscribers.PUT("/:id", h.subscriberHandler.UpdateSubscriber)
			subscribers.DELETE("/:id", h.subscriberHandler.DeleteSubscriber)
		}

		// Subscription routes
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", h.subscriptionHandler.Subscribe)
			subscriptions.GET("/:id", h.subscriptionHandler.GetSubscription)
			subscriptions.DELETE("/:subscriber_id/:topic_id", h.subscriptionHandler.Unsubscribe)
		}

		// Subscriber-specific subscription routes
		v1.GET("/subscribers/:subscriber_id/topics", h.subscriptionHandler.ListSubscriberTopics)
		
		// Topic-specific subscription routes
		v1.GET("/topics/:topic_id/subscribers", h.subscriptionHandler.ListTopicSubscribers)
	}

	return router
}

func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
