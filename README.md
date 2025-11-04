# Habit Tracker

A full-stack habit tracking application with a React Native mobile app and a RESTful API backend built in Go. Track your daily habits, build streaks, and maintain positive routines with a modern, minimal interface.

## Features

### Mobile App
- **Beautiful UI** - Modern, minimal design with pastel colors and rounded corners
- **Habit Tracking** - Create and manage multiple habits
- **Streak Counting** - Track consecutive days (allows 1-day skip before reset)
- **Daily Logging** - Log habit completions with optional notes
- **Recent Activity** - View the last 7 days of activity for each habit
- **User Authentication** - Secure login and registration

### Backend API
- **JWT-based authentication** - Secure token-based authentication
- **Input validation** - Comprehensive request validation for all endpoints
- **Error logging** - All errors are logged with appropriate context
- **HTTP status codes** - Proper error responses with meaningful messages
- **Modular architecture** - Separation of concerns with repository and service layers
- **Dependency injection** - Clean, testable code structure
- **Security headers** - XSS protection, clickjacking prevention, CSP
- **RBAC** - Role-based access control (users can only access their own data)
- **SQLite database** - Lightweight database for local development

## Prerequisites

### For Backend
- Go 1.21 or higher
- SQLite3

### For Mobile App
- Node.js 20 or higher
- React Native development environment ([Setup Guide](https://reactnative.dev/docs/set-up-your-environment))
- For iOS: Xcode and CocoaPods
- For Android: Android Studio and Android SDK

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

### Backend Setup

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

### Mobile App Setup

5. Install npm dependencies:
```bash
npm install
```

6. For iOS, install CocoaPods dependencies:
```bash
# First time only
bundle install

# Install pods
cd ios && bundle exec pod install && cd ..
```

7. Start the Metro bundler:
```bash
npm start
```

8. In a new terminal, run the app:

**For iOS:**
```bash
npm run ios
```

**For Android:**
```bash
npm run android
```

## Configuration

### Backend API URL

By default, the mobile app connects to `http://localhost:8080`. To connect to a different backend:

1. Open `src/services/api.ts`
2. Update the `API_BASE_URL` constant:
   ```typescript
   const API_BASE_URL = 'http://your-backend-url:8080';
   ```

For iOS simulator, use `http://localhost:8080`
For Android emulator, use `http://10.0.2.2:8080`
For physical devices, use your computer's IP address (e.g., `http://192.168.1.100:8080`)

## Usage

### Getting Started

1. Start the backend server (see Backend Setup)
2. Launch the mobile app (see Mobile App Setup)
3. Register a new account or login with existing credentials
4. Add your first habit using the "+" button
5. Track your daily progress by logging completions
6. Build and maintain your streaks!

### Understanding Streaks

- A streak continues if you log a habit daily
- You can skip 1 day without breaking your streak
- Missing 2 or more consecutive days resets the streak to 0
- The streak counter shows your current consecutive days

## REST API

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
This is a new [**React Native**](https://reactnative.dev) project, bootstrapped using [`@react-native-community/cli`](https://github.com/react-native-community/cli).

# Getting Started

> **Note**: Make sure you have completed the [Set Up Your Environment](https://reactnative.dev/docs/set-up-your-environment) guide before proceeding.

## Step 1: Start Metro

First, you will need to run **Metro**, the JavaScript build tool for React Native.

To start the Metro dev server, run the following command from the root of your React Native project:

```sh
# Using npm
npm start

# OR using Yarn
yarn start
```

## Step 2: Build and run your app

With Metro running, open a new terminal window/pane from the root of your React Native project, and use one of the following commands to build and run your Android or iOS app:

### Android

```sh
# Using npm
npm run android

# OR using Yarn
yarn android
```

### iOS

For iOS, remember to install CocoaPods dependencies (this only needs to be run on first clone or after updating native deps).

The first time you create a new project, run the Ruby bundler to install CocoaPods itself:

```sh
bundle install
```

Then, and every time you update your native dependencies, run:

```sh
bundle exec pod install
```

For more information, please visit [CocoaPods Getting Started guide](https://guides.cocoapods.org/using/getting-started.html).

```sh
# Using npm
npm run ios

# OR using Yarn
yarn ios
```

If everything is set up correctly, you should see your new app running in the Android Emulator, iOS Simulator, or your connected device.

This is one way to run your app — you can also build it directly from Android Studio or Xcode.

## Step 3: Modify your app

Now that you have successfully run the app, let's make changes!

Open `App.tsx` in your text editor of choice and make some changes. When you save, your app will automatically update and reflect these changes — this is powered by [Fast Refresh](https://reactnative.dev/docs/fast-refresh).

When you want to forcefully reload, for example to reset the state of your app, you can perform a full reload:

- **Android**: Press the <kbd>R</kbd> key twice or select **"Reload"** from the **Dev Menu**, accessed via <kbd>Ctrl</kbd> + <kbd>M</kbd> (Windows/Linux) or <kbd>Cmd ⌘</kbd> + <kbd>M</kbd> (macOS).
- **iOS**: Press <kbd>R</kbd> in iOS Simulator.

## Congratulations! :tada:

You've successfully run and modified your React Native App. :partying_face:

### Now what?

- If you want to add this new React Native code to an existing application, check out the [Integration guide](https://reactnative.dev/docs/integration-with-existing-apps).
- If you're curious to learn more about React Native, check out the [docs](https://reactnative.dev/docs/getting-started).

# Troubleshooting

If you're having issues getting the above steps to work, see the [Troubleshooting](https://reactnative.dev/docs/troubleshooting) page.

# Learn More

To learn more about React Native, take a look at the following resources:

- [React Native Website](https://reactnative.dev) - learn more about React Native.
- [Getting Started](https://reactnative.dev/docs/environment-setup) - an **overview** of React Native and how setup your environment.
- [Learn the Basics](https://reactnative.dev/docs/getting-started) - a **guided tour** of the React Native **basics**.
- [Blog](https://reactnative.dev/blog) - read the latest official React Native **Blog** posts.
- [`@facebook/react-native`](https://github.com/facebook/react-native) - the Open Source; GitHub **repository** for React Native.
