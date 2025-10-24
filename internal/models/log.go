package models

import "time"

// Log represents a log entry for a habit
type Log struct {
	ID              string    `json:"id"`
	HabitID         string    `json:"habit_id"`
	Notes           string    `json:"notes,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateLogRequest represents the request to create a new log
type CreateLogRequest struct {
	Notes           string `json:"notes,omitempty"`
	DurationSeconds *int   `json:"duration_seconds,omitempty"`
}

// UpdateLogRequest represents the request to update a log
type UpdateLogRequest struct {
	Notes           *string `json:"notes,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}
