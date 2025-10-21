# Go REST API with SQLite

This is a RESTful API built with Go, using the standard library for HTTP handling and SQLite for the database.

## Features

- User, Habit, and Log CRUD operations
- JWT-based authentication
- Modular architecture (repository and service layers)
- Dependency injection
- SQL migrations

## Prerequisites

- Go (version 1.18 or higher)

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Run the application:**
   ```sh
   go run cmd/api/main.go
   ```
   The server will start on `http://localhost:8080`.

## API Endpoints

All endpoints are prefixed with `/`.

### Users

- `POST /users`: Create a new user
- `GET /users?id={id}`: Get a user by ID
- `PUT /users`: Update a user
- `DELETE /users?id={id}`: Delete a user

### Habits (Protected)

- `POST /habits`: Create a new habit
- `GET /habits?id={id}`: Get a habit by ID
- `GET /habits`: Get all habits for the authenticated user
- `PUT /habits`: Update a habit
- `DELETE /habits?id={id}`: Delete a habit

### Logs (Protected)

- `POST /logs`: Create a new log
- `GET /logs?id={id}`: Get a log by ID
- `GET /logs?habit_id={habit_id}`: Get all logs for a habit
- `PUT /logs`: Update a log
- `DELETE /logs?id={id}`: Delete a log

## Authentication

To access protected endpoints, you need to include a JWT in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```
