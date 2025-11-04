# Quick Start Guide

## Get Up and Running in 5 Minutes

### 1. Install Dependencies

```bash
npm install
```

For iOS:
```bash
cd ios && pod install && cd ..
```

### 2. Start the Backend API

In a new terminal:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### 3. Configure for Your Platform

**For Android Emulator**: Update `src/services/api.ts`:
```typescript
const API_BASE_URL = 'http://10.0.2.2:8080';
```

**For Physical Device**: Use your computer's local IP:
```typescript
const API_BASE_URL = 'http://192.168.1.100:8080'; // Replace with your IP
```

To find your IP:
- **Mac/Linux**: `ifconfig | grep "inet "`
- **Windows**: `ipconfig`

### 4. Start Metro

```bash
npm start
```

### 5. Run the App

In a new terminal:

**iOS:**
```bash
npm run ios
```

**Android:**
```bash
npm run android
```

### 6. Create Your First Account

1. Tap "Register" on the login screen
2. Fill in:
   - Name: Your name
   - Phone: Any unique phone number (e.g., +1234567890)
   - Time Zone: e.g., America/New_York
   - Password: At least 8 characters
3. Tap "Register"

### 7. Start Tracking Habits!

1. Tap the "+" button
2. Create your first habit (e.g., "Morning Exercise")
3. Tap "+ Log Today" to start your streak
4. Watch your streak counter grow! 🔥

## Common Issues

### "Network request failed"
- Ensure backend API is running
- Check API_BASE_URL matches your setup
- For Android emulator, use `10.0.2.2` instead of `localhost`

### "Cannot connect to Metro"
```bash
npm start -- --reset-cache
```

### Build errors
```bash
# Clean and rebuild
cd android && ./gradlew clean && cd ..
npm run android
```

## Test Users

You can create multiple test users with different phone numbers to test the app.

Each user's data is completely isolated - they can only see their own habits and logs.

## What's Next?

Check out the full [README.md](./README.md) for:
- Complete feature documentation
- Project structure
- API endpoints
- Troubleshooting guide
- Development guidelines

Happy habit tracking! 🎯

