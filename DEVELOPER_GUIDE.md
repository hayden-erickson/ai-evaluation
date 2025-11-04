# Developer Quick Reference

## Project Structure

```
ai-evaluation/
├── App.tsx                    # Root React Native component
├── src/
│   ├── components/           # UI components
│   │   ├── Habit.tsx        # Habit card with streak
│   │   ├── HabitList.tsx    # List of habits
│   │   ├── HabitDetailsModal.tsx  # Create/edit habit
│   │   ├── LogDetailsModal.tsx    # Create/edit log
│   │   ├── LoginScreen.tsx  # Authentication
│   │   └── HomeScreen.tsx   # Main app screen
│   ├── contexts/
│   │   └── AuthContext.tsx  # Auth state management
│   ├── services/
│   │   └── api.ts          # API client
│   ├── types/
│   │   └── index.ts        # TypeScript types
│   ├── utils/
│   │   └── helpers.ts      # Utility functions
│   └── styles/
│       └── theme.ts        # Design system
├── handlers/                # Go HTTP handlers
├── service/                 # Go business logic
├── repository/             # Go data access
├── models/                 # Go data models
└── main.go                 # Go server entry point
```

## Key API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/users/register` | Register new user | No |
| POST | `/users/login` | Login user | No |
| GET | `/habits` | List user's habits | Yes |
| POST | `/habits` | Create habit | Yes |
| PUT | `/habits/:id` | Update habit | Yes |
| DELETE | `/habits/:id` | Delete habit | Yes |
| GET | `/habits/:id/logs` | List habit logs | Yes |
| POST | `/habits/:id/logs` | Create log | Yes |
| PUT | `/logs/:id` | Update log | Yes |
| DELETE | `/logs/:id` | Delete log | Yes |

## Environment Variables

### Backend (Go)
```bash
PORT=8080                    # Server port
JWT_SECRET=your-secret       # JWT signing key
DB_PATH=habits.db            # SQLite database path
```

### Frontend (React Native)
Edit `src/services/api.ts`:
```typescript
const API_BASE_URL = 'http://localhost:8080';
```

## Common Commands

### Backend
```bash
# Run server
go run main.go

# Build binary
go build -o habit-server main.go

# Run tests
go test ./...

# Format code
go fmt ./...
```

### Frontend
```bash
# Install dependencies
npm install

# Start Metro
npm start

# Run iOS
npm run ios

# Run Android
npm run android

# Lint
npm run lint

# Type check
npx tsc --noEmit

# Clear cache
npm start -- --reset-cache
```

## Streak Calculation Logic

Located in `src/utils/helpers.ts`:

```typescript
calculateStreak(logs: Log[]): StreakInfo
```

Rules:
1. Count consecutive days with logs
2. Allow 1-day skip without breaking streak
3. Reset to 0 if more than 1 day gap
4. Return current streak count and continuation status

Example:
```
Day 1: ✓ → Streak: 1
Day 2: ✓ → Streak: 2
Day 3: - → Streak: 2 (skip allowed)
Day 4: ✓ → Streak: 3
Day 5: - → Streak: 3 (skip allowed)
Day 6: - → Streak: 0 (broken)
```

## Design Tokens

From `src/styles/theme.ts`:

```typescript
// Colors
colors.primary = '#A8DADC'    // Pastel cyan
colors.accent = '#E63946'     // Muted red
colors.success = '#B8E6B8'    // Pastel green
colors.text = '#2D3436'       // Dark gray

// Spacing
spacing.sm = 8
spacing.md = 16
spacing.lg = 24

// Border Radius
borderRadius.sm = 8
borderRadius.md = 12
borderRadius.lg = 16
```

## Authentication Flow

1. User opens app
2. `AuthContext` checks `AsyncStorage` for token
3. If token exists → `HomeScreen`
4. If no token → `LoginScreen`
5. After login/register → Save token to `AsyncStorage`
6. All API calls include `Authorization: Bearer <token>` header

## Adding New Features

### Adding a New Component
1. Create file in `src/components/`
2. Use theme from `src/styles/theme.ts`
3. Follow existing component patterns
4. Add TypeScript types
5. Document with JSDoc comments

### Adding a New API Endpoint
1. Add method to appropriate service in `src/services/api.ts`
2. Add types to `src/types/index.ts`
3. Handle errors with `ApiError`
4. Use in components via async/await

### Adding Backend Endpoint
1. Add handler in `handlers/`
2. Add service method in `service/`
3. Add repository method in `repository/`
4. Add model/validation in `models/`
5. Register route in `main.go`

## Testing Checklist

- [ ] Backend compiles: `go build main.go`
- [ ] Frontend compiles: `npx tsc --noEmit`
- [ ] Linting passes: `npm run lint`
- [ ] Backend runs: `go run main.go`
- [ ] Health endpoint works: `curl http://localhost:8080/health`
- [ ] Can register user
- [ ] Can login
- [ ] Can create habit
- [ ] Can log habit
- [ ] Streak calculates correctly
- [ ] UI looks good on both iOS and Android

## Troubleshooting

### "Cannot connect to API"
- Check backend is running
- Verify API_BASE_URL is correct
- For Android emulator, use `10.0.2.2` not `localhost`

### "Pod install failed" (iOS)
```bash
cd ios
rm -rf Pods Podfile.lock
pod install
cd ..
```

### "Metro bundler error"
```bash
watchman watch-del-all
rm -rf node_modules
npm install
npm start -- --reset-cache
```

### "Database locked"
- Only one backend instance can run at a time
- Kill existing process: `pkill -f habit-server`

## Code Style

### TypeScript
- Use functional components
- Use hooks (useState, useEffect)
- Use TypeScript strict mode
- Document complex functions
- Handle all error cases

### Go
- Use standard project layout
- Dependency injection
- Interface-based design
- Proper error handling
- Structured logging

## Security Considerations

1. **Passwords**: Hashed with Argon2id
2. **JWT**: 24-hour expiration
3. **Authorization**: Users can only access their own data
4. **Input Validation**: Client and server-side
5. **Token Storage**: AsyncStorage (encrypted on device)

## Performance Tips

1. Use pull-to-refresh instead of auto-polling
2. Cache user data in AuthContext
3. Lazy load images (if added)
4. Debounce search inputs (if added)
5. Use FlatList for long lists

## Future Enhancements

Potential features to add:
- [ ] Profile pictures
- [ ] Push notifications for reminders
- [ ] Social features (share streaks)
- [ ] Charts and analytics
- [ ] Custom habit frequencies (weekly, etc.)
- [ ] Dark mode
- [ ] Habit categories/tags
- [ ] Export data
- [ ] Multiple languages
- [ ] Offline mode with sync

## Resources

- [React Native Docs](https://reactnative.dev/docs/getting-started)
- [Go Documentation](https://golang.org/doc/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [SQLite Docs](https://www.sqlite.org/docs.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
