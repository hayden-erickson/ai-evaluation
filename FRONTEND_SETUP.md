# Habit Tracker - Full Stack Setup Guide

This guide will help you set up and run both the backend API and the React frontend.

## Prerequisites

- **Go** (version 1.16 or higher) - for the backend
- **Node.js** (version 14 or higher) - for the frontend
- **npm** (comes with Node.js)

## Quick Start

### 1. Start the Backend API

Open a terminal in the project root directory and run:

```bash
# Build and run the Go backend
go run main.go
```

The backend API will start on `http://localhost:8080`

You should see output like:
```
Database initialized successfully
Server starting on port 8080
```

### 2. Start the Frontend

Open a **new terminal** window, navigate to the frontend directory, and run:

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies (first time only)
npm install

# Start the React development server
npm start
```

The React app will automatically open in your browser at `http://localhost:3000`

## What You Should See

1. **Login/Register Screen** - Create a new account or login
2. **Habit Dashboard** - After logging in, you'll see your habits list
3. **Add Habits** - Click "Add New Habit" to create your first habit
4. **Track Progress** - Log daily completions and watch your streaks grow!

## Architecture Overview

### Backend (Go)
- **Port:** 8080
- **Database:** SQLite (`habits.db`)
- **Authentication:** JWT tokens
- **API Endpoints:**
  - `/users/register` - Create account
  - `/users/login` - Login
  - `/habits` - Manage habits
  - `/habits/:id/logs` - Manage habit logs

### Frontend (React)
- **Port:** 3000 (development)
- **Framework:** React 18
- **Styling:** Custom CSS with pastel theme
- **State Management:** React hooks
- **API Communication:** Fetch API with JWT authentication

### Communication Flow

```
React Frontend (localhost:3000)
        ↓
    Fetch API
        ↓
CORS Middleware (allows localhost:3000)
        ↓
Go Backend API (localhost:8080)
        ↓
    SQLite Database
```

## Features Implemented

### ✅ Authentication
- User registration with validation
- Secure login with JWT tokens
- Password hashing
- Auto-logout on token expiration

### ✅ Habit Management
- Create, read, update, delete habits
- Habit name and description
- Multiple habits per user

### ✅ Streak Tracking
- Automatic streak calculation
- 1-day skip allowance (won't break streak)
- Visual 30-day calendar view
- Real-time streak updates

### ✅ Daily Logging
- Log habit completions
- Add notes to each log
- Edit or delete past logs
- Click any day to view/edit

### ✅ Modern UI
- Pastel color scheme
- Rounded corners
- Smooth animations
- Responsive design
- Mobile-friendly

## Troubleshooting

### Backend Issues

**Problem:** `go: cannot find main module`
```bash
# Solution: Make sure you're in the project root directory
cd /path/to/new-ui
go run main.go
```

**Problem:** Database errors
```bash
# Solution: Delete the database and restart
rm habits.db
go run main.go
```

### Frontend Issues

**Problem:** `npm: command not found`
```bash
# Solution: Install Node.js from https://nodejs.org/
```

**Problem:** Port 3000 already in use
```bash
# Solution: Kill the process or use a different port
# On Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or set a different port:
set PORT=3001 && npm start
```

**Problem:** "Failed to load habits"
```bash
# Solution: Make sure backend is running on port 8080
# Check if you can access: http://localhost:8080/health
```

### CORS Issues

If you see CORS errors in the browser console:

1. Verify backend is running on port 8080
2. Verify frontend is running on port 3000
3. Check that `middleware/cors.go` allows `http://localhost:3000`
4. Restart both servers

## Development Workflow

### Making Changes to Backend

1. Edit Go files
2. Stop the server (Ctrl+C)
3. Restart with `go run main.go`

### Making Changes to Frontend

1. Edit React files
2. Changes will auto-reload in the browser
3. No restart needed (hot reload enabled)

### Testing the API Directly

You can test API endpoints using curl:

```bash
# Register a user
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","phone_number":"1234567890","time_zone":"America/New_York","password":"password123"}'

# Login
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"1234567890","password":"password123"}'

# Get habits (replace TOKEN with actual JWT token from login)
curl -X GET http://localhost:8080/habits \
  -H "Authorization: Bearer TOKEN"
```

## Project Structure

```
new-ui/
├── frontend/                 # React application
│   ├── public/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── services/        # API service layer
│   │   ├── utils/           # Utility functions
│   │   ├── App.js           # Main app component
│   │   ├── index.js         # Entry point
│   │   └── index.css        # Global styles
│   ├── package.json
│   └── README.md
├── config/                   # Database configuration
├── handlers/                 # HTTP request handlers
├── middleware/               # HTTP middleware
│   ├── auth.go              # JWT authentication
│   ├── cors.go              # CORS headers
│   ├── logging.go           # Request logging
│   └── security.go          # Security headers
├── models/                   # Data models
├── repository/               # Database operations
├── service/                  # Business logic
├── utils/                    # Utility functions
├── main.go                   # Backend entry point
└── habits.db                 # SQLite database (created on first run)
```

## Next Steps

1. **Create an account** - Register with your phone number
2. **Add your first habit** - Click "Add New Habit"
3. **Log your progress** - Click "Log Today" to record completions
4. **Build your streak** - Keep logging daily to build your streak!
5. **Explore the UI** - Click on past days to view or edit logs

## Security Notes

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- All authenticated endpoints require valid tokens
- CORS is configured for localhost development
- Input validation on both frontend and backend

## Production Deployment

For production deployment:

1. Build the React app: `cd frontend && npm run build`
2. Serve the build folder with the Go backend
3. Update CORS settings to allow your production domain
4. Enable HTTPS and update security headers
5. Use environment variables for sensitive configuration
6. Use a production-grade database (PostgreSQL, MySQL)

## Support

If you encounter any issues:

1. Check both terminal windows for error messages
2. Verify both servers are running
3. Check browser console for frontend errors
4. Review the troubleshooting section above
5. Ensure all prerequisites are installed

Happy habit tracking! 🎯
