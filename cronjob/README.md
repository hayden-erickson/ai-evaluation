# Habit Notification CronJob

A production-ready Kubernetes CronJob that sends Twilio SMS notifications to users based on their habit logging activity over the past 2 days.

## Features

- ✅ Queries MySQL database for habit logs across all users
- ✅ Sends Twilio SMS notifications based on configurable rules
- ✅ Secure credential management using Kubernetes Secrets
- ✅ RBAC with least-privilege ServiceAccount
- ✅ Resource limits and health probes
- ✅ Comprehensive logging and error handling
- ✅ Production-ready security hardening
- ✅ Multi-stage Docker build for minimal image size

## Notification Rules

The CronJob sends notifications to users in the following scenarios:

1. **No logs in the last 2 days**: User receives a reminder to maintain their habit streak
2. **Log on day 1 but no log on day 2**: User receives a reminder not to break their streak
3. **Logs on both days**: No notification sent

## Architecture

```
┌─────────────────┐
│  Kubernetes     │
│  CronJob        │
│  (Daily @ 8AM)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      ┌─────────────────┐
│  Go Application │─────▶│  MySQL Database │
│  (Container)    │      │  (Habits Data)  │
└────────┬────────┘      └─────────────────┘
         │
         ▼
┌─────────────────┐
│  Twilio API     │
│  (SMS Gateway)  │
└─────────────────┘
```

## Prerequisites

- Kubernetes cluster (v1.21+)
- kubectl configured to access your cluster
- Docker for building container images
- MySQL database with the following schema:
  - `users` table (id, name, phone_number, time_zone, created_at)
  - `habits` table (id, user_id, name, description, created_at)
  - `logs` table (id, habit_id, created_at, notes)
- Twilio account with:
  - Account SID
  - Auth Token
  - Verified phone number

## Quick Start

### 1. Build the Docker Image

```bash
cd cronjob

# Build the image
docker build -t your-registry/habit-notifier:latest .

# Push to your registry
docker push your-registry/habit-notifier:latest
```

### 2. Configure Secrets

**Option A: Using kubectl (Recommended)**

```bash
kubectl create secret generic habit-notifier-secret \
  --from-literal=DB_USER=your_db_user \
  --from-literal=DB_PASSWORD=your_db_password \
  --from-literal=TWILIO_ACCOUNT_SID=your_twilio_account_sid \
  --from-literal=TWILIO_AUTH_TOKEN=your_twilio_auth_token \
  --from-literal=TWILIO_FROM_NUMBER=+1234567890 \
  --namespace=default
```

**Option B: Using YAML file**

Edit `k8s/secret.yaml` with your actual credentials and apply:

```bash
kubectl apply -f k8s/secret.yaml
```

⚠️ **Security Note**: Never commit secrets to version control. Use sealed-secrets, external-secrets, or a secrets management solution in production.

### 3. Configure Non-Sensitive Settings

Edit `k8s/configmap.yaml` to match your environment:

```yaml
data:
  DB_HOST: "your-mysql-host.svc.cluster.local"
  DB_PORT: "3306"
  DB_NAME: "your_database_name"
```

Apply the ConfigMap:

```bash
kubectl apply -f k8s/configmap.yaml
```

### 4. Create ServiceAccount and RBAC

```bash
kubectl apply -f k8s/serviceaccount.yaml
```

### 5. Deploy the CronJob

Edit `k8s/cronjob.yaml` to set your container image:

```yaml
containers:
- name: notifier
  image: your-registry/habit-notifier:latest
```

Optionally adjust the schedule (default is 8:00 AM UTC daily):

```yaml
spec:
  schedule: "0 8 * * *"  # Cron format: min hour day month weekday
  timeZone: "UTC"
```

Deploy:

```bash
kubectl apply -f k8s/cronjob.yaml
```

### 6. Deploy All Resources at Once

Alternatively, deploy everything at once:

```bash
kubectl apply -f k8s/
```

## Deployment Using Kustomize

For easier management across environments:

```bash
kubectl apply -k k8s/
```

## Verification

### Check CronJob Status

```bash
# View CronJob
kubectl get cronjob habit-notifier

# View scheduled Jobs
kubectl get jobs -l app=habit-notifier

# View Job logs
kubectl logs -l app=habit-notifier --tail=100
```

### Test the CronJob Manually

Create a one-time Job from the CronJob:

```bash
kubectl create job --from=cronjob/habit-notifier habit-notifier-manual-test
```

Watch the job:

```bash
kubectl get job habit-notifier-manual-test -w
```

View logs:

```bash
kubectl logs job/habit-notifier-manual-test -f
```

Clean up test job:

```bash
kubectl delete job habit-notifier-manual-test
```

## Configuration

### Environment Variables

#### From ConfigMap (Non-Sensitive)

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL database hostname | Required |
| `DB_PORT` | MySQL database port | `3306` |
| `DB_NAME` | Database name | Required |
| `LOG_LEVEL` | Application log level | `info` |
| `TZ` | Timezone | `UTC` |

#### From Secret (Sensitive)

| Variable | Description |
|----------|-------------|
| `DB_USER` | Database username |
| `DB_PASSWORD` | Database password |
| `TWILIO_ACCOUNT_SID` | Twilio Account SID |
| `TWILIO_AUTH_TOKEN` | Twilio Auth Token |
| `TWILIO_FROM_NUMBER` | Twilio phone number (E.164 format) |

### CronJob Schedule

Edit the schedule in `k8s/cronjob.yaml`:

```yaml
spec:
  schedule: "0 8 * * *"  # Daily at 8:00 AM
  timeZone: "UTC"
```

Common schedules:
- `0 8 * * *` - Daily at 8:00 AM
- `0 */6 * * *` - Every 6 hours
- `0 8 * * 1-5` - Weekdays at 8:00 AM
- `0 20 * * *` - Daily at 8:00 PM

### Resource Limits

Adjust resources in `k8s/cronjob.yaml`:

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "500m"
```

## Monitoring

### View Logs

```bash
# Stream logs from latest job
kubectl logs -l app=habit-notifier -f

# View logs from all jobs
kubectl logs -l app=habit-notifier --all-containers=true --tail=1000

# View logs from specific job
kubectl logs job/habit-notifier-<job-id>
```

### Check Job History

```bash
# View successful jobs
kubectl get jobs -l app=habit-notifier --field-selector status.successful=1

# View failed jobs
kubectl get jobs -l app=habit-notifier --field-selector status.failed=1
```

### Prometheus Metrics (if enabled)

The CronJob can be monitored using kube-state-metrics:

```promql
# CronJob status
kube_cronjob_status_active{cronjob="habit-notifier"}

# Job success/failure
kube_job_status_succeeded{job_name=~"habit-notifier-.*"}
kube_job_status_failed{job_name=~"habit-notifier-.*"}

# Job duration
kube_job_complete{job_name=~"habit-notifier-.*"}
```

## Troubleshooting

### Job Fails to Start

```bash
# Check CronJob events
kubectl describe cronjob habit-notifier

# Check pod events
kubectl describe pod -l app=habit-notifier
```

### Database Connection Issues

1. Verify database credentials in Secret:
```bash
kubectl get secret habit-notifier-secret -o yaml
```

2. Check database host is reachable from pod:
```bash
kubectl run -it --rm debug --image=mysql:8 --restart=Never -- \
  mysql -h <DB_HOST> -u <DB_USER> -p
```

3. Verify network policies allow egress to database

### Twilio API Issues

1. Verify Twilio credentials are correct
2. Check Twilio phone number is in E.164 format: `+1234567890`
3. Ensure outbound HTTPS (port 443) is allowed
4. Check Twilio account balance and sending limits

### Permission Errors

Check ServiceAccount and RBAC:

```bash
kubectl get serviceaccount habit-notifier-sa
kubectl get role habit-notifier-role
kubectl get rolebinding habit-notifier-rolebinding
```

### Pod Security Issues

If pod fails with security errors:

```bash
kubectl describe pod -l app=habit-notifier
```

Check your cluster's PodSecurityPolicy or Pod Security Standards.

## Security Best Practices

### ✅ Implemented

- ✅ Non-root user (UID 1000)
- ✅ Read-only root filesystem
- ✅ No privilege escalation
- ✅ Dropped all capabilities
- ✅ SecComp profile
- ✅ Secrets stored in Kubernetes Secrets (not in code)
- ✅ RBAC with least privilege
- ✅ Multi-stage Docker build
- ✅ Minimal Alpine base image
- ✅ No hardcoded credentials

### 🔒 Additional Recommendations

1. **Use External Secrets Operator** or similar for secrets management
2. **Enable Network Policies** to restrict egress/ingress
3. **Use Private Container Registry** with image scanning
4. **Implement Pod Security Standards** (restricted mode)
5. **Enable Audit Logging** for compliance
6. **Use Sealed Secrets** for GitOps workflows
7. **Rotate credentials** regularly
8. **Monitor for CVEs** in dependencies

## Maintenance

### Update the Application

```bash
# Build new image with version tag
docker build -t your-registry/habit-notifier:v1.1.0 .
docker push your-registry/habit-notifier:v1.1.0

# Update CronJob
kubectl set image cronjob/habit-notifier \
  notifier=your-registry/habit-notifier:v1.1.0

# Or update the YAML and reapply
kubectl apply -f k8s/cronjob.yaml
```

### Clean Up Old Jobs

The CronJob automatically retains:
- 3 successful jobs
- 3 failed jobs

To manually clean up:

```bash
# Delete completed jobs
kubectl delete jobs -l app=habit-notifier --field-selector status.successful=1

# Delete old jobs older than 7 days
kubectl get jobs -l app=habit-notifier -o json | \
  jq -r '.items[] | select(.status.completionTime < (now - 604800 | strftime("%Y-%m-%dT%H:%M:%SZ"))) | .metadata.name' | \
  xargs -I {} kubectl delete job {}
```

### Suspend CronJob

To temporarily disable without deleting:

```bash
kubectl patch cronjob habit-notifier -p '{"spec":{"suspend":true}}'
```

Resume:

```bash
kubectl patch cronjob habit-notifier -p '{"spec":{"suspend":false}}'
```

## Uninstallation

```bash
# Delete all resources
kubectl delete -f k8s/

# Or individually
kubectl delete cronjob habit-notifier
kubectl delete secret habit-notifier-secret
kubectl delete configmap habit-notifier-config
kubectl delete serviceaccount habit-notifier-sa
kubectl delete role habit-notifier-role
kubectl delete rolebinding habit-notifier-rolebinding
```

## Development

### Local Testing

1. **Set environment variables:**

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=habits_db
export TWILIO_ACCOUNT_SID=your_sid
export TWILIO_AUTH_TOKEN=your_token
export TWILIO_FROM_NUMBER=+1234567890
```

2. **Run locally:**

```bash
cd cronjob
go run .
```

3. **Build:**

```bash
go build -o habit-notifier
./habit-notifier
```

### Testing with Docker

```bash
docker build -t habit-notifier:dev .

docker run --rm \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=password \
  -e DB_NAME=habits_db \
  -e TWILIO_ACCOUNT_SID=your_sid \
  -e TWILIO_AUTH_TOKEN=your_token \
  -e TWILIO_FROM_NUMBER=+1234567890 \
  habit-notifier:dev
```

## Support

For issues or questions:
1. Check logs: `kubectl logs -l app=habit-notifier`
2. Review Kubernetes events: `kubectl get events --sort-by='.lastTimestamp'`
3. Verify database schema matches expected structure
4. Test Twilio credentials manually

## License

[Your License Here]

## Contributors

[Your Team Here]

