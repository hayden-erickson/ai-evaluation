# Development Checklist

## Before Running the App

### Backend Setup
- [ ] Go backend server is running on port 8080
- [ ] Database is initialized (habits.db)
- [ ] JWT_SECRET environment variable is set (optional, has default)
- [ ] Backend health check responds: `curl http://localhost:8080/health`

### Frontend Setup
- [ ] Node.js 20+ is installed
- [ ] Dependencies installed: `npm install`
- [ ] For iOS: CocoaPods installed: `cd ios && pod install && cd ..`
- [ ] Metro bundler is running: `npm start`

### Platform-Specific
**Android:**
- [ ] Android Studio installed
- [ ] Android SDK configured
- [ ] Emulator running OR device connected
- [ ] USB debugging enabled (for physical device)

**iOS (macOS only):**
- [ ] Xcode installed
- [ ] iOS Simulator available OR device connected
- [ ] Development certificate configured (for physical device)

## Running the App

### Start Backend
```bash
# Terminal 1
go run main.go
```

### Start Metro
```bash
# Terminal 2
npm start
```

### Launch App
```bash
# Terminal 3 - Android
npm run android

# OR Terminal 3 - iOS
npm run ios
```

## First Run Testing

### Authentication Flow
- [ ] Register screen loads
- [ ] Can create account with valid credentials
- [ ] Validation works (8+ char password, required fields)
- [ ] After registration, automatically logged in
- [ ] Can logout
- [ ] Can login with created credentials
- [ ] Invalid credentials show error
- [ ] Token persists (close and reopen app)

### Habit Management
- [ ] Home screen shows empty state initially
- [ ] Can tap "+" button to open create modal
- [ ] Can create habit with name only
- [ ] Can create habit with name and description
- [ ] Habit appears in list after creation
- [ ] Can edit habit (tap pencil icon)
- [ ] Changes save correctly
- [ ] Can delete habit (tap trash icon)
- [ ] Confirmation dialog appears before delete
- [ ] Habit removed from list after delete

### Logging & Streaks
- [ ] New habit shows 0-day streak
- [ ] Can tap "+ Log Today" button
- [ ] Log modal opens
- [ ] Can save log with notes
- [ ] Can save log without notes
- [ ] Button changes to "✓ Logged Today" after logging
- [ ] Streak increases to 1 after first log
- [ ] Can't log again same day (button disabled)
- [ ] Can tap "Show History" to expand
- [ ] Last 30 days display correctly
- [ ] Today shows as logged (green with checkmark)
- [ ] Can tap logged day to edit notes
- [ ] Changes to log save correctly

### Streak Logic Testing
**Test 1: Consecutive Days**
- [ ] Log habit today → Streak = 1
- [ ] Change device date to tomorrow, log → Streak = 2
- [ ] Change device date to next day, log → Streak = 3

**Test 2: One-Day Skip**
- [ ] Log habit today → Streak = 1
- [ ] Skip one day (don't log)
- [ ] Log next day → Streak = 3 (skip allowed)

**Test 3: Two-Day Skip (Streak Break)**
- [ ] Log habit today → Streak = 1
- [ ] Skip two days
- [ ] Log on day 4 → Streak = 1 (reset)

### UI/UX Testing
- [ ] All text is readable
- [ ] Colors are pastel and pleasant
- [ ] Rounded corners on all cards/buttons
- [ ] Buttons respond to touch
- [ ] Loading indicators show during API calls
- [ ] Error messages are clear and helpful
- [ ] Modals can be dismissed by tapping outside
- [ ] Keyboard doesn't cover input fields
- [ ] Pull-to-refresh works on habit list
- [ ] Scrolling is smooth
- [ ] No UI elements cut off screen

### Error Handling
- [ ] Network error shows helpful message
- [ ] Backend down shows error (not crash)
- [ ] Invalid token logs user out
- [ ] Duplicate habit name handled gracefully
- [ ] Empty form submissions prevented
- [ ] API errors display to user

## Platform-Specific Checks

### Android
- [ ] App works on emulator
- [ ] App works on physical device
- [ ] API connects using 10.0.2.2 (emulator)
- [ ] Back button works correctly
- [ ] Status bar color is correct
- [ ] Keyboard behavior is correct

### iOS
- [ ] App works on simulator
- [ ] App works on physical device
- [ ] API connects using localhost (simulator)
- [ ] Safe area respected (notch devices)
- [ ] Status bar style is correct
- [ ] Keyboard behavior is correct

## Performance Checks
- [ ] App launches quickly (< 3 seconds)
- [ ] Screen transitions are smooth
- [ ] No lag when scrolling habit list
- [ ] API calls complete in reasonable time
- [ ] No memory leaks (check with profiler)
- [ ] App doesn't crash under normal use

## Code Quality Checks
- [ ] No TypeScript errors
- [ ] No ESLint warnings (run `npm run lint`)
- [ ] All functions have comments
- [ ] No console.log statements in production code
- [ ] Proper error handling throughout
- [ ] Input validation on all forms

## Documentation Checks
- [ ] HABIT_TRACKER_README.md is complete
- [ ] QUICKSTART.md has accurate instructions
- [ ] APP_SUMMARY.md reflects implementation
- [ ] Code comments are clear
- [ ] API endpoints documented

## Pre-Production Checklist
- [ ] Change JWT_SECRET to secure value
- [ ] Update API_BASE_URL to production URL
- [ ] Remove debug logging
- [ ] Test on multiple devices
- [ ] Test on different screen sizes
- [ ] Test with slow network
- [ ] Test with no network
- [ ] Add app icon
- [ ] Add splash screen
- [ ] Configure app name
- [ ] Set up error tracking (e.g., Sentry)
- [ ] Set up analytics (optional)
- [ ] Create privacy policy
- [ ] Create terms of service

## Known Issues to Check
- [ ] Date-fns timezone handling
- [ ] AsyncStorage quota limits
- [ ] JWT token expiration handling
- [ ] Concurrent API requests
- [ ] Race conditions in state updates

## Troubleshooting Steps

### If app won't start:
1. Clear Metro cache: `npm start -- --reset-cache`
2. Clean build: Android `cd android && ./gradlew clean` or iOS `rm -rf ios/build`
3. Reinstall dependencies: `rm -rf node_modules && npm install`
4. For iOS: `cd ios && rm -rf Pods && pod install`

### If can't connect to backend:
1. Check backend is running: `curl http://localhost:8080/health`
2. Check API URL in `src/config/config.ts`
3. For Android emulator, ensure using `10.0.2.2`
4. For physical device, use computer's local IP
5. Check firewall settings

### If authentication fails:
1. Check backend logs for errors
2. Verify JWT_SECRET is consistent
3. Clear AsyncStorage: Uninstall and reinstall app
4. Check network requests in debugger

## Success Criteria

The app is ready when:
- ✅ All authentication flows work
- ✅ All CRUD operations work
- ✅ Streak calculation is accurate
- ✅ UI is polished and responsive
- ✅ Error handling is comprehensive
- ✅ Documentation is complete
- ✅ No critical bugs
- ✅ Performance is acceptable

## Next Steps After Checklist

1. User acceptance testing
2. Beta testing with real users
3. Performance optimization
4. Additional features from roadmap
5. App store preparation
6. Production deployment

---

**Note**: This checklist should be completed before considering the app production-ready. Check off items as you verify them.
