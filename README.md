# Habit Tracker API

A production-ready REST API built with Go for tracking personal habits. Features JWT authentication, SQLite database for local development, and complete REST endpoints for Users, Habits, and Logs.

## Features

- ✅ **RESTful API** with full CRUD operations for Users, Habits, and Logs
- 🔐 **JWT Authentication** with secure token-based access control
- 🏗️ **Clean Architecture** with separation of concerns (handlers → services → repositories)
- 🔒 **Security Best Practices** including RBAC, input validation, secure headers, and HTTPS
- 📊 **SQLite Database** with automatic migrations, proper indexing and foreign key constraints
- 🚀 **Easy Local Development** - No external database required, automatic setup on first run
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
- SQLite3 (included with most operating systems)

### Setup and Run

1. **Clone the repository**
   ```bash
   git clone https://github.com/hayden-erickson/ai-evaluation.git
   cd ai-evaluation
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

The server will:
- Create a SQLite database file `habits.db` in the current directory
- Automatically run all migrations
- Start listening on port 8080 (or the port specified in `PORT` environment variable)

### Configuration

The application can be configured using environment variables:

- `PORT` - Server port (default: 8080)
- `DB_PATH` - SQLite database file path (default: habits.db)
- `JWT_SECRET` - Secret key for JWT token generation (default: default-secret-change-in-production)

Example:
```bash
export PORT=3000
export DB_PATH=my_habits.db
export JWT_SECRET=your-secret-key-here
go run main.go
```

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 with full REST API endpoints for Users, Habits, and Logs.
