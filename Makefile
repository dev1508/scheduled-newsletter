.PHONY: run-api run-worker build-api build-worker clean test deps

# Development commands
run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

# Build commands
build-api:
	go build -o bin/api cmd/api/main.go

build-worker:
	go build -o bin/worker cmd/worker/main.go

build: build-api build-worker

# Utility commands
deps:
	go mod tidy
	go mod download

test:
	go test -v ./...

clean:
	rm -rf bin/

# Docker commands (for future use)
docker-build:
	docker build -t newsletter-assignment .

# Database commands (for future use)
db-up:
	docker-compose up -d postgres redis

db-down:
	docker-compose down

# Help
help:
	@echo "Available commands:"
	@echo "  run-api      - Run the API server"
	@echo "  run-worker   - Run the worker"
	@echo "  build-api    - Build the API binary"
	@echo "  build-worker - Build the worker binary"
	@echo "  build        - Build both binaries"
	@echo "  deps         - Download dependencies"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
