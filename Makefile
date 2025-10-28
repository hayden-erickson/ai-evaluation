# Makefile for Habit Notifier CronJob

# Variables
IMAGE_NAME ?= your-registry/habit-notifier
IMAGE_TAG ?= latest
NAMESPACE ?= habit-tracker

.PHONY: help build push deploy test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Built image: $(IMAGE_NAME):$(IMAGE_TAG)"

push: ## Push Docker image to registry
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	@echo "Pushed image: $(IMAGE_NAME):$(IMAGE_TAG)"

build-push: build push ## Build and push Docker image

deploy: ## Deploy to Kubernetes
	kubectl apply -k k8s/
	@echo "Deployed to namespace: $(NAMESPACE)"

deploy-raw: ## Deploy using raw kubectl apply
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/serviceaccount.yaml
	kubectl apply -f k8s/networkpolicy.yaml
	kubectl apply -f k8s/cronjob.yaml

test: ## Run a manual test job
	kubectl create job --from=cronjob/habit-notifier habit-notifier-test-$$(date +%s) -n $(NAMESPACE)
	@echo "Test job created. Monitor with: make logs"

logs: ## View logs from most recent job
	@POD=$$(kubectl get pods -n $(NAMESPACE) -l app=habit-notifier --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}' 2>/dev/null); \
	if [ -n "$$POD" ]; then \
		kubectl logs -n $(NAMESPACE) $$POD -f; \
	else \
		echo "No pods found"; \
	fi

status: ## Check CronJob and Job status
	@echo "=== CronJob Status ==="
	kubectl get cronjob -n $(NAMESPACE)
	@echo ""
	@echo "=== Recent Jobs ==="
	kubectl get jobs -n $(NAMESPACE) --sort-by=.metadata.creationTimestamp
	@echo ""
	@echo "=== Recent Pods ==="
	kubectl get pods -n $(NAMESPACE) --sort-by=.metadata.creationTimestamp

describe: ## Describe CronJob
	kubectl describe cronjob habit-notifier -n $(NAMESPACE)

suspend: ## Suspend CronJob
	kubectl patch cronjob habit-notifier -n $(NAMESPACE) -p '{"spec":{"suspend":true}}'
	@echo "CronJob suspended"

resume: ## Resume CronJob
	kubectl patch cronjob habit-notifier -n $(NAMESPACE) -p '{"spec":{"suspend":false}}'
	@echo "CronJob resumed"

clean: ## Delete all Kubernetes resources
	kubectl delete -k k8s/ || true
	@echo "Resources deleted"

clean-jobs: ## Delete completed jobs
	kubectl delete jobs -n $(NAMESPACE) --field-selector status.successful=1
	@echo "Completed jobs deleted"

update-image: ## Update CronJob image
	kubectl set image cronjob/habit-notifier habit-notifier=$(IMAGE_NAME):$(IMAGE_TAG) -n $(NAMESPACE)
	@echo "Image updated to: $(IMAGE_NAME):$(IMAGE_TAG)"

shell: ## Get shell in a debug pod
	kubectl run -it --rm debug --image=python:3.11-slim --restart=Never -n $(NAMESPACE) -- /bin/bash

events: ## Show recent events
	kubectl get events -n $(NAMESPACE) --sort-by='.lastTimestamp' | tail -20

validate: ## Validate Kubernetes manifests
	kubectl apply -k k8s/ --dry-run=client

lint: ## Lint Python code
	docker run --rm -v $$(pwd)/app:/app python:3.11-slim sh -c "pip install pylint && pylint /app/*.py"

security-scan: ## Scan Docker image for vulnerabilities
	docker scan $(IMAGE_NAME):$(IMAGE_TAG) || echo "Install docker scan or use trivy"
