.PHONY: help build run test docker-build docker-up docker-down migrate clean

# Variables
APP_NAME=habits-api
DOCKER_IMAGE=$(APP_NAME):latest
GO_FILES=$(shell find . -name '*.go' -type f)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/api/main.go

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-down: ## Stop services
	@echo "Stopping services..."
	docker-compose down

docker-logs: ## View docker-compose logs
	docker-compose logs -f

migrate: ## Run database migrations
	@echo "Running migrations..."
	@bash scripts/run-migrations.sh

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...

# Development workflow
dev-setup: deps docker-build docker-up ## Set up local development environment
	@echo "Waiting for MySQL to be ready..."
	@sleep 10
	$(MAKE) migrate
	@echo "Development environment ready!"

dev-start: docker-up ## Start development environment
	@echo "Development environment started"
	@echo "API will be available at http://localhost:8080"

dev-stop: docker-down ## Stop development environment

# AWS deployment
aws-setup: ## Set up AWS EKS cluster
	@bash scripts/setup-eks.sh

aws-deploy: ## Deploy to AWS
	@bash scripts/deploy-aws.sh

# Kubernetes commands
k8s-status: ## Check Kubernetes deployment status
	kubectl get all -n habits-app

k8s-logs: ## View application logs
	kubectl logs -f deployment/habits-api -n habits-app

k8s-describe: ## Describe API deployment
	kubectl describe deployment habits-api -n habits-app

k8s-restart: ## Restart API deployment
	kubectl rollout restart deployment/habits-api -n habits-app
