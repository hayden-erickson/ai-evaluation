# Habit Tracker App

A full-stack habit tracking application with a Go REST API backend and React Native mobile frontend.

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

### Frontend (React Native)
- **Habit tracking** - Create, edit, and delete habits
- **Streak tracking** - Track habit streaks with 1-day skip allowance
- **Log entries** - Add notes and details for each habit completion
- **Modern UI** - Minimal design with pastel colors and rounded corners
- **Authentication** - Secure login and registration
- **Error handling** - Graceful error handling with user-friendly messages

## Architecture

### Backend Package Structure

- **`/models`** - Data structures and validation logic (User, Habit, Log)
- **`/repository`** - Database operations and data access layer
- **`/service`** - Business logic and service layer
- **`/handlers`** - HTTP request handlers
- **`/middleware`** - Authentication, logging, and security middleware
- **`/config`** - Database configuration and migrations
- **`/utils`** - Utility functions (JWT, password hashing)
- **`/migrations`** - SQL migration files

### Frontend Structure

- **`/src/components`** - Reusable UI components
- **`/src/screens`** - Screen components (Login, Register, Home)
- **`/src/services`** - API client and authentication service
- **`/src/utils`** - Utility functions and helpers
- **`/src/types`** - TypeScript type definitions

## Prerequisites

### Backend
- Go 1.21 or higher
- SQLite3

### Frontend
- Node.js 20 or higher
- React Native development environment (see [React Native Environment Setup](https://reactnative.dev/docs/set-up-your-environment))
- For iOS: macOS with Xcode and CocoaPods
- For Android: Android Studio and Android SDK

## Installation

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

2. Install Go dependencies:
```bash
go mod download
```

3. (Optional) Set environment variables:
```bash
export PORT=8080                    # Default: 8080
export JWT_SECRET=your-secret-key   # Default: "default-secret-key-change-in-production"
export DB_PATH=habits.db            # Default: habits.db
```

4. Run the backend server:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

### Frontend Setup

1. Install Node.js dependencies:
```bash
npm install
```

2. For iOS, install CocoaPods dependencies:
```bash
cd ios
bundle install
bundle exec pod install
cd ..
```

3. Configure the API URL:
   - The app is configured to connect to `http://localhost:8080` by default
   - To connect to a different server, update the `API_URL` in `src/services/api.ts`

4. Run the React Native app:

   **For iOS:**
   ```bash
   npm run ios
   ```

   **For Android:**
   ```bash
   npm run android
   ```

   **Start Metro bundler separately:**
   ```bash
   npm start
   ```

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

### Protected Endpoints (Requires Authentication)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <token>
```

#### Habits
- `GET /habits` - Get all user habits
- `POST /habits` - Create a new habit
- `GET /habits/{id}` - Get a specific habit
- `PUT /habits/{id}` - Update a habit
- `DELETE /habits/{id}` - Delete a habit

#### Logs
- `GET /habits/{habit_id}/logs` - Get all logs for a habit
- `POST /habits/{habit_id}/logs` - Create a log for a habit
- `GET /logs/{id}` - Get a specific log
- `PUT /logs/{id}` - Update a log
- `DELETE /logs/{id}` - Delete a log

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

### Backend
```bash
go test ./...
```

### Frontend
```bash
npm test
```

## License

MIT
