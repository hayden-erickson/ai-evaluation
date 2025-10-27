# Kubernetes CronJob Implementation - Quick Start

This repository now includes a production-ready Kubernetes CronJob for sending habit log notifications.

## What Was Added

### Core Application
- **`cmd/notification-job/main.go`** - Go application that:
  - Connects to MySQL database
  - Queries all users and their habit logs for the past 2 days
  - Determines which users need notifications based on logging patterns
  - Sends SMS notifications via Twilio API

### Docker
- **`cmd/notification-job/Dockerfile`** - Multi-stage Docker build
- **`.dockerignore`** - Optimized build context

### Kubernetes Manifests (`k8s/`)
- **`cronjob.yaml`** - CronJob definition with security best practices
- **`configmap.yaml`** - Database configuration
- **`secret.yaml`** - Template for sensitive credentials
- **`networkpolicy.yaml`** - Optional network security policy
- **`kustomization.yaml`** - Kustomize configuration for easier deployment
- **`README.md`** - Comprehensive deployment and troubleshooting guide

### Tooling
- **`Makefile`** - Common operations (build, deploy, test, logs, etc.)

## Quick Deployment

1. **Build and push the Docker image:**
   ```bash
   make build-docker DOCKER_REGISTRY=your-registry
   make push-docker DOCKER_REGISTRY=your-registry
   ```

2. **Update secrets in `k8s/secret.yaml`:**
   ```yaml
   DB_PASSWORD: "your-actual-password"
   TWILIO_ACCOUNT_SID: "ACxxxxx..."
   TWILIO_AUTH_TOKEN: "your-token"
   TWILIO_PHONE_NUMBER: "+1234567890"
   ```

3. **Deploy to Kubernetes:**
   ```bash
   make deploy-k8s
   ```

4. **Test manually:**
   ```bash
   make test-job
   make logs
   ```

## Notification Logic

- **No logs for 2 days** → User receives reminder
- **Log on day 1, none on day 2** → User receives encouragement
- **Logs on both days** → No notification sent

## Schedule

Default: Daily at 6 PM UTC (`0 18 * * *`)
Edit `k8s/cronjob.yaml` to change the schedule.

## Security

✅ No vulnerabilities in dependencies
✅ CodeQL security scan passed
✅ Follows Kubernetes security best practices:
  - Non-root user execution
  - Read-only root filesystem
  - Resource limits
  - Dropped capabilities
  - Network policies (optional)

## Documentation

See `k8s/README.md` for complete documentation including:
- Detailed deployment instructions
- Monitoring and troubleshooting
- Security considerations
- Customization options

## Dependencies Added

- `github.com/go-sql-driver/mysql` v1.9.3
- `github.com/twilio/twilio-go` v1.28.4
