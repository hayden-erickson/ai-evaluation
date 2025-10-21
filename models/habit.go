package models

import (
	"time"
)

// Habit represents a habit in the system
type Habit struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// HabitCreateRequest represents the request to create a new habit
type HabitCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// HabitUpdateRequest represents the request to update a habit
type HabitUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
