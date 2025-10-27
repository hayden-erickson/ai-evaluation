# Quick Start Deployment Guide

This guide will help you deploy the Habit Tracking Notification CronJob to a Kubernetes cluster in under 10 minutes.

## Prerequisites Checklist

- [ ] Kubernetes cluster (version 1.23+)
- [ ] `kubectl` configured to access your cluster
- [ ] Docker installed (for building the image)
- [ ] MySQL database accessible from the cluster
- [ ] Twilio account with:
  - [ ] Account SID
  - [ ] Auth Token
  - [ ] Phone number for sending SMS

## Step-by-Step Deployment

### 1. Build and Push the Docker Image

```bash
# Clone the repository if you haven't already
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation

# Build the Docker image
# Replace 'your-dockerhub-username' with your actual Docker Hub username
# or your private registry URL
docker build -t your-dockerhub-username/habit-notification-cron:v1.0.0 \
  -f cmd/notification-cron/Dockerfile .

# Log in to Docker Hub (or your registry)
docker login

# Push the image
docker push your-dockerhub-username/habit-notification-cron:v1.0.0
```

### 2. Update Configuration Files

#### Edit `k8s/configmap.yaml`:
```yaml
data:
  DB_HOST: "your-mysql-host"      # e.g., "mysql-service" or "10.0.1.5"
  DB_PORT: "3306"
  DB_USER: "your-db-username"     # e.g., "habits_app"
  DB_NAME: "habits"               # your database name
```

#### Edit `k8s/cronjob.yaml`:
Find the image line and update it:
```yaml
image: your-dockerhub-username/habit-notification-cron:v1.0.0
```

Optionally adjust the schedule (default is 8:00 AM UTC daily):
```yaml
schedule: "0 8 * * *"  # minute hour day month weekday
```

### 3. Create Kubernetes Secret

**Important:** Never commit real credentials to git!

```bash
# Create the secret with your actual credentials
kubectl create secret generic notification-cron-secrets \
  --from-literal=DB_PASSWORD='your-actual-database-password' \
  --from-literal=TWILIO_ACCOUNT_SID='ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  --from-literal=TWILIO_AUTH_TOKEN='your-actual-twilio-auth-token' \
  --from-literal=TWILIO_FROM_NUMBER='+15551234567' \
  -n default

# Verify
kubectl get secret notification-cron-secrets -n default
```

### 4. Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/cronjob.yaml

# Verify the deployment
kubectl get cronjob habit-notification-cron -n default
```

Expected output:
```
NAME                       SCHEDULE    SUSPEND   ACTIVE   LAST SCHEDULE   AGE
habit-notification-cron    0 8 * * *   False     0        <none>          10s
```

### 5. Test the CronJob

Don't wait for the schedule! Test it immediately:

```bash
# Create a manual job for testing
kubectl create job --from=cronjob/habit-notification-cron \
  habit-notification-test -n default

# Watch the job status
kubectl get jobs -n default -w

# View logs
kubectl logs -n default job/habit-notification-test

# Clean up test job
kubectl delete job habit-notification-test -n default
```

## Verification

### Check CronJob is Running
```bash
kubectl describe cronjob habit-notification-cron -n default
```

Look for:
- `Schedule: 0 8 * * *` (or your custom schedule)
- `Suspend: False`
- `Active: 0` (when not running)
- No error events

### Monitor Logs
```bash
# View logs from the most recent job
kubectl logs -l app=habit-notification-cron -n default --tail=50

# Follow logs in real-time during execution
kubectl logs -f -l app=habit-notification-cron -n default
```

Expected log output:
```
Starting habit tracking notification cron job
Successfully connected to database
Found 5 users who need notifications
Successfully sent notification to user 1 (+15551234567)
Successfully sent notification to user 2 (+15559876543)
...
Notification job completed. Success: 5, Failures: 0
```

## Common Issues and Solutions

### Issue: "Failed to connect to database"

**Solution:**
1. Verify database host is correct:
   ```bash
   kubectl get configmap notification-cron-config -o yaml
   ```
2. Test database connection from within the cluster:
   ```bash
   kubectl run -it --rm debug --image=mysql:8 --restart=Never -- \
     mysql -h your-mysql-host -u your-db-user -p
   ```
3. Check if database password is correct:
   ```bash
   # Don't run this in production - it will display the password!
   kubectl get secret notification-cron-secrets -o jsonpath='{.data.DB_PASSWORD}' | base64 -d
   ```

### Issue: "TWILIO_ACCOUNT_SID environment variable is required"

**Solution:**
Check that the secret exists and has all required fields:
```bash
kubectl describe secret notification-cron-secrets -n default
```

Expected output should show all four keys:
- `DB_PASSWORD`
- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_FROM_NUMBER`

### Issue: Twilio API errors (401 Unauthorized)

**Solution:**
Verify Twilio credentials:
1. Log in to [Twilio Console](https://console.twilio.com/)
2. Verify your Account SID and Auth Token
3. Recreate the secret with correct credentials:
   ```bash
   kubectl delete secret notification-cron-secrets -n default
   kubectl create secret generic notification-cron-secrets \
     --from-literal=DB_PASSWORD='...' \
     --from-literal=TWILIO_ACCOUNT_SID='AC...' \
     --from-literal=TWILIO_AUTH_TOKEN='...' \
     --from-literal=TWILIO_FROM_NUMBER='+1...' \
     -n default
   ```

### Issue: Image pull error

**Solution:**
1. Verify the image exists and is accessible:
   ```bash
   docker pull your-dockerhub-username/habit-notification-cron:v1.0.0
   ```
2. If using a private registry, create an image pull secret:
   ```bash
   kubectl create secret docker-registry regcred \
     --docker-server=https://index.docker.io/v1/ \
     --docker-username=your-username \
     --docker-password=your-password \
     --docker-email=your-email \
     -n default
   ```
3. Update `cronjob.yaml` to use the secret:
   ```yaml
   imagePullSecrets:
   - name: regcred
   ```

## Next Steps

### 1. Set Up Monitoring

Add Prometheus monitoring:
```bash
# Install Prometheus (if not already installed)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack

# Create a ServiceMonitor for the CronJob
kubectl apply -f k8s/servicemonitor.yaml
```

### 2. Configure Alerting

Create alerts for job failures:
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: habit-notification-alerts
spec:
  groups:
  - name: cronjob
    interval: 30s
    rules:
    - alert: CronJobFailed
      expr: kube_job_status_failed{job_name=~"habit-notification-cron.*"} > 0
      for: 5m
      annotations:
        summary: "Habit notification job failed"
```

### 3. Schedule Optimization

Consider different schedules based on user time zones:
- Morning reminder: `0 8 * * *` (8 AM UTC)
- Evening reminder: `0 20 * * *` (8 PM UTC)
- Multiple times: Deploy two CronJobs with different schedules

### 4. Scale for Large User Base

For > 10,000 users:
1. Increase resource limits in `cronjob.yaml`
2. Consider batch processing
3. Add retries for Twilio API rate limits
4. Use a message queue (e.g., Redis) for better handling

## Production Checklist

Before going to production:

- [ ] Database indexes created (see `k8s/README.md`)
- [ ] Secrets stored securely (consider HashiCorp Vault or AWS Secrets Manager)
- [ ] Image scanned for vulnerabilities
- [ ] Resource limits adjusted based on user count
- [ ] Network policies applied (restrict egress)
- [ ] Monitoring and alerting configured
- [ ] Backup CronJob schedule configured (failover)
- [ ] Tested manual job execution successfully
- [ ] Reviewed Twilio rate limits and adjusted batch sizes
- [ ] Set up log aggregation (e.g., ELK stack, CloudWatch)
- [ ] Documented incident response procedures

## Clean Up

To remove everything:

```bash
# Delete the CronJob
kubectl delete cronjob habit-notification-cron -n default

# Delete ConfigMap
kubectl delete configmap notification-cron-config -n default

# Delete Secret
kubectl delete secret notification-cron-secrets -n default

# Delete ServiceAccount and RBAC
kubectl delete serviceaccount notification-cron-sa -n default

# Delete any remaining jobs
kubectl delete jobs -l app=habit-notification-cron -n default
```

## Getting Help

- Check the main [README.md](./README.md) for detailed documentation
- Review logs: `kubectl logs -l app=habit-notification-cron --tail=100`
- Describe resources: `kubectl describe cronjob habit-notification-cron`
- View events: `kubectl get events --sort-by='.lastTimestamp'`

## Summary

You should now have:
✅ A Docker image built and pushed to a registry
✅ Kubernetes manifests configured for your environment
✅ Secrets created with your credentials
✅ CronJob deployed and scheduled
✅ Successful test run completed

The CronJob will now run automatically on schedule, sending notifications to users who haven't logged their habits.
