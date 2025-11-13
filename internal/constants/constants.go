package constants

// Environment constants
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// Default configuration values
const (
	DefaultPort     = "8080"
	DefaultLogLevel = "info"
	DefaultEnv      = EnvDevelopment
)

// Database configuration defaults
const (
	DefaultDBHost     = "localhost"
	DefaultDBPort     = "5432"
	DefaultDBUser     = "newsletter"
	DefaultDBPassword = "password"
	DefaultDBName     = "newsletter_db"
)

// Redis configuration defaults
const (
	DefaultRedisHost = "localhost"
	DefaultRedisPort = "6379"
)

// SMTP configuration defaults
const (
	DefaultSMTPHost      = "smtp-relay.brevo.com"
	DefaultSMTPPort      = "587"
	DefaultSMTPFromEmail = "noreply@yourapp.com"
	DefaultSMTPFromName  = "Newsletter App"
)

// Email HTTP API configuration defaults
const (
	DefaultEmailAPIBaseURL = "https://api.brevo.com"
	DefaultEmailUseHTTP    = true
)

// Content status constants
const (
	ContentStatusScheduled = "scheduled"
	ContentStatusSent      = "sent"
	ContentStatusFailed    = "failed"
	ContentStatusCancelled = "cancelled"
)

// Delivery status constants
const (
	DeliveryStatusPending = "pending"
	DeliveryStatusSent    = "sent"
	DeliveryStatusFailed  = "failed"
	DeliveryStatusBounced = "bounced"
)

// Job status constants
const (
	JobStatusPending   = "pending"
	JobStatusEnqueued  = "enqueued"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
)

// Job types
const (
	JobTypeSendNewsletter = "send_newsletter"
	JobTypeCleanupOldJobs = "cleanup_old_jobs"
)

// Pagination defaults
const (
	DefaultLimit  = 10
	MaxLimit      = 100
	DefaultOffset = 0
)

// Validation constants
const (
	MaxTopicNameLength        = 255
	MaxTopicDescriptionLength = 1000
	MaxSubjectLength          = 500
	MaxSubscriberNameLength   = 255
	MaxEmailLength            = 255
)

// Database connection pool settings
const (
	DefaultMaxConnections = 10
	DefaultMinConnections = 2
)

// Job processing settings
const (
	DefaultMaxAttempts = 3
	DefaultJobLimit    = 100
)

// Scheduler settings
const (
	DefaultSchedulerInterval  = "30s"
	DefaultSchedulerBatchSize = 100
)

// Environment variable keys
const (
	EnvKeyPort        = "PORT"
	EnvKeyEnvironment = "ENV"
	EnvKeyLogLevel    = "LOG_LEVEL"
	EnvKeyDatabaseURL = "DATABASE_URL"
)

// Database environment variable keys
const (
	EnvKeyDBHost     = "DB_HOST"
	EnvKeyDBPort     = "DB_PORT"
	EnvKeyDBUser     = "DB_USER"
	EnvKeyDBPassword = "DB_PASSWORD"
	EnvKeyDBName     = "DB_NAME"
)

// Redis environment variable keys
const (
	EnvKeyRedisHost     = "REDIS_HOST"
	EnvKeyRedisPort     = "REDIS_PORT"
	EnvKeyRedisPassword = "REDIS_PASSWORD"
)

// SMTP environment variable keys
const (
	EnvKeySMTPHost      = "SMTP_HOST"
	EnvKeySMTPPort      = "SMTP_PORT"
	EnvKeySMTPUsername  = "SMTP_USERNAME"
	EnvKeySMTPPassword  = "SMTP_PASSWORD"
	EnvKeySMTPFromEmail = "SMTP_FROM_EMAIL"
	EnvKeySMTPFromName  = "SMTP_FROM_NAME"
)

// Email HTTP API environment variable keys
const (
	EnvKeyEmailAPIKey     = "EMAIL_API_KEY"
	EnvKeyEmailAPIBaseURL = "EMAIL_API_BASE_URL"
	EnvKeyEmailUseHTTP    = "EMAIL_USE_HTTP"
)

// Scheduler environment variable keys
const (
	EnvKeySchedulerInterval  = "SCHEDULER_INTERVAL"
	EnvKeySchedulerBatchSize = "SCHEDULER_BATCH_SIZE"
)

// Asynq environment variable keys
const (
	EnvKeyAsynqRedisAddr       = "ASYNQ_REDIS_ADDR"
	EnvKeyAsynqRedisPassword   = "ASYNQ_REDIS_PASSWORD"
	EnvKeyAsynqRedisDB         = "ASYNQ_REDIS_DB"
	EnvKeyAsynqTLSConfigNeeded = "ASYNQ_TLS_CONFIG_NEEDED"
)
