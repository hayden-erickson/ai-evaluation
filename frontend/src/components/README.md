# Components Directory

This directory contains all React components for the Habit Tracker UI.

## Structure

```
components/
├── Auth/              # Authentication components
├── Habits/            # Habit management components
└── Modals/            # Modal dialog components
```

## Component Overview

### Auth Components

#### Login.js
- **Purpose**: User authentication form
- **Props**: 
  - `onLogin(userData)` - Called on successful login
  - `onSwitchToRegister()` - Switch to registration view
- **Features**:
  - Phone number and password inputs
  - Form validation
  - Error handling
  - Loading state

#### Register.js
- **Purpose**: New user registration form
- **Props**:
  - `onRegister(userData)` - Called on successful registration
  - `onSwitchToLogin()` - Switch to login view
- **Features**:
  - Name, phone, timezone, password inputs
  - Password confirmation
  - Form validation
  - Error handling

### Habit Components

#### HabitList.js
- **Purpose**: Container for all user habits
- **Props**: None (fetches data directly)
- **Features**:
  - Fetches all habits from API
  - Displays habit cards
  - "Add New Habit" button
  - Manages HabitDetailsModal
  - Loading and error states

#### Habit.js
- **Purpose**: Individual habit card with streak tracking
- **Props**:
  - `habit` - Habit object
  - `onEdit(habit)` - Edit callback
  - `onDelete(habitId)` - Delete callback
  - `onUpdate()` - Refresh callback
- **Features**:
  - Displays habit name and description
  - Shows current streak
  - "Log Today" button
  - 30-day calendar grid
  - Edit and delete buttons
  - Manages LogDetailsModal

### Modal Components

#### HabitDetailsModal.js
- **Purpose**: Create or edit habit
- **Props**:
  - `isOpen` - Modal visibility
  - `onClose()` - Close callback
  - `onSave(habitData)` - Save callback
  - `habit` - Habit to edit (null for new)
- **Features**:
  - Name input (required)
  - Description textarea (optional)
  - Form validation
  - Save and cancel buttons

#### LogDetailsModal.js
- **Purpose**: Create or edit habit log
- **Props**:
  - `isOpen` - Modal visibility
  - `onClose()` - Close callback
  - `onSave(logData)` - Save callback
  - `log` - Log to edit (null for new)
  - `date` - Date for the log
- **Features**:
  - Notes textarea
  - Date display
  - Save and cancel buttons

## Component Relationships

```
App
├── Login
├── Register
└── HabitList
    ├── HabitDetailsModal
    └── Habit (multiple)
        └── LogDetailsModal
```

## Usage Examples

### Using Login Component
```javascript
import Login from './components/Auth/Login';

<Login 
  onLogin={(userData) => console.log('Logged in:', userData)}
  onSwitchToRegister={() => setShowRegister(true)}
/>
```

### Using HabitList Component
```javascript
import HabitList from './components/Habits/HabitList';

<HabitList />
```

### Using Modals
```javascript
import HabitDetailsModal from './components/Modals/HabitDetailsModal';

<HabitDetailsModal
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  onSave={handleSave}
  habit={selectedHabit}
/>
```

## Styling

All components use classes from `src/index.css`. Key classes:

- `.btn`, `.btn-primary`, `.btn-secondary` - Buttons
- `.form-input`, `.form-textarea` - Form fields
- `.modal`, `.modal-overlay` - Modals
- `.habit-card` - Habit cards
- `.streak-badge` - Streak display
- `.days-grid`, `.day-container` - Calendar

## State Management

Components manage their own state using React hooks:
- `useState` - Local component state
- `useEffect` - Side effects (API calls, lifecycle)

No global state management library is used.

## Error Handling

All components implement error handling:
1. Try/catch blocks for async operations
2. Error state variable
3. Error message display
4. User-friendly error messages

## Best Practices

1. **Controlled Components** - All form inputs are controlled
2. **PropTypes** - Could be added for type checking
3. **Accessibility** - Labels, semantic HTML
4. **Loading States** - Show feedback during async operations
5. **Error Boundaries** - Graceful error handling
6. **Clean Code** - Comments, clear naming, modular design
