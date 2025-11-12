package models

import (
	"time"

	"github.com/google/uuid"
)

// Topic represents a newsletter topic
type Topic struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}


// Subscriber represents an email subscriber
type Subscriber struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      *string   `json:"name" db:"name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription represents a subscriber's subscription to a topic
type Subscription struct {
	ID           uuid.UUID `json:"id" db:"id"`
	SubscriberID uuid.UUID `json:"subscriber_id" db:"subscriber_id"`
	TopicID      uuid.UUID `json:"topic_id" db:"topic_id"`
	SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
}

// Content represents scheduled newsletter content
type Content struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TopicID   uuid.UUID `json:"topic_id" db:"topic_id"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	SendAt    time.Time `json:"send_at" db:"send_at"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Delivery represents an individual email delivery
type Delivery struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ContentID    uuid.UUID  `json:"content_id" db:"content_id"`
	SubscriberID uuid.UUID  `json:"subscriber_id" db:"subscriber_id"`
	Email        string     `json:"email" db:"email"`
	Status       string     `json:"status" db:"status"`
	SentAt       *time.Time `json:"sent_at" db:"sent_at"`
	ErrorMessage *string    `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// JobScheduler represents a scheduled job
type JobScheduler struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ContentID   uuid.UUID  `json:"content_id" db:"content_id"`
	JobType     string     `json:"job_type" db:"job_type"`
	ScheduledAt time.Time  `json:"scheduled_at" db:"scheduled_at"`
	Status      string     `json:"status" db:"status"`
	Attempts    int        `json:"attempts" db:"attempts"`
	MaxAttempts int        `json:"max_attempts" db:"max_attempts"`
	ErrorMessage *string   `json:"error_message" db:"error_message"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
