package models

import "time"

// Log represents a habit log entry in the system
type Log struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateLogRequest represents the request body for creating a log
type CreateLogRequest struct {
	HabitID string `json:"habit_id"`
	Notes   string `json:"notes,omitempty"`
}

// UpdateLogRequest represents the request body for updating a log
type UpdateLogRequest struct {
	Notes *string `json:"notes,omitempty"`
}
