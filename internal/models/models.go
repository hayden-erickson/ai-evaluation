package models

import "time"

// User represents a single application user.
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // only used for write payloads; not stored
	PasswordHash string `json:"-"`
	ProfileImageURL string `json:"profile_image_url"`
	Name      string    `json:"name"`
	TimeZone  string    `json:"time_zone"`
	Phone     string    `json:"phone_number"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Habit represents a habit owned by a user.
type Habit struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Description string  `json:"description"`
	DurationSeconds *int64   `json:"duration_seconds,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// LogEntry represents a habit log entry.
type LogEntry struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Notes     string    `json:"notes"`
	DurationSeconds *int64   `json:"duration_seconds,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
