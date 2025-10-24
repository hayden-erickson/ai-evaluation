package models

import "time"

// Habit represents a habit tracked by a user
type Habit struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateHabitRequest represents the request to create a new habit
type CreateHabitRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	DurationSeconds *int   `json:"duration_seconds,omitempty"`
}

// UpdateHabitRequest represents the request to update a habit
type UpdateHabitRequest struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}
