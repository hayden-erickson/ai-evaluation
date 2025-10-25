package models

import (
	"errors"
	"time"
)

// Log represents a log entry for a habit
type Log struct {
	ID              int64     `json:"id"`
	HabitID         int64     `json:"habit_id"`
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

// Validate validates the CreateLogRequest
func (r *CreateLogRequest) Validate() error {
	// No required fields for log creation
	return nil
}

// Validate validates the UpdateLogRequest
func (r *UpdateLogRequest) Validate() error {
	// Validate duration if provided
	if r.DurationSeconds != nil && *r.DurationSeconds < 0 {
		return errors.New("duration_seconds cannot be negative")
	}
	return nil
}
