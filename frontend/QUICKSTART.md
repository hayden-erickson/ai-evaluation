# 🚀 Quick Start Guide

Get up and running with the Habit Tracker in 3 simple steps!

## Step 1: Start the Backend

Open a terminal in the project root and run:

```bash
go run main.go
```

✅ Backend running on `http://localhost:8080`

## Step 2: Start the Frontend

Open a **new** terminal and run:

```bash
cd frontend
npm run dev
```

✅ Frontend running on `http://localhost:3000`

## Step 3: Open Your Browser

Navigate to: **http://localhost:3000**

## First Time Use

1. **Create Account** → Click "Register"
2. **Create Habit** → Click "+ Add New Habit"
3. **Log Progress** → Click "✓ Log for Today"
4. **Build Streaks** → Keep logging daily! 🔥

## What You'll See

### Login/Register Screen
- Clean, modern authentication
- Pastel gradient background
- Mobile-responsive design

### Dashboard
- All your habits in a card grid
- Streak counters with emojis
- Quick log buttons
- 30-day history view

### Features
- ✅ Create and edit habits
- 📝 Add notes to logs
- ⏱️ Track duration
- 🔥 Build streaks (1 skip day allowed)
- 📊 View 30-day history
- 🗑️ Delete habits and logs

## Need Help?

- **Full Documentation**: See `README.md` in this directory
- **Setup Guide**: See `FRONTEND_SETUP.md` in project root
- **API Docs**: Check the Go handlers in `../handlers/`

## Troubleshooting

**Can't connect?**
- Make sure backend is running on port 8080
- Check browser console (F12) for errors

**CORS errors?**
- Verify `../middleware/security.go` allows port 3000

**Want to reset?**
- Clear browser localStorage: `localStorage.clear()`
- Delete database: Remove `../data/habits.db` and restart backend

---

Happy habit tracking! 🎯

