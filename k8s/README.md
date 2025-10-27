# Habit Log Notification CronJob

This directory contains a production-ready Kubernetes CronJob that sends SMS notifications via Twilio to users who haven't logged their habits recently.

## Overview

The CronJob performs the following:
1. Queries a MySQL database for habit logs across all users over the past 2 days
2. Sends a notification via Twilio to each user's phone number if:
   - The user has no logs over the last 2 days, OR
   - The user has 1 log on day 1, but no log on day 2
3. If both days have logs, no notification is sent

## Architecture

### Components

- **CronJob**: Kubernetes CronJob that runs daily at 8:00 AM UTC
- **ConfigMap**: Non-sensitive configuration (database host, port, etc.)
- **Secret**: Sensitive credentials (database password, Twilio credentials)
- **ServiceAccount & RBAC**: Minimal permissions for the job to run

### Application Logic

The Go application (`cmd/cronjob/main.go`):
- Connects to MySQL database
- Queries all users and their habit logs
- Determines which users need notifications based on logging patterns
- Sends personalized SMS via Twilio API
- Logs all operations for monitoring

## Prerequisites

1. **Kubernetes Cluster**: A running Kubernetes cluster (v1.19+)
2. **MySQL Database**: With the following schema:
   - `users` table: id, phone_number, time_zone, name
   - `habits` table: id, user_id, name, description
   - `logs` table: id, habit_id, created_at, notes
3. **Twilio Account**: 
   - Account SID
   - Auth Token
   - Phone Number for sending SMS
4. **Docker Registry**: To store the container image

## Setup Instructions

### 1. Build and Push Docker Image

```bash
# Build the Docker image
docker build -f Dockerfile.cronjob -t your-registry/habit-cronjob:latest .

# Push to your registry
docker push your-registry/habit-cronjob:latest
```

### 2. Configure Secrets

Create base64-encoded values for your secrets:

```bash
# Database password
echo -n 'your-db-password' | base64

# Twilio Account SID
echo -n 'your-twilio-account-sid' | base64

# Twilio Auth Token
echo -n 'your-twilio-auth-token' | base64
```

Edit `k8s/secret.yaml` and add the base64-encoded values:

```yaml
data:
  db-password: "eW91ci1kYi1wYXNzd29yZA=="
  twilio-account-sid: "eW91ci10d2lsaW8tYWNjb3VudC1zaWQ="
  twilio-auth-token: "eW91ci10d2lsaW8tYXV0aC10b2tlbg=="
```

### 3. Configure ConfigMap

Edit `k8s/configmap.yaml` with your actual values:

```yaml
data:
  DB_HOST: "mysql-service.default.svc.cluster.local"
  DB_PORT: "3306"
  DB_USER: "habit_user"
  DB_NAME: "habits"
  TWILIO_PHONE_NUMBER: "+1234567890"
```

### 4. Update CronJob Image

Edit `k8s/cronjob.yaml` and update the image reference:

```yaml
containers:
- name: notification-job
  image: your-registry/habit-cronjob:v1.0.0  # Update this
```

### 5. Deploy to Kubernetes

Deploy all components in order:

```bash
# Create namespace (optional)
kubectl create namespace habit-tracker

# Deploy RBAC (ServiceAccount, Role, RoleBinding)
kubectl apply -f k8s/rbac.yaml

# Deploy ConfigMap
kubectl apply -f k8s/configmap.yaml

# Deploy Secret
kubectl apply -f k8s/secret.yaml

# Deploy CronJob
kubectl apply -f k8s/cronjob.yaml
```

## Configuration

### Environment Variables

#### Database Configuration
- `DB_HOST`: MySQL host (default: localhost)
- `DB_PORT`: MySQL port (default: 3306)
- `DB_USER`: Database user (default: root)
- `DB_PASSWORD`: Database password (required, from Secret)
- `DB_NAME`: Database name (default: habits)

#### Twilio Configuration
- `TWILIO_ACCOUNT_SID`: Twilio Account SID (required, from Secret)
- `TWILIO_AUTH_TOKEN`: Twilio Auth Token (required, from Secret)
- `TWILIO_PHONE_NUMBER`: Twilio phone number for sending SMS (required)

### CronJob Schedule

The default schedule is `0 8 * * *` (daily at 8:00 AM UTC).

To change the schedule, edit `k8s/cronjob.yaml`:

```yaml
spec:
  schedule: "0 8 * * *"  # Cron format: minute hour day month weekday
```

Examples:
- `0 8 * * *`: Daily at 8:00 AM
- `0 */6 * * *`: Every 6 hours
- `0 8 * * 1-5`: Weekdays at 8:00 AM
- `0 20 * * *`: Daily at 8:00 PM

## Monitoring and Debugging

### View CronJob Status

```bash
# List CronJobs
kubectl get cronjobs

# Describe CronJob
kubectl describe cronjob habit-log-notification

# List Jobs created by CronJob
kubectl get jobs --selector=app=habit-log-notification
```

### View Logs

```bash
# Get the most recent job
JOB_NAME=$(kubectl get jobs --selector=app=habit-log-notification --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}')

# View logs
kubectl logs job/$JOB_NAME

# Follow logs in real-time (if job is running)
kubectl logs job/$JOB_NAME -f
```

### Manual Trigger

To manually trigger the CronJob for testing:

```bash
# Create a one-time Job from the CronJob
kubectl create job --from=cronjob/habit-log-notification manual-test-$(date +%s)

# View the job
kubectl get jobs

# View logs
kubectl logs job/manual-test-<timestamp>
```

### Debug Failed Jobs

```bash
# List failed jobs
kubectl get jobs --field-selector status.successful=0

# Describe a failed job
kubectl describe job <job-name>

# Get pod logs from failed job
kubectl logs -l job-name=<job-name>
```

## Resource Usage

The CronJob is configured with the following resource limits:

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "500m"
```

Adjust these based on your workload:
- For large user bases (10,000+ users), increase memory to 512Mi
- For slower database connections, increase activeDeadlineSeconds

## Security Features

1. **Non-root User**: Container runs as user 1000 (non-root)
2. **Read-only Filesystem**: Root filesystem is read-only
3. **Dropped Capabilities**: All Linux capabilities are dropped
4. **No Privilege Escalation**: `allowPrivilegeEscalation: false`
5. **Secret Management**: Sensitive data stored in Kubernetes Secrets
6. **RBAC**: Minimal permissions via ServiceAccount

## Troubleshooting

### Common Issues

#### 1. "Failed to connect to database"
- Verify database host and port in ConfigMap
- Check database credentials in Secret
- Ensure database is accessible from the cluster
- Test network connectivity: `kubectl run -it --rm debug --image=mysql:8 --restart=Never -- mysql -h <DB_HOST> -u <DB_USER> -p`

#### 2. "Twilio error: 21211"
- Invalid phone number format
- Ensure phone numbers in database include country code (e.g., +1234567890)

#### 3. "Twilio error: 21408"
- Permission denied to send to this number
- Check Twilio account trial limitations
- Verify phone numbers are verified in Twilio console (for trial accounts)

#### 4. CronJob not running at scheduled time
- Check CronJob status: `kubectl get cronjob habit-log-notification`
- Verify schedule format is correct
- Check for `startingDeadlineSeconds` expiration

#### 5. Job timeout
- Increase `activeDeadlineSeconds` in cronjob.yaml
- Optimize database queries
- Check database performance

## Testing

### Local Testing

Test the application locally before deploying:

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your-password
export DB_NAME=habits
export TWILIO_ACCOUNT_SID=your-sid
export TWILIO_AUTH_TOKEN=your-token
export TWILIO_PHONE_NUMBER=+1234567890

# Run the application
go run cmd/cronjob/main.go
```

### Kubernetes Testing

1. Create a test job manually (see "Manual Trigger" above)
2. Monitor logs for errors
3. Verify notifications are sent
4. Check Twilio dashboard for message delivery status

## Maintenance

### Updating the CronJob

1. Build new image with updated code
2. Push to registry with new tag
3. Update image in `k8s/cronjob.yaml`
4. Apply changes: `kubectl apply -f k8s/cronjob.yaml`

### Suspending the CronJob

To temporarily disable the CronJob:

```bash
kubectl patch cronjob habit-log-notification -p '{"spec":{"suspend":true}}'
```

To resume:

```bash
kubectl patch cronjob habit-log-notification -p '{"spec":{"suspend":false}}'
```

### Cleaning Up

To remove all components:

```bash
kubectl delete -f k8s/cronjob.yaml
kubectl delete -f k8s/secret.yaml
kubectl delete -f k8s/configmap.yaml
kubectl delete -f k8s/rbac.yaml
```

## Database Schema

The application expects the following MySQL schema:

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    profile_image_url TEXT,
    name VARCHAR(255) NOT NULL,
    time_zone VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    habit_id BIGINT NOT NULL,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_habits_user_id ON habits(user_id);
CREATE INDEX idx_logs_habit_id ON logs(habit_id);
CREATE INDEX idx_logs_created_at ON logs(created_at);
```

## Performance Considerations

- **Batch Processing**: For large user bases (100,000+), consider implementing batch processing
- **Database Indexing**: Ensure `created_at` is indexed on the `logs` table
- **Connection Pooling**: MySQL driver uses connection pooling by default
- **Rate Limiting**: Twilio has rate limits - consider implementing delays for large batches

## License

MIT
