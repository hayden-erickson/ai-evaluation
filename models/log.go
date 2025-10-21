package models

import "time"

// Log represents a log entry for a habit
type Log struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateLogRequest represents the request body for creating a log
type CreateLogRequest struct {
	HabitID int64  `json:"habit_id"`
	Notes   string `json:"notes,omitempty"`
}

// UpdateLogRequest represents the request body for updating a log
type UpdateLogRequest struct {
	Notes *string `json:"notes,omitempty"`
}
