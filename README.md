# Habit Tracker React Native App

This is a React Native application for tracking habits, built to interact with a Go backend API.

## Features

- User authentication (Sign up and Login)
- Create, Read, Update, and Delete habits
- Log daily progress for each habit
- View habit streaks

## Getting Started

### Prerequisites

- Node.js and npm
- React Native CLI
- Android Studio or Xcode for running on an emulator/simulator or a physical device
- Go (for running the backend API)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-directory>
    ```

2.  **Install frontend dependencies:**
    ```bash
    npm install
    ```

3.  **Install backend dependencies:**
    ```bash
    go mod tidy
    ```

### Running the Application

1.  **Start the backend server:**
    Open a terminal in the project root and run:
    ```bash
    go run main.go
    ```
    The server should start on `http://localhost:8080`.

2.  **Start the React Native Metro bundler:**
    Open another terminal in the project root and run:
    ```bash
    npx react-native start
    ```

3.  **Run on Android or iOS:**
    -   **For Android:**
        ```bash
        npx react-native run-android
        ```
    -   **For iOS:**
        ```bash
        npx react-native run-ios
        ```

## Project Structure

```
.
├── android/
├── ios/
├── src/
│   ├── api/
│   │   ├── api.ts
│   │   ├── habitService.ts
│   │   └── logService.ts
│   ├── components/
│   │   ├── HabitDetailsModal.tsx
│   │   ├── HabitList.tsx
│   │   ├── HabitListItem.tsx
│   │   └── LogDetailsModal.tsx
│   ├── contexts/
│   │   └── AuthContext.tsx
│   ├── navigation/
│   │   └── AppNavigator.tsx
│   ├── screens/
│   │   ├── HabitListScreen.tsx
│   │   ├── LoginScreen.tsx
│   │   ├── SignUpScreen.tsx
│   │   └── SplashScreen.tsx
│   ├── styles/
│   │   ├── colors.ts
│   │   └── theme.ts
│   └── types/
│       └── types.ts
├── App.tsx
└── ... (other project files)
```
