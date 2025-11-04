# Habit Tracker App

A modern, minimal React Native application for tracking daily habits and building streaks. Users can create habits, log daily progress, and maintain streaks with a one-day skip allowance.

## Features

- **User Authentication**: Secure registration and login with JWT tokens
- **Habit Management**: Create, edit, and delete habits with names and descriptions
- **Streak Tracking**: Visual streak counter with one-day skip allowance
- **Daily Logging**: Log habits daily with optional notes
- **30-Day History**: View the last 30 days of habit logs at a glance
- **Modern UI**: Clean, minimal design with pastel colors and rounded corners
- **Pull to Refresh**: Easily refresh your habit list
- **Error Handling**: Graceful error messages for all operations

## Tech Stack

- **Frontend**: React Native 0.82.1 with TypeScript
- **Navigation**: React Navigation 7.x
- **Storage**: AsyncStorage for token persistence
- **Date Handling**: date-fns for date calculations
- **Backend**: Go REST API with SQLite database
- **Authentication**: JWT-based authentication

## Prerequisites

Before running the app, ensure you have the following installed:

- **Node.js**: Version 20 or higher
- **npm** or **yarn**: Package manager
- **React Native CLI**: For running the app
- **Go**: Version 1.19 or higher (for backend)
- **Android Studio** (for Android development) or **Xcode** (for iOS development)

### Platform-Specific Requirements

#### Android
- Android Studio with Android SDK
- Android Virtual Device (AVD) or physical Android device
- Java Development Kit (JDK) 11 or higher

#### iOS (macOS only)
- Xcode 14 or higher
- CocoaPods
- iOS Simulator or physical iOS device

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd claude-sonnet-4.5
```

### 2. Install Dependencies

```bash
npm install
```

For iOS, also install CocoaPods dependencies:

```bash
cd ios
pod install
cd ..
```

### 3. Start the Backend Server

The app requires the Go backend to be running. In a separate terminal:

```bash
# Set environment variables (optional)
export PORT=8080
export JWT_SECRET=your-secret-key
export DB_PATH=habits.db

# Run the server
go run main.go
```

The server will start on `http://localhost:8080`.

### 4. Configure API Endpoint

If your backend is running on a different host/port, update the API base URL in:

```typescript
// src/services/api.ts
const API_BASE_URL = 'http://localhost:8080';
```

**Note for Android Emulator**: Use `http://10.0.2.2:8080` instead of `localhost`.

**Note for Physical Devices**: Use your computer's local IP address (e.g., `http://192.168.1.100:8080`).

## Running the App

### Start Metro Bundler

```bash
npm start
```

### Run on Android

In a new terminal:

```bash
npm run android
```

Or use Android Studio to build and run the app.

### Run on iOS (macOS only)

In a new terminal:

```bash
npm run ios
```

Or use Xcode to build and run the app.

## Usage

### First Time Setup

1. **Register**: Create a new account with your name, phone number, and password
2. **Login**: Sign in with your credentials

### Managing Habits

1. **Create a Habit**: Tap the pink "+" button at the bottom right
2. **Edit a Habit**: Tap the pencil icon (✏️) on any habit card
3. **Delete a Habit**: Tap the trash icon (🗑️) on any habit card

### Logging Habits

1. **Log Today**: Tap the "+ Log Today" button on a habit card
2. **Add Notes**: Optionally add notes about your progress
3. **View History**: Tap "Show History" to see the last 30 days
4. **Edit a Log**: Tap on any logged day in the history to edit notes

### Understanding Streaks

- Your streak increases each day you log the habit
- You can skip **one day** without breaking your streak
- Skipping two or more consecutive days resets the streak to 0
- The streak badge shows your current streak count

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Habit.tsx       # Individual habit card with streak display
│   ├── HabitDetailsModal.tsx  # Modal for creating/editing habits
│   └── LogDetailsModal.tsx    # Modal for creating/editing logs
├── screens/            # Screen components
│   ├── LoginScreen.tsx
│   ├── RegisterScreen.tsx
│   └── HomeScreen.tsx
├── services/           # API service layer
│   └── api.ts         # API client with authentication
├── utils/             # Utility functions
│   └── streakCalculator.ts  # Streak calculation logic
└── styles/            # Theme and styling
    └── theme.ts       # Colors, typography, spacing
```

## API Endpoints

The app communicates with the following backend endpoints:

### Authentication
- `POST /users/register` - Create a new user account
- `POST /users/login` - Authenticate and receive JWT token

### Habits
- `GET /habits` - Get all habits for authenticated user
- `POST /habits` - Create a new habit
- `GET /habits/:id` - Get a specific habit
- `PUT /habits/:id` - Update a habit
- `DELETE /habits/:id` - Delete a habit

### Logs
- `GET /habits/:id/logs` - Get all logs for a habit
- `POST /habits/:id/logs` - Create a new log
- `GET /logs/:id` - Get a specific log
- `PUT /logs/:id` - Update a log
- `DELETE /logs/:id` - Delete a log

All protected endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <token>
```

## Troubleshooting

### Metro Bundler Issues

If you encounter bundler issues:

```bash
npm start -- --reset-cache
```

### Android Build Errors

Clean the build:

```bash
cd android
./gradlew clean
cd ..
npm run android
```

### iOS Build Errors

Clean the build:

```bash
cd ios
rm -rf Pods Podfile.lock
pod install
cd ..
npm run ios
```

### Network Connection Issues

- Ensure the backend server is running
- Check the API_BASE_URL in `src/services/api.ts`
- For Android emulator, use `10.0.2.2` instead of `localhost`
- For physical devices, use your computer's local IP address
- Ensure your device and computer are on the same network

### Authentication Issues

If you're logged out unexpectedly:
- Check that the backend server is running
- Verify the JWT_SECRET matches between app restarts
- Clear app data and try logging in again

## Security Considerations

- **JWT Secret**: Change the default JWT secret in production
- **HTTPS**: Use HTTPS for production API endpoints
- **Password Storage**: Passwords are hashed using bcrypt on the backend
- **Token Storage**: JWT tokens are stored securely in AsyncStorage
- **Input Validation**: All user inputs are validated on both client and server

## Design Philosophy

The app follows a modern, minimal design approach:

- **Pastel Colors**: Soft, easy-on-the-eyes color palette
- **Sans-Serif Fonts**: Clean, readable typography
- **Rounded Corners**: Friendly, approachable UI elements
- **Generous Spacing**: Comfortable layout with breathing room
- **Clear Hierarchy**: Important information stands out
- **Intuitive Icons**: Emoji-based icons for quick recognition

## Future Enhancements

Potential features for future versions:

- Push notifications for habit reminders
- Habit categories and tags
- Statistics and analytics dashboard
- Social features (share streaks, compete with friends)
- Custom habit colors and icons
- Export data to CSV
- Dark mode support
- Habit templates and suggestions

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

## Support

For issues or questions, please [add contact information or issue tracker link].
