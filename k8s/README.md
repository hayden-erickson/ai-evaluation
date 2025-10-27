# Habit Tracking Notification CronJob

This directory contains Kubernetes manifests for deploying a CronJob that sends SMS notifications to users who haven't logged their habits.

## Overview

The notification CronJob performs the following tasks:
1. Connects to a MySQL database to query habit logs
2. Identifies users who meet notification criteria:
   - Users with no logs in the past 2 days
   - Users who logged habits yesterday but not today
3. Sends personalized SMS notifications via Twilio
4. Logs all activities and errors for monitoring

## Architecture

```
┌─────────────────┐
│  Kubernetes     │
│  CronJob        │
│  (Scheduled)    │
└────────┬────────┘
         │
         ├──► MySQL Database (queries habit logs)
         │
         └──► Twilio API (sends SMS notifications)
```

## Prerequisites

1. **Kubernetes Cluster** (version 1.23+)
   - Required for CronJob support
   - Version 1.25+ recommended for timezone support

2. **MySQL Database**
   - Database with users, habits, and logs tables
   - Network accessible from the Kubernetes cluster
   
3. **Twilio Account**
   - Account SID
   - Auth Token
   - Twilio phone number (for sending SMS)

4. **Container Registry** (optional)
   - Docker Hub, Google Container Registry, or similar
   - For storing the Docker image

## Files

- `cronjob.yaml` - CronJob resource definition with scheduling and configuration
- `configmap.yaml` - Non-sensitive configuration (database host, port, etc.)
- `secret.yaml.template` - Template for sensitive credentials (NOT for git)
- `rbac.yaml` - ServiceAccount for the CronJob (minimal permissions)
- `Dockerfile` - Multi-stage Docker build for the cron job application
- `README.md` - This file

## Setup Instructions

### Step 1: Build and Push Docker Image

```bash
# Navigate to the repository root
cd /path/to/ai-evaluation

# Build the Docker image
docker build -t your-registry/habit-notification-cron:v1.0.0 -f cmd/notification-cron/Dockerfile .

# Push to your container registry
docker push your-registry/habit-notification-cron:v1.0.0
```

### Step 2: Configure Kubernetes Manifests

1. **Update ConfigMap** (`configmap.yaml`):
   ```yaml
   DB_HOST: "your-mysql-host"      # MySQL service name or IP
   DB_PORT: "3306"                  # MySQL port
   DB_USER: "your-db-user"          # Database username
   DB_NAME: "habits"                # Database name
   ```

2. **Update CronJob** (`cronjob.yaml`):
   - Set the correct image path:
     ```yaml
     image: your-registry/habit-notification-cron:v1.0.0
     ```
   - Adjust the schedule if needed (default: 8:00 AM UTC daily):
     ```yaml
     schedule: "0 8 * * *"  # minute hour day month weekday
     ```

### Step 3: Create Secrets

**Important:** Never commit secrets to git. Create them directly in Kubernetes:

```bash
# Create the secret with your actual credentials
kubectl create secret generic notification-cron-secrets \
  --from-literal=DB_PASSWORD='your-database-password' \
  --from-literal=TWILIO_ACCOUNT_SID='ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  --from-literal=TWILIO_AUTH_TOKEN='your-twilio-auth-token' \
  --from-literal=TWILIO_FROM_NUMBER='+1234567890' \
  -n default

# Verify the secret was created
kubectl get secret notification-cron-secrets -n default
```

### Step 4: Deploy to Kubernetes

```bash
# Apply the manifests in order
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/cronjob.yaml

# Verify deployment
kubectl get cronjob habit-notification-cron -n default
```

## Verification

### Check CronJob Status
```bash
# View the CronJob
kubectl get cronjob habit-notification-cron -n default

# View CronJob details
kubectl describe cronjob habit-notification-cron -n default
```

### Monitor Job Execution
```bash
# List jobs created by the CronJob
kubectl get jobs -n default -l app=habit-notification-cron

# View logs from the most recent job
kubectl logs -n default -l app=habit-notification-cron --tail=100
```

### Manual Trigger (for testing)
```bash
# Create a one-time job from the CronJob for testing
kubectl create job --from=cronjob/habit-notification-cron habit-notification-manual -n default

# Watch the job
kubectl get jobs -n default -w

# View logs
kubectl logs -n default job/habit-notification-manual
```

## Configuration

### Schedule Format

The CronJob uses standard cron syntax:
```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

Examples:
- `0 8 * * *` - Daily at 8:00 AM UTC
- `0 9 * * 1-5` - Weekdays at 9:00 AM UTC
- `0 */6 * * *` - Every 6 hours
- `30 7 * * *` - Daily at 7:30 AM UTC

### Resource Limits

Current settings:
```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "128Mi"
    cpu: "200m"
```

Adjust based on your user base:
- Small (< 1000 users): Default is fine
- Medium (1000-10000 users): Increase memory to 256Mi
- Large (> 10000 users): Increase memory to 512Mi and cpu to 500m

### Environment Variables

**From ConfigMap:**
- `DB_HOST` - MySQL host
- `DB_PORT` - MySQL port
- `DB_USER` - MySQL username
- `DB_NAME` - Database name

**From Secret:**
- `DB_PASSWORD` - MySQL password
- `TWILIO_ACCOUNT_SID` - Twilio account identifier
- `TWILIO_AUTH_TOKEN` - Twilio authentication token
- `TWILIO_FROM_NUMBER` - Twilio phone number for sending SMS

## Security Best Practices

### Implemented Security Features

1. **Non-root User**: Container runs as user 1000 (non-root)
2. **Read-only Filesystem**: Container filesystem is read-only
3. **No Privilege Escalation**: `allowPrivilegeEscalation: false`
4. **Dropped Capabilities**: All Linux capabilities dropped
5. **Seccomp Profile**: Runtime default seccomp profile applied
6. **Secret Management**: Credentials stored in Kubernetes Secrets (not in code)
7. **Minimal RBAC**: ServiceAccount has no special permissions

### Additional Recommendations

1. **Use a Secret Management System**
   - Consider using tools like HashiCorp Vault, AWS Secrets Manager, or Google Secret Manager
   - Integrate with Kubernetes using the External Secrets Operator

2. **Network Policies**
   - Restrict network access to only MySQL and Twilio API endpoints
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: notification-cron-netpol
   spec:
     podSelector:
       matchLabels:
         app: habit-notification-cron
     policyTypes:
     - Egress
     egress:
     - to:
       - podSelector:
           matchLabels:
             app: mysql  # Your MySQL pod label
       ports:
       - protocol: TCP
         port: 3306
     - to:
       - namespaceSelector: {}
       ports:
       - protocol: TCP
         port: 443  # For Twilio HTTPS API
   ```

3. **Image Scanning**
   - Scan Docker images for vulnerabilities before deployment
   - Use tools like Trivy, Clair, or cloud provider scanning

4. **Audit Logging**
   - Enable Kubernetes audit logging to track CronJob executions
   - Monitor for unexpected behavior

## Troubleshooting

### CronJob Not Running

```bash
# Check if CronJob is suspended
kubectl get cronjob habit-notification-cron -o jsonpath='{.spec.suspend}'

# View CronJob events
kubectl describe cronjob habit-notification-cron

# Check for scheduling issues
kubectl get events --sort-by='.lastTimestamp' | grep habit-notification
```

### Job Failed

```bash
# View failed jobs
kubectl get jobs -l app=habit-notification-cron --field-selector status.successful=0

# Get logs from failed job
kubectl logs -l app=habit-notification-cron --tail=200

# Describe the job to see error details
kubectl describe job <job-name>
```

### Database Connection Issues

```bash
# Test database connectivity from a debug pod
kubectl run -it --rm debug --image=mysql:8 --restart=Never -- \
  mysql -h mysql-service -u habits_app -p

# Check if ConfigMap values are correct
kubectl get configmap notification-cron-config -o yaml

# Verify secret exists
kubectl get secret notification-cron-secrets
```

### Twilio Issues

Common issues:
- Invalid credentials: Check `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN`
- Invalid phone number: Ensure `TWILIO_FROM_NUMBER` is in E.164 format (+1234567890)
- Rate limiting: Twilio has rate limits; consider adding delays or batching
- Unverified numbers: In trial accounts, you can only send to verified numbers

### View Logs

```bash
# Get logs from the most recent successful job
kubectl logs -l app=habit-notification-cron --tail=100

# Follow logs in real-time
kubectl logs -f -l app=habit-notification-cron

# Get logs from a specific job
kubectl logs job/<job-name>
```

## Monitoring and Alerts

### Metrics to Monitor

1. **Job Success Rate**: Track successful vs. failed job runs
2. **Execution Duration**: Monitor how long each job takes
3. **Notification Success**: Track Twilio API success/failure
4. **Resource Usage**: Monitor memory and CPU consumption

### Setting Up Alerts

Example Prometheus alert:
```yaml
- alert: HabitNotificationCronJobFailed
  expr: kube_job_status_failed{job_name=~"habit-notification-cron.*"} > 0
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Habit notification cron job failed"
    description: "CronJob {{ $labels.job_name }} has failed"
```

## Database Schema Requirements

The CronJob expects the following database schema:

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    time_zone VARCHAR(50) NOT NULL,
    profile_image_url VARCHAR(500),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    habit_id BIGINT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id)
);

-- Recommended indexes for performance
CREATE INDEX idx_logs_habit_id_created_at ON logs(habit_id, created_at);
CREATE INDEX idx_habits_user_id ON habits(user_id);
```

## Notification Logic

Users receive notifications if:

1. **No logs in past 2 days**
   - Query finds no log entries for any of the user's habits in the last 48 hours
   - Message: "We noticed you haven't logged any habits in the past 2 days..."

2. **Logged yesterday but not today**
   - User has at least one log between 48-24 hours ago
   - User has no logs in the last 24 hours
   - Message: "You logged your habits yesterday but not today..."

Users who have logged habits in the last 24 hours receive no notification.

## Maintenance

### Updating the CronJob

```bash
# Update the image version
kubectl set image cronjob/habit-notification-cron \
  notification-cron=your-registry/habit-notification-cron:v1.1.0

# Or edit the manifest and reapply
kubectl apply -f k8s/cronjob.yaml
```

### Suspending the CronJob

```bash
# Suspend (stop creating new jobs)
kubectl patch cronjob habit-notification-cron -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob habit-notification-cron -p '{"spec":{"suspend":false}}'
```

### Cleanup

```bash
# Delete the CronJob (keeps existing jobs)
kubectl delete cronjob habit-notification-cron

# Delete all resources
kubectl delete -f k8s/

# Delete the secret
kubectl delete secret notification-cron-secrets
```

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review job logs for error messages
3. Verify all configuration values are correct
4. Check Twilio console for API errors

## License

MIT
