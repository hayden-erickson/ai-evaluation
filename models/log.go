package models

import (
	"time"
)

// Log represents a log entry for a habit
type Log struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateLogRequest represents the request to create a new log
type CreateLogRequest struct {
	Notes string `json:"notes,omitempty"`
}

// UpdateLogRequest represents the request to update a log
type UpdateLogRequest struct {
	Notes *string `json:"notes,omitempty"`
}

// Validate validates the CreateLogRequest
func (r *CreateLogRequest) Validate() error {
	// No required fields for log creation
	return nil
}

// Validate validates the UpdateLogRequest
func (r *UpdateLogRequest) Validate() error {
	// Notes field is optional
	return nil
}
