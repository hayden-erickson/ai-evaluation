# Habit Tracker UI

A modern, minimal React application for tracking daily habits and building lasting streaks.

## Features

- 🎯 **Habit Tracking**: Create and manage your daily habits
- 🔥 **Streak Counting**: Track consecutive days with visual indicators
- 📊 **History View**: See your last 30 days of habit logs
- ✅ **Easy Logging**: Quick log creation with notes and duration
- 🎨 **Modern Design**: Clean UI with pastel colors and rounded corners
- 🔒 **Secure Authentication**: JWT-based user authentication
- 📱 **Responsive**: Works on desktop, tablet, and mobile devices

## Prerequisites

Before running this application, ensure you have:

- **Node.js** (version 16 or higher)
- **npm** (comes with Node.js)
- **Backend API** running on `http://localhost:8080`

## Installation

1. **Install dependencies:**

```bash
npm install
```

## Running Locally

1. **Start the development server:**

```bash
npm run dev
```

2. **Open your browser:**

Navigate to [http://localhost:3000](http://localhost:3000)

3. **Ensure the backend API is running:**

The backend Go API should be running on `http://localhost:8080`. The frontend will proxy API requests to this address.

## Building for Production

To create a production build:

```bash
npm run build
```

The build output will be in the `dist` directory.

To preview the production build:

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── Auth/           # Authentication components
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   └── Auth.css
│   │   ├── Habits/         # Habit tracking components
│   │   │   ├── HabitList.jsx
│   │   │   ├── Habit.jsx
│   │   │   ├── StreakCount.jsx
│   │   │   ├── DayContainer.jsx
│   │   │   └── Habits.css
│   │   └── Modals/         # Modal components
│   │       ├── HabitDetailsModal.jsx
│   │       ├── LogDetailsModal.jsx
│   │       └── Modals.css
│   ├── services/
│   │   └── api.js          # API service layer
│   ├── utils/
│   │   └── streakUtils.js  # Streak calculation utilities
│   ├── App.jsx             # Main application component
│   ├── App.css
│   ├── main.jsx            # Application entry point
│   └── index.css           # Global styles
├── index.html
├── package.json
└── vite.config.js
```

## Usage Guide

### Getting Started

1. **Create an Account:**
   - Click "Register" on the login page
   - Enter your name, phone number, and password
   - You'll be automatically logged in

2. **Create Your First Habit:**
   - Click the "+ Add New Habit" button
   - Enter a name (required) and optional description
   - Optionally set a target duration in seconds
   - Click "Save Habit"

3. **Log Your Progress:**
   - Click "✓ Log for Today" on any habit
   - Add optional notes about your session
   - Add optional duration if you want to track time
   - Click "Save Log"

4. **View Your Streak:**
   - Your current streak is displayed prominently
   - Click "▼ Show History" to see the last 30 days
   - Days with logs are highlighted in green
   - Click any logged day to view or edit that log

### Streak Rules

- **Active Streak**: Log at least once per day to maintain your streak
- **Grace Period**: You can skip **one day** without breaking your streak
- **Broken Streak**: Missing more than one consecutive day resets your streak to 0

### Features

#### Habit Management
- ✏️ **Edit**: Modify habit name, description, or duration
- 🗑️ **Delete**: Remove a habit (requires confirmation)

#### Log Management
- 📝 **Add Notes**: Document how your session went
- ⏱️ **Track Duration**: Record how long you spent
- ✏️ **Edit Logs**: Click any logged day to modify
- 🗑️ **Delete Logs**: Remove logs you don't want

## API Configuration

The application communicates with a backend API. By default, it proxies requests from `/api` to `http://localhost:8080`.

To change the backend URL, modify `vite.config.js`:

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://your-backend-url:port',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

## Technologies Used

- **React 18**: UI library
- **Vite**: Build tool and development server
- **Vanilla CSS**: Styling (no CSS frameworks for minimal dependencies)
- **Standard Fetch API**: HTTP requests

## Security Features

- JWT-based authentication
- Token stored securely in localStorage
- Automatic token refresh on API errors
- Input validation on all forms
- XSS protection through React's built-in escaping
- CORS configured on backend

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Troubleshooting

### Backend Connection Issues

If you see "Network error" messages:
1. Verify the backend API is running on `http://localhost:8080`
2. Check the browser console for CORS errors
3. Ensure the backend has CORS headers configured

### Authentication Issues

If you're logged out unexpectedly:
1. Check if your JWT token has expired
2. Try logging in again
3. Clear localStorage and try again: `localStorage.clear()`

### Build Issues

If `npm install` fails:
1. Delete `node_modules` and `package-lock.json`
2. Run `npm install` again
3. Ensure you're using Node.js version 16 or higher

## Contributing

This is a demonstration project. For production use, consider adding:

- Unit tests (Jest, React Testing Library)
- E2E tests (Cypress, Playwright)
- State management (Redux, Zustand)
- Form validation library (React Hook Form)
- Error boundary components
- Analytics tracking
- Offline support with Service Workers

## License

MIT License - feel free to use this code for your own projects.

## Support

For issues or questions, please refer to the backend API documentation or contact the development team.

