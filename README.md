# Habit Streaks App (React Native + Go API)

Modern, minimal habit tracker focused on streaks. Users log daily habits; streaks allow skipping one day before resetting. The React Native app communicates with the Go API in this repository.

## Features

- JWT auth (login/register) with secure token storage
- Create, edit, delete habits
- Create, edit, delete logs per habit
- Streak calculation with one allowed skip day
- Pastel, rounded UI with accessible, responsive layout

## Prerequisites

- Node 20+, Android Studio and/or Xcode (see React Native environment setup)
- Go 1.21+

## Run the API (backend)

```bash
# From repo root
go mod download

# Optional env vars
# set PORT=8080 (Windows PowerShell: $env:PORT="8080")
# set JWT_SECRET=change-me
# set DB_PATH=habits.db

go run main.go
```

The API listens on `http://localhost:8080` (Android emulator uses `http://10.0.2.2:8080`).

## Run the App (frontend)

```bash
# Install JS deps
npm install

# Start Metro
npm start

# In another terminal, run a platform
npm run android
# or
npm run ios
```

The app will connect to the API automatically. On Android emulator, it targets `10.0.2.2` per default configuration.

## Usage

1. Register with name, phone number, password, and time zone
2. You will be logged in automatically
3. Create a habit (optionally specify duration in seconds)
4. Tap “New Log” to add a log for today; tap a day in the grid to edit an existing log
5. Streak increases for consecutive days; one missed day is allowed before reset

## API Reference (Summary)

- `POST /users/register` – Register
- `POST /users/login` – Login (returns JWT)
- `GET /habits` – List habits
- `POST /habits` – Create habit
- `PUT /habits/{id}` – Update habit
- `DELETE /habits/{id}` – Delete habit
- `GET /habits/{habit_id}/logs` – List logs for habit
- `POST /habits/{habit_id}/logs` – Create log
- `PUT /logs/{id}` – Update log
- `DELETE /logs/{id}` – Delete log

All protected endpoints require `Authorization: Bearer <token>`.

## Security & Validation

- Tokens stored in AsyncStorage; Authorization header sent on each request
- Client-side validation for required fields and basic constraints
- Server enforces RBAC: users can only access their own resources

## Troubleshooting

- Android cannot reach `localhost`: the app uses `10.0.2.2` automatically
- API port different than 8080: set `PORT` for the Go server and adjust the app if needed

## License

MIT
