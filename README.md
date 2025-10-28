# Habit Tracker REST API

A RESTful API built in Go for tracking user habits with JWT-based authentication, SQLite database, and comprehensive security features.

## Features

- **JWT-based authentication** - Secure token-based authentication
- **Input validation** - Comprehensive request validation for all endpoints
- **Error logging** - All errors are logged with appropriate context
- **HTTP status codes** - Proper error responses with meaningful messages
- **Modular architecture** - Separation of concerns with repository and service layers
- **Dependency injection** - Clean, testable code structure
- **Security headers** - XSS protection, clickjacking prevention, CSP
- **RBAC** - Role-based access control (users can only access their own data)
- **SQLite database** - Lightweight database for local development

## Architecture

The application follows a clean, modular architecture:

### Package Structure

- **`/models`** - Data structures and validation logic (User, Habit, Log)
- **`/repository`** - Database operations and data access layer
- **`/service`** - Business logic and service layer
- **`/handlers`** - HTTP request handlers
- **`/middleware`** - Authentication, logging, and security middleware
- **`/config`** - Database configuration and migrations
- **`/utils`** - Utility functions (JWT, password hashing)
- **`/migrations`** - SQL migration files

## Prerequisites

- Go 1.21 or higher
- SQLite3

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

2. Install dependencies:
```bash
go mod download
```

3. (Optional) Set environment variables:
```bash
export PORT=8080                    # Default: 8080
export JWT_SECRET=your-secret-key   # Default: "default-secret-key-change-in-production"
export DB_PATH=habits.db            # Default: habits.db
```

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

## API Endpoints

### Authentication

#### Register a new user
```http
POST /users/register
Content-Type: application/json

{
  "name": "John Doe",
  "phone_number": "+1234567890",
  "password": "securepassword123",
  "time_zone": "America/New_York",
  "profile_image_url": "https://example.com/avatar.jpg"
}
```

#### Login
```http
POST /users/login
Content-Type: application/json

{
  "phone_number": "+1234567890",
  "password": "securepassword123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "phone_number": "+1234567890",
    "time_zone": "America/New_York",
    "profile_image_url": "https://example.com/avatar.jpg",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### User Endpoints (Requires Authentication)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <token>
```

#### Get user details
```http
GET /users/{id}
```

#### Update user
```http
PUT /users/{id}
Content-Type: application/json

{
  "name": "Jane Doe",
  "time_zone": "America/Los_Angeles"
}
```

#### Delete user
```http
DELETE /users/{id}
```

### Habit Endpoints (Requires Authentication)

#### Create a habit
```http
POST /habits
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio"
}
```

#### Get all user habits
```http
GET /habits
Authorization: Bearer <token>
```

#### Get a specific habit
```http
GET /habits/{id}
Authorization: Bearer <token>
```

#### Update a habit
```http
PUT /habits/{id}
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Evening Exercise",
  "description": "45 minutes of yoga"
}
```

#### Delete a habit
```http
DELETE /habits/{id}
Authorization: Bearer <token>
```

### Log Endpoints (Requires Authentication)

#### Create a log for a habit
```http
POST /habits/{habit_id}/logs
Content-Type: application/json
Authorization: Bearer <token>

{
  "notes": "Completed 30 minutes of running"
}
```

#### Get all logs for a habit
```http
GET /habits/{habit_id}/logs
Authorization: Bearer <token>
```

#### Get a specific log
```http
GET /logs/{id}
Authorization: Bearer <token>
```

#### Update a log
```http
PUT /logs/{id}
Content-Type: application/json
Authorization: Bearer <token>

{
  "notes": "Updated: Completed 45 minutes of running"
}
```

#### Delete a log
```http
DELETE /logs/{id}
Authorization: Bearer <token>
```

### Health Check

```http
GET /health

Response: OK
```

## Database Schema

### Users Table
- `id` - INTEGER PRIMARY KEY AUTOINCREMENT
- `profile_image_url` - TEXT
- `name` - TEXT NOT NULL
- `time_zone` - TEXT NOT NULL
- `phone_number` - TEXT NOT NULL (indexed)
- `password_hash` - TEXT NOT NULL
- `created_at` - DATETIME DEFAULT CURRENT_TIMESTAMP

### Habits Table
- `id` - INTEGER PRIMARY KEY AUTOINCREMENT
- `user_id` - INTEGER NOT NULL (foreign key to users)
- `name` - TEXT NOT NULL
- `description` - TEXT
- `created_at` - DATETIME DEFAULT CURRENT_TIMESTAMP

### Logs Table
- `id` - INTEGER PRIMARY KEY AUTOINCREMENT
- `habit_id` - INTEGER NOT NULL (foreign key to habits)
- `notes` - TEXT
- `created_at` - DATETIME DEFAULT CURRENT_TIMESTAMP

## Security Features

- **Password Hashing** - Argon2id algorithm for secure password storage
- **JWT Authentication** - Token-based authentication with expiration
- **RBAC** - Users can only access their own resources
- **Input Validation** - All requests are validated before processing
- **Security Headers** - XSS protection, CSP, clickjacking prevention
- **Error Logging** - Comprehensive error logging for debugging

## Testing

To test the API, you can use tools like:
- `curl`
- Postman
- HTTPie

Example with curl:

```bash
# Register a user
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "phone_number": "+1234567890",
    "password": "securepassword123",
    "time_zone": "America/New_York"
  }'

# Login
TOKEN=$(curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890",
    "password": "securepassword123"
  }' | jq -r '.token')

# Create a habit
curl -X POST http://localhost:8080/habits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Morning Exercise",
    "description": "30 minutes of cardio"
  }'
```

## License

MIT

---

# Habit Notifier Kubernetes CronJob

This repository also includes a production-ready Kubernetes CronJob that:

- Queries a MySQL database for each user's habit logs over the past 2 days (per-user timezone).
- Sends an SMS via Twilio if either:
  - No logs in the last two days, or
  - A log exists on day 1 but no log on day 2.

If both days have logs, no notification is sent.

## App Source

- Python app under `app/` with modules: `config.py`, `db.py`, `notifier.py`, `logic.py`, `health.py`, `main.py`.
- Container spec: `Dockerfile`
- Manifests: `k8s/namespace.yaml`, `k8s/serviceaccount.yaml`, `k8s/configmap.yaml`, `k8s/secret.yaml` (placeholders), `k8s/cronjob.yaml`, `k8s/mysql-ca-secret.example.yaml`.

## Configuration (env)

- `LOG_LEVEL` (default INFO)
- `MYSQL_HOST` (required)
- `MYSQL_PORT` (default 3306)
- `MYSQL_DB` (required)
- `MYSQL_USER` (required, Secret)
- `MYSQL_PASSWORD` (required, Secret)
- `TWILIO_ACCOUNT_SID` (required, Secret)
- `TWILIO_AUTH_TOKEN` (required, Secret)
- `TWILIO_FROM` (required; E.164 sender number)
- `MYSQL_SSL_DISABLED` (default false). Set to `true` to disable TLS (not recommended).
- `MYSQL_SSL_CA_PATH` (optional path, e.g. `/etc/mysql/certs/ca.crt`) to verify server cert.

Note: Database tables are assumed to be `User`, `Habit`, and `Log` with fields described in the task. Time zones are read from `User.time_zone` (IANA tz).

## Build and Push the Image

```bash
docker build -t your-registry/habit-notifier:latest .
docker push your-registry/habit-notifier:latest
```

Update `k8s/cronjob.yaml` image to your pushed image.

## Deploy

```bash
kubectl apply -f k8s/namespace.yaml
kubectl -n habit-cron apply -f k8s/serviceaccount.yaml

# Create runtime secrets securely (preferred over committing secret.yaml)
kubectl -n habit-cron create secret generic habit-cron-secrets \
  --from-literal=MYSQL_USER=youruser \
  --from-literal=MYSQL_PASSWORD=yourpass \
  --from-literal=TWILIO_ACCOUNT_SID=ACxxxxxxxx \
  --from-literal=TWILIO_AUTH_TOKEN=xxxxxxxx

# Optional: add a MySQL CA bundle if your server uses a custom CA
# Provide a PEM at ./mysql-ca.crt and create the secret:
kubectl -n habit-cron create secret generic mysql-ca \
  --from-file=ca.crt=./mysql-ca.crt

# Edit k8s/configmap.yaml to set MYSQL_* and TWILIO_FROM values as needed
kubectl -n habit-cron apply -f k8s/configmap.yaml

# Deploy the CronJob
kubectl -n habit-cron apply -f k8s/cronjob.yaml
```

### Schedule

- The manifest schedules the job at `0 14 * * *` (14:00 UTC). Adjust as needed.
- The app evaluates the previous two local dates per user using their stored time zone.

## Verify

```bash
kubectl -n habit-cron get cronjobs
kubectl -n habit-cron get jobs
kubectl -n habit-cron logs job/<job-name>
```

Exit code is non-zero if any per-user processing error occurred.

## Security

- Runs as non-root with seccomp `RuntimeDefault`, read-only root FS, and dropped Linux capabilities.
- Secrets are mounted via environment variables. Do not commit real credentials to Git.
- MySQL TLS supported with CA bundle mounted at `/etc/mysql/certs/ca.crt` and `MYSQL_SSL_CA_PATH` set accordingly. Prefer TLS over disabling it.
- ServiceAccount has no RBAC permissions beyond default; the job does not call the Kubernetes API.

## Uninstall

```bash
kubectl -n habit-cron delete -f k8s/cronjob.yaml
kubectl -n habit-cron delete secret habit-cron-secrets mysql-ca --ignore-not-found
kubectl -n habit-cron delete -f k8s/configmap.yaml
kubectl -n habit-cron delete -f k8s/serviceaccount.yaml
kubectl delete -f k8s/namespace.yaml
```
