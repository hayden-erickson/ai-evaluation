package models

import (
	"time"
)

// Log represents a log entry for a habit
type Log struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Notes     string    `json:"notes"`
	Duration  *int64    `json:"duration,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// LogCreateRequest represents the request to create a new log
type LogCreateRequest struct {
	HabitID  int64  `json:"habit_id" validate:"required"`
	Notes    string `json:"notes"`
	Duration *int64 `json:"duration"`
}

// LogUpdateRequest represents the request to update a log
type LogUpdateRequest struct {
	Notes    *string `json:"notes,omitempty"`
	Duration *int64  `json:"duration,omitempty"`
}
