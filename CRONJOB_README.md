# Habit Tracker Notification CronJob

A production-ready Kubernetes CronJob that sends SMS notifications via Twilio to users who haven't logged their habits in the past 2 days.

## Features

- ✅ **Production-ready**: Security hardened with non-root user, read-only filesystem, dropped capabilities
- ✅ **Secure**: RBAC, NetworkPolicies, Secret management
- ✅ **Observable**: Comprehensive logging and error handling
- ✅ **Maintainable**: Modular code, clear separation of concerns
- ✅ **Configurable**: Environment-based configuration via ConfigMaps and Secrets
- ✅ **Resource-limited**: CPU and memory limits prevent resource exhaustion
- ✅ **Timezone-aware**: Respects user timezones for accurate day calculations

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes CronJob                        │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                     Job Pod                            │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │          habit_notifier.py                       │  │  │
│  │  │  1. Connect to MySQL                             │  │  │
│  │  │  2. Query users and habit logs                   │  │  │
│  │  │  3. Identify users needing notifications         │  │  │
│  │  │  4. Send SMS via Twilio                          │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
         │                                      │
         │ MySQL                                │ HTTPS
         ▼                                      ▼
    ┌─────────┐                          ┌──────────┐
    │ MySQL   │                          │ Twilio   │
    │ Database│                          │ API      │
    └─────────┘                          └──────────┘
```

## Notification Logic

Users receive notifications if:

1. **No logs in past 2 days**: User has 0 logs on both day 1 and day 2
2. **Missed yesterday**: User logged on day 1 but not on day 2

Users do NOT receive notifications if:
- Both days have logs
- Only day 2 has a log (they're back on track)

## Project Structure

```
new-cron/
├── app/
│   ├── habit_notifier.py      # Main application code
│   └── requirements.txt        # Python dependencies
├── k8s/
│   ├── namespace.yaml          # Namespace definition
│   ├── configmap.yaml          # Non-sensitive config
│   ├── secret.yaml             # Sensitive credentials
│   ├── serviceaccount.yaml     # RBAC configuration
│   ├── cronjob.yaml            # CronJob definition
│   ├── networkpolicy.yaml      # Network security
│   ├── kustomization.yaml      # Kustomize config
│   └── README.md               # K8s-specific docs
├── Dockerfile                  # Container definition
├── .dockerignore               # Docker build exclusions
├── Makefile                    # Automation commands
├── DEPLOYMENT.md               # Detailed deployment guide
└── CRONJOB_README.md          # This file
```

## Quick Start

### 1. Build and push Docker image

```bash
make build IMAGE_NAME=your-registry/habit-notifier IMAGE_TAG=v1.0.0
make push IMAGE_NAME=your-registry/habit-notifier IMAGE_TAG=v1.0.0
```

### 2. Configure secrets

```bash
# Create secret with your credentials
kubectl create secret generic habit-notifier-secrets \
  --namespace=habit-tracker \
  --from-literal=DB_USER='your-db-user' \
  --from-literal=DB_PASSWORD='your-db-password' \
  --from-literal=TWILIO_ACCOUNT_SID='your-sid' \
  --from-literal=TWILIO_AUTH_TOKEN='your-token'
```

### 3. Update configuration

Edit `k8s/configmap.yaml` with your values:
- Database host
- Twilio phone number

Edit `k8s/cronjob.yaml`:
- Docker image reference
- Schedule (default: daily at 9 AM UTC)

### 4. Deploy

```bash
make deploy
```

### 5. Test

```bash
make test
make logs
```

## Configuration

### Environment Variables

**From ConfigMap (non-sensitive):**
- `DB_HOST` - MySQL host
- `DB_PORT` - MySQL port (default: 3306)
- `DB_NAME` - Database name
- `TWILIO_FROM_NUMBER` - Twilio phone number

**From Secret (sensitive):**
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `TWILIO_ACCOUNT_SID` - Twilio account SID
- `TWILIO_AUTH_TOKEN` - Twilio auth token

### CronJob Schedule

The default schedule is `0 9 * * *` (daily at 9 AM UTC).

Common schedules:
- `0 9 * * *` - Daily at 9 AM
- `0 */6 * * *` - Every 6 hours
- `0 9 * * 1-5` - Weekdays at 9 AM
- `0 20 * * *` - Daily at 8 PM

## Database Schema

The application expects these tables:

### User
- `id` (PRIMARY KEY)
- `name`
- `phone_number`
- `time_zone`
- `created_at`

### Habit
- `id` (PRIMARY KEY)
- `user_id` (FOREIGN KEY → User)
- `name`
- `description`
- `created_at`

### Log
- `id` (PRIMARY KEY)
- `habit_id` (FOREIGN KEY → Habit)
- `notes`
- `created_at`

## Security Features

### Container Security
- ✅ Runs as non-root user (UID 1000)
- ✅ Read-only root filesystem
- ✅ No privilege escalation
- ✅ All Linux capabilities dropped
- ✅ Seccomp profile enabled

### Kubernetes Security
- ✅ RBAC with least-privilege ServiceAccount
- ✅ NetworkPolicy restricts egress
- ✅ Secrets for credential management
- ✅ Resource limits prevent DoS
- ✅ Pod Security Standards compatible

### Application Security
- ✅ Input validation
- ✅ Parameterized SQL queries (prevents SQL injection)
- ✅ Error handling and logging
- ✅ Connection timeouts

## Monitoring

### View CronJob status

```bash
kubectl get cronjob -n habit-tracker
```

### View job history

```bash
kubectl get jobs -n habit-tracker --sort-by=.metadata.creationTimestamp
```

### View logs

```bash
# Latest job logs
kubectl logs -n habit-tracker -l app=habit-notifier --tail=100

# Follow logs
kubectl logs -n habit-tracker -l app=habit-notifier -f

# Specific job
kubectl logs -n habit-tracker -l job-name=habit-notifier-28123456
```

### Metrics

The application logs include:
- Number of users checked
- Number of notifications sent
- Success/failure counts
- Error details

Example log output:
```
2024-01-15 09:00:01 - __main__ - INFO - Starting habit notification job
2024-01-15 09:00:02 - __main__ - INFO - Database connection established
2024-01-15 09:00:02 - __main__ - INFO - Twilio client initialized
2024-01-15 09:00:03 - __main__ - INFO - Found 150 users with phone numbers
2024-01-15 09:00:05 - __main__ - INFO - User John Doe (ID: 42) needs notification: No logs in the past 2 days
2024-01-15 09:00:06 - __main__ - INFO - Notification sent to John Doe (+1234567890). Message SID: SM123abc
2024-01-15 09:00:10 - __main__ - INFO - Job completed. Notifications sent: 12, Failed: 0
```

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

### Common Issues

**Database connection fails:**
- Verify `DB_HOST` in ConfigMap
- Check database credentials in Secret
- Ensure NetworkPolicy allows egress to database
- Test connectivity: `make shell` then try connecting to MySQL

**Twilio errors:**
- Verify credentials in Secret
- Check account balance
- Ensure phone numbers are in E.164 format (+1234567890)
- Check Twilio status: https://status.twilio.com/

**No notifications sent:**
- Check if users have phone numbers in database
- Verify timezone handling (users must have valid timezone)
- Check notification logic in logs

### Manual Testing

Run a one-off job:

```bash
make test
```

Or manually:

```bash
kubectl create job --from=cronjob/habit-notifier habit-notifier-manual -n habit-tracker
kubectl logs -n habit-tracker -l job-name=habit-notifier-manual -f
```

## Maintenance

### Suspend CronJob

```bash
make suspend
```

### Resume CronJob

```bash
make resume
```

### Update image

```bash
make build-push IMAGE_TAG=v1.1.0
make update-image IMAGE_TAG=v1.1.0
```

### Clean up old jobs

```bash
make clean-jobs
```

## Development

### Local testing

```bash
# Install dependencies
pip install -r app/requirements.txt

# Set environment variables
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=habit_tracker
export DB_USER=root
export DB_PASSWORD=password
export TWILIO_ACCOUNT_SID=your-sid
export TWILIO_AUTH_TOKEN=your-token
export TWILIO_FROM_NUMBER=+1234567890

# Run application
python app/habit_notifier.py
```

### Code linting

```bash
make lint
```

### Security scanning

```bash
make security-scan
```

## Production Checklist

Before deploying to production:

- [ ] Update all placeholder values in `secret.yaml`
- [ ] Configure correct database host in `configmap.yaml`
- [ ] Set Twilio phone number in `configmap.yaml`
- [ ] Update Docker image reference in `cronjob.yaml`
- [ ] Set appropriate CronJob schedule
- [ ] Configure resource limits based on workload
- [ ] Set up monitoring and alerting
- [ ] Test with manual job run
- [ ] Verify notifications are sent correctly
- [ ] Review and adjust NetworkPolicy if needed
- [ ] Consider using External Secrets Operator
- [ ] Set up log aggregation (ELK, Loki, etc.)
- [ ] Configure backup for job logs

## License

MIT

## Support

For issues or questions:
1. Check logs: `make logs`
2. Review events: `make events`
3. Check status: `make status`
4. See detailed deployment guide: `DEPLOYMENT.md`
