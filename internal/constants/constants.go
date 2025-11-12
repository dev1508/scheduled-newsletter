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
	DefaultLimit    = 10
	MaxLimit        = 100
	DefaultOffset   = 0
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
