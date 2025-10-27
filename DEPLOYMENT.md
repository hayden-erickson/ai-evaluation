# Quick Start Deployment Guide

This guide provides step-by-step instructions to deploy the Habit Log Notification CronJob to your Kubernetes cluster.

## Prerequisites Checklist

- [ ] Kubernetes cluster (v1.19+) with `kubectl` configured
- [ ] Docker installed and configured
- [ ] Access to a Docker registry (Docker Hub, GCR, ECR, etc.)
- [ ] MySQL database accessible from the cluster
- [ ] Twilio account with:
  - [ ] Account SID
  - [ ] Auth Token
  - [ ] Phone number for sending SMS

## Step-by-Step Deployment

### 1. Clone and Build

```bash
# Clone the repository
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation

# Download dependencies
make deps

# Test local build
make build
make clean
```

### 2. Build and Push Docker Image

```bash
# Set your Docker registry
export DOCKER_REGISTRY=your-docker-username  # e.g., myusername for Docker Hub
export IMAGE_TAG=v1.0.0

# Build the Docker image
make docker-build

# Login to your Docker registry (if needed)
docker login

# Push the image
make docker-push
```

**Alternative: Using Docker commands directly**
```bash
docker build -f Dockerfile.cronjob -t your-registry/habit-cronjob:v1.0.0 .
docker push your-registry/habit-cronjob:v1.0.0
```

### 3. Configure Kubernetes Secrets

**Option A: Using kubectl (Recommended)**
```bash
kubectl create secret generic habit-cronjob-secrets \
  --from-literal=db-password='your-database-password' \
  --from-literal=twilio-account-sid='ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  --from-literal=twilio-auth-token='your-twilio-auth-token'
```

**Option B: Using YAML file**
```bash
# Copy the example file
cp k8s/secret.example.yaml k8s/secret.yaml

# Encode your secrets
echo -n 'your-db-password' | base64
echo -n 'ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' | base64
echo -n 'your-twilio-auth-token' | base64

# Edit k8s/secret.yaml and paste the base64-encoded values
nano k8s/secret.yaml

# Apply the secret
kubectl apply -f k8s/secret.yaml
```

### 4. Configure Settings

Edit `k8s/configmap.yaml`:
```yaml
data:
  DB_HOST: "mysql-service.default.svc.cluster.local"  # Your MySQL service
  DB_PORT: "3306"
  DB_USER: "habit_user"           # Your database user
  DB_NAME: "habits"               # Your database name
  TWILIO_PHONE_NUMBER: "+1234567890"  # Your Twilio phone number
```

Edit `k8s/cronjob.yaml`:
```yaml
spec:
  schedule: "0 8 * * *"  # Adjust schedule as needed
  # ...
  containers:
  - name: notification-job
    image: your-registry/habit-cronjob:v1.0.0  # Update with your image
```

### 5. Deploy to Kubernetes

```bash
# Deploy all resources
make k8s-apply

# Or manually:
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/cronjob.yaml
```

### 6. Verify Deployment

```bash
# Check CronJob status
kubectl get cronjob habit-log-notification

# Expected output:
# NAME                      SCHEDULE    SUSPEND   ACTIVE   LAST SCHEDULE   AGE
# habit-log-notification    0 8 * * *   False     0        <none>          10s
```

### 7. Test the CronJob

**Manually trigger a test run:**
```bash
# Create a test job
make k8s-trigger

# Or manually:
kubectl create job --from=cronjob/habit-log-notification test-run-$(date +%s)

# Wait a few seconds, then check job status
kubectl get jobs

# View logs
make k8s-logs

# Or manually:
kubectl logs job/test-run-<timestamp>
```

### 8. Monitor Logs

```bash
# View most recent job logs
make k8s-logs

# Follow logs in real-time (if job is running)
JOB_NAME=$(kubectl get jobs --selector=app=habit-tracker --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}')
kubectl logs job/$JOB_NAME -f

# Check job history
kubectl get jobs --selector=app=habit-tracker
```

## Common Configuration Scenarios

### Different Time Zones

The CronJob schedule is in UTC. To run at 8 AM in your local timezone:

| Timezone | Schedule | Description |
|----------|----------|-------------|
| UTC | `0 8 * * *` | 8:00 AM UTC |
| EST (UTC-5) | `0 13 * * *` | 8:00 AM EST |
| PST (UTC-8) | `0 16 * * *` | 8:00 AM PST |
| CST (UTC-6) | `0 14 * * *` | 8:00 AM CST |

### Multiple Runs Per Day

```yaml
# Every 6 hours
schedule: "0 */6 * * *"

# Twice daily (8 AM and 8 PM UTC)
# Create two separate CronJobs or use: "0 8,20 * * *"
```

### Weekdays Only

```yaml
# Monday through Friday at 8 AM UTC
schedule: "0 8 * * 1-5"
```

## Troubleshooting

### Issue: CronJob created but no jobs running

**Solution:**
```bash
# Check if CronJob is suspended
kubectl get cronjob habit-log-notification -o jsonpath='{.spec.suspend}'

# If true, resume it
kubectl patch cronjob habit-log-notification -p '{"spec":{"suspend":false}}'

# Check for recent events
kubectl describe cronjob habit-log-notification
```

### Issue: Job fails with "ImagePullBackOff"

**Solution:**
```bash
# Verify image exists and is accessible
docker pull your-registry/habit-cronjob:v1.0.0

# Check if registry credentials are needed
kubectl create secret docker-registry registry-credentials \
  --docker-server=your-registry \
  --docker-username=your-username \
  --docker-password=your-password

# Add to cronjob.yaml:
# imagePullSecrets:
# - name: registry-credentials
```

### Issue: Job fails with database connection error

**Solution:**
```bash
# Test database connectivity from a pod
kubectl run -it --rm debug --image=mysql:8 --restart=Never -- \
  mysql -h mysql-service -u habit_user -p

# Check if database host is correct in ConfigMap
kubectl get configmap habit-cronjob-config -o yaml

# Verify database credentials in Secret
kubectl get secret habit-cronjob-secrets -o jsonpath='{.data.db-password}' | base64 -d
```

### Issue: Twilio errors

**Common Twilio error codes:**
- `21211`: Invalid phone number - ensure format includes country code (+1234567890)
- `21408`: Permission denied - verify phone number in Twilio console (for trial accounts)
- `20003`: Authentication error - check Account SID and Auth Token

## Maintenance

### Update Application

```bash
# Build new version
export IMAGE_TAG=v1.1.0
make docker-build-push

# Update cronjob.yaml with new image tag
# Then apply changes
kubectl apply -f k8s/cronjob.yaml

# Kubernetes will use the new image for next scheduled run
```

### Suspend CronJob Temporarily

```bash
kubectl patch cronjob habit-log-notification -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob habit-log-notification -p '{"spec":{"suspend":false}}'
```

### Clean Up

```bash
# Remove all resources
make k8s-delete

# Or manually
kubectl delete -f k8s/cronjob.yaml
kubectl delete -f k8s/secret.yaml
kubectl delete -f k8s/configmap.yaml
kubectl delete -f k8s/rbac.yaml
```

## Security Best Practices

1. **Never commit secrets**: The actual `k8s/secret.yaml` is gitignored
2. **Use RBAC**: The ServiceAccount has minimal permissions
3. **Limit resources**: Prevent resource exhaustion with limits
4. **Run as non-root**: Container runs as user 1000
5. **Read-only filesystem**: Container cannot write to its filesystem
6. **Rotate credentials**: Regularly update Twilio and database credentials

## Support

For issues and questions:
- Check the main [README](k8s/README.md) for detailed documentation
- Review logs: `make k8s-logs`
- Check job status: `make k8s-status`
- Verify configuration: `kubectl get configmap habit-cronjob-config -o yaml`
