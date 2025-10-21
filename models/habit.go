package models

import "time"

// Habit represents a habit in the system
type Habit struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateHabitRequest represents the request body for creating a habit
type CreateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateHabitRequest represents the request body for updating a habit
type UpdateHabitRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
