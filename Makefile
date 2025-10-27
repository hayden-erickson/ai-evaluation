# Makefile for Habit Log Notification CronJob

# Variables
DOCKER_REGISTRY ?= your-registry
IMAGE_NAME ?= habit-cronjob
IMAGE_TAG ?= latest
FULL_IMAGE = $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Go variables
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_CLEAN = $(GO_CMD) clean
GO_TEST = $(GO_CMD) test
GO_GET = $(GO_CMD) get
BINARY_NAME = cronjob
BINARY_PATH = ./cmd/cronjob

.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the cronjob binary
	@echo "Building cronjob binary..."
	$(GO_BUILD) -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "Binary built: $(BINARY_NAME)"

.PHONY: build-linux
build-linux: ## Build the cronjob binary for Linux
	@echo "Building cronjob binary for Linux..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_NAME)-linux $(BINARY_PATH)
	@echo "Binary built: $(BINARY_NAME)-linux"

.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	$(GO_CLEAN)
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux
	@echo "Clean complete"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GO_TEST) -v ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image: $(FULL_IMAGE)"
	docker build -f Dockerfile.cronjob -t $(FULL_IMAGE) .
	@echo "Docker image built successfully"

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	@echo "Pushing Docker image: $(FULL_IMAGE)"
	docker push $(FULL_IMAGE)
	@echo "Docker image pushed successfully"

.PHONY: docker-build-push
docker-build-push: docker-build docker-push ## Build and push Docker image

.PHONY: k8s-apply
k8s-apply: ## Apply Kubernetes manifests
	@echo "Applying Kubernetes manifests..."
	kubectl apply -f k8s/rbac.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/cronjob.yaml
	@echo "Kubernetes resources applied"

.PHONY: k8s-delete
k8s-delete: ## Delete Kubernetes resources
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/cronjob.yaml --ignore-not-found
	kubectl delete -f k8s/secret.yaml --ignore-not-found
	kubectl delete -f k8s/configmap.yaml --ignore-not-found
	kubectl delete -f k8s/rbac.yaml --ignore-not-found
	@echo "Kubernetes resources deleted"

.PHONY: k8s-status
k8s-status: ## Check status of CronJob
	@echo "CronJob status:"
	kubectl get cronjob habit-log-notification
	@echo "\nRecent jobs:"
	kubectl get jobs --selector=app=habit-tracker -o wide

.PHONY: k8s-logs
k8s-logs: ## View logs from the most recent job
	@echo "Fetching logs from most recent job..."
	@JOB_NAME=$$(kubectl get jobs --selector=app=habit-tracker --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}' 2>/dev/null); \
	if [ -z "$$JOB_NAME" ]; then \
		echo "No jobs found"; \
	else \
		echo "Logs from job: $$JOB_NAME"; \
		kubectl logs job/$$JOB_NAME; \
	fi

.PHONY: k8s-trigger
k8s-trigger: ## Manually trigger the CronJob
	@echo "Manually triggering CronJob..."
	kubectl create job --from=cronjob/habit-log-notification manual-test-$$(date +%s)
	@echo "Job created. Check status with: make k8s-status"

.PHONY: deps
deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	$(GO_GET) -v ./...
	@echo "Dependencies downloaded"

.PHONY: run-local
run-local: ## Run the cronjob locally (requires env vars)
	@echo "Running cronjob locally..."
	@if [ -z "$$DB_PASSWORD" ]; then echo "Error: DB_PASSWORD not set"; exit 1; fi
	@if [ -z "$$TWILIO_ACCOUNT_SID" ]; then echo "Error: TWILIO_ACCOUNT_SID not set"; exit 1; fi
	@if [ -z "$$TWILIO_AUTH_TOKEN" ]; then echo "Error: TWILIO_AUTH_TOKEN not set"; exit 1; fi
	@if [ -z "$$TWILIO_PHONE_NUMBER" ]; then echo "Error: TWILIO_PHONE_NUMBER not set"; exit 1; fi
	$(GO_CMD) run $(BINARY_PATH)

.PHONY: verify-env
verify-env: ## Verify required environment variables are set
	@echo "Verifying environment variables..."
	@if [ -z "$(DOCKER_REGISTRY)" ] || [ "$(DOCKER_REGISTRY)" = "your-registry" ]; then \
		echo "Warning: DOCKER_REGISTRY is not set or using default value"; \
	fi
	@echo "Image will be: $(FULL_IMAGE)"
