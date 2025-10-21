# Habit Tracker REST API

A RESTful API built with Go and SQLite for tracking user habits with JWT-based authentication.

## Features

- **User Management**: Create, read, update, and delete users
- **Habit Tracking**: Manage habits with user ownership
- **Logging**: Track habit completion with notes
- **JWT Authentication**: Secure API endpoints with token-based auth
- **Input Validation**: Comprehensive request validation
- **Error Logging**: All errors are logged with appropriate HTTP status codes
- **Security**: RBAC, secure headers, input validation
- **Modular Architecture**: Separation of concerns with repository and service layers
- **Dependency Injection**: Clean, testable code structure

## Project Structure

```
.
├── main.go                 # Application entry point
├── migrations/             # SQL migration files
│   ├── 001_create_users_table.sql
│   ├── 002_create_habits_table.sql
│   └── 003_create_logs_table.sql
├── internal/
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic layer
│   └── handler/           # HTTP handlers and middleware
└── pkg/
    ├── auth/              # Authentication utilities
    ├── database/          # Database initialization
    └── validator/         # Input validation

```

## Prerequisites

- Go 1.21 or higher
- GCC (for SQLite CGO compilation)

## Installation

1. Clone the repository
2. Install dependencies:

```bash
go mod download
```

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Create User
```
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg",
  "password": "securepassword"
}
```

#### Login
```
POST /api/login
Content-Type: application/json

{
  "id": "user-id-here",
  "password": "securepassword"
}

Response:
{
  "token": "jwt-token-here",
  "user": { ... }
}
```

#### Get User
```
GET /api/users/{id}
```

#### Update User
```
PUT /api/users/{id}
Content-Type: application/json

{
  "name": "Jane Doe",
  "time_zone": "America/Los_Angeles"
}
```

#### Delete User
```
DELETE /api/users/{id}
```

#### List Users
```
GET /api/users
```

### Protected Endpoints (Authentication Required)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer {jwt-token}
```

#### Create Habit
```
POST /api/habits
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio"
}
```

#### Get Habit
```
GET /api/habits/{id}
Authorization: Bearer {token}
```

#### Update Habit
```
PUT /api/habits/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Evening Exercise",
  "description": "45 minutes of strength training"
}
```

#### Delete Habit
```
DELETE /api/habits/{id}
Authorization: Bearer {token}
```

#### List User's Habits
```
GET /api/habits
Authorization: Bearer {token}
```

#### Create Log
```
POST /api/logs
Content-Type: application/json
Authorization: Bearer {token}

{
  "habit_id": "habit-id-here",
  "notes": "Felt great today!"
}
```

#### Get Log
```
GET /api/logs/{id}
Authorization: Bearer {token}
```

#### Update Log
```
PUT /api/logs/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "notes": "Updated notes"
}
```

#### Delete Log
```
DELETE /api/logs/{id}
Authorization: Bearer {token}
```

#### List Logs by Habit
```
GET /api/logs?habit_id={habit-id}
Authorization: Bearer {token}
```

## Security Features

- **JWT Authentication**: Token-based authentication for protected endpoints
- **Password Hashing**: Bcrypt for secure password storage
- **Input Validation**: All requests are validated before processing
- **RBAC**: Users can only access their own habits and logs
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, HSTS, CSP
- **Error Logging**: All errors are logged for monitoring

## Database

The application uses SQLite for local development. The database file is created at `./data/habits.db`.

### Schema

- **users**: User accounts with authentication
- **habits**: User habits with foreign key to users
- **logs**: Habit completion logs with foreign key to habits

Foreign key constraints ensure referential integrity with cascade deletes.

## Development

### Adding New Endpoints

1. Define models in `internal/models/`
2. Create repository interface and implementation in `internal/repository/`
3. Create service interface and implementation in `internal/service/`
4. Create HTTP handler in `internal/handler/`
5. Register routes in `main.go`

### Running Tests

```bash
go test ./...
```

## Error Handling

All errors return appropriate HTTP status codes:

- `400 Bad Request`: Invalid input or validation errors
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: User doesn't have permission to access resource
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side errors

Error responses follow this format:
```json
{
  "error": "Error message here"
}
```

## License

MIT
