# Habit Tracker App

A full-stack application for tracking user habits, featuring a Go backend and a React Native mobile app.

## Features

- **JWT-based authentication** - Secure token-based authentication.
- **Modern UI** - A minimal, modern user interface built with React Native.
- **Streak Tracking** - Calculates and displays habit streaks, allowing for one skipped day.
- **Modular Architecture** - A clean, modular architecture for both the backend and frontend.
- **SQLite Database** - Lightweight database for local development.

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Node.js (LTS version recommended)
- npm or yarn
- A configured React Native development environment (see the [official guide](https://reactnative.dev/docs/set-up-your-environment)).

## Getting Started

Follow these steps to get the application running locally.

### 1. Backend Setup (Go API)

First, set up and run the Go backend server.

```bash
# 1. Install Go dependencies
go mod download

# 2. (Optional) Set environment variables for port, JWT secret, and DB path
export PORT=8080
export JWT_SECRET=your-super-secret-key
export DB_PATH=habits.db

# 3. Run the server
go run main.go
```

The API server will start on port `8080`.

### 2. Frontend Setup (React Native App)

With the backend running, set up and run the mobile app in a separate terminal.

```bash
# 1. Install Node.js dependencies
npm install
# or
yarn install

# 2. Start the Metro bundler
npm start
# or
yarn start

# 3. Run the app on your desired platform
# For Android
npm run android
# or
yarn android

# For iOS (make sure to install pods first)
bundle install && bundle exec pod install
npm run ios
# or
yarn ios
```

## API Endpoints

(The API endpoint documentation remains the same as in the previous version of this file.)

## License

MIT
