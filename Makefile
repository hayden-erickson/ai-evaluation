# Makefile for Habit Notification CronJob

# Variables
DOCKER_REGISTRY ?= your-registry
IMAGE_NAME ?= notification-job
IMAGE_TAG ?= latest
FULL_IMAGE := $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build-local    - Build the Go binary locally"
	@echo "  build-docker   - Build the Docker image"
	@echo "  push-docker    - Push the Docker image to registry"
	@echo "  deploy-k8s     - Deploy to Kubernetes"
	@echo "  delete-k8s     - Delete from Kubernetes"
	@echo "  test-job       - Manually trigger a test job"
	@echo "  logs           - View logs from the latest job"
	@echo "  clean          - Clean up local build artifacts"

# Build the Go binary locally
.PHONY: build-local
build-local:
	@echo "Building notification job locally..."
	go build -o bin/notification-job ./cmd/notification-job
	@echo "Binary built: bin/notification-job"

# Build the Docker image
.PHONY: build-docker
build-docker:
	@echo "Building Docker image $(FULL_IMAGE)..."
	docker build -f cmd/notification-job/Dockerfile -t $(FULL_IMAGE) .
	@echo "Docker image built: $(FULL_IMAGE)"

# Push the Docker image
.PHONY: push-docker
push-docker:
	@echo "Pushing Docker image $(FULL_IMAGE)..."
	docker push $(FULL_IMAGE)
	@echo "Docker image pushed successfully"

# Deploy to Kubernetes
.PHONY: deploy-k8s
deploy-k8s:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/cronjob.yaml
	@echo "Deployment complete"

# Deploy using kustomize
.PHONY: deploy-kustomize
deploy-kustomize:
	@echo "Deploying with kustomize..."
	kubectl apply -k k8s/
	@echo "Deployment complete"

# Delete from Kubernetes
.PHONY: delete-k8s
delete-k8s:
	@echo "Deleting from Kubernetes..."
	kubectl delete -f k8s/cronjob.yaml --ignore-not-found
	kubectl delete -f k8s/secret.yaml --ignore-not-found
	kubectl delete -f k8s/configmap.yaml --ignore-not-found
	@echo "Resources deleted"

# Manually trigger a test job
.PHONY: test-job
test-job:
	@echo "Creating manual test job..."
	kubectl create job --from=cronjob/habit-notification-job manual-test-$(shell date +%s)
	@echo "Test job created. Use 'make logs' to view logs"

# View logs from the latest job
.PHONY: logs
logs:
	@echo "Fetching logs from latest job..."
	@JOB_NAME=$$(kubectl get jobs --selector=app=habit-notification-job --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}' 2>/dev/null) && \
	if [ -z "$$JOB_NAME" ]; then \
		echo "No jobs found"; \
	else \
		echo "Job: $$JOB_NAME"; \
		kubectl logs job/$$JOB_NAME --follow; \
	fi

# Check CronJob status
.PHONY: status
status:
	@echo "CronJob status:"
	@kubectl get cronjob habit-notification-job 2>/dev/null || echo "CronJob not found"
	@echo ""
	@echo "Recent jobs:"
	@kubectl get jobs --selector=app=habit-notification-job --sort-by=.metadata.creationTimestamp 2>/dev/null || echo "No jobs found"

# Clean up local build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	@echo "Clean complete"

# Full build and push pipeline
.PHONY: release
release: build-docker push-docker
	@echo "Release complete: $(FULL_IMAGE)"
