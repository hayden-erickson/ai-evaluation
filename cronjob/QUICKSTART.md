# Quick Start Guide

Get the Habit Notifier CronJob up and running in 5 minutes.

## Prerequisites

- Kubernetes cluster access
- Docker installed
- MySQL database with users, habits, and logs tables
- Twilio account credentials

## Step 1: Build and Push Docker Image

```bash
cd cronjob

# Update registry in Makefile or use directly:
docker build -t your-registry/habit-notifier:latest .
docker push your-registry/habit-notifier:latest
```

## Step 2: Create Kubernetes Secret

```bash
kubectl create secret generic habit-notifier-secret \
  --from-literal=DB_USER=your_db_user \
  --from-literal=DB_PASSWORD=your_db_password \
  --from-literal=TWILIO_ACCOUNT_SID=ACxxxx \
  --from-literal=TWILIO_AUTH_TOKEN=your_token \
  --from-literal=TWILIO_FROM_NUMBER=+1234567890 \
  --namespace=default
```

## Step 3: Update ConfigMap

Edit `k8s/configmap.yaml`:

```yaml
data:
  DB_HOST: "your-mysql-host.svc.cluster.local"
  DB_PORT: "3306"
  DB_NAME: "your_database_name"
```

## Step 4: Update CronJob Image

Edit `k8s/cronjob.yaml` line 90:

```yaml
image: your-registry/habit-notifier:latest
```

## Step 5: Deploy

```bash
# Apply all resources
kubectl apply -f k8s/

# Verify deployment
kubectl get cronjob habit-notifier
```

## Step 6: Test Manually

```bash
# Create a test job
kubectl create job --from=cronjob/habit-notifier test-run

# Watch logs
kubectl logs -l app=habit-notifier -f
```

## Verify It Works

1. Check job completed successfully:
   ```bash
   kubectl get jobs -l app=habit-notifier
   ```

2. View logs:
   ```bash
   kubectl logs -l app=habit-notifier --tail=50
   ```

3. Verify users received SMS notifications

## Troubleshooting

### Can't connect to database
- Verify `DB_HOST` in ConfigMap
- Check network policies allow egress
- Test connection: `kubectl run -it --rm debug --image=mysql:8 --restart=Never -- mysql -h <DB_HOST> -u <DB_USER> -p`

### Twilio errors
- Verify credentials in secret
- Check phone number format: `+1234567890`
- Ensure HTTPS egress allowed (port 443)

### Pod won't start
- Check image exists: `docker pull your-registry/habit-notifier:latest`
- View pod events: `kubectl describe pod -l app=habit-notifier`
- Check logs: `kubectl logs -l app=habit-notifier`

## Next Steps

- Adjust schedule in `k8s/cronjob.yaml` (default: 8 AM UTC daily)
- Set up monitoring and alerts
- Configure resource limits based on load
- Review security settings

For detailed documentation, see [README.md](README.md)

