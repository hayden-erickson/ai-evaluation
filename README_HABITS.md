# Habit Tracker REST API

A robust, production-ready REST API for tracking user habits with JWT authentication, built using Go and SQLite.

## Features

- **JWT-based Authentication**: Secure token-based authentication
- **User Management**: Register, login, update, and delete user accounts
- **Habit Tracking**: Create, read, update, and delete habits
- **Activity Logging**: Track habit completion with timestamped logs
- **Role-Based Access Control**: Users can only access their own data
- **Input Validation**: Comprehensive validation of all user inputs
- **Security Headers**: Protection against common web vulnerabilities
- **Error Logging**: All errors are logged for debugging
- **Modular Architecture**: Clear separation of concerns with repository and service layers
- **Dependency Injection**: Clean, testable code structure

## Architecture

The application follows a layered architecture:

- **Models**: Data structures and DTOs
- **Repository**: Data access layer (database operations)
- **Service**: Business logic layer
- **Handlers**: HTTP request handlers
- **Middleware**: Authentication, logging, security headers
- **Utils**: Validation, cryptography, JWT utilities

## Prerequisites

- Go 1.21 or higher
- SQLite3 (included via go-sqlite3 driver)

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

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## Configuration

The application can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path (default: ./habits.db)
- `JWT_SECRET`: Secret key for JWT signing (auto-generated if not provided)

Example:
```bash
export PORT=3000
export DB_PATH=/path/to/database.db
export JWT_SECRET=your-secret-key-here
go run main.go
```

## Database

The application uses SQLite for local development. The database file is created automatically on first run, and migrations are executed automatically.

### Database Schema

**Users Table**:
- `id`: Integer (Primary Key)
- `profile_image_url`: Text
- `name`: Text (Required)
- `time_zone`: Text
- `phone_number`: Text (Unique, Required)
- `password_hash`: Text (Required)
- `created_at`: DateTime

**Habits Table**:
- `id`: Integer (Primary Key)
- `user_id`: Integer (Foreign Key → users.id)
- `name`: Text (Required)
- `description`: Text
- `created_at`: DateTime

**Logs Table**:
- `id`: Integer (Primary Key)
- `habit_id`: Integer (Foreign Key → habits.id)
- `notes`: Text
- `created_at`: DateTime

## API Endpoints

### Authentication

#### Register a New User
```http
POST /api/register
Content-Type: application/json

{
  "name": "John Doe",
  "phone_number": "+1234567890",
  "password": "securepassword123",
  "time_zone": "America/New_York",
  "profile_image_url": "https://example.com/photo.jpg"
}
```

Response:
```json
{
  "id": 1,
  "name": "John Doe",
  "phone_number": "+1234567890",
  "time_zone": "America/New_York",
  "profile_image_url": "https://example.com/photo.jpg",
  "created_at": "2024-01-01T10:00:00Z"
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "phone_number": "+1234567890",
  "password": "securepassword123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "phone_number": "+1234567890",
    "time_zone": "America/New_York",
    "profile_image_url": "https://example.com/photo.jpg",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### Users (Protected)

All user endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

#### Get User
```http
GET /api/users/{id}
Authorization: Bearer <token>
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
```

#### Delete User
```http
DELETE /api/users/{id}
Authorization: Bearer <token>
```

### Habits (Protected)

#### Create Habit
```http
POST /api/habits
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio"
}
```

#### List User's Habits
```http
GET /api/habits
Authorization: Bearer <token>
```

#### Get Specific Habit
```http
GET /api/habits/{id}
Authorization: Bearer <token>
```

#### Update Habit
```http
PUT /api/habits/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Evening Exercise",
  "description": "45 minutes of cardio and weights"
}
```

#### Delete Habit
```http
DELETE /api/habits/{id}
Authorization: Bearer <token>
```

### Logs (Protected)

#### Create Log Entry
```http
POST /api/logs
Authorization: Bearer <token>
Content-Type: application/json

{
  "habit_id": 1,
  "notes": "Completed 5km run"
}
```

#### List Logs for a Habit
```http
GET /api/logs?habit_id=1
Authorization: Bearer <token>
```

#### Get Specific Log
```http
GET /api/logs/{id}
Authorization: Bearer <token>
```

#### Update Log
```http
PUT /api/logs/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "notes": "Completed 5km run in 30 minutes"
}
```

#### Delete Log
```http
DELETE /api/logs/{id}
Authorization: Bearer <token>
```

### Health Check

```http
GET /health
```

Response:
```
OK
```

## Error Handling

All errors are returned with appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid input or validation errors
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User doesn't have permission to access the resource
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate user)
- `500 Internal Server Error`: Server-side errors

Example error response:
```json
{
  "error": "Invalid phone number format"
}
```

## Security

The application implements multiple security best practices:

- **Password Hashing**: Bcrypt for secure password storage
- **JWT Authentication**: Secure token-based authentication with 24-hour expiration
- **Input Validation**: All inputs are validated and sanitized
- **RBAC**: Users can only access their own data
- **Security Headers**: X-Frame-Options, X-Content-Type-Options, CSP, etc.
- **HTTPS Enforcement**: Strict-Transport-Security header (configure reverse proxy for production)
- **CORS**: Configurable CORS policy
- **SQL Injection Protection**: Parameterized queries

## Testing

To test the API, you can use curl, Postman, or any HTTP client.

Example workflow:

1. Register a user:
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","phone_number":"+1234567890","password":"password123"}'
```

2. Login:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"+1234567890","password":"password123"}'
```

3. Create a habit (use token from login response):
```bash
curl -X POST http://localhost:8080/api/habits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{"name":"Exercise","description":"Daily workout"}'
```

## Development

### Project Structure

```
.
├── main.go              # Application entry point
├── models/              # Data models
│   ├── user.go
│   ├── habit.go
│   └── log.go
├── repository/          # Data access layer
│   ├── user_repository.go
│   ├── habit_repository.go
│   └── log_repository.go
├── service/             # Business logic layer
│   ├── user_service.go
│   ├── habit_service.go
│   └── log_service.go
├── handlers/            # HTTP handlers
│   ├── user_handler.go
│   ├── habit_handler.go
│   └── log_handler.go
├── middleware/          # HTTP middleware
│   ├── auth.go
│   ├── logging.go
│   └── security.go
├── utils/               # Utility functions
│   ├── validation.go
│   ├── crypto.go
│   └── jwt.go
├── migrations/          # Database migrations
│   ├── 001_create_users_table.sql
│   ├── 002_create_habits_table.sql
│   └── 003_create_logs_table.sql
└── README.md
```

### Adding New Features

1. Add model in `models/` package
2. Create repository interface and implementation in `repository/` package
3. Create service interface and implementation in `service/` package
4. Create HTTP handler in `handlers/` package
5. Register routes in `main.go`

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
