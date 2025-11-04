# Testing the Habit Tracker App

## Prerequisites

Before testing, ensure you have:
1. Go 1.21+ installed
2. Node.js 20+ installed
3. React Native development environment set up
4. For iOS: Xcode and CocoaPods (macOS only)
5. For Android: Android Studio and Android SDK

## Quick Start Testing

### 1. Start the Backend Server

```bash
# From the project root
go run main.go
```

Expected output:
```
Database initialized successfully
Server starting on port 8080
```

### 2. Install Mobile Dependencies

```bash
npm install

# For iOS only
cd ios
bundle install
bundle exec pod install
cd ..
```

### 3. Configure API Endpoint

Edit `src/services/api.ts` and set the correct API URL:

**For iOS Simulator:**
```typescript
const API_BASE_URL = 'http://localhost:8080';
```

**For Android Emulator:**
```typescript
const API_BASE_URL = 'http://10.0.2.2:8080';
```

**For Physical Device:**
```typescript
const API_BASE_URL = 'http://YOUR_LOCAL_IP:8080';
```
(Replace YOUR_LOCAL_IP with your computer's IP on the local network)

### 4. Start Metro Bundler

```bash
npm start
```

### 5. Run the App

**iOS:**
```bash
npm run ios
```

**Android:**
```bash
npm run android
```

## Manual Test Scenarios

### Test 1: User Registration ✅
1. Launch the app
2. Tap "Don't have an account? Register"
3. Enter:
   - Name: "Test User"
   - Phone: "+1234567890"
   - Password: "testpass123"
4. Tap "Create Account"
5. ✅ Should auto-login and navigate to home screen

### Test 2: User Login ✅
1. Close and reopen the app
2. Enter credentials from Test 1
3. Tap "Sign In"
4. ✅ Should navigate to home screen

### Test 3: Create First Habit ✅
1. Tap "+ Create Your First Habit"
2. Enter:
   - Name: "Morning Exercise"
   - Description: "30 minutes of cardio"
3. Tap "Save"
4. ✅ Habit should appear in the list

### Test 4: Log Habit Completion ✅
1. Find "Morning Exercise" habit
2. Tap "+ Log Today"
3. Add notes: "Completed 5km run"
4. Tap "Save"
5. ✅ Streak counter should show "1 Day Streak 🔥"

### Test 5: View Streak History ✅
1. On the same habit, tap "Show History ▼"
2. ✅ Should show last 7 days
3. ✅ Today should be highlighted in green with checkmark
4. ✅ Other days should be empty

### Test 6: Edit Habit ✅
1. Tap the ✏️ icon on any habit
2. Change the name or description
3. Tap "Save"
4. ✅ Changes should be reflected immediately

### Test 7: Edit Log ✅
1. Tap "Show History" on a habit with logs
2. Tap on a green log entry
3. Modify the notes
4. Tap "Save"
5. ✅ Notes should be updated

### Test 8: Delete Log ✅
1. Long press on a log entry
2. Confirm deletion
3. ✅ Log should be removed
4. ✅ Streak counter should update

### Test 9: Streak with Skip ✅
This requires time manipulation or multiple days:
- Day 1: Log ✓ (Streak: 1)
- Day 2: Skip - (Streak still valid)
- Day 3: Log ✓ (Streak: 2)
- Day 4: Skip - (Streak still valid)
- Day 5: Skip - (Streak breaks, resets to 0)

### Test 10: Delete Habit ✅
1. Tap the 🗑️ icon on any habit
2. Confirm deletion
3. ✅ Habit should be removed from the list

### Test 11: Create Multiple Habits ✅
1. Tap the "+" FAB (floating action button)
2. Create 3-5 different habits
3. ✅ All should appear in scrollable list

### Test 12: Logout ✅
1. Tap "Logout" in the top right
2. Confirm logout
3. ✅ Should return to login screen
4. ✅ Token should be cleared

### Test 13: Session Persistence ✅
1. Login to the app
2. Close the app completely (force quit)
3. Reopen the app
4. ✅ Should still be logged in

### Test 14: Input Validation ✅
1. Try to create a habit with empty name
2. ✅ Should show error message
3. Try to register with invalid phone (no +)
4. ✅ Should show validation alert
5. Try to register with password < 8 chars
6. ✅ Should show validation alert

### Test 15: Pull to Refresh ✅
1. On the home screen with habits
2. Pull down the list
3. ✅ Loading indicator should appear
4. ✅ Habits should refresh

## Expected UI Appearance

### Colors
- Background: Off-white (#FEFEFE)
- Cards: White with light shadows
- Primary buttons: Pastel cyan (#A8DADC)
- Accent (FAB): Muted red (#E63946)
- Success indicators: Pastel green (#A8E6CF)
- Text: Dark gray (#2D3436)

### Typography
- System sans-serif font throughout
- Clear hierarchy (12px-36px)
- Readable and clean

### Layout
- Rounded corners on all interactive elements
- Consistent spacing (8px, 16px, 24px)
- No content overflow
- Proper safe area handling

## Common Issues and Solutions

### Backend Connection Failed
- Ensure backend is running: `curl http://localhost:8080/health`
- Check API_BASE_URL in `src/services/api.ts`
- For Android emulator, use `http://10.0.2.2:8080`
- For physical device, use your local IP

### iOS Pod Install Errors
```bash
cd ios
pod deintegrate
pod install
cd ..
```

### Metro Bundler Cache Issues
```bash
npm start -- --reset-cache
```

### TypeScript Errors
```bash
npx tsc --noEmit
```

### Linting Errors
```bash
npm run lint
```

## API Testing (Optional)

Test the backend independently using curl:

```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "phone_number": "+1234567890",
    "password": "testpass123",
    "time_zone": "America/New_York"
  }'

# Extract token from response and use it
TOKEN="your_token_here"

# Create habit
curl -X POST http://localhost:8080/habits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Habit",
    "description": "Test description"
  }'
```

## Success Criteria

The app is working correctly if:
- ✅ All registration and login flows work
- ✅ Habits can be created, edited, and deleted
- ✅ Logs can be created, edited, and deleted
- ✅ Streak calculation is accurate
- ✅ UI is clean and responsive
- ✅ No crashes or errors during normal use
- ✅ Session persists after app restart
- ✅ All validations work properly
- ✅ Error messages are user-friendly
