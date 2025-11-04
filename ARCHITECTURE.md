# React Native Habit Tracker - Developer Documentation

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Habit.tsx       # Individual habit card with streak display
│   ├── HabitList.tsx   # List of habits with empty state
│   ├── HabitDetailsModal.tsx  # Modal for creating/editing habits
│   └── LogDetailsModal.tsx    # Modal for creating/editing logs
│
├── screens/            # Full-screen views
│   ├── LoginScreen.tsx
│   ├── RegisterScreen.tsx
│   └── HomeScreen.tsx
│
├── services/           # External service integrations
│   └── api.ts         # API client for backend communication
│
├── types/             # TypeScript type definitions
│   └── index.ts       # Shared types (User, Habit, Log, etc.)
│
├── utils/             # Helper functions
│   └── streakUtils.ts # Streak calculation logic
│
└── styles/            # Theme and styling
    └── theme.ts       # Global colors, typography, spacing
```

## Component Hierarchy

```
App
├── LoginScreen
│   └── (forms and buttons)
├── RegisterScreen
│   └── (forms and buttons)
└── HomeScreen
    ├── HabitList
    │   ├── Habit (multiple)
    │   │   ├── LogDetailsModal
    │   │   └── DayContainer (for each day in StreakList)
    │   └── FloatingAddButton
    └── HabitDetailsModal
```

## Key Features

### 1. Authentication
- JWT-based authentication with secure token storage
- Phone number + password login
- User registration with validation
- Auto-logout on token expiration

### 2. Habit Management
- Create habits with name and description
- Edit existing habits
- Delete habits (with confirmation)
- View all habits in a list

### 3. Habit Logging
- Log habit completion for today
- Add optional notes to each log
- Edit existing logs
- Delete logs (long press on day container)
- View last 7 days of activity

### 4. Streak Calculation
- Counts consecutive days with logs
- Allows 1-day skip without breaking streak
- Resets to 0 after 2+ days without logging
- Visual indicator with colored badges

## Design System

### Color Palette (Pastel Theme)
- Primary: `#A8DADC` (Soft cyan)
- Secondary: `#F1FAEE` (Off white)
- Accent: `#E63946` (Soft red)
- Success: `#B8E6D5` (Soft mint green)
- Background: `#F8F9FA` (Light gray)

### Typography
- Font Family: System (sans-serif)
- Heading sizes: 32px, 24px, 20px, 18px
- Body: 16px
- Small: 14px

### Spacing
- xs: 4px
- sm: 8px
- md: 16px
- lg: 24px
- xl: 32px
- xxl: 48px

### Border Radius
- sm: 8px
- md: 12px
- lg: 16px
- xl: 24px
- full: 9999px (circular)

## API Integration

### Base URL Configuration
Default: `http://localhost:8080`

For different environments:
- iOS Simulator: `http://localhost:8080`
- Android Emulator: `http://10.0.2.2:8080`
- Physical Device: `http://<your-ip>:8080`

### Endpoints Used
1. `POST /users/register` - Create new account
2. `POST /users/login` - Authenticate user
3. `GET /habits` - Fetch all user habits
4. `POST /habits` - Create new habit
5. `PUT /habits/:id` - Update habit
6. `DELETE /habits/:id` - Delete habit
7. `GET /habits/:id/logs` - Fetch habit logs
8. `POST /habits/:id/logs` - Create log
9. `PUT /logs/:id` - Update log
10. `DELETE /logs/:id` - Delete log

### Error Handling
All API calls include:
- Try-catch blocks for error handling
- User-friendly error messages via Alert
- Network error detection
- Token expiration handling

## State Management

The app uses React hooks for state management:
- `useState` for local component state
- `useEffect` for side effects (API calls, data loading)
- Props drilling for passing data between components

### Key State Variables
- `currentUser`: Logged-in user object
- `currentScreen`: Current view ('login' | 'register' | 'home')
- `habits`: Array of user's habits
- `logs`: Array of logs for each habit
- `isLoading`: Loading state for async operations

## Testing

### Manual Testing Checklist
1. **Authentication**
   - [ ] Register new user
   - [ ] Login with valid credentials
   - [ ] Login with invalid credentials (should fail)
   - [ ] Logout

2. **Habit Management**
   - [ ] Create new habit
   - [ ] Edit habit
   - [ ] Delete habit
   - [ ] View empty state (no habits)
   - [ ] Pull to refresh

3. **Habit Logging**
   - [ ] Log habit for today
   - [ ] Prevent duplicate logs for same day
   - [ ] Edit existing log
   - [ ] Delete log (long press)
   - [ ] View 7-day activity

4. **Streak Calculation**
   - [ ] Verify streak counts correctly
   - [ ] Verify 1-day skip doesn't break streak
   - [ ] Verify 2+ day gap resets streak

### Unit Tests (Future)
Key functions to test:
- `calculateStreak()`
- `getLogsGroupedByDate()`
- `formatDate()`
- API service methods

## Performance Considerations

1. **API Calls**
   - Refresh only when needed (user action, pull-to-refresh)
   - No auto-refresh intervals to save battery/data

2. **Rendering**
   - FlatList for large lists (not needed for current scale)
   - Memoization opportunities (not implemented yet)

3. **Bundle Size**
   - Minimal dependencies to keep app size small
   - Standard React Native components only

## Security Best Practices

1. **Authentication**
   - Passwords never stored in plain text
   - JWT token stored in memory only (not persisted)
   - Token included in Authorization header

2. **Input Validation**
   - Client-side validation for all forms
   - Server-side validation in backend
   - XSS prevention via React Native's built-in escaping

3. **API Security**
   - HTTPS recommended for production
   - CORS configured on backend
   - Rate limiting recommended for production

## Future Enhancements

1. **Features**
   - Push notifications for habit reminders
   - Habit categories/tags
   - Statistics and charts
   - Social sharing of streaks
   - Dark mode support

2. **Technical**
   - Persistent token storage (AsyncStorage)
   - Offline mode support
   - Redux/Context API for global state
   - Automated tests (Jest, React Native Testing Library)
   - CI/CD pipeline

3. **UI/UX**
   - Animations and transitions
   - Haptic feedback
   - Customizable themes
   - Accessibility improvements (screen reader support)

## Troubleshooting

### Common Issues

1. **Cannot connect to backend**
   - Verify backend is running on correct port
   - Check API_BASE_URL in `src/services/api.ts`
   - For Android emulator, use `http://10.0.2.2:8080`

2. **Build errors**
   - Run `npm install` to ensure dependencies are installed
   - Clear Metro cache: `npm start -- --reset-cache`
   - For iOS: `cd ios && pod install`

3. **Authentication fails**
   - Verify backend database is accessible
   - Check network connectivity
   - Verify JWT_SECRET is consistent

4. **Streak calculation incorrect**
   - Check timezone settings
   - Verify log dates are correct
   - Review streak calculation logic in `streakUtils.ts`

## Contributing

When adding new features:
1. Follow existing code style and structure
2. Add TypeScript types for new data structures
3. Update this documentation
4. Test on both iOS and Android
5. Ensure backward compatibility with API

## License

MIT
