# Habit Tracker UI

This is a React frontend for the Habit Tracker application.

## Running the Application

To run the application, you need to have both the Go backend and the React frontend running concurrently.

### Backend

1.  Navigate to the root directory of the project.
2.  Run the backend server:
    ```sh
    go run main.go
    ```

The backend will be running on `http://localhost:8080`.

### Frontend

1.  Navigate to the `ui` directory.
2.  Install the dependencies:
    ```sh
    npm install
    ```
3.  Start the frontend development server:
    ```sh
    npm start
    ```

The frontend will be running on `http://localhost:3000`.

## Proxy Setup

To allow the React frontend to communicate with the Go backend, a proxy has been set up in `package.json`. All API requests from the frontend will be forwarded to `http://localhost:8080`.
