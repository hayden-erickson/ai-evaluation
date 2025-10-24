package models

import "time"

// Log represents a log entry for a habit
type Log struct {
	ID      int64  `json:"id"`
	HabitID int64  `json:"habit_id"`
	Notes   string `json:"notes,omitempty"`
	// DurationSeconds is optional unless the associated habit has a duration
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateLogRequest represents the request body for creating a log
type CreateLogRequest struct {
	HabitID int64  `json:"habit_id"`
	Notes   string `json:"notes,omitempty"`
	// Optional unless the habit specifies a duration; must be > 0 seconds when provided
	DurationSeconds *int `json:"duration_seconds,omitempty"`
}

// UpdateLogRequest represents the request body for updating a log
type UpdateLogRequest struct {
	Notes *string `json:"notes,omitempty"`
	// Optional; when provided, must be > 0 seconds
	DurationSeconds *int `json:"duration_seconds,omitempty"`
}
