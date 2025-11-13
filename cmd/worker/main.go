package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"newsletter-assignment/internal/config"
	"newsletter-assignment/internal/constants"
	"newsletter-assignment/internal/db"
	"newsletter-assignment/internal/email"
	"newsletter-assignment/internal/log"
	"newsletter-assignment/internal/queue"
	"newsletter-assignment/internal/repo"
	"newsletter-assignment/internal/version"
	"newsletter-assignment/internal/worker"

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

	logger.Info("Newsletter Worker starting",
		zap.String("version", version.Version),
		zap.String("build", version.Build),
		zap.String("env", cfg.Env),
	)

	// Connect to database
	database, err := db.New(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Initialize repositories
	contentRepo := repo.NewContentRepository(database)
	subscriptionRepo := repo.NewSubscriptionRepository(database)
	subscriberRepo := repo.NewSubscriberRepository(database)
	jobRepo := repo.NewJobRepository(database)
	deliveryRepo := repo.NewDeliveryRepository(database)

	// Initialize unified email sender (supports both SMTP and HTTP API)
	smtpConfig := &email.SMTPConfig{
		Host:      cfg.SMTP.Host,
		Port:      cfg.SMTP.Port,
		Username:  cfg.SMTP.Username,
		Password:  cfg.SMTP.Password,
		FromEmail: cfg.SMTP.FromEmail,
		FromName:  cfg.SMTP.FromName,
	}

	httpConfig := &email.HTTPConfig{
		APIKey:    cfg.Email.APIKey,
		FromEmail: cfg.Email.FromEmail,
		FromName:  cfg.Email.FromName,
		BaseURL:   cfg.Email.BaseURL,
	}

	unifiedConfig := &email.UnifiedConfig{
		SMTP:    smtpConfig,
		HTTP:    httpConfig,
		UseHTTP: cfg.Email.UseHTTP,
	}

	emailSender := email.NewUnifiedEmailSender(unifiedConfig, logger)

	// Initialize worker
	sendContentWorker := worker.NewSendContentWorker(
		contentRepo,
		subscriptionRepo,
		subscriberRepo,
		jobRepo,
		deliveryRepo,
		emailSender,
		logger,
	)

	// Initialize queue
	jobQueue := queue.NewAsynqQueue(
		cfg.Asynq.RedisAddr,
		cfg.Asynq.RedisPassword,
		cfg.Asynq.RedisDB,
		cfg.Asynq.TLSConfigNeeded,
		logger,
	)

	// Register task handlers
	jobQueue.RegisterHandler(constants.JobTypeSendNewsletter, sendContentWorker.HandleSendContent)

	// Start health check server for Render
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Worker is healthy"))
		})
		logger.Info("Starting health check server", zap.String("port", "8081"))
		if err := http.ListenAndServe(":8081", nil); err != nil {
			logger.Error("Health check server failed", zap.Error(err))
		}
	}()

	// Start server in goroutine
	go func() {
		if err := jobQueue.Start(); err != nil {
			logger.Fatal("Failed to start Asynq worker server", zap.Error(err))
		}
	}()

	logger.Info("Asynq worker server started successfully")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	jobQueue.Shutdown()
	logger.Info("Worker exited")
}
