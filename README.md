# Habit Tracker App

A full-stack habit tracking application with a React Native mobile frontend and Go REST API backend. Track your daily habits, build streaks, and maintain consistency with an intuitive, modern interface.

## Overview

This project consists of two main components:
1. **React Native Mobile App** - Modern, minimal UI for tracking habits on iOS and Android
2. **Go REST API** - Backend server with JWT authentication and SQLite database

## Features

### Mobile App Features
- **JWT-based authentication** - Secure login and registration
- **Habit tracking** - Create, edit, and delete habits
- **Streak calculation** - Automatically track habit streaks with 1-day skip allowance
- **Log management** - Add notes to daily habit completions
- **Modern UI** - Pastel colors, rounded corners, and clean sans-serif fonts
- **Responsive design** - Works seamlessly on different screen sizes
- **Offline storage** - Authentication tokens persist between sessions
- **Error handling** - User-friendly error messages throughout

### Backend API Features
- **RESTful endpoints** - Complete CRUD operations for users, habits, and logs
- **JWT authentication** - Secure token-based authentication
- **Input validation** - Comprehensive request validation
- **Error logging** - All errors logged with context
- **Security headers** - XSS protection, clickjacking prevention, CSP
- **RBAC** - Role-based access control (users can only access their own data)
- **SQLite database** - Lightweight database for local development

## Prerequisites

### For Mobile App
- Node.js 20 or higher
- npm or yarn
- React Native development environment set up
  - For iOS: Xcode (macOS only)
  - For Android: Android Studio and Android SDK
- CocoaPods (for iOS dependencies)

### For Backend API
- Go 1.21 or higher
- SQLite3

## Installation

### 1. Clone the repository
```bash
git clone https://github.com/hayden-erickson/ai-evaluation.git
cd ai-evaluation
```

### 2. Set up the Backend API

Install Go dependencies:
```bash
go mod download
```

(Optional) Set environment variables:
```bash
export PORT=8080                    # Default: 8080
export JWT_SECRET=your-secret-key   # Default: "default-secret-key-change-in-production"
export DB_PATH=habits.db            # Default: habits.db
```

Start the backend server:
```bash
go run main.go
```

The API server will start on `http://localhost:8080` (or your configured PORT).

### 3. Set up the Mobile App

Install npm dependencies:
```bash
npm install
```

For iOS, install CocoaPods dependencies:
```bash
# First time setup - install Ruby dependencies
bundle install

# Install iOS pods
cd ios
bundle exec pod install
cd ..
```

### 4. Configure API Endpoint

Update the API base URL in `src/services/api.ts` if your backend is running on a different address:

```typescript
// For iOS Simulator
const API_BASE_URL = 'http://localhost:8080';

// For Android Emulator
const API_BASE_URL = 'http://10.0.2.2:8080';

// For physical device on same network
const API_BASE_URL = 'http://YOUR_LOCAL_IP:8080';
```

## Running the Application

### Start the Backend API

```bash
go run main.go
```

The server will display:
```
Database initialized successfully
Server starting on port 8080
```

### Start Metro Bundler

In a new terminal:
```bash
npm start
```

### Run on iOS

In a new terminal:
```bash
npm run ios
# OR
yarn ios
```

### Run on Android

In a new terminal:
```bash
npm run android
# OR
yarn android
```

Make sure you have:
- An iOS Simulator running (for iOS)
- An Android Emulator running or a physical device connected (for Android)

## Usage

### First Time Setup

1. Launch the app on your device/simulator
2. Register a new account:
   - Enter your name
   - Enter a phone number (format: +1234567890)
   - Create a password (minimum 8 characters)
3. After registration, you'll be automatically logged in

### Creating Habits

1. Tap the "+" button (or "Create Your First Habit" if you have no habits)
2. Enter a habit name (required)
3. Add an optional description
4. Tap "Save"

### Logging Habit Completion

1. Find the habit you want to log
2. Tap "+ Log Today" button
3. Optionally add notes about your completion
4. Tap "Save"

### Viewing Streak History

1. Tap "Show History ▼" on any habit
2. View the last 7 days of logs
3. Days with logs are highlighted in green
4. Days without logs show a dash

### Managing Habits

- **Edit**: Tap the ✏️ icon on any habit
- **Delete**: Tap the 🗑️ icon on any habit
- **Edit Log**: Tap on any green log entry in the history
- **Delete Log**: Long press on any log entry

## Streak Logic

The app tracks habit streaks with the following rules:

1. **Streak starts** when you complete a habit for the first time
2. **Streak continues** when you log the habit the next day
3. **One-day skip allowed** - You can miss one day and still maintain your streak
4. **Streak breaks** if you miss more than one consecutive day
5. **Streak counter** shows your current consecutive days

Example:
- Day 1: ✓ (Streak: 1)
- Day 2: ✓ (Streak: 2)
- Day 3: - (Skip allowed)
- Day 4: ✓ (Streak: 3)
- Day 5: - (Skip used)
- Day 6: - (Streak broken, resets to 0)

## Architecture

### Mobile App Structure

```
src/
├── components/          # React Native components
│   ├── Habit.tsx       # Individual habit card with streak
│   ├── HabitList.tsx   # List of all habits
│   ├── HabitDetailsModal.tsx  # Create/edit habit modal
│   ├── LogDetailsModal.tsx    # Create/edit log modal
│   ├── LoginScreen.tsx        # Authentication screen
│   └── HomeScreen.tsx         # Main app screen
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state management
├── services/           # API communication
│   └── api.ts         # API service with all endpoints
├── types/             # TypeScript type definitions
│   └── index.ts      # User, Habit, Log types
├── utils/            # Utility functions
│   └── helpers.ts    # Streak calculation, validation, date formatting
└── styles/           # Styling
    └── theme.ts      # Colors, spacing, typography
```

### Backend Package Structure

- **`/models`** - Data structures and validation logic
- **`/repository`** - Database operations and data access layer
- **`/service`** - Business logic and service layer
- **`/handlers`** - HTTP request handlers
- **`/middleware`** - Authentication, logging, and security middleware
- **`/config`** - Database configuration and migrations
- **`/utils`** - Utility functions (JWT, password hashing)

## Design System

### Colors (Pastel Palette)
- Primary: #A8DADC (Pastel cyan)
- Secondary: #F1FAEE (Pastel cream)
- Accent: #E63946 (Muted red)
- Success: #B8E6B8 (Pastel green)
- Error: #FFB3BA (Pastel red)

### Typography
- Font Family: System sans-serif
- Sizes: 12px - 36px
- Weights: Regular (400), Medium (500), Bold (600)

### Spacing
- Uses a consistent 4px base unit
- Scale: 4, 8, 16, 24, 32, 48

### Border Radius
- Small: 8px
- Medium: 12px
- Large: 16px
- Round: 999px (for circular elements)

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
  "time_zone": "America/New_York"
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
```

### Habit Endpoints (Requires Authentication)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <token>
```

#### Create a habit
```http
POST /habits
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

#### Update a habit
```http
PUT /habits/{id}
Authorization: Bearer <token>

{
  "name": "Evening Exercise"
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

#### Update a log
```http
PUT /logs/{id}
Authorization: Bearer <token>

{
  "notes": "Updated notes"
}
```

#### Delete a log
```http
DELETE /logs/{id}
Authorization: Bearer <token>
```

## Database Schema

### Users Table
- `id` - INTEGER PRIMARY KEY AUTOINCREMENT
- `name` - TEXT NOT NULL
- `phone_number` - TEXT NOT NULL (indexed)
- `password_hash` - TEXT NOT NULL
- `time_zone` - TEXT NOT NULL
- `profile_image_url` - TEXT
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
- **Secure Token Storage** - AsyncStorage for encrypted token persistence
- **RBAC** - Users can only access their own resources
- **Input Validation** - Client and server-side validation
- **Security Headers** - XSS protection, CSP, clickjacking prevention
- **HTTPS Support** - Ready for production deployment with HTTPS

## Testing

### Backend Tests
```bash
go test ./...
```

### Mobile App Tests
```bash
npm test
# OR
yarn test
```

## Troubleshooting

### iOS Build Issues

If you get CocoaPods errors:
```bash
cd ios
pod deintegrate
pod install
cd ..
```

### Android Build Issues

If you get Gradle errors:
```bash
cd android
./gradlew clean
cd ..
```

### Metro Bundler Issues

Clear the cache:
```bash
npm start -- --reset-cache
```

### API Connection Issues

**On iOS Simulator**: Use `http://localhost:8080`

**On Android Emulator**: Use `http://10.0.2.2:8080` (Android emulator's special alias for host machine)

**On Physical Device**: Use your computer's local IP address (e.g., `http://192.168.1.100:8080`)
- Both device and computer must be on the same network
- Make sure firewall allows connections on port 8080

## License

MIT
