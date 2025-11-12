package main

import (
	"fmt"

	"newsletter-assignment/internal/config"
	"newsletter-assignment/internal/log"
	"newsletter-assignment/internal/version"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	logger, err := log.NewLogger(cfg.Env, cfg.LogLevel)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	logger.Info("Starting newsletter worker",
		zap.String("version", version.Version),
		zap.String("build", version.Build),
		zap.String("env", cfg.Env),
	)

	logger.Info("Worker started successfully - ready to process jobs")
}
