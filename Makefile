# Makefile for VYOM ERP LMS Backend

.PHONY: help run test build fmt lint clean docker-build docker-up docker-down migrate prisma-generate

# Variables
BINARY_NAME=vyom-lms
GO_FILES=$(shell find . -type f -name '*.go')
DOCKER_IMAGE=vyom-lms:latest

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)VYOM ERP LMS - Makefile$(NC)"
	@echo "========================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

run: ## Run the application
	@echo "$(BLUE)Starting LMS backend...$(NC)"
	go run ./cmd/api/main.go

dev: ## Run with hot reload (requires air)
	@echo "$(BLUE)Starting LMS backend with hot reload...$(NC)"
	air

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	go test -v -tags=unit ./internal/...

test-integration: ## Run integration tests only
	@echo "$(BLUE)Running integration tests...$(NC)"
	go test -v -tags=integration ./internal/...

bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	gofmt -s -w .
	goimports -w .

lint: ## Lint code
	@echo "$(BLUE)Linting code...$(NC)"
	golangci-lint run ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	go vet ./...

tidy: ## Clean up dependencies
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	go mod tidy
	go mod verify

build: ## Build binary
	@echo "$(BLUE)Building binary...$(NC)"
	go build -o bin/$(BINARY_NAME) ./cmd/api/main.go
	@echo "$(GREEN)Binary built: bin/$(BINARY_NAME)$(NC)"

release: ## Build release binary (optimized)
	@echo "$(BLUE)Building release binary...$(NC)"
	go build -ldflags="-w -s" -o bin/$(BINARY_NAME) ./cmd/api/main.go
	@echo "$(GREEN)Release binary built: bin/$(BINARY_NAME)$(NC)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning up...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Prisma commands
prisma-generate: ## Generate Prisma client
	@echo "$(BLUE)Generating Prisma client...$(NC)"
	go run github.com/prisma/prisma-client-go generate

prisma-migrate: ## Run migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	go run github.com/prisma/prisma-client-go migrate deploy

prisma-migrate-dev: ## Run migrations in dev mode
	@echo "$(BLUE)Running migrations (dev mode)...$(NC)"
	go run github.com/prisma/prisma-client-go migrate dev

prisma-seed: ## Seed database
	@echo "$(BLUE)Seeding database...$(NC)"
	go run ./prisma/seed.ts

# Database commands
migrate: ## Run database migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	@prisma migrate deploy

migrate-dev: ## Run migrations in dev
	@echo "$(BLUE)Running migrations (dev)...$(NC)"
	@prisma migrate dev

migrate-create: ## Create new migration (NAME=migration_name)
	@echo "$(BLUE)Creating migration: $(NAME)$(NC)"
	@prisma migrate dev --name $(NAME)

db-reset: ## Reset database (careful!)
	@echo "$(RED)WARNING: This will reset your database!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		prisma migrate reset; \
	fi

db-push: ## Push schema to database
	@echo "$(BLUE)Pushing schema to database...$(NC)"
	@prisma db push

db-seed-sample: ## Seed sample data
	@echo "$(BLUE)Seeding sample data...$(NC)"
	go run ./cmd/seed/main.go

# Docker commands
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-up: ## Start Docker containers
	@echo "$(BLUE)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Containers started$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)Containers stopped$(NC)"

docker-logs: ## View Docker logs
	@echo "$(BLUE)Docker logs...$(NC)"
	docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "$(BLUE)Cleaning Docker resources...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)Docker resources cleaned$(NC)"

# Development environment
env-setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	cp .env.example .env
	go mod download
	go mod tidy
	docker-compose up -d
	sleep 5
	make migrate
	@echo "$(GREEN)Development environment setup complete$(NC)"

env-clean: ## Clean development environment
	@echo "$(BLUE)Cleaning development environment...$(NC)"
	docker-compose down -v
	rm -f .env
	rm -rf bin/
	@echo "$(GREEN)Development environment cleaned$(NC)"

# Code generation
generate: ## Generate code (mocks, etc.)
	@echo "$(BLUE)Generating code...$(NC)"
	go generate ./...

# Proto commands
proto-generate: ## Generate gRPC code from proto files
	@echo "$(BLUE)Generating gRPC code...$(NC)"
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Documentation
docs: ## Generate API documentation
	@echo "$(BLUE)Generating API documentation...$(NC)"
	swag init -g cmd/api/main.go
	@echo "$(GREEN)API docs generated$(NC)"

# Development helpers
check: fmt vet lint ## Run all checks
	@echo "$(GREEN)All checks passed!$(NC)"

setup-hooks: ## Setup git hooks
	@echo "$(BLUE)Setting up git hooks...$(NC)"
	@mkdir -p .git/hooks
	@echo "#!/bin/bash\nmake check" > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Git hooks setup complete$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install gotest.tools/gotestsum@latest
	@echo "$(GREEN)Development tools installed$(NC)"

# Debugging
debug: ## Run with debugger (requires dlv)
	@echo "$(BLUE)Running with debugger...$(NC)"
	dlv debug ./cmd/api/main.go

# Release
release-publish: ## Publish release (requires goreleaser)
	@echo "$(BLUE)Publishing release...$(NC)"
	goreleaser release --rm-dist

.DEFAULT_GOAL := help
