## New API - Habits Tracker (Go + SQLite)

### Overview
This is a minimal REST API implementing Users, Habits, and Logs with:
- SQLite for local development
- JWT authentication
- Input validation and consistent error responses
- Security headers and simple request logging

Endpoints are implemented in `main.go` and database migrations are in `migrations/`.

### Prerequisites
- Go 1.21+

### Run
Set a JWT secret (for dev you can skip it; a default will be used):

```bash
$env:JWT_SECRET="your-dev-secret"  # PowerShell
# export JWT_SECRET=your-dev-secret  # bash/zsh

go run main.go
```

The server starts on port 8080. Health check: `GET /healthz`.

### Database
SQLite file is created at `./new-api.db`. Migrations run automatically from `./migrations`.

### Auth Flow
1) Register a user and receive a JWT:

```bash
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"profileImageUrl":"","name":"Alice","timeZone":"UTC","phoneNumber":"+10000000000"}'
```

Copy `token` from the response. For subsequent requests, include:
`Authorization: Bearer <token>`

### Habits
- GET `/habits` — list your habits
- POST `/habits` — create habit
  - body: `{ "name": "Water", "description": "Drink water" }`
- GET `/habits/{id}` — get habit
- PUT `/habits/{id}` — update habit (any of `name`, `description`)
- DELETE `/habits/{id}` — delete habit

### Logs
- GET `/logs` — list logs across your habits
- GET `/logs?habitId={id}` — list logs for a habit
- POST `/logs` — create log
  - body: `{ "habitId": 1, "notes": "Did it" }`
- GET `/logs/{id}` — get log
- PUT `/logs/{id}` — update log (`notes`)
- DELETE `/logs/{id}` — delete log

Include header: `Authorization: Bearer <token>` for all except `/register` and `/healthz`.
