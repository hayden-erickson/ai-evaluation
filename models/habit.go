package models

import (
	"errors"
	"time"
)

// Habit represents a habit in the system
type Habit struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateHabitRequest represents the request to create a new habit
type CreateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateHabitRequest represents the request to update a habit
type UpdateHabitRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Validate validates the CreateHabitRequest
func (r *CreateHabitRequest) Validate() error {
	// Validate required fields
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// Validate validates the UpdateHabitRequest
func (r *UpdateHabitRequest) Validate() error {
	// At least one field must be provided
	if r.Name == nil && r.Description == nil {
		return errors.New("at least one field must be provided for update")
	}
	// Validate name if provided
	if r.Name != nil && *r.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}
