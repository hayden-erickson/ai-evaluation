# Makefile for Habit Notifier CronJob

# Variables
APP_NAME := habit-notifier
REGISTRY := your-registry
IMAGE_NAME := $(REGISTRY)/$(APP_NAME)
VERSION := latest
NAMESPACE := default

# Go variables
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_CLEAN := $(GO_CMD) clean
GO_TEST := $(GO_CMD) test
GO_GET := $(GO_CMD) get
GO_MOD := $(GO_CMD) mod

.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download Go dependencies
	$(GO_MOD) download
	$(GO_MOD) tidy

.PHONY: build
build: ## Build the Go binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -ldflags="-w -s" -o $(APP_NAME) .

.PHONY: clean
clean: ## Clean build artifacts
	$(GO_CLEAN)
	rm -f $(APP_NAME)

.PHONY: test
test: ## Run tests
	$(GO_TEST) -v ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(VERSION) .

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	docker push $(IMAGE_NAME):$(VERSION)

.PHONY: docker-all
docker-all: docker-build docker-push ## Build and push Docker image

.PHONY: k8s-create-secret
k8s-create-secret: ## Create Kubernetes secret (requires env vars)
	kubectl create secret generic habit-notifier-secret \
		--from-literal=DB_USER=${DB_USER} \
		--from-literal=DB_PASSWORD=${DB_PASSWORD} \
		--from-literal=TWILIO_ACCOUNT_SID=${TWILIO_ACCOUNT_SID} \
		--from-literal=TWILIO_AUTH_TOKEN=${TWILIO_AUTH_TOKEN} \
		--from-literal=TWILIO_FROM_NUMBER=${TWILIO_FROM_NUMBER} \
		--namespace=$(NAMESPACE) \
		--dry-run=client -o yaml | kubectl apply -f -

.PHONY: k8s-deploy
k8s-deploy: ## Deploy all Kubernetes resources
	kubectl apply -f k8s/ --namespace=$(NAMESPACE)

.PHONY: k8s-delete
k8s-delete: ## Delete all Kubernetes resources
	kubectl delete -f k8s/ --namespace=$(NAMESPACE) --ignore-not-found=true

.PHONY: k8s-logs
k8s-logs: ## View logs from the latest job
	kubectl logs -l app=habit-notifier --namespace=$(NAMESPACE) --tail=100 -f

.PHONY: k8s-test
k8s-test: ## Create a manual test job
	kubectl create job --from=cronjob/habit-notifier habit-notifier-manual-test --namespace=$(NAMESPACE)
	@echo "Watch the job with: make k8s-logs"

.PHONY: k8s-status
k8s-status: ## Check CronJob and Job status
	@echo "=== CronJob Status ==="
	kubectl get cronjob habit-notifier --namespace=$(NAMESPACE)
	@echo ""
	@echo "=== Recent Jobs ==="
	kubectl get jobs -l app=habit-notifier --namespace=$(NAMESPACE) --sort-by=.metadata.creationTimestamp

.PHONY: k8s-describe
k8s-describe: ## Describe CronJob
	kubectl describe cronjob habit-notifier --namespace=$(NAMESPACE)

.PHONY: deploy-all
deploy-all: docker-all k8s-deploy ## Build, push, and deploy everything

.PHONY: run-local
run-local: ## Run the application locally (requires env vars)
	$(GO_CMD) run .

.PHONY: fmt
fmt: ## Format Go code
	$(GO_CMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO_CMD) vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint installed)
	golangci-lint run

.PHONY: security-scan
security-scan: ## Run security scan on Docker image
	docker scan $(IMAGE_NAME):$(VERSION)

# Default target
.DEFAULT_GOAL := help

