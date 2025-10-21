# Habit Tracker API

A RESTful API built with Go for tracking user habits, featuring JWT authentication, SQLite database, and clean architecture.

## Features

- **JWT-based Authentication**: Secure token-based authentication
- **User Management**: Create, read, update, and delete user accounts
- **Habit Tracking**: Manage habits with full CRUD operations
- **Activity Logging**: Track habit completion with notes
- **Role-Based Access Control (RBAC)**: Users can only access their own data
- **Input Validation**: Comprehensive request validation
- **Error Logging**: All errors are logged with appropriate context
- **Security Headers**: Implements secure HTTP headers
- **Modular Architecture**: Clean separation of concerns with repository and service layers
- **Dependency Injection**: Loosely coupled components

## Tech Stack

- **Go 1.21**: Programming language
- **SQLite**: Database for local development
- **Standard Library**: Uses `net/http`, `context`, and `log`
- **JWT**: golang-jwt/jwt for authentication
- **bcrypt**: Secure password hashing

## Project Structure

```
.
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── migrations/                  # SQL migration files
│   ├── 001_create_users_table.sql
│   ├── 002_create_habits_table.sql
│   └── 003_create_logs_table.sql
├── internal/
│   ├── database/               # Database initialization and migrations
│   ├── models/                 # Data models and DTOs
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic layer
│   ├── handlers/               # HTTP handlers
│   └── middleware/             # Authentication and middleware
└── data/                       # SQLite database (created at runtime)
```

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd new-api
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Health Check
- **GET** `/health` - Check API health status

#### User Registration
- **POST** `/users` - Create a new user account
  ```json
  {
    "name": "John Doe",
    "time_zone": "America/New_York",
    "phone_number": "1234567890",
    "password": "securepassword123",
    "profile_image_url": "https://example.com/image.jpg"
  }
  ```

#### Login
- **POST** `/login` - Authenticate and receive JWT token
  ```json
  {
    "phone_number": "1234567890",
    "password": "securepassword123"
  }
  ```

### Protected Endpoints (Requires Authentication)

All protected endpoints require the `Authorization` header with a Bearer token:
```
Authorization: Bearer <your-jwt-token>
```

#### Users
- **GET** `/users/{id}` - Get user details
- **PUT** `/users/{id}` - Update user information
- **DELETE** `/users/{id}` - Delete user account

#### Habits
- **POST** `/habits` - Create a new habit
  ```json
  {
    "name": "Daily Exercise",
    "description": "30 minutes of exercise"
  }
  ```
- **GET** `/habits` - Get all habits for authenticated user
- **GET** `/habits/{id}` - Get specific habit details
- **PUT** `/habits/{id}` - Update a habit
- **DELETE** `/habits/{id}` - Delete a habit

#### Logs
- **POST** `/habits/{habit_id}/logs` - Create a new log entry
  ```json
  {
    "notes": "Completed morning workout"
  }
  ```
- **GET** `/habits/{habit_id}/logs` - Get all logs for a habit
- **GET** `/logs/{id}` - Get specific log details
- **PUT** `/logs/{id}` - Update a log entry
- **DELETE** `/logs/{id}` - Delete a log entry

## Example Usage

### 1. Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "time_zone": "America/New_York",
    "phone_number": "1234567890",
    "password": "securepass123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "1234567890",
    "password": "securepass123"
  }'
```

### 3. Create a Habit (with token)
```bash
curl -X POST http://localhost:8080/habits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "name": "Morning Meditation",
    "description": "10 minutes of meditation"
  }'
```

### 4. Log a Habit Completion
```bash
curl -X POST http://localhost:8080/habits/{habit-id}/logs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "notes": "Felt very relaxed today"
  }'
```

## Security Features

- **Password Hashing**: All passwords are hashed using bcrypt
- **JWT Authentication**: Stateless authentication with 24-hour token expiration
- **RBAC**: Users can only access their own data
- **Input Validation**: All inputs are validated before processing
- **Secure Headers**: Implements security headers (X-Frame-Options, CSP, etc.)
- **Foreign Key Constraints**: Ensures data integrity at database level

## Database Schema

### Users Table
- `id` (TEXT, PRIMARY KEY)
- `profile_image_url` (TEXT)
- `name` (TEXT, NOT NULL)
- `time_zone` (TEXT, NOT NULL)
- `phone_number` (TEXT)
- `password_hash` (TEXT, NOT NULL)
- `created_at` (DATETIME, NOT NULL)

### Habits Table
- `id` (TEXT, PRIMARY KEY)
- `user_id` (TEXT, FOREIGN KEY → users.id)
- `name` (TEXT, NOT NULL)
- `description` (TEXT)
- `created_at` (DATETIME, NOT NULL)

### Logs Table
- `id` (TEXT, PRIMARY KEY)
- `habit_id` (TEXT, FOREIGN KEY → habits.id)
- `notes` (TEXT)
- `created_at` (DATETIME, NOT NULL)

## Error Handling

All errors return appropriate HTTP status codes with JSON error messages:

```json
{
  "error": "Bad Request",
  "message": "name is required"
}
```

Common status codes:
- `200` OK
- `201` Created
- `204` No Content
- `400` Bad Request
- `401` Unauthorized
- `403` Forbidden
- `404` Not Found
- `405` Method Not Allowed
- `500` Internal Server Error

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o habit-tracker main.go
./habit-tracker
```

## Configuration

The application uses the following default configuration:
- **Server Port**: 8080
- **Database Path**: `./data/habits.db`
- **JWT Secret**: Defined in `internal/middleware/auth.go` (change in production)
- **Token Expiration**: 24 hours

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
