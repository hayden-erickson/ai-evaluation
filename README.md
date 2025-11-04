# Habit Tracker App

A beautiful, modern React Native mobile application for tracking daily habits and building streaks. This app helps users maintain consistency by visualizing their progress and allowing them to skip one day before their streak resets.

## Features

- 🎯 **Habit Management**: Create, edit, and delete habits with names and descriptions
- 🔥 **Streak Tracking**: Visual streak counter with one-day skip allowance
- 📅 **Daily Logging**: Log habit completions with optional notes
- 📊 **Activity View**: See your last 14 days of activity at a glance
- 🔐 **Secure Authentication**: User registration and login with JWT tokens
- 🎨 **Modern UI**: Pastel color scheme with rounded corners and clean typography
- ✨ **Smooth UX**: Intuitive modals, pull-to-refresh, and confirmation dialogs

## Screenshots

The app features a clean, minimal design with:
- Pastel color palette (lavender, soft blue, soft pink, soft green)
- Sans-serif fonts throughout
- Rounded corners on all components
- Clear visual feedback for logged days
- Expandable habit cards with streak history

## Tech Stack

- **Frontend**: React Native with TypeScript
- **State Management**: React Context API
- **Storage**: AsyncStorage for authentication tokens
- **Backend**: Go REST API (included in this repository)
- **Database**: SQLite (for the backend)

## Prerequisites

Before running this app, make sure you have:

- Node.js (v18 or higher)
- npm or yarn
- React Native development environment set up
  - For iOS: Xcode (Mac only)
  - For Android: Android Studio and Android SDK
- Go (v1.21 or higher) for running the backend API

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd claude-sonnet-4.5
```

### 2. Install Dependencies

```bash
npm install
# or
yarn install
```

### 3. Install iOS Dependencies (iOS only)

```bash
cd ios
pod install
cd ..
```

## Running the App

### Step 1: Start the Backend API

The app requires the Go backend API to be running. In a terminal window:

```bash
# Set environment variables (optional)
export PORT=8080
export JWT_SECRET=your-secret-key
export DB_PATH=habits.db

# Run the API server
go run main.go
```

The API will start on `http://localhost:8080` by default.

### Step 2: Configure API Endpoint

If your backend is running on a different host/port, update the API base URL in:

```typescript
// src/services/api.ts
const API_BASE_URL = 'http://localhost:8080'; // Change this as needed
```

**Important for Android Emulator**: If testing on an Android emulator, use `http://10.0.2.2:8080` instead of `http://localhost:8080`.

**Important for Physical Devices**: Use your computer's local IP address (e.g., `http://192.168.1.100:8080`).

### Step 3: Start Metro Bundler

```bash
npm start
# or
yarn start
```

### Step 4: Run the App

In a separate terminal:

**For iOS:**
```bash
npm run ios
# or
yarn ios
```

**For Android:**
```bash
npm run android
# or
yarn android
```

## Usage

### First Time Setup

1. **Register**: When you first open the app, tap "Register" to create an account
   - Enter your name, phone number, time zone, and password (min 8 characters)
   - Phone number is used as your unique identifier

2. **Login**: After registration, you'll be automatically logged in
   - Use your phone number and password to login on subsequent sessions

### Managing Habits

1. **Create a Habit**: 
   - Tap the "+" floating action button
   - Enter a name (required) and optional description
   - Tap "Create"

2. **Edit a Habit**:
   - Tap "✏️ Edit" on any habit card
   - Update the details
   - Tap "Update"

3. **Delete a Habit**:
   - Tap "🗑️ Delete" on any habit card
   - Confirm the deletion (this also deletes all logs)

### Logging Habits

1. **Log Today**:
   - Tap "+ Log Today" button on a habit
   - Add optional notes about your progress
   - Tap "Save"

2. **Edit a Log**:
   - Tap to expand a habit to see recent activity
   - Tap on any day with a log (green square)
   - Update the notes
   - Tap "Update"

### Understanding Streaks

- **Streak Counter**: Shows the number of consecutive days you've logged this habit
- **One-Day Grace**: You can skip one day without breaking your streak
- **Two Days**: Missing two consecutive days will reset your streak to 0
- **Visual Feedback**: 
  - Green squares = Days with logs
  - Gray squares = Days without logs
  - Blue border = Today

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Habit.tsx       # Individual habit display
│   ├── HabitDetailsModal.tsx
│   ├── LogDetailsModal.tsx
│   └── StreakList.tsx  # Visual streak calendar
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state
├── screens/            # Main app screens
│   ├── HomeScreen.tsx  # Main habit list
│   ├── LoginScreen.tsx
│   └── RegisterScreen.tsx
├── services/           # API services
│   └── api.ts         # HTTP client & API methods
├── styles/            # Styling
│   ├── theme.ts       # Color palette & design tokens
│   └── commonStyles.ts
├── types/             # TypeScript types
│   └── models.ts      # Data models & interfaces
└── utils/             # Utility functions
    └── streakCalculator.ts
```

## API Endpoints

The app communicates with these backend endpoints:

### Authentication
- `POST /users/register` - Register new user
- `POST /users/login` - Login and get JWT token
- `GET /users/{id}` - Get user details
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Habits
- `GET /habits` - Get all user habits
- `POST /habits` - Create new habit
- `GET /habits/{id}` - Get specific habit
- `PUT /habits/{id}` - Update habit
- `DELETE /habits/{id}` - Delete habit

### Logs
- `GET /habits/{habit_id}/logs` - Get all logs for a habit
- `POST /habits/{habit_id}/logs` - Create new log
- `GET /logs/{id}` - Get specific log
- `PUT /logs/{id}` - Update log
- `DELETE /logs/{id}` - Delete log

## Troubleshooting

### Cannot Connect to Backend

- Verify the backend API is running (`curl http://localhost:8080/health`)
- Check the API_BASE_URL in `src/services/api.ts`
- For Android emulator, use `10.0.2.2` instead of `localhost`
- For physical devices, use your computer's local IP address
- Ensure your firewall allows the connection

### Build Errors

**iOS:**
```bash
cd ios
pod deintegrate
pod install
cd ..
npm run ios
```

**Android:**
```bash
cd android
./gradlew clean
cd ..
npm run android
```

### Metro Bundler Issues

```bash
npm start -- --reset-cache
```

### Missing Dependencies

```bash
npm install @react-native-async-storage/async-storage
```

## Development

### Running Tests

```bash
npm test
# or
yarn test
```

### Linting

```bash
npm run lint
# or
yarn lint
```

### Type Checking

```bash
npx tsc --noEmit
```

## Security Considerations

- **JWT Tokens**: Stored securely in AsyncStorage
- **Password Requirements**: Minimum 8 characters
- **Input Validation**: All inputs are validated on both client and server
- **Authentication**: All habit and log endpoints require valid JWT token
- **Authorization**: Users can only access their own data

## Design Principles

1. **Modern & Minimal**: Clean interface with plenty of white space
2. **Pastel Colors**: Soft, easy-on-the-eyes color palette
3. **Rounded Corners**: All buttons and cards have rounded corners
4. **Sans-serif Typography**: Clear, readable system fonts
5. **Intuitive UX**: Clear visual feedback and confirmation dialogs
6. **Mobile-First**: Optimized for touch interactions

## Future Enhancements

Potential features for future versions:
- Push notifications for habit reminders
- Habit categories and filtering
- Statistics and analytics dashboard
- Social features (share streaks with friends)
- Dark mode support
- Custom habit icons
- Export data functionality

## License

This project is for evaluation purposes.

## Support

For issues or questions, please refer to the project documentation or create an issue in the repository.

---

Built with ❤️ using React Native and Go
