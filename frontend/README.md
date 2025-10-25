# Habit Tracker UI

A modern, minimal React application for tracking daily habits and maintaining streaks.

## Features

- 🔐 **User Authentication** - Secure login and registration
- 📊 **Habit Tracking** - Create and manage multiple habits
- 🔥 **Streak Counting** - Track consecutive days with 1-day skip allowance
- 📝 **Daily Logs** - Add notes for each habit completion
- 🎨 **Modern UI** - Clean, pastel-themed interface with rounded corners
- 📱 **Responsive Design** - Works on desktop and mobile devices

## Prerequisites

Before running this application, ensure you have:

- **Node.js** (version 14 or higher)
- **npm** (comes with Node.js)
- **Backend API** running on `http://localhost:8080`

## Installation

1. **Navigate to the frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

## Running the Application

### Development Mode

1. **Start the backend API server** (in a separate terminal):
   ```bash
   # From the project root directory
   go run main.go
   ```
   The backend should be running on `http://localhost:8080`

2. **Start the React development server:**
   ```bash
   npm start
   ```
   The app will automatically open in your browser at `http://localhost:3000`

### Production Build

To create an optimized production build:

```bash
npm run build
```

The build files will be in the `build/` directory.

## Usage Guide

### First Time Setup

1. **Register an Account**
   - Click "Register" on the login page
   - Fill in your name, phone number, time zone, and password
   - Password must be at least 8 characters long

2. **Login**
   - Enter your phone number and password
   - Click "Login"

### Managing Habits

1. **Create a New Habit**
   - Click the "➕ Add New Habit" button
   - Enter a name (required) and optional description
   - Click "Save"

2. **Edit a Habit**
   - Click the ✏️ (edit) button on any habit card
   - Update the name or description
   - Click "Save"

3. **Delete a Habit**
   - Click the 🗑️ (delete) button on any habit card
   - Confirm the deletion
   - Note: This will also delete all associated logs

### Logging Habit Completions

1. **Log Today's Completion**
   - Click the "✅ Log Today" button on a habit card
   - Add optional notes about your progress
   - Click "Save Log"

2. **View Past Logs**
   - The last 30 days are displayed as a grid below each habit
   - Green squares indicate days with logs
   - Click any day to view or edit its log

3. **Understanding Streaks**
   - Your current streak is displayed prominently on each habit
   - You can skip one day without breaking your streak
   - Skipping more than one day will reset the streak to 0

## Project Structure

```
frontend/
├── public/
│   └── index.html          # HTML template
├── src/
│   ├── components/
│   │   ├── Auth/
│   │   │   ├── Login.js           # Login component
│   │   │   └── Register.js        # Registration component
│   │   ├── Habits/
│   │   │   ├── Habit.js           # Individual habit card
│   │   │   └── HabitList.js       # List of all habits
│   │   └── Modals/
│   │       ├── HabitDetailsModal.js  # Create/edit habit modal
│   │       └── LogDetailsModal.js    # Create/edit log modal
│   ├── services/
│   │   └── api.js          # API service layer
│   ├── utils/
│   │   └── streakCalculator.js  # Streak calculation logic
│   ├── App.js              # Main application component
│   ├── index.js            # React entry point
│   └── index.css           # Global styles
├── package.json            # Dependencies and scripts
└── README.md              # This file
```

## API Integration

The application communicates with the backend API using the following endpoints:

### Authentication
- `POST /users/register` - Register new user
- `POST /users/login` - Login user

### Habits
- `GET /habits` - Get all user habits
- `POST /habits` - Create new habit
- `GET /habits/:id` - Get specific habit
- `PUT /habits/:id` - Update habit
- `DELETE /habits/:id` - Delete habit

### Logs
- `GET /habits/:id/logs` - Get logs for a habit
- `POST /habits/:id/logs` - Create new log
- `PUT /logs/:id` - Update log
- `DELETE /logs/:id` - Delete log

All authenticated requests include a JWT token in the `Authorization` header.

## Security Features

- **JWT Authentication** - Secure token-based authentication
- **Password Validation** - Minimum 8 characters required
- **Input Validation** - All user inputs are validated
- **Secure Storage** - Auth tokens stored in localStorage
- **Auto-logout** - Automatic logout on 401 responses

## Styling

The application uses a modern pastel color scheme with:
- Sans-serif fonts for clean readability
- Rounded corners on all UI elements
- Soft shadows for depth
- Responsive grid layouts
- Smooth transitions and hover effects

## Troubleshooting

### Backend Connection Issues

If you see "Failed to load habits" or authentication errors:

1. Ensure the backend server is running on `http://localhost:8080`
2. Check that the `proxy` setting in `package.json` is correct
3. Verify your network connection

### Build Errors

If you encounter build errors:

1. Delete `node_modules/` and `package-lock.json`
2. Run `npm install` again
3. Clear your browser cache

### Authentication Issues

If you can't log in:

1. Verify your credentials are correct
2. Check browser console for error messages
3. Ensure the backend database is properly initialized

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

This project is part of the AI Evaluation codebase.
