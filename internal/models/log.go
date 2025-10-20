package models

import "time"

// Log represents a log entry for a habit
type Log struct {
	ID        string    `json:"id" db:"id"`
	HabitID   string    `json:"habit_id" db:"habit_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Notes     string    `json:"notes" db:"notes"`
}

// CreateLogRequest represents the request body for creating a log entry
type CreateLogRequest struct {
	HabitID string `json:"habit_id" validate:"required,uuid"`
	Notes   string `json:"notes" validate:"omitempty,max=1000"`
}

// UpdateLogRequest represents the request body for updating a log entry
type UpdateLogRequest struct {
	Notes *string `json:"notes" validate:"omitempty,max=1000"`
}
