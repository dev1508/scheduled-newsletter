# Newsletter Scheduler System

A newsletter scheduler system built with Go that allows admins to create topics, schedule content, and automatically send emails to subscribers.

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
│   ├── repo/                # Repository layer
│   ├── service/             # Business logic layer
│   ├── handler/             # HTTP handlers
│   ├── http/                # HTTP router setup
│   ├── log/                 # Logging utilities
│   ├── queue/               # Queue management (future)
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
   # Edit .env with your database and Redis credentials
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

### Example Usage

**Create a topic**:
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech News",
    "description": "Latest technology news and updates"
  }'
```

**List topics**:
```bash
curl http://localhost:8080/api/v1/topics?limit=10&offset=0
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

See `.env.example` for all available configuration options.

## Database Schema

The system uses the following main tables:

- **topics** - Newsletter topics
- **subscribers** - Email subscribers  
- **subscriptions** - Subscriber-topic relationships
- **content** - Scheduled newsletter content
- **deliveries** - Individual email delivery tracking
- **job_scheduler** - Durable job scheduling

## Deployment

The system is designed to be free-tier friendly:

- **Database**: Neon PostgreSQL
- **Redis**: Upstash Redis
- **Hosting**: Render or similar

## Contributing

1. Follow the existing code structure and patterns
2. Write tests for new functionality
3. Update documentation as needed
4. Make small, focused commits

## License

MIT License
