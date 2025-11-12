# Newsletter Scheduler System

A production-ready newsletter scheduler system built with Go that allows admins to create topics, manage subscribers, schedule content, and automatically send emails with concurrent processing and delivery tracking.

## ✨ Features

- 📧 **Real SMTP Email Delivery** - Send emails via Brevo with TLS encryption
- ⚡ **High-Performance Concurrent Processing** - 20x faster with 20 concurrent goroutines
- 📊 **Complete Delivery Tracking** - Track sent/failed status for every email
- 🔄 **Automated Job Scheduling** - Background processing with Redis/Asynq
- 🎯 **Topic-Based Subscriptions** - Organize content by topics
- 👥 **Subscriber Management** - Full CRUD operations for subscribers
- 📝 **Content Scheduling** - Schedule newsletters for future delivery
- 🛡️ **Error Handling & Logging** - Comprehensive error tracking and recovery
- 🏗️ **Clean Architecture** - Hexagonal architecture with clear separation of concerns

## Tech Stack

- **Backend**: Go (Gin + pgx + Asynq + net/smtp)
- **Database**: PostgreSQL for persistence
- **Queue**: Redis (Asynq) for background job execution
- **Email**: Brevo SMTP for email delivery
- **Architecture**: Clean hexagonal layers (handler → service → repository)

## Project Structure

```
newsletter-assignment/
├── cmd/
│   ├── api/main.go          # API server entrypoint
│   └── worker/main.go       # Background worker entrypoint
├── internal/
│   ├── config/              # Configuration management
│   ├── db/                  # Database connection
│   ├── models/              # Domain entities
│   ├── repo/                # Repository layer (topics, subscribers, content, deliveries)
│   ├── service/             # Business logic layer
│   ├── handler/             # HTTP handlers
│   ├── http/                # HTTP router setup
│   ├── log/                 # Logging utilities
│   ├── email/               # SMTP email service
│   ├── worker/              # Background job workers
│   ├── scheduler/           # Job scheduling service
│   ├── queue/               # Queue management (Asynq)
│   └── version/             # Version constants
├── migrations/              # Database migrations
├── .env.example            # Environment variables template
├── Makefile               # Build and run commands
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Redis 6+

### Setup

1. **Clone and setup environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your database, Redis, and SMTP credentials
   ```

2. **Run database migrations**:
   ```bash
   # Connect to your PostgreSQL database and run:
   psql -d your_database -f migrations/001_init.sql
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the API server**:
   ```bash
   make run-api
   # or
   go run cmd/api/main.go
   ```

5. **Run the worker** (in another terminal):
   ```bash
   make run-worker
   # or
   go run cmd/worker/main.go
   ```

### API Endpoints

#### Health Check
- `GET /healthz` - Health check endpoint

#### Topics
- `POST /api/v1/topics` - Create a new topic
- `GET /api/v1/topics` - List all topics (with pagination)
- `GET /api/v1/topics/:id` - Get topic by ID
- `PUT /api/v1/topics/:id` - Update topic
- `DELETE /api/v1/topics/:id` - Delete topic
- `GET /api/v1/topics/:id/subscribers` - Get topic subscribers
- `GET /api/v1/topics/:id/content` - Get topic content

#### Subscribers
- `POST /api/v1/subscribers` - Create a new subscriber
- `GET /api/v1/subscribers` - List all subscribers (with pagination)
- `GET /api/v1/subscribers/:id` - Get subscriber by ID
- `PUT /api/v1/subscribers/:id` - Update subscriber
- `DELETE /api/v1/subscribers/:id` - Delete subscriber
- `GET /api/v1/subscribers/:id/topics` - Get subscriber topics

#### Subscriptions
- `POST /api/v1/subscriptions` - Subscribe user to topic
- `GET /api/v1/subscriptions/:id` - Get subscription details
- `DELETE /api/v1/subscriptions/:subscriber_id/:topic_id` - Unsubscribe

#### Content & Newsletters
- `POST /api/v1/content` - Create and schedule newsletter content
- `GET /api/v1/content` - List all content (with pagination)
- `GET /api/v1/content/:id` - Get content by ID
- `PUT /api/v1/content/:id` - Update content
- `DELETE /api/v1/content/:id` - Delete content

### Example Usage

#### Complete Newsletter Workflow

**1. Create a topic**:
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech News",
    "description": "Latest technology news and updates"
  }'
```

**2. Create a subscriber**:
```bash
curl -X POST http://localhost:8080/api/v1/subscribers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

**3. Subscribe user to topic**:
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "subscriber_id": "SUBSCRIBER_UUID",
    "topic_id": "TOPIC_UUID"
  }'
```

**4. Create and schedule newsletter content**:
```bash
curl -X POST http://localhost:8080/api/v1/content \
  -H "Content-Type: application/json" \
  -d '{
    "topic_id": "TOPIC_UUID",
    "subject": "Weekly Tech Update",
    "body": "<h1>Hello!</h1><p>This weeks tech news...</p>",
    "send_at": "2025-11-13T15:30:00+05:30"
  }'
```

**5. Monitor delivery status**:
```bash
curl http://localhost:8080/api/v1/content/CONTENT_UUID
# Check "status" field: "scheduled" → "sent"
```

## Development

### Available Make Commands

- `make run-api` - Run the API server
- `make run-worker` - Run the background worker
- `make build` - Build both API and worker binaries
- `make test` - Run tests
- `make deps` - Download dependencies
- `make clean` - Clean build artifacts

### Environment Variables

Key configuration options in `.env`:

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/newsletter_db

# Redis (for job queue)
REDIS_ADDR=localhost:6379

# SMTP (Brevo configuration)
SMTP_HOST=smtp-relay.brevo.com
SMTP_PORT=587
SMTP_USERNAME=your_brevo_login@smtp-brevo.com
SMTP_PASSWORD=your_brevo_smtp_key
SMTP_FROM_EMAIL=your_email@example.com
SMTP_FROM_NAME=Newsletter App

# Scheduler
SCHEDULER_INTERVAL=30s
SCHEDULER_BATCH_SIZE=100
```

See `.env.example` for all available configuration options.

## Database Schema

The system uses the following main tables:

- **topics** - Newsletter topics
- **subscribers** - Email subscribers  
- **subscriptions** - Subscriber-topic relationships
- **content** - Scheduled newsletter content
- **deliveries** - Individual email delivery tracking
- **job_scheduler** - Durable job scheduling

## How It Works

### Architecture Overview

1. **API Server** (`cmd/api/main.go`):
   - Handles HTTP requests for CRUD operations
   - Runs job scheduler every 30 seconds
   - Enqueues newsletter jobs to Redis

2. **Background Worker** (`cmd/worker/main.go`):
   - Processes newsletter sending jobs from Redis queue
   - Sends emails concurrently (20 goroutines)
   - Tracks delivery status in database

3. **Email Flow**:
   ```
   Content Created → Job Scheduled → Worker Picks Up → 
   Concurrent Email Sending → Delivery Tracking → Status Update
   ```

### Performance Features

- **Concurrent Processing**: 20 parallel email sends per batch
- **Delivery Tracking**: Individual status for each email (pending/sent/failed)
- **Error Handling**: Failed emails are logged with error messages
- **Job Persistence**: Durable job scheduling with Redis/Asynq

## Deployment

The system is designed to be free-tier friendly:

- **Database**: Neon PostgreSQL (free tier)
- **Redis**: Upstash Redis (free tier)
- **Email**: Brevo SMTP (300 emails/day free)
- **Hosting**: Render, Railway, or similar (free tier)

## Contributing

1. Follow the existing code structure and patterns
2. Write tests for new functionality
3. Update documentation as needed
4. Make small, focused commits

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
