package models

import "time"

// User represents a user in the database
type User struct {
	ID              int64     `json:"id"`
	ProfileImageURL string    `json:"profile_image_url"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	PhoneNumber     string    `json:"phone_number"`
	CreatedAt       time.Time `json:"created_at"`
}

// Habit represents a habit in the database
type Habit struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    *int      `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Log represents a log entry for a habit
type Log struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Notes     string    `json:"notes"`
	Duration  *int      `json:"duration,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
