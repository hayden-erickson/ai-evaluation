# Habit Tracker API

A production-ready REST API built with Go for tracking personal habits. Features JWT authentication, MySQL database, and complete AWS Kubernetes deployment setup.

## Features

- ✅ **RESTful API** with full CRUD operations for Users, Habits, and Logs
- 🔐 **JWT Authentication** with secure token-based access control
- 🏗️ **Clean Architecture** with separation of concerns (handlers → services → repositories)
- 🔒 **Security Best Practices** including RBAC, input validation, secure headers, and HTTPS
- 📊 **MySQL Database** with proper indexing and foreign key constraints
- 🐳 **Docker Support** for local development and production builds
- ☸️ **Kubernetes Deployment** with auto-scaling, health checks, and monitoring
- 📈 **AWS Ready** with EKS deployment scripts and CloudWatch integration
- 📝 **Comprehensive Logging** of all errors and requests
- ✨ **Idiomatic Go Code** using standard library (net/http, context, log)

## Architecture

### Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Authentication, logging, security
│   ├── models/                  # Data structures and DTOs
│   ├── repository/              # Database access layer
│   └── service/                 # Business logic layer
├── migrations/                  # SQL database migrations
├── k8s/                        # Kubernetes manifests
├── scripts/                    # Deployment scripts
├── docker-compose.yml          # Local development setup
└── Dockerfile                  # Production container image
```

### Layer Architecture

1. **Handlers Layer** - HTTP request/response handling, input validation
2. **Service Layer** - Business logic, authorization, data transformation
3. **Repository Layer** - Database operations, SQL queries
4. **Models Layer** - Data structures, request/response types

## Data Models

### User
- ID (UUID)
- Email (unique)
- Password (bcrypt hashed)
- Profile Image URL
- Name
- Time Zone
- Phone Number
- Created At

### Habit
- ID (UUID)
- User ID (foreign key)
- Name
- Description
- Created At

### Log
- ID (UUID)
- Habit ID (foreign key)
- Created At
- Notes

## API Endpoints

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { ... }
}
```

### Users (Protected)

#### Get User
```http
GET /api/v1/users/{id}
Authorization: Bearer {token}
```

#### Update User
```http
PUT /api/v1/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Jane Doe",
  "time_zone": "America/Los_Angeles"
}
```

#### Delete User
```http
DELETE /api/v1/users/{id}
Authorization: Bearer {token}
```

### Habits (Protected)

#### Create Habit
```http
POST /api/v1/habits
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio"
}
```

#### List User's Habits
```http
GET /api/v1/habits
Authorization: Bearer {token}
```

#### Get Habit
```http
GET /api/v1/habits/{id}
Authorization: Bearer {token}
```

#### Update Habit
```http
PUT /api/v1/habits/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Updated Habit Name"
}
```

#### Delete Habit
```http
DELETE /api/v1/habits/{id}
Authorization: Bearer {token}
```

### Logs (Protected)

#### Create Log Entry
```http
POST /api/v1/logs
Authorization: Bearer {token}
Content-Type: application/json

{
  "habit_id": "uuid-here",
  "notes": "Completed 45 minutes today"
}
```

#### List Logs for Habit
```http
GET /api/v1/logs?habit_id={habit_id}
Authorization: Bearer {token}
```

#### Get Log
```http
GET /api/v1/logs/{id}
Authorization: Bearer {token}
```

#### Update Log
```http
PUT /api/v1/logs/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "notes": "Updated notes"
}
```

#### Delete Log
```http
DELETE /api/v1/logs/{id}
Authorization: Bearer {token}
```

## Local Development

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- MySQL client (optional, for manual migrations)

### Setup and Run

1. **Clone the repository**
2. **Improved testability** - Packages can be tested in isolation
3. **Clear dependencies** - Import structure shows relationships between components
4. **Easier navigation** - Related code is grouped together

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 with the access code edit endpoint available at `/api/access-code/edit`.
