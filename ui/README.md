# Habit Streaks UI

A modern minimal React (Vite + TypeScript) UI to interact with the Habits API in this repository.

## Prerequisites
- Node 18+
- The Go API running locally on port 8080 (or set `VITE_API_BASE`)

## Run locally
```bash
# in repo root, start API (uses PORT=8080 by default)
go run main.go

# in another terminal, start UI
cd ui
# optional: set API base URL
# PowerShell
$env:VITE_API_BASE="http://localhost:8080"
# bash/zsh
# export VITE_API_BASE="http://localhost:8080"
npm run dev
```

Open `http://localhost:5173`.

## Environment
- `VITE_API_BASE` (default `http://localhost:8080`)
- Allow the UI origin in the API via `ALLOWED_ORIGINS`. For local dev it defaults to `http://localhost:5173`.

## Features
- Login via phone number + password; token stored in `localStorage`
- View habits, create/edit/delete habits
- Add logs and edit logs
- 30-day streak grid with one-day skip rule
- Minimal, pastel theme; accessible form validation and error messages

## Notes
- Ensure a user exists in the API to log in, or register via `POST /users/register`.
