# Makefile for Habit Notification CronJob

# Configuration
IMAGE_REGISTRY ?= docker.io
IMAGE_NAME ?= habit-notification-cron
IMAGE_TAG ?= v1.0.0
FULL_IMAGE_NAME = $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
NAMESPACE ?= default

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=notification-cron
BINARY_PATH=./cmd/notification-cron

.PHONY: help
help: ## Display this help message
	@echo "Habit Notification CronJob - Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download Go dependencies
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: build
build: ## Build the notification cron binary
	$(GOBUILD) -v -o $(BINARY_NAME) $(BINARY_PATH)

.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: test
test: ## Run tests
	$(GOTEST) -v ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(FULL_IMAGE_NAME) -f cmd/notification-cron/Dockerfile .
	@echo "Built image: $(FULL_IMAGE_NAME)"

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	docker push $(FULL_IMAGE_NAME)
	@echo "Pushed image: $(FULL_IMAGE_NAME)"

.PHONY: docker-build-push
docker-build-push: docker-build docker-push ## Build and push Docker image

.PHONY: k8s-secret
k8s-secret: ## Create Kubernetes secret (requires env vars)
	@if [ -z "$(DB_PASSWORD)" ]; then echo "Error: DB_PASSWORD not set"; exit 1; fi
	@if [ -z "$(TWILIO_ACCOUNT_SID)" ]; then echo "Error: TWILIO_ACCOUNT_SID not set"; exit 1; fi
	@if [ -z "$(TWILIO_AUTH_TOKEN)" ]; then echo "Error: TWILIO_AUTH_TOKEN not set"; exit 1; fi
	@if [ -z "$(TWILIO_FROM_NUMBER)" ]; then echo "Error: TWILIO_FROM_NUMBER not set"; exit 1; fi
	kubectl create secret generic notification-cron-secrets \
		--from-literal=DB_PASSWORD='$(DB_PASSWORD)' \
		--from-literal=TWILIO_ACCOUNT_SID='$(TWILIO_ACCOUNT_SID)' \
		--from-literal=TWILIO_AUTH_TOKEN='$(TWILIO_AUTH_TOKEN)' \
		--from-literal=TWILIO_FROM_NUMBER='$(TWILIO_FROM_NUMBER)' \
		-n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@echo "Secret created/updated successfully"

.PHONY: k8s-deploy
k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f k8s/configmap.yaml -n $(NAMESPACE)
	kubectl apply -f k8s/rbac.yaml -n $(NAMESPACE)
	kubectl apply -f k8s/cronjob.yaml -n $(NAMESPACE)
	@echo "CronJob deployed successfully"

.PHONY: k8s-delete
k8s-delete: ## Delete from Kubernetes
	kubectl delete -f k8s/cronjob.yaml -n $(NAMESPACE) --ignore-not-found
	kubectl delete -f k8s/rbac.yaml -n $(NAMESPACE) --ignore-not-found
	kubectl delete -f k8s/configmap.yaml -n $(NAMESPACE) --ignore-not-found
	@echo "CronJob deleted successfully"

.PHONY: k8s-test
k8s-test: ## Create a manual test job
	kubectl create job --from=cronjob/habit-notification-cron habit-notification-test-$$(date +%s) -n $(NAMESPACE)
	@echo "Test job created. View logs with: make k8s-logs"

.PHONY: k8s-logs
k8s-logs: ## Show logs from the latest job
	kubectl logs -l app=habit-notification-cron -n $(NAMESPACE) --tail=100

.PHONY: k8s-status
k8s-status: ## Show CronJob status
	kubectl get cronjob habit-notification-cron -n $(NAMESPACE)
	@echo ""
	kubectl get jobs -l app=habit-notification-cron -n $(NAMESPACE)

.PHONY: k8s-describe
k8s-describe: ## Describe CronJob
	kubectl describe cronjob habit-notification-cron -n $(NAMESPACE)

.PHONY: all
all: deps build ## Build everything

.PHONY: deploy-all
deploy-all: docker-build-push k8s-deploy ## Build, push image, and deploy to Kubernetes
	@echo "Deployment complete!"

# Example usage:
# make build
# make docker-build IMAGE_REGISTRY=your-registry IMAGE_NAME=your-name
# make k8s-secret DB_PASSWORD=xxx TWILIO_ACCOUNT_SID=xxx TWILIO_AUTH_TOKEN=xxx TWILIO_FROM_NUMBER=xxx
# make k8s-deploy
# make k8s-test
# make k8s-logs
