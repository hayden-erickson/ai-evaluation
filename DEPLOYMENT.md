# Kubernetes CronJob Deployment Guide

## Overview

This guide provides step-by-step instructions for deploying the Habit Tracker Notification CronJob to a production Kubernetes cluster.

## Architecture

The solution consists of:
- **Python Application**: Queries MySQL database and sends Twilio SMS notifications
- **Docker Container**: Containerized application with security best practices
- **Kubernetes CronJob**: Scheduled job that runs daily
- **ConfigMap**: Non-sensitive configuration
- **Secret**: Sensitive credentials (database, Twilio)
- **ServiceAccount & RBAC**: Least-privilege access control
- **NetworkPolicy**: Network segmentation and security

## Prerequisites

- Kubernetes cluster (1.25+ recommended for timezone support)
- `kubectl` configured to access your cluster
- Docker registry access (Docker Hub, GCR, ECR, etc.)
- MySQL database accessible from the cluster
- Twilio account with SMS capabilities

## Step 1: Build and Push Docker Image

### 1.1 Build the Docker image

```bash
cd new-cron
docker build -t your-registry/habit-notifier:v1.0.0 .
```

### 1.2 Test the image locally (optional)

```bash
docker run --rm \
  -e DB_HOST=your-mysql-host \
  -e DB_PORT=3306 \
  -e DB_NAME=habit_tracker \
  -e DB_USER=your-user \
  -e DB_PASSWORD=your-password \
  -e TWILIO_ACCOUNT_SID=your-sid \
  -e TWILIO_AUTH_TOKEN=your-token \
  -e TWILIO_FROM_NUMBER=+1234567890 \
  your-registry/habit-notifier:v1.0.0
```

### 1.3 Push to registry

```bash
docker push your-registry/habit-notifier:v1.0.0
```

## Step 2: Prepare Kubernetes Secrets

### 2.1 Create base64-encoded secrets

**IMPORTANT**: Never commit actual secrets to version control!

```bash
# Database credentials
echo -n 'your-db-user' | base64
echo -n 'your-db-password' | base64

# Twilio credentials
echo -n 'your-twilio-account-sid' | base64
echo -n 'your-twilio-auth-token' | base64
```

### 2.2 Update secret.yaml

Edit `k8s/secret.yaml` and replace the placeholder values with your actual base64-encoded credentials.

**Alternative: Use kubectl to create secret directly**

```bash
kubectl create secret generic habit-notifier-secrets \
  --namespace=habit-tracker \
  --from-literal=DB_USER='your-db-user' \
  --from-literal=DB_PASSWORD='your-db-password' \
  --from-literal=TWILIO_ACCOUNT_SID='your-twilio-sid' \
  --from-literal=TWILIO_AUTH_TOKEN='your-twilio-token' \
  --dry-run=client -o yaml > k8s/secret.yaml
```

### 2.3 (Recommended) Use External Secrets Operator

For production environments, consider using [External Secrets Operator](https://external-secrets.io/) with AWS Secrets Manager, HashiCorp Vault, or similar:

```bash
# Install External Secrets Operator
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets-system --create-namespace

# See commented example in k8s/secret.yaml
```

## Step 3: Configure Application Settings

### 3.1 Update ConfigMap

Edit `k8s/configmap.yaml`:

```yaml
data:
  DB_HOST: "your-mysql-host.database.svc.cluster.local"
  DB_PORT: "3306"
  DB_NAME: "habit_tracker"
  TWILIO_FROM_NUMBER: "+1234567890"  # Your Twilio phone number
```

### 3.2 Update CronJob schedule

Edit `k8s/cronjob.yaml` to set your desired schedule:

```yaml
spec:
  # Examples:
  schedule: "0 9 * * *"      # Daily at 9 AM UTC
  # schedule: "0 */6 * * *"  # Every 6 hours
  # schedule: "*/30 * * * *" # Every 30 minutes
  
  timeZone: "UTC"  # Or your preferred timezone
```

### 3.3 Update image reference

Edit `k8s/cronjob.yaml` and `k8s/kustomization.yaml` to reference your actual Docker image:

```yaml
image: your-registry/habit-notifier:v1.0.0
```

## Step 4: Deploy to Kubernetes

### 4.1 Deploy using kubectl

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Deploy all resources
kubectl apply -f k8s/
```

### 4.2 Deploy using Kustomize (recommended)

```bash
kubectl apply -k k8s/
```

### 4.3 Verify deployment

```bash
# Check CronJob
kubectl get cronjob -n habit-tracker

# Check if job has run
kubectl get jobs -n habit-tracker

# View logs
kubectl logs -n habit-tracker -l app=habit-notifier --tail=100
```

## Step 5: Testing

### 5.1 Trigger a manual job run

```bash
kubectl create job --from=cronjob/habit-notifier habit-notifier-manual-test -n habit-tracker
```

### 5.2 Monitor the job

```bash
# Watch job status
kubectl get jobs -n habit-tracker -w

# View pod logs
kubectl logs -n habit-tracker -l job-name=habit-notifier-manual-test -f
```

### 5.3 Check for errors

```bash
# Describe the job
kubectl describe job habit-notifier-manual-test -n habit-tracker

# Check pod events
kubectl get events -n habit-tracker --sort-by='.lastTimestamp'
```

## Step 6: Monitoring and Maintenance

### 6.1 View CronJob history

```bash
# List all jobs created by the CronJob
kubectl get jobs -n habit-tracker

# View successful jobs
kubectl get jobs -n habit-tracker --field-selector status.successful=1
```

### 6.2 Access logs

```bash
# Recent logs from all job runs
kubectl logs -n habit-tracker -l app=habit-notifier --tail=500

# Logs from specific job
kubectl logs -n habit-tracker -l job-name=habit-notifier-28123456
```

### 6.3 Suspend CronJob (if needed)

```bash
# Suspend
kubectl patch cronjob habit-notifier -n habit-tracker -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob habit-notifier -n habit-tracker -p '{"spec":{"suspend":false}}'
```

## Security Best Practices

### ✅ Implemented Security Features

1. **Non-root user**: Container runs as user 1000
2. **Read-only root filesystem**: Prevents tampering
3. **No privilege escalation**: Security context enforced
4. **Dropped capabilities**: All Linux capabilities dropped
5. **Network policies**: Restricts egress to only necessary services
6. **RBAC**: Least-privilege service account
7. **Secret management**: Credentials stored in Kubernetes Secrets
8. **Resource limits**: CPU and memory limits prevent resource exhaustion

### 🔒 Additional Recommendations

1. **Use External Secrets**: Integrate with AWS Secrets Manager or HashiCorp Vault
2. **Enable Pod Security Standards**: Use `restricted` policy
3. **Image scanning**: Scan Docker images for vulnerabilities
4. **Private registry**: Use private container registry with authentication
5. **Audit logging**: Enable Kubernetes audit logs
6. **Monitoring**: Set up Prometheus/Grafana for metrics
7. **Alerting**: Configure alerts for job failures

## Troubleshooting

### Job fails immediately

```bash
# Check pod status
kubectl get pods -n habit-tracker

# View pod details
kubectl describe pod <pod-name> -n habit-tracker

# Check logs
kubectl logs <pod-name> -n habit-tracker
```

**Common issues:**
- Missing or incorrect secrets
- Database connection issues
- Twilio authentication errors
- Image pull errors

### Database connection fails

1. Verify database host is accessible from cluster
2. Check database credentials in secret
3. Verify network policies allow egress to database
4. Test connection from a debug pod:

```bash
kubectl run -it --rm debug --image=mysql:8 --restart=Never -n habit-tracker -- \
  mysql -h your-mysql-host -u your-user -p
```

### Twilio API errors

1. Verify Twilio credentials are correct
2. Check Twilio account balance
3. Verify phone numbers are in E.164 format (+1234567890)
4. Check Twilio API status: https://status.twilio.com/

### Job doesn't run at scheduled time

1. Check CronJob schedule syntax
2. Verify timezone setting (requires K8s 1.25+)
3. Check if `startingDeadlineSeconds` is too restrictive
4. View CronJob events:

```bash
kubectl describe cronjob habit-notifier -n habit-tracker
```

## Updating the Application

### Update image version

```bash
# Build new version
docker build -t your-registry/habit-notifier:v1.1.0 .
docker push your-registry/habit-notifier:v1.1.0

# Update CronJob
kubectl set image cronjob/habit-notifier \
  habit-notifier=your-registry/habit-notifier:v1.1.0 \
  -n habit-tracker
```

### Update configuration

```bash
# Edit ConfigMap
kubectl edit configmap habit-notifier-config -n habit-tracker

# Edit Secret
kubectl edit secret habit-notifier-secrets -n habit-tracker
```

## Cleanup

### Remove all resources

```bash
# Delete namespace (removes everything)
kubectl delete namespace habit-tracker

# Or delete individual resources
kubectl delete -k k8s/
```

## Production Checklist

- [ ] Docker image built and pushed to registry
- [ ] Secrets created with actual credentials (not placeholders)
- [ ] ConfigMap updated with correct values
- [ ] CronJob schedule configured
- [ ] Image reference updated in manifests
- [ ] Database accessible from cluster
- [ ] Twilio account configured and funded
- [ ] RBAC permissions reviewed
- [ ] Network policies configured
- [ ] Resource limits appropriate for workload
- [ ] Monitoring and alerting configured
- [ ] Manual test job successful
- [ ] Documentation updated with environment-specific details

## Support

For issues or questions:
1. Check application logs: `kubectl logs -n habit-tracker -l app=habit-notifier`
2. Review Kubernetes events: `kubectl get events -n habit-tracker`
3. Verify configuration in ConfigMap and Secret
4. Test database and Twilio connectivity manually

## License

MIT
