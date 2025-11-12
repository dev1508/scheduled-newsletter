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
	contentHandler      *handler.ContentHandler
}

func NewHandler(
	topicHandler *handler.TopicHandler,
	subscriberHandler *handler.SubscriberHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	contentHandler *handler.ContentHandler,
) *Handler {
	return &Handler{
		topicHandler:        topicHandler,
		subscriberHandler:   subscriberHandler,
		subscriptionHandler: subscriptionHandler,
		contentHandler:      contentHandler,
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
			
			// Topic-specific subscription routes
			topics.GET("/:id/subscribers", h.subscriptionHandler.ListTopicSubscribers)
			
			// Topic-specific content routes
			topics.GET("/:id/content", h.contentHandler.ListContentByTopic)
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
			
			// Subscriber-specific subscription routes
			subscribers.GET("/:id/topics", h.subscriptionHandler.ListSubscriberTopics)
		}

		// Subscription routes
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", h.subscriptionHandler.Subscribe)
			subscriptions.GET("/:id", h.subscriptionHandler.GetSubscription)
			subscriptions.DELETE("/:subscriber_id/:topic_id", h.subscriptionHandler.Unsubscribe)
		}

		// Content routes
		content := v1.Group("/content")
		{
			content.POST("", h.contentHandler.CreateContent)
			content.GET("", h.contentHandler.ListContent)
			content.GET("/:id", h.contentHandler.GetContent)
			content.PUT("/:id", h.contentHandler.UpdateContent)
			content.DELETE("/:id", h.contentHandler.DeleteContent)
			content.POST("/:id/schedule", h.contentHandler.ScheduleContent)
		}
	}

	return router
}

func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
