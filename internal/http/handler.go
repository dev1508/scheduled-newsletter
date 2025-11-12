package http

import (
	"net/http"

	"newsletter-assignment/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	topicHandler *handler.TopicHandler
}

func NewHandler(topicHandler *handler.TopicHandler) *Handler {
	return &Handler{
		topicHandler: topicHandler,
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
	}

	return router
}

func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
