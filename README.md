# Habit Tracking REST API

A REST API built with Go for tracking user habits with SQLite database, JWT authentication, and comprehensive CRUD operations.

## Features

- **JWT-based Authentication** - Secure user authentication with token-based access
- **User Management** - Full CRUD operations for user accounts
- **Habit Tracking** - Create, read, update, and delete habits
- **Log Entries** - Track habit completion with timestamped logs and notes
- **Role-Based Access Control** - Users can only access their own data
- **Input Validation** - Comprehensive request validation
- **Error Handling** - All errors logged and returned with appropriate HTTP status codes
- **Security Headers** - XSS protection, HSTS, CSP, and more
- **Modular Architecture** - Clean separation of concerns with repository and service layers
- **Database Migrations** - Automatic schema creation on startup

## Architecture

The application follows a layered architecture with dependency injection:

```
┌─────────────────────────────────────────────┐
│              HTTP Handlers                   │
│  (User, Habit, Log request handling)        │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│            Service Layer                     │
│  (Business logic & validation)              │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│          Repository Layer                    │
│  (Database operations)                      │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│            SQLite Database                   │
│  (Users, Habits, Logs tables)               │
└─────────────────────────────────────────────┘
```

### Package Structure

- **`/models`** - Data models and DTOs (User, Habit, Log, requests/responses)
- **`/repository`** - Database operations and data access layer
- **`/service`** - Business logic and validation
- **`/handlers`** - HTTP request handlers
- **`/middleware`** - JWT authentication and security middleware
- **`/config`** - Application configuration
- **`/migrations`** - SQL schema migration files

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified by the `PORT` environment variable).

### Configuration

The application can be configured using environment variables:

- `DB_PATH` - Path to SQLite database file (default: `./habits.db`)
- `JWT_SECRET` - Secret key for JWT token signing (default: `your-secret-key-change-in-production`)
- `PORT` - Server port (default: `8080`)

Example:
```bash
export DB_PATH=/data/habits.db
export JWT_SECRET=my-super-secret-key
export PORT=3000
go run main.go
```

## API Endpoints

### Authentication

#### Register User
```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "password": "password123",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg"
}

Response: 201 Created
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "time_zone": "America/New_York",
    "phone_number": "+1234567890",
    "profile_image_url": "https://example.com/image.jpg",
    "created_at": "2025-10-21T00:00:00Z"
  }
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "name": "John Doe",
  "password": "password123"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { ... }
}
```

### Users (Protected)

All user endpoints require the `Authorization: Bearer <token>` header.

#### Get User
```http
GET /api/users/{id}
Authorization: Bearer <token>

Response: 200 OK
{
  "id": 1,
  "name": "John Doe",
  ...
}
```

#### Update User
```http
PUT /api/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe",
  "time_zone": "America/Los_Angeles"
}

Response: 200 OK
```

#### Delete User
```http
DELETE /api/users/{id}
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "user deleted successfully"
}
```

### Habits (Protected)

#### List Habits
```http
GET /api/habits
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Morning Exercise",
    "description": "30 minutes of cardio",
    "created_at": "2025-10-21T00:00:00Z"
  }
]
```

#### Create Habit
```http
POST /api/habits
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio"
}

Response: 201 Created
```

#### Get Habit
```http
GET /api/habits/{id}
Authorization: Bearer <token>

Response: 200 OK
```

#### Update Habit
```http
PUT /api/habits/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Evening Exercise",
  "description": "45 minutes of yoga"
}

Response: 200 OK
```

#### Delete Habit
```http
DELETE /api/habits/{id}
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "habit deleted successfully"
}
```

### Logs (Protected)

#### List Logs for Habit
```http
GET /api/habits/{habit_id}/logs
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "habit_id": 1,
    "notes": "Completed 5k run",
    "created_at": "2025-10-21T08:00:00Z"
  }
]
```

#### Create Log
```http
POST /api/logs
Authorization: Bearer <token>
Content-Type: application/json

{
  "habit_id": 1,
  "notes": "Completed 5k run"
}

Response: 201 Created
```

#### Get Log
```http
GET /api/logs/{id}
Authorization: Bearer <token>

Response: 200 OK
```

#### Update Log
```http
PUT /api/logs/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "notes": "Completed 10k run instead"
}

Response: 200 OK
```

#### Delete Log
```http
DELETE /api/logs/{id}
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "log deleted successfully"
}
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_image_url TEXT,
    name TEXT NOT NULL,
    time_zone TEXT NOT NULL,
    phone_number TEXT,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Habits Table
```sql
CREATE TABLE habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Logs Table
```sql
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);
```

## Security Features

- **Password Hashing** - Passwords are hashed using bcrypt
- **JWT Authentication** - Stateless token-based authentication
- **RBAC** - Users can only access their own resources
- **Input Validation** - All requests are validated before processing
- **Security Headers** - XSS Protection, HSTS, CSP, X-Frame-Options
- **SQL Injection Protection** - Parameterized queries prevent SQL injection
- **Error Sanitization** - Sensitive errors are logged but not exposed to clients

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o habit-api main.go
./habit-api
```

### Code Structure Guidelines

The application follows these principles:
- **Single Responsibility** - Each package/struct has one clear purpose
- **Dependency Injection** - Dependencies are injected, not created
- **Interface Segregation** - Services use interfaces for testability
- **Error Handling** - All errors are logged and returned appropriately
- **Idiomatic Go** - Uses standard library where possible

## License

This project is licensed under the MIT License.
