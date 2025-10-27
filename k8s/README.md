# Habit Log Notification CronJob

This directory contains the Kubernetes CronJob implementation for sending habit log notifications to users via Twilio SMS.

## Overview

The notification job runs daily and performs the following tasks:

1. Queries the MySQL database for all users and their habit logs over the past 2 days
2. Determines which users need notifications based on the following criteria:
   - **No logs for both days**: User receives a notification about not logging any habits
   - **Log on day 1 but not day 2**: User receives a notification to maintain momentum
   - **Logs on both days**: No notification sent
3. Sends SMS notifications via Twilio to eligible users

## Architecture

### Components

- **Go Application** (`cmd/notification-job/main.go`): The notification job application
- **Dockerfile** (`cmd/notification-job/Dockerfile`): Multi-stage Docker build for the application
- **Kubernetes Manifests** (`k8s/`):
  - `configmap.yaml`: Configuration for database connection settings
  - `secret.yaml`: Sensitive credentials (database password, Twilio credentials)
  - `cronjob.yaml`: CronJob specification with resource limits and security settings

### Database Schema

The application expects the following MySQL tables:

- **users**: Contains user information (id, name, phone_number, time_zone, created_at)
- **habits**: Contains habit information (id, user_id, name, description, created_at)
- **logs**: Contains habit log entries (id, habit_id, created_at, notes)

## Prerequisites

Before deploying this CronJob, ensure you have:

1. A running Kubernetes cluster
2. A MySQL database with the required schema
3. Twilio account credentials (Account SID, Auth Token, and Phone Number)
4. Docker registry access to push the container image

## Building the Container Image

1. Build the Docker image:

```bash
cd /home/runner/work/ai-evaluation/ai-evaluation
docker build -f cmd/notification-job/Dockerfile -t your-registry/notification-job:latest .
```

2. Push the image to your registry:

```bash
docker push your-registry/notification-job:latest
```

## Deployment

### Step 1: Update the Secret

Edit `k8s/secret.yaml` and replace the placeholder values with your actual credentials:

```yaml
stringData:
  DB_PASSWORD: "your-actual-database-password"
  TWILIO_ACCOUNT_SID: "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  TWILIO_AUTH_TOKEN: "your-twilio-auth-token"
  TWILIO_PHONE_NUMBER: "+15551234567"
```

**Important**: For production, use a secure secret management solution like:
- Kubernetes External Secrets Operator
- HashiCorp Vault
- Cloud provider secret managers (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault)

### Step 2: Update the ConfigMap (Optional)

If your database connection details differ from the defaults, edit `k8s/configmap.yaml`:

```yaml
data:
  DB_HOST: "your-mysql-host"
  DB_PORT: "3306"
  DB_USER: "your-db-user"
  DB_NAME: "your-database-name"
```

### Step 3: Update the CronJob Image

Edit `k8s/cronjob.yaml` and update the image reference:

```yaml
containers:
- name: notification-job
  image: your-registry/notification-job:latest
```

Also adjust the schedule if needed (currently set to 6 PM UTC daily):

```yaml
spec:
  schedule: "0 18 * * *"  # Cron format: minute hour day month weekday
```

### Step 4: Deploy to Kubernetes

Apply the manifests to your cluster:

```bash
# Create the ConfigMap
kubectl apply -f k8s/configmap.yaml

# Create the Secret
kubectl apply -f k8s/secret.yaml

# Create the CronJob
kubectl apply -f k8s/cronjob.yaml
```

## Monitoring and Troubleshooting

### View CronJob Status

```bash
kubectl get cronjob habit-notification-job
```

### List Jobs Created by CronJob

```bash
kubectl get jobs --selector=app=habit-notification-job
```

### View Job Logs

```bash
# Get the most recent job
JOB_NAME=$(kubectl get jobs --selector=app=habit-notification-job --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}')

# View logs
kubectl logs job/$JOB_NAME
```

### Manual Job Trigger (Testing)

To manually trigger the job for testing:

```bash
kubectl create job --from=cronjob/habit-notification-job manual-test-$(date +%s)
```

### Common Issues

1. **Database Connection Failures**
   - Verify the database host and credentials in ConfigMap and Secret
   - Ensure the database is accessible from the Kubernetes cluster
   - Check network policies and firewall rules

2. **Twilio API Errors**
   - Verify Twilio credentials are correct
   - Ensure the Twilio phone number is verified and active
   - Check Twilio account balance and SMS capabilities

3. **Permission Issues**
   - The container runs as non-root user (UID 1000)
   - Ensure the image is built correctly with appropriate permissions

## Security Considerations

The CronJob implements several security best practices:

1. **Secret Management**: Sensitive credentials stored in Kubernetes Secrets
2. **Non-root User**: Container runs as UID 1000 (non-root)
3. **Read-only Root Filesystem**: Prevents runtime modifications
4. **Dropped Capabilities**: All Linux capabilities are dropped
5. **Resource Limits**: CPU and memory limits prevent resource exhaustion
6. **No Privilege Escalation**: Prevents container from gaining additional privileges

## Resource Configuration

The CronJob is configured with the following resource limits:

- **CPU Request**: 100m (0.1 CPU cores)
- **CPU Limit**: 200m (0.2 CPU cores)
- **Memory Request**: 64Mi
- **Memory Limit**: 128Mi

Adjust these values based on your actual usage patterns and cluster capacity.

## Customization

### Adjusting the Schedule

The schedule uses standard cron format:

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
│ │ │ │ │
* * * * *
```

Examples:
- `0 18 * * *` - Daily at 6 PM UTC
- `0 9,18 * * *` - Twice daily at 9 AM and 6 PM UTC
- `0 18 * * 1-5` - Weekdays at 6 PM UTC
- `*/30 * * * *` - Every 30 minutes

### Notification Messages

To customize notification messages, edit the message templates in `cmd/notification-job/main.go`:

```go
// No logs for 2 days
message = fmt.Sprintf("Hi %s! Custom message here...", user.Name)

// Has day 1 log but not day 2
message = fmt.Sprintf("Hi %s! Different message...", user.Name)
```

## Cleanup

To remove the CronJob and related resources:

```bash
kubectl delete cronjob habit-notification-job
kubectl delete secret notification-job-secret
kubectl delete configmap notification-job-config
```

## Support

For issues or questions:
1. Check the application logs
2. Verify database connectivity
3. Confirm Twilio credentials and configuration
4. Review Kubernetes events: `kubectl describe cronjob habit-notification-job`
