# API Documentation

## Base URL

**Local Development:** `http://localhost:8080`  
**Production:** `https://api.yourdomain.com`

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens expire after 24 hours.

## Response Format

### Success Response
```json
{
  "id": "uuid",
  "field": "value",
  ...
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

## HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Authentication Endpoints

### Register New User

**POST** `/api/v1/auth/register`

Creates a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg"
}
```

**Validation Rules:**
- `name`: required, 1-100 characters
- `email`: required, valid email format
- `password`: required, minimum 8 characters
- `time_zone`: required
- `phone_number`: optional
- `profile_image_url`: optional, valid URL

**Success Response (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg",
  "created_at": "2025-10-20T12:00:00Z"
}
```

**Error Response (400):**
```json
{
  "error": "Invalid request body"
}
```

---

### Login

**POST** `/api/v1/auth/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "time_zone": "America/New_York",
    "phone_number": "+1234567890",
    "profile_image_url": "https://example.com/image.jpg",
    "created_at": "2025-10-20T12:00:00Z"
  }
}
```

**Error Response (401):**
```json
{
  "error": "Invalid credentials"
}
```

---

## User Endpoints

All user endpoints require authentication.

### Get User

**GET** `/api/v1/users/{id}`

Retrieves user information. Users can only access their own data.

**Path Parameters:**
- `id` - User UUID

**Success Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "time_zone": "America/New_York",
  "phone_number": "+1234567890",
  "profile_image_url": "https://example.com/image.jpg",
  "created_at": "2025-10-20T12:00:00Z"
}
```

---

### Update User

**PUT** `/api/v1/users/{id}`

Updates user information. Users can only update their own data.

**Path Parameters:**
- `id` - User UUID

**Request Body:**
All fields are optional. Only provided fields will be updated.

```json
{
  "name": "Jane Doe",
  "profile_image_url": "https://example.com/newimage.jpg",
  "time_zone": "America/Los_Angeles",
  "phone_number": "+9876543210"
}
```

**Success Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Jane Doe",
  "email": "john@example.com",
  "time_zone": "America/Los_Angeles",
  "phone_number": "+9876543210",
  "profile_image_url": "https://example.com/newimage.jpg",
  "created_at": "2025-10-20T12:00:00Z"
}
```

---

### Delete User

**DELETE** `/api/v1/users/{id}`

Deletes a user account. Users can only delete their own account.

**Path Parameters:**
- `id` - User UUID

**Success Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

---

## Habit Endpoints

All habit endpoints require authentication.

### Create Habit

**POST** `/api/v1/habits`

Creates a new habit for the authenticated user.

**Request Body:**
```json
{
  "name": "Morning Exercise",
  "description": "30 minutes of cardio every morning"
}
```

**Validation Rules:**
- `name`: required, 1-100 characters
- `description`: optional, max 500 characters

**Success Response (201):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Morning Exercise",
  "description": "30 minutes of cardio every morning",
  "created_at": "2025-10-20T12:00:00Z"
}
```

---

### List User's Habits

**GET** `/api/v1/habits`

Retrieves all habits for the authenticated user, ordered by creation date (newest first).

**Success Response (200):**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Morning Exercise",
    "description": "30 minutes of cardio",
    "created_at": "2025-10-20T12:00:00Z"
  },
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Read Daily",
    "description": "Read for 20 minutes",
    "created_at": "2025-10-19T12:00:00Z"
  }
]
```

---

### Get Habit

**GET** `/api/v1/habits/{id}`

Retrieves a specific habit. Users can only access their own habits.

**Path Parameters:**
- `id` - Habit UUID

**Success Response (200):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Morning Exercise",
  "description": "30 minutes of cardio",
  "created_at": "2025-10-20T12:00:00Z"
}
```

---

### Update Habit

**PUT** `/api/v1/habits/{id}`

Updates a habit. Users can only update their own habits.

**Path Parameters:**
- `id` - Habit UUID

**Request Body:**
All fields are optional. Only provided fields will be updated.

```json
{
  "name": "Morning Workout",
  "description": "45 minutes of mixed cardio and strength training"
}
```

**Success Response (200):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Morning Workout",
  "description": "45 minutes of mixed cardio and strength training",
  "created_at": "2025-10-20T12:00:00Z"
}
```

---

### Delete Habit

**DELETE** `/api/v1/habits/{id}`

Deletes a habit and all associated log entries. Users can only delete their own habits.

**Path Parameters:**
- `id` - Habit UUID

**Success Response (200):**
```json
{
  "message": "Habit deleted successfully"
}
```

---

## Log Endpoints

All log endpoints require authentication.

### Create Log Entry

**POST** `/api/v1/logs`

Creates a new log entry for a habit.

**Request Body:**
```json
{
  "habit_id": "660e8400-e29b-41d4-a716-446655440000",
  "notes": "Completed 45 minutes today, felt great!"
}
```

**Validation Rules:**
- `habit_id`: required, valid UUID
- `notes`: optional, max 1000 characters

**Success Response (201):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "habit_id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-10-20T12:00:00Z",
  "notes": "Completed 45 minutes today, felt great!"
}
```

---

### List Logs for Habit

**GET** `/api/v1/logs?habit_id={habit_id}`

Retrieves all log entries for a specific habit, ordered by creation date (newest first).

**Query Parameters:**
- `habit_id` - Habit UUID (required)

**Success Response (200):**
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "habit_id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-10-20T12:00:00Z",
    "notes": "Completed 45 minutes"
  },
  {
    "id": "990e8400-e29b-41d4-a716-446655440000",
    "habit_id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-10-19T12:00:00Z",
    "notes": "Completed 30 minutes"
  }
]
```

---

### Get Log Entry

**GET** `/api/v1/logs/{id}`

Retrieves a specific log entry. Users can only access logs for their own habits.

**Path Parameters:**
- `id` - Log UUID

**Success Response (200):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "habit_id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-10-20T12:00:00Z",
  "notes": "Completed 45 minutes today"
}
```

---

### Update Log Entry

**PUT** `/api/v1/logs/{id}`

Updates a log entry. Users can only update logs for their own habits.

**Path Parameters:**
- `id` - Log UUID

**Request Body:**
```json
{
  "notes": "Completed 60 minutes today, new personal record!"
}
```

**Success Response (200):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "habit_id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-10-20T12:00:00Z",
  "notes": "Completed 60 minutes today, new personal record!"
}
```

---

### Delete Log Entry

**DELETE** `/api/v1/logs/{id}`

Deletes a log entry. Users can only delete logs for their own habits.

**Path Parameters:**
- `id` - Log UUID

**Success Response (200):**
```json
{
  "message": "Log deleted successfully"
}
```

---

## Example cURL Commands

### Register and Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "time_zone": "America/New_York"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'

# Save the token
TOKEN="your-token-here"
```

### Working with Habits
```bash
# Create habit
curl -X POST http://localhost:8080/api/v1/habits \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Exercise", "description": "Daily workout"}'

# List habits
curl http://localhost:8080/api/v1/habits \
  -H "Authorization: Bearer $TOKEN"

# Get specific habit
curl http://localhost:8080/api/v1/habits/{habit-id} \
  -H "Authorization: Bearer $TOKEN"

# Update habit
curl -X PUT http://localhost:8080/api/v1/habits/{habit-id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Morning Exercise"}'

# Delete habit
curl -X DELETE http://localhost:8080/api/v1/habits/{habit-id} \
  -H "Authorization: Bearer $TOKEN"
```

### Working with Logs
```bash
# Create log
curl -X POST http://localhost:8080/api/v1/logs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"habit_id": "{habit-id}", "notes": "Completed today"}'

# List logs for habit
curl "http://localhost:8080/api/v1/logs?habit_id={habit-id}" \
  -H "Authorization: Bearer $TOKEN"

# Update log
curl -X PUT http://localhost:8080/api/v1/logs/{log-id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Updated notes"}'

# Delete log
curl -X DELETE http://localhost:8080/api/v1/logs/{log-id} \
  -H "Authorization: Bearer $TOKEN"
```

---

## Rate Limiting

Currently no rate limiting is implemented. For production deployment, consider adding rate limiting middleware or using AWS WAF.

## Pagination

Currently no pagination is implemented for list endpoints. For large datasets, consider adding pagination with query parameters like `?page=1&limit=20`.

## Versioning

The API is versioned via the URL path (`/api/v1/`). Breaking changes will be introduced in new versions (`/api/v2/`).
