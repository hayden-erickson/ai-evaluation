# Habit Tracker Frontend

A modern, minimal React-based UI for the Habit Tracker API. Track your daily habits and maintain streaks with a clean, beautiful interface.

## Features

- **User Authentication** - Secure login and registration with JWT tokens
- **Habit Management** - Create, edit, and delete habits
- **Streak Tracking** - Visual streak counter with 1-day grace period
- **Daily Logs** - Log your habit completions with notes
- **Modern UI** - Clean design with pastel colors, rounded corners, and smooth animations
- **Responsive** - Works on desktop and mobile devices
- **Real-time Updates** - Instant feedback on all actions

## Prerequisites

- Node.js 14 or higher
- npm or yarn

## Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

## Running Locally

### Development Mode

Start the development server with hot reload:

```bash
npm start
```

The app will open at [http://localhost:3000](http://localhost:3000)

The development server is configured to proxy API requests to `http://localhost:8080`, so make sure your backend is running on port 8080.

### Production Build

Build the app for production:

```bash
npm run build
```

This creates an optimized build in the `build/` directory, which the Go backend will serve.

## Running the Full Stack

1. Build the frontend:
```bash
cd frontend
npm install
npm run build
cd ..
```

2. Start the Go backend:
```bash
go run main.go
```

3. Open your browser to [http://localhost:8080](http://localhost:8080)

The Go backend will serve both the API and the React frontend.

## Architecture

### Component Structure

```
App
├── Auth (Login/Register)
└── HabitList
    ├── HabitDetailsModal
    └── Habit (for each habit)
        ├── LogDetailsModal
        └── StreakList
```

### Key Components

- **App.js** - Main component handling authentication state
- **Auth.js** - Login and registration forms
- **HabitList.js** - List of all user habits
- **Habit.js** - Individual habit card with streak counter and log list
- **HabitDetailsModal.js** - Modal for creating/editing habits
- **LogDetailsModal.js** - Modal for creating/editing logs
- **api.js** - API client for backend communication

### Styling

The app uses a modern, minimal design with:
- Sans-serif fonts (system font stack)
- Pastel color gradients
- Rounded corners on all elements
- Smooth animations and transitions
- Responsive layout that adapts to screen size

### Streak Calculation

The app calculates streaks with a 1-day grace period:
- If you log a habit today, your streak continues
- If you skip one day, your streak continues (grace period)
- If you skip two days in a row, your streak resets to 0

### Security

- JWT tokens stored in localStorage
- All API requests include authentication headers
- Input validation on all forms
- Error handling with user-friendly messages

## API Integration

The frontend communicates with the backend API:

- **Authentication**
  - POST /users/register - Register new user
  - POST /users/login - Login user

- **Habits**
  - GET /habits - Get all user habits
  - POST /habits - Create new habit
  - PUT /habits/:id - Update habit
  - DELETE /habits/:id - Delete habit

- **Logs**
  - GET /habits/:id/logs - Get all logs for a habit
  - POST /habits/:id/logs - Create new log
  - PUT /logs/:id - Update log
  - DELETE /logs/:id - Delete log

## Troubleshooting

### API Connection Issues

If you see "Failed to load habits" or similar errors:
1. Ensure the backend is running on port 8080
2. Check browser console for CORS errors
3. Verify JWT token is stored (check localStorage in DevTools)

### Build Issues

If `npm run build` fails:
1. Delete `node_modules` and `package-lock.json`
2. Run `npm install` again
3. Try building again

### Authentication Issues

If you can't login:
1. Check backend logs for errors
2. Verify phone number format (include country code)
3. Ensure password is at least 6 characters

## Development

### Adding New Features

1. Create new components in `src/`
2. Update API client in `src/api.js` if needed
3. Add styles to `src/index.css`
4. Test thoroughly in development mode

### Code Style

- Use functional components with hooks
- Add JSDoc comments to all functions
- Use descriptive variable names
- Keep components focused and modular
- Handle errors gracefully with user feedback

## License

MIT
