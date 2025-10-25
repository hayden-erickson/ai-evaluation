# Component Hierarchy and Data Flow

## Visual Component Tree

```
App (Main Application)
│
├── [Not Authenticated]
│   │
│   ├── Login
│   │   ├── Form (phone_number, password)
│   │   ├── Submit Button → authAPI.login()
│   │   └── Switch to Register Link
│   │
│   └── Register
│       ├── Form (name, phone_number, time_zone, password)
│       ├── Submit Button → authAPI.register()
│       └── Switch to Login Link
│
└── [Authenticated]
    │
    ├── Header
    │   ├── Title
    │   ├── Welcome Message
    │   └── Logout Button → authAPI.logout()
    │
    └── HabitList
        │
        ├── Add New Habit Button → Opens HabitDetailsModal
        │
        ├── Habit (for each habit)
        │   │
        │   ├── Habit Header
        │   │   ├── Habit Info
        │   │   │   ├── Name
        │   │   │   └── Description
        │   │   │
        │   │   └── Habit Actions
        │   │       ├── Edit Button → Opens HabitDetailsModal
        │   │       └── Delete Button → habitsAPI.delete()
        │   │
        │   ├── Streak Badge
        │   │   └── Streak Number (calculated from logs)
        │   │
        │   ├── Log Today Button → Opens LogDetailsModal
        │   │
        │   ├── Days Grid (30-day calendar)
        │   │   └── DayContainer (×30)
        │   │       ├── [Has Log] → Click to edit
        │   │       └── [No Log] → Click to create
        │   │
        │   └── LogDetailsModal (conditional)
        │       ├── Date Display
        │       ├── Notes Field
        │       ├── Cancel Button
        │       └── Save Button → logsAPI.create() or logsAPI.update()
        │
        └── HabitDetailsModal (conditional)
            ├── Modal Title (New/Edit)
            ├── Name Field
            ├── Description Field
            ├── Cancel Button
            └── Save Button → habitsAPI.create() or habitsAPI.update()
```

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         User Actions                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    React Components                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Login   │  │ Register │  │ HabitList│  │  Habit   │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Service Layer                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │ authAPI  │  │habitsAPI │  │ logsAPI  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
│                                                              │
│  • Handles HTTP requests                                    │
│  • Manages JWT tokens                                       │
│  • Error handling                                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend API (Go)                          │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              CORS Middleware                         │  │
│  │  • Allows localhost:3000                             │  │
│  │  • Handles preflight OPTIONS requests               │  │
│  └──────────────────────────────────────────────────────┘  │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Authentication Middleware                    │  │
│  │  • Validates JWT tokens                              │  │
│  │  • Extracts user ID                                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Route Handlers                          │  │
│  │  • UserHandler                                       │  │
│  │  • HabitHandler                                      │  │
│  │  • LogHandler                                        │  │
│  └──────────────────────────────────────────────────────┘  │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Service Layer                           │  │
│  │  • Business logic                                    │  │
│  │  • Validation                                        │  │
│  └──────────────────────────────────────────────────────┘  │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │            Repository Layer                          │  │
│  │  • Database operations                               │  │
│  │  • SQL queries                                       │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    SQLite Database                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │  users   │  │  habits  │  │   logs   │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

## State Management

### App Component State
```javascript
{
  isAuthenticated: boolean,    // User login status
  showRegister: boolean,        // Toggle login/register view
  user: Object | null,          // Current user data
  loading: boolean              // Initial auth check
}
```

### HabitList Component State
```javascript
{
  habits: Array,                // List of all user habits
  loading: boolean,             // Loading state
  error: string,                // Error message
  isModalOpen: boolean,         // Habit modal visibility
  selectedHabit: Object | null  // Habit being edited
}
```

### Habit Component State
```javascript
{
  logs: Array,                  // Habit logs
  loading: boolean,             // Loading state
  error: string,                // Error message
  isLogModalOpen: boolean,      // Log modal visibility
  selectedLog: Object | null,   // Log being edited
  selectedDate: string | null   // Date for new/edit log
}
```

### Modal Component State
```javascript
{
  formData: Object,             // Form field values
  error: string,                // Validation/API errors
  loading: boolean              // Submit state
}
```

## API Request Flow

### Example: Creating a New Habit

```
1. User clicks "Add New Habit"
   └─> HabitList sets isModalOpen = true

2. User fills form and clicks "Save"
   └─> HabitDetailsModal.handleSubmit()
       └─> Validates input
       └─> Calls onSave(formData)
           └─> HabitList.handleSaveHabit()
               └─> habitsAPI.create(habitData)
                   └─> apiRequest('/habits', { method: 'POST', body: JSON.stringify(habitData) })
                       └─> Adds Authorization header with JWT token
                       └─> fetch('http://localhost:8080/habits', ...)
                           
3. Backend receives request
   └─> CORS Middleware (allows request)
   └─> Auth Middleware (validates JWT, extracts user_id)
   └─> HabitHandler.CreateHabit()
       └─> HabitService.CreateHabit()
           └─> Validates habit data
           └─> HabitRepository.Create()
               └─> INSERT INTO habits (user_id, name, description, ...)
               
4. Backend sends response
   └─> Returns created habit with ID
   
5. Frontend receives response
   └─> habitsAPI.create() returns habit
   └─> HabitList.handleSaveHabit() calls loadHabits()
       └─> Refreshes habit list
   └─> Modal closes
   └─> New habit appears in UI
```

## Streak Calculation Flow

```
1. Habit component loads
   └─> useEffect() calls loadLogs()
       └─> habitsAPI.getLogs(habitId)
           └─> GET /habits/:id/logs
           
2. Logs received from API
   └─> setLogs(fetchedLogs)
   
3. Component renders
   └─> calculateStreak(logs) called
       └─> Sorts logs by date (newest first)
       └─> Removes duplicate dates
       └─> Checks if last log is within 2 days
       └─> Counts consecutive days (allowing 1 skip)
       └─> Returns streak number
       
4. Streak displayed in UI
   └─> <div className="streak-badge">
       └─> Shows streak number with fire emoji
       
5. 30-day calendar rendered
   └─> getLast30Days(logs) called
       └─> Creates array of last 30 days
       └─> Maps logs to dates
       └─> Returns array with hasLog flags
       
6. Days rendered in grid
   └─> Green background if hasLog
   └─> Border highlight if isToday
   └─> Click handler for each day
```

## Authentication Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Initial Page Load                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    Check localStorage for token
                              │
                ┌─────────────┴─────────────┐
                │                           │
           Token exists                No token
                │                           │
                ▼                           ▼
        Show main app                Show login page
                │                           │
                │                           ▼
                │                   User logs in/registers
                │                           │
                │                           ▼
                │                   API returns JWT token
                │                           │
                │                           ▼
                │                   Save token to localStorage
                │                           │
                └───────────────────────────┘
                              │
                              ▼
                    All API requests include:
                    Authorization: Bearer <token>
                              │
                ┌─────────────┴─────────────┐
                │                           │
           200 OK                      401 Unauthorized
                │                           │
        Continue normally                   ▼
                                    Remove token
                                    Redirect to login
```

## Component Communication

### Parent → Child (Props)
```
App → Login
  - onLogin: function
  - onSwitchToRegister: function

App → HabitList
  - (no props, uses API directly)

HabitList → Habit
  - habit: object
  - onEdit: function
  - onDelete: function
  - onUpdate: function

HabitList → HabitDetailsModal
  - isOpen: boolean
  - onClose: function
  - onSave: function
  - habit: object | null

Habit → LogDetailsModal
  - isOpen: boolean
  - onClose: function
  - onSave: function
  - log: object | null
  - date: string | null
```

### Child → Parent (Callbacks)
```
Login → App
  - onLogin(userData) → Sets user and isAuthenticated

Habit → HabitList
  - onDelete(habitId) → Removes habit from list
  - onEdit(habit) → Opens edit modal
  - onUpdate() → Refreshes habit list

HabitDetailsModal → HabitList
  - onSave(formData) → Creates/updates habit
  - onClose() → Closes modal

LogDetailsModal → Habit
  - onSave(formData) → Creates/updates log
  - onClose() → Closes modal
```

## File Dependencies

```
index.js
  └─> App.js
      ├─> services/api.js
      ├─> components/Auth/Login.js
      │   └─> services/api.js
      ├─> components/Auth/Register.js
      │   └─> services/api.js
      └─> components/Habits/HabitList.js
          ├─> services/api.js
          ├─> components/Modals/HabitDetailsModal.js
          └─> components/Habits/Habit.js
              ├─> services/api.js
              ├─> utils/streakCalculator.js
              └─> components/Modals/LogDetailsModal.js
```

## Key Design Patterns

### 1. Container/Presentational Pattern
- **Containers**: HabitList, Habit (manage state, API calls)
- **Presentational**: Modals (receive props, render UI)

### 2. Controlled Components
- All forms use controlled inputs
- State managed in component
- onChange handlers update state

### 3. Composition
- Small, focused components
- Composed into larger features
- Reusable modals

### 4. Separation of Concerns
- **Components**: UI and user interaction
- **Services**: API communication
- **Utils**: Business logic (streak calculation)
- **Styles**: CSS in separate file

### 5. Error Boundaries
- Try/catch in async functions
- Error state in components
- User-friendly error messages
