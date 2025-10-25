package models

import (
	"errors"
	"time"
)

// Habit represents a habit in the system
type Habit struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
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

// Validate validates the CreateHabitRequest
func (r *CreateHabitRequest) Validate() error {
	// Validate required fields
	if r.Name == "" {
		return errors.New("name is required")
	}
	// Validate duration if provided
	if r.DurationSeconds != nil && *r.DurationSeconds < 0 {
		return errors.New("duration_seconds cannot be negative")
	}
	return nil
}

// Validate validates the UpdateHabitRequest
func (r *UpdateHabitRequest) Validate() error {
	// At least one field must be provided
	if r.Name == nil && r.Description == nil && r.DurationSeconds == nil {
		return errors.New("at least one field must be provided for update")
	}
	// Validate name if provided
	if r.Name != nil && *r.Name == "" {
		return errors.New("name cannot be empty")
	}
	// Validate duration if provided
	if r.DurationSeconds != nil && *r.DurationSeconds < 0 {
		return errors.New("duration_seconds cannot be negative")
	}
	return nil
}
