# Habit Tracker UI

This project is a React application for tracking habits.

## Prerequisites

- Node.js and npm installed.
- The backend API server must be running.

## Getting Started

1.  **Navigate to the UI directory:**
    ```bash
    cd ui
    ```

2.  **Install dependencies:**
    ```bash
    npm install
    ```

3.  **Run the application:**
    ```bash
    npm start
    ```
    The application will be available at [http://localhost:3000](http://localhost:3000).

## Proxy Setup

To ensure that API requests from the React application are correctly routed to the backend server, a proxy has been configured in `package.json`. All requests to `/users`, `/habits`, and `/logs` will be forwarded to the backend server running on port 8080.
