# Habit Tracker Frontend

A modern, minimal React-based UI for tracking habits with streak functionality. Built with TypeScript and designed with a focus on usability and clean aesthetics.

## Features

- **User Authentication** - Secure login and registration with JWT tokens
- **Habit Management** - Create, edit, and delete habits with descriptions
- **Streak Tracking** - Track your consistency with current and longest streaks
- **Smart Streak Logic** - Allows one day skip before breaking a streak
- **Daily Logs** - Add notes to your daily habit completions
- **Visual Timeline** - See your last 14 days of progress at a glance
- **Modern UI** - Clean, minimal design with pastel colors and rounded corners
- **Responsive Design** - Works seamlessly on desktop and mobile devices
- **Error Handling** - Graceful error messages and input validation

## Tech Stack

- **React** - UI library
- **TypeScript** - Type-safe JavaScript
- **Context API** - State management for authentication
- **CSS3** - Modern styling with gradients and animations
- **Fetch API** - HTTP requests to backend

## Prerequisites

- Node.js 14 or higher
- npm or yarn
- Running instance of the Habit Tracker API (see main README)

## Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. (Optional) Configure the API URL:
Create a `.env` file in the frontend directory:
```
REACT_APP_API_URL=http://localhost:8080
```

## Development

Start the development server:
```bash
npm start
```

The app will open at [http://localhost:3000](http://localhost:3000).

## Building for Production

Build the optimized production bundle:
```bash
npm run build
```

This creates a `build` directory with static files ready to be served by the Go backend.

## Usage

### First Time Setup

1. **Register an Account**
   - Click "Sign up" on the login page
   - Enter your name, phone number, and password
   - Password must be at least 8 characters

2. **Create Your First Habit**
   - Click "+ Add Habit"
   - Enter a name (required) and description (optional)
   - Click "Save"

3. **Log Your Progress**
   - Click "+ Log Today" on any habit card
   - Optionally add notes about your progress
   - Click "Save"

### Understanding Streaks

- **Current Streak**: Days you've logged consecutively (allowing one skip)
- **Longest Streak**: Your best streak for this habit
- **One Day Skip Rule**: You can skip one day and maintain your streak
- **Streak Reset**: Missing two consecutive days breaks your streak

### Managing Habits

- **Edit Habit**: Click the pencil (✏️) icon
- **Delete Habit**: Click the trash (🗑️) icon
- **View History**: Click the "Last 14 Days" dropdown
- **Edit Log**: Click on any completed day in the timeline
- **Delete Log**: Click the "×" button when hovering over a log

## Project Structure

```
frontend/
├── public/              # Static files
├── src/
│   ├── components/      # React components
│   │   ├── Dashboard.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── HabitList.tsx
│   │   ├── Habit.tsx
│   │   ├── HabitDetailsModal.tsx
│   │   └── LogDetailsModal.tsx
│   ├── contexts/        # React contexts
│   │   └── AuthContext.tsx
│   ├── services/        # API client
│   │   └── api.ts
│   ├── types/           # TypeScript types
│   │   └── index.ts
│   ├── utils/           # Utility functions
│   │   └── streakUtils.ts
│   ├── styles/          # CSS files
│   │   ├── Auth.css
│   │   ├── Dashboard.css
│   │   ├── HabitList.css
│   │   ├── Habit.css
│   │   └── Modal.css
│   ├── App.tsx          # Main app component
│   ├── App.css          # App styles
│   └── index.tsx        # Entry point
└── package.json
```

## Component Architecture

### Authentication Flow
- `AuthProvider` manages authentication state
- `Login` and `Register` components handle user auth
- JWT tokens stored in localStorage
- Protected routes check authentication status

### Habit Management
- `HabitList` displays all user habits
- `Habit` component shows individual habit with streak data
- `HabitDetailsModal` for creating/editing habits
- `LogDetailsModal` for creating/editing logs

### Streak Calculation
- Calculated client-side from log data
- Allows one day skip between logs
- Shows current and longest streaks
- Visual timeline shows last 14 days

## Design System

### Colors
- **Primary**: `#a8d5e2` (light blue) - Main actions
- **Secondary**: `#fad0c4` (light pink) - Add buttons
- **Success**: `#ffeaa7` (light yellow) - Streak info
- **Background**: `#f5f7fa` (light gray)
- **Text**: `#2c3e50` (dark blue-gray)

### Typography
- Font family: System sans-serif stack
- Sizes: 12px - 32px
- Weights: 400 (regular), 500 (medium), 600 (semibold), 700 (bold)

### Layout
- Border radius: 10px - 20px
- Padding: 12px - 24px
- Gaps: 8px - 32px
- Max content width: 1200px

## API Integration

The frontend communicates with the backend API at `/api/*` endpoints:

- `POST /users/register` - Register new user
- `POST /users/login` - Login user
- `GET /habits` - Get all habits
- `POST /habits` - Create habit
- `PUT /habits/:id` - Update habit
- `DELETE /habits/:id` - Delete habit
- `GET /habits/:id/logs` - Get habit logs
- `POST /habits/:id/logs` - Create log
- `PUT /logs/:id` - Update log
- `DELETE /logs/:id` - Delete log

## Security

- JWT tokens for authentication
- Tokens stored in localStorage
- Authorization header on all protected requests
- Input validation on all forms
- Password minimum length enforcement
- XSS protection via React's automatic escaping

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Troubleshooting

### API Connection Issues
- Ensure the backend is running on the correct port
- Check REACT_APP_API_URL in .env file
- Verify CORS settings on the backend

### Build Errors
- Delete node_modules and package-lock.json
- Run `npm install` again
- Check Node.js version (14+)

### Authentication Issues
- Clear localStorage in browser dev tools
- Check JWT_SECRET matches between frontend and backend
- Verify token expiration hasn't passed

## Contributing

1. Create a new branch for your feature
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## License

MIT
