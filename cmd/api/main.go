package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"newsletter-assignment/internal/config"
	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/handler"
	httphandler "newsletter-assignment/internal/http"
	"newsletter-assignment/internal/log"
	"newsletter-assignment/internal/repo"
	"newsletter-assignment/internal/queue"
	"newsletter-assignment/internal/scheduler"
	"newsletter-assignment/internal/service"
	"newsletter-assignment/internal/version"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := log.NewLogger(cfg.Env, cfg.LogLevel)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting newsletter API server",
		zap.String("version", version.Version),
		zap.String("build", version.Build),
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
	)

	// Initialize database
	database, err := db.New(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Initialize repositories
	topicRepo := repo.NewTopicRepository(database)
	subscriberRepo := repo.NewSubscriberRepository(database)
	subscriptionRepo := repo.NewSubscriptionRepository(database)
	contentRepo := repo.NewContentRepository(database)
	jobRepo := repo.NewJobRepository(database)

	// Initialize services
	topicService := service.NewTopicService(topicRepo, logger)
	subscriberService := service.NewSubscriberService(subscriberRepo, logger)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, subscriberRepo, topicRepo, logger)
	contentService := service.NewContentService(contentRepo, topicRepo, jobRepo, database, logger)

	// Initialize handlers
	topicHandler := handler.NewTopicHandler(topicService, logger)
	subscriberHandler := handler.NewSubscriberHandler(subscriberService, logger)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, logger)
	contentHandler := handler.NewContentHandler(contentService, logger)

	// Initialize queue
	jobQueue := queue.NewAsynqQueue(
		cfg.Asynq.RedisAddr,
		cfg.Asynq.RedisPassword,
		cfg.Asynq.RedisDB,
		logger,
	)
	defer jobQueue.Close()

	// Parse scheduler interval
	schedulerInterval, err := time.ParseDuration(cfg.Scheduler.Interval)
	if err != nil {
		logger.Fatal("Invalid scheduler interval", zap.String("interval", cfg.Scheduler.Interval), zap.Error(err))
	}

	// Initialize scheduler
	jobScheduler := scheduler.NewScheduler(jobRepo, jobQueue, logger, schedulerInterval, cfg.Scheduler.BatchSize)

	// Initialize HTTP handler with dependencies
	httpHandler := httphandler.NewHandler(topicHandler, subscriberHandler, subscriptionHandler, contentHandler)
	router := httpHandler.SetupRoutes()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start scheduler in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go func() {
		logger.Info("Starting job scheduler")
		jobScheduler.Start(ctx)
	}()

	go func() {
		logger.Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
