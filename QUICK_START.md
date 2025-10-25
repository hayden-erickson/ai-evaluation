# 🚀 Quick Start Guide - Habit Tracker

Get up and running in 2 minutes!

## Prerequisites Check

Before starting, verify you have:

- ✅ **Go** installed - Run `go version` (need 1.16+)
- ✅ **Node.js** installed - Run `node --version` (need 14+)
- ✅ **npm** installed - Run `npm --version`

If any are missing, install them first.

## Option 1: Using Batch Scripts (Windows - Easiest)

### Step 1: Start Backend
Double-click `start-backend.bat` or run in terminal:
```bash
start-backend.bat
```

### Step 2: Start Frontend
Open a new terminal and double-click `start-frontend.bat` or run:
```bash
start-frontend.bat
```

That's it! The app will open in your browser automatically.

## Option 2: Manual Start (All Platforms)

### Terminal 1 - Backend
```bash
# From project root
go run main.go
```

Wait for: `Server starting on port 8080`

### Terminal 2 - Frontend
```bash
# From project root
cd frontend
npm install    # First time only
npm start
```

Wait for: `Compiled successfully!`

## First Time Setup

1. **Browser opens** to `http://localhost:3000`
2. **Click "Register"** to create an account
3. **Fill in the form:**
   - Name: Your name
   - Phone: Any phone number (e.g., 1234567890)
   - Time Zone: America/New_York (or your timezone)
   - Password: At least 8 characters
4. **Click "Register"** - You're in!

## Create Your First Habit

1. **Click "➕ Add New Habit"**
2. **Enter details:**
   - Name: "Morning Exercise" (or whatever you want)
   - Description: Optional
3. **Click "Save"**
4. **Click "✅ Log Today"** to record your first completion
5. **Watch your streak grow!** 🔥

## What You'll See

```
┌─────────────────────────────────────────────────────┐
│  🎯 Habit Tracker                                   │
│  Welcome back, [Your Name]!                         │
│  [Logout]                                           │
└─────────────────────────────────────────────────────┘

[➕ Add New Habit]

┌─────────────────────────────────────────────────────┐
│  Morning Exercise                    [✏️] [🗑️]      │
│  Stay healthy and energized                         │
│                                                      │
│  🔥 1 day streak                                    │
│                                                      │
│  [✅ Log Today]                                     │
│                                                      │
│  Last 30 Days:                                      │
│  ┌─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┐                  │
│  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │✓│                  │
│  └─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┘                  │
│  ┌─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─┐                  │
│  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │                  │
│  └─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┘                  │
└─────────────────────────────────────────────────────┘
```

## Common Issues & Solutions

### ❌ "Port 8080 already in use"
```bash
# Kill the process using port 8080
# Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Or change the port:
set PORT=8081
go run main.go
```

### ❌ "Port 3000 already in use"
```bash
# The terminal will ask if you want to use another port
# Type 'y' and press Enter
```

### ❌ "Failed to load habits"
- Make sure backend is running (check Terminal 1)
- Visit http://localhost:8080/health - should show "OK"
- Check for errors in backend terminal

### ❌ "npm: command not found"
- Install Node.js from https://nodejs.org/
- Restart your terminal after installation

### ❌ CORS errors in browser console
- Make sure frontend is on port 3000
- Make sure backend is on port 8080
- Restart both servers

## Understanding Streaks

### How Streaks Work
- ✅ Log a habit → Streak increases
- ⚠️ Skip 1 day → Streak stays (you get one free skip!)
- ❌ Skip 2+ days → Streak resets to 0

### Example
```
Monday:    Log ✅  → Streak: 1
Tuesday:   Skip ⚠️ → Streak: 1 (skip allowed)
Wednesday: Log ✅  → Streak: 2
Thursday:  Skip ❌ → Streak: 0 (too many skips)
```

## Tips for Success

### 🎯 Best Practices
1. **Log daily** - Keep your streaks alive
2. **Add notes** - Reflect on your progress
3. **Start small** - Don't create too many habits at once
4. **Be consistent** - Use the 1-day skip wisely

### 📱 Daily Workflow
1. Open app each morning
2. Click "✅ Log Today" for completed habits
3. Add notes about how it went
4. Check your streak progress
5. Plan for tomorrow

### 🎨 UI Tips
- **Green squares** = Days you logged
- **Highlighted border** = Today
- **Click any day** = View/edit that day's log
- **Fire emoji** = Your streak count

## Features Overview

### ✅ What You Can Do
- Create unlimited habits
- Track daily completions
- Build and maintain streaks
- Add notes to each log
- View 30-day history
- Edit past logs
- Delete habits/logs

### 🔒 Security
- Passwords are encrypted
- Secure JWT authentication
- Auto-logout on session expiry
- Input validation

## Keyboard Shortcuts

- **Esc** - Close any modal
- **Enter** - Submit forms (when focused)
- **Tab** - Navigate between fields

## Mobile Usage

The app works great on mobile browsers:
1. Open `http://localhost:3000` on your phone
2. Make sure phone is on same network as computer
3. Or use your computer's IP instead of localhost

## Next Steps

### 📚 Learn More
- Read `FRONTEND_SETUP.md` for detailed setup
- Read `UI_IMPLEMENTATION_SUMMARY.md` for features
- Read `COMPONENT_HIERARCHY.md` for architecture

### 🛠️ Customize
- Edit `frontend/src/index.css` to change colors
- Modify components in `frontend/src/components/`
- Add new features to the API

### 🚀 Deploy
- Build for production: `npm run build`
- Deploy backend to a server
- Serve frontend build folder
- Update CORS settings for production domain

## Stopping the Application

### Stop Backend
- Go to Terminal 1
- Press `Ctrl+C`

### Stop Frontend
- Go to Terminal 2
- Press `Ctrl+C`

## Restarting

Just run the same start commands again:
- Backend: `go run main.go` or `start-backend.bat`
- Frontend: `npm start` or `start-frontend.bat`

Your data is saved in `habits.db` and will persist between restarts.

## Getting Help

### Check Logs
- **Backend logs**: Terminal 1 shows all API requests
- **Frontend logs**: Browser DevTools → Console (F12)

### Verify Setup
1. Backend health: http://localhost:8080/health
2. Frontend running: http://localhost:3000
3. Database exists: Check for `habits.db` file

### Reset Everything
```bash
# Delete database (loses all data!)
rm habits.db

# Clear frontend cache
cd frontend
rm -rf node_modules package-lock.json
npm install
```

## Success Checklist

- ✅ Backend running on port 8080
- ✅ Frontend running on port 3000
- ✅ Browser opened automatically
- ✅ Can register a new account
- ✅ Can login successfully
- ✅ Can create a habit
- ✅ Can log a habit completion
- ✅ Can see streak count
- ✅ Can view 30-day calendar

If all checked, you're ready to start tracking habits! 🎉

## Support

If you encounter issues:
1. Check the "Common Issues" section above
2. Review terminal output for errors
3. Check browser console (F12) for frontend errors
4. Ensure all prerequisites are installed
5. Try restarting both servers

Happy habit tracking! 🎯🔥
