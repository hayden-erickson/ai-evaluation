# Habit Tracker with React UI

A full-stack habit tracking application with a Go REST API backend and modern React frontend. Track your habits with streak functionality, allowing you to maintain momentum while being forgiving of occasional skips.

## Features

### Backend (Go API)
- **JWT-based authentication** - Secure token-based authentication
- **Input validation** - Comprehensive request validation for all endpoints
- **Error logging** - All errors are logged with appropriate context
- **HTTP status codes** - Proper error responses with meaningful messages
- **Modular architecture** - Separation of concerns with repository and service layers
- **Dependency injection** - Clean, testable code structure
- **Security headers** - XSS protection, clickjacking prevention, CSP
- **RBAC** - Role-based access control (users can only access their own data)
- **SQLite database** - Lightweight database for local development

### Frontend (React UI)
- **Modern UI** - Clean, minimal design with pastel colors and rounded corners
- **Streak tracking** - Visual representation of habit consistency
- **Smart streak logic** - Allows one day skip before breaking a streak
- **Daily logs** - Add notes to track your progress
- **Responsive design** - Works on desktop and mobile
- **Error handling** - Graceful error messages and validation

## Quick Start

### Running Locally

1. **Clone the repository:**
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

2. **Install and build the frontend:**
```bash
cd frontend
npm install
npm run build
cd ..
```

3. **Run the server:**
```bash
go run main.go
```

4. **Access the application:**
Open your browser and navigate to [http://localhost:8080](http://localhost:8080)

The server serves both the API endpoints and the React frontend.

## Architecture

The application follows a clean, modular architecture:

### Package Structure

**Backend:**
- **`/models`** - Data structures and validation logic (User, Habit, Log)
- **`/repository`** - Database operations and data access layer
- **`/service`** - Business logic and service layer
- **`/handlers`** - HTTP request handlers
- **`/middleware`** - Authentication, logging, and security middleware
- **`/config`** - Database configuration and migrations
- **`/utils`** - Utility functions (JWT, password hashing)
- **`/migrations`** - SQL migration files

**Frontend:**
- **`/frontend`** - React application
  - **`/src/components`** - React components (Dashboard, Login, HabitList, etc.)
  - **`/src/contexts`** - React contexts (AuthContext)
  - **`/src/services`** - API client
  - **`/src/types`** - TypeScript type definitions
  - **`/src/utils`** - Utility functions (streak calculation)
  - **`/src/styles`** - CSS stylesheets

## Prerequisites

- Go 1.21 or higher
- Node.js 14 or higher
- npm or yarn
- SQLite3

## Installation

### Full Stack Setup

1. Clone the repository:
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install frontend dependencies:
```bash
cd frontend
npm install
cd ..
```

4. (Optional) Set environment variables:
```bash
export PORT=8080                    # Default: 8080
export JWT_SECRET=your-secret-key   # Default: "default-secret-key-change-in-production"
export DB_PATH=habits.db            # Default: habits.db
```

## Running the Application

### Development Mode (Frontend with Hot Reload)

1. Start the Go backend:
```bash
go run main.go
```

2. In a separate terminal, start the React development server:
```bash
cd frontend
npm start
```

The backend runs on port 8080 and the frontend on port 3000. The frontend will proxy API requests to the backend.

### Production Mode (Integrated)

1. Build the frontend:
```bash
cd frontend
npm run build
cd ..
```

2. Start the server:
```bash
go run main.go
```

The server will serve both the API and the React app on port 8080 (or the port specified in the PORT environment variable).

Access the application at [http://localhost:8080](http://localhost:8080)

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
