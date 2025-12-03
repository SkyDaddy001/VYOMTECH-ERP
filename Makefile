.PHONY: help setup-dev dev-up dev-down build run test test-all lint format clean docker-build docker-run migrate

help:
	@echo "Multi-Tenant AI Call Center - Development Commands"
	@echo ""
	@echo "Setup & Development:"
	@echo "  make setup-dev     - Install dependencies"
	@echo "  make dev-up        - Start development environment (requires Docker)"
	@echo "  make dev-down      - Stop development environment"
	@echo ""
	@echo "Building & Running:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo ""
	@echo "Testing:"
	@echo "  make test          - Run unit tests"
	@echo "  make test-all      - Run all tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make format        - Format code with gofmt"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo ""
	@echo "Database:"
	@echo "  make migrate       - Run database migrations"
	@echo "  make clean         - Clean build artifacts"

# Setup and dependencies
setup-dev:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

# Development environment
dev-up:
	@echo "Starting development environment with Docker Compose..."
	docker-compose up -d
	@echo "Waiting for database to be ready..."
	sleep 5
	make migrate
	@echo "Development environment is ready!"

dev-down:
	@echo "Stopping development environment..."
	docker-compose down
	@echo "Development environment stopped"

# Building
build:
	@echo "Building application..."
	go build -o bin/main ./cmd/main.go
	@echo "Build complete: bin/main"

run: build
	@echo "Running application..."
	./bin/main

# Testing
test:
	@echo "Running tests..."
	go test -v ./...

test-all:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Code quality
lint:
	@echo "Running linter..."
	golangci-lint run ./...

format:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t vyomtech-backend:latest .
	@echo "Docker image built successfully"

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 \
		--env-file .env \
		vyomtech-backend:latest

# Database
migrate:
	@echo "Running database migrations..."
	@echo "Note: Ensure MySQL is running and accessible"
	# Run migration script here
	@echo "Migrations complete"

# Cleanup
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "Cleanup complete"

# Development with hot reload
dev-watch:
	@echo "Starting development with hot reload..."
	@echo "Make sure to install air: go install github.com/cosmtrek/air@latest"
	air

# Database setup (local)
db-setup:
	@echo "Setting up local database..."
	mysql -u root -p < migrations/001_initial_schema.sql
	@echo "Database setup complete"

# Docker database setup
db-docker:
	docker run --name mysql-callcenter \
		-e MYSQL_ROOT_PASSWORD=rootpass \
		-e MYSQL_DATABASE=callcenter \
		-e MYSQL_USER=callcenter_user \
		-e MYSQL_PASSWORD=secure_app_pass \
		-p 3306:3306 \
		-d mysql:8.0

db-stop:
	docker stop mysql-callcenter || true
	docker rm mysql-callcenter || true
