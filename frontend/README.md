# Habit Tracker - React Frontend

A modern, minimal React UI for tracking habits and building streaks. Built with React, Vite, and React Router.

## Features

- **User Authentication** - Secure JWT-based login and registration
- **Habit Management** - Create, edit, and delete habits
- **Streak Tracking** - Visual calendar showing daily progress
- **Flexible Streak Rules** - Users can skip one day before their streak resets
- **Log Management** - Add notes to daily habit logs
- **Modern UI** - Pastel color scheme with rounded corners and smooth animations
- **Responsive Design** - Works seamlessly on desktop and mobile devices

## Prerequisites

- Node.js 18.x or higher
- npm or yarn
- Backend API running on port 8080 (see main README.md)

## Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Configure environment variables (optional):
```bash
# Create .env file if you need to change the API URL
# Default is http://localhost:8080
echo "VITE_API_URL=http://localhost:8080" > .env
```

## Running Locally

### Development Mode

Start the development server with hot reload:

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Production Build

Build the application for production:

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

## Usage Guide

### First Time Setup

1. **Register an Account**
   - Navigate to the registration page
   - Enter your name, phone number, and password
   - Password must be at least 8 characters

2. **Login**
   - Use your phone number and password to login
   - Your session will be saved in localStorage

### Managing Habits

1. **Create a Habit**
   - Click the "Add New Habit" button
   - Enter a name and optional description
   - Click "Save"

2. **Edit a Habit**
   - Click the "Edit" button on any habit card
   - Update the name or description
   - Click "Save"

3. **Delete a Habit**
   - Click the "Delete" button on any habit card
   - Confirm the deletion

### Tracking Progress

1. **Log Today's Activity**
   - Click the "Log Today" button on any habit
   - Add optional notes about your progress
   - Click "Save"

2. **View Streak Calendar**
   - Each habit shows a 30-day calendar
   - Green squares indicate completed days
   - Click any day to view or edit the log

3. **Understanding Streaks**
   - Your streak counts consecutive days with logs
   - You can skip ONE day before your streak resets

## Security Features

- JWT Authentication
- Token Storage in localStorage
- Protected Routes
- Input Validation
- Error Handling

## License

MIT
