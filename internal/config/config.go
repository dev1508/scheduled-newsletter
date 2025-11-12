package config

import (
	"os"
	"strconv"

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
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return error if it doesn't exist
	}

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "newsletter")
	cfg.DB.Password = getEnv("DB_PASSWORD", "password")
	cfg.DB.Name = getEnv("DB_NAME", "newsletter_db")

	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")

	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp-relay.brevo.com")
	cfg.SMTP.Port = getEnv("SMTP_PORT", "587")
	cfg.SMTP.Username = getEnv("SMTP_USERNAME", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.FromEmail = getEnv("SMTP_FROM_EMAIL", "noreply@yourapp.com")
	cfg.SMTP.FromName = getEnv("SMTP_FROM_NAME", "Newsletter App")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
