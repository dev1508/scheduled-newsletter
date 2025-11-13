package config

import (
	"os"
	"strconv"

	"newsletter-assignment/internal/constants"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	LogLevel    string
	DatabaseURL string

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}

	SMTP struct {
		Host      string
		Port      string
		Username  string
		Password  string
		FromEmail string
		FromName  string
	}

	Email struct {
		APIKey    string
		BaseURL   string
		UseHTTP   bool
		FromEmail string
		FromName  string
	}

	Scheduler struct {
		Interval  string
		BatchSize int
	}

	Asynq struct {
		RedisAddr       string
		RedisPassword   string
		RedisDB         int
		TLSConfigNeeded bool
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return error if it doesn't exist
	}

	cfg := &Config{
		Port:        getEnv(constants.EnvKeyPort, constants.DefaultPort),
		Env:         getEnv(constants.EnvKeyEnvironment, constants.DefaultEnv),
		LogLevel:    getEnv(constants.EnvKeyLogLevel, constants.DefaultLogLevel),
		DatabaseURL: getEnv(constants.EnvKeyDatabaseURL, ""),
	}

	cfg.DB.Host = getEnv(constants.EnvKeyDBHost, constants.DefaultDBHost)
	cfg.DB.Port = getEnv(constants.EnvKeyDBPort, constants.DefaultDBPort)
	cfg.DB.User = getEnv(constants.EnvKeyDBUser, constants.DefaultDBUser)
	cfg.DB.Password = getEnv(constants.EnvKeyDBPassword, constants.DefaultDBPassword)
	cfg.DB.Name = getEnv(constants.EnvKeyDBName, constants.DefaultDBName)

	cfg.Redis.Host = getEnv(constants.EnvKeyRedisHost, constants.DefaultRedisHost)
	cfg.Redis.Port = getEnv(constants.EnvKeyRedisPort, constants.DefaultRedisPort)
	cfg.Redis.Password = getEnv(constants.EnvKeyRedisPassword, "")

	cfg.SMTP.Host = getEnv(constants.EnvKeySMTPHost, constants.DefaultSMTPHost)
	cfg.SMTP.Port = getEnv(constants.EnvKeySMTPPort, constants.DefaultSMTPPort)
	cfg.SMTP.Username = getEnv(constants.EnvKeySMTPUsername, "")
	cfg.SMTP.Password = getEnv(constants.EnvKeySMTPPassword, "")
	cfg.SMTP.FromEmail = getEnv(constants.EnvKeySMTPFromEmail, constants.DefaultSMTPFromEmail)
	cfg.SMTP.FromName = getEnv(constants.EnvKeySMTPFromName, constants.DefaultSMTPFromName)

	cfg.Email.APIKey = getEnv(constants.EnvKeyEmailAPIKey, "")
	cfg.Email.BaseURL = getEnv(constants.EnvKeyEmailAPIBaseURL, constants.DefaultEmailAPIBaseURL)
	cfg.Email.UseHTTP = getEnvBool(constants.EnvKeyEmailUseHTTP, constants.DefaultEmailUseHTTP)
	cfg.Email.FromEmail = getEnv(constants.EnvKeySMTPFromEmail, constants.DefaultSMTPFromEmail) // Reuse SMTP from email
	cfg.Email.FromName = getEnv(constants.EnvKeySMTPFromName, constants.DefaultSMTPFromName)    // Reuse SMTP from name

	cfg.Scheduler.Interval = getEnv(constants.EnvKeySchedulerInterval, constants.DefaultSchedulerInterval)
	cfg.Scheduler.BatchSize = getEnvInt(constants.EnvKeySchedulerBatchSize, constants.DefaultSchedulerBatchSize)

	cfg.Asynq.RedisAddr = getEnv(constants.EnvKeyAsynqRedisAddr, constants.DefaultRedisHost+":"+constants.DefaultRedisPort)
	cfg.Asynq.RedisPassword = getEnv(constants.EnvKeyAsynqRedisPassword, "")
	cfg.Asynq.RedisDB = getEnvInt(constants.EnvKeyAsynqRedisDB, 0)
	cfg.Asynq.TLSConfigNeeded = getEnvBool(constants.EnvKeyAsynqTLSConfigNeeded, false)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
