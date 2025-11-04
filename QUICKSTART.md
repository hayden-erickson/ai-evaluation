# Quick Start Guide - Habit Tracker App

This guide will help you get the Habit Tracker app running quickly.

## Prerequisites

- Node.js 20+
- Go 1.19+
- React Native development environment set up
- Android Studio (for Android) or Xcode (for iOS/macOS)

## 5-Minute Setup

### 1. Install Dependencies

```bash
npm install
```

For iOS (macOS only):
```bash
cd ios && pod install && cd ..
```

### 2. Start the Backend Server

Open a new terminal and run:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### 3. Start Metro Bundler

Open another terminal:

```bash
npm start
```

### 4. Run the App

**For Android:**
```bash
npm run android
```

**For iOS (macOS only):**
```bash
npm run ios
```

## First Time Usage

1. **Register a new account** with:
   - Name: Your name
   - Phone: Any phone number (e.g., "1234567890")
   - Password: At least 8 characters

2. **Create your first habit**:
   - Tap the pink "+" button
   - Enter habit name (e.g., "Morning Exercise")
   - Add optional description
   - Tap "Save"

3. **Log your habit**:
   - Tap "+ Log Today" on the habit card
   - Add optional notes
   - Tap "Save"

4. **View your streak**:
   - The green badge shows your current streak
   - Tap "Show History" to see the last 30 days

## Troubleshooting

### Can't connect to backend?

**Android Emulator:**
The app is configured to use `10.0.2.2:8080` for Android emulator. Make sure:
- Backend is running on port 8080
- No firewall blocking the connection

**iOS Simulator:**
Uses `localhost:8080`. Make sure the backend is running.

**Physical Device:**
Edit `src/config/config.ts` and change the API URL to your computer's local IP:
```typescript
return 'http://192.168.1.XXX:8080'; // Replace with your IP
```

Find your IP:
- **Windows**: `ipconfig` (look for IPv4 Address)
- **macOS/Linux**: `ifconfig` or `ip addr`

### Metro bundler cache issues?

```bash
npm start -- --reset-cache
```

### Build errors?

**Android:**
```bash
cd android && ./gradlew clean && cd ..
```

**iOS:**
```bash
cd ios && rm -rf Pods Podfile.lock && pod install && cd ..
```

## Key Features

- ✅ **One-day skip allowance**: Miss one day without breaking your streak
- 📊 **30-day history**: Visual calendar of your progress
- 📝 **Notes**: Add notes to each log entry
- 🎨 **Modern UI**: Pastel colors and clean design
- 🔄 **Pull to refresh**: Swipe down to refresh habits

## API Endpoints

All authenticated endpoints require a JWT token in the Authorization header.

**Public:**
- `POST /users/register` - Create account
- `POST /users/login` - Login

**Protected:**
- `GET /habits` - List habits
- `POST /habits` - Create habit
- `PUT /habits/:id` - Update habit
- `DELETE /habits/:id` - Delete habit
- `GET /habits/:id/logs` - Get habit logs
- `POST /habits/:id/logs` - Create log
- `PUT /logs/:id` - Update log
- `DELETE /logs/:id` - Delete log

## Project Structure

```
src/
├── components/     # UI components (Habit, Modals)
├── screens/        # Screen components (Login, Register, Home)
├── services/       # API client
├── utils/          # Streak calculation logic
├── styles/         # Theme and colors
└── config/         # Configuration
```

## Need Help?

See the full [HABIT_TRACKER_README.md](./HABIT_TRACKER_README.md) for detailed documentation.

## Development Tips

1. **Hot Reload**: Press `r` in Metro terminal to reload
2. **Debug Menu**: Shake device or press `Cmd+D` (iOS) / `Cmd+M` (Android)
3. **Console Logs**: Check Metro terminal for console output
4. **Network Requests**: Use React Native Debugger or Flipper

## Next Steps

- Customize colors in `src/styles/theme.ts`
- Add more habit fields in `src/services/api.ts`
- Implement push notifications
- Add habit categories
- Create statistics dashboard

Happy habit tracking! 🎯
