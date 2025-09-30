package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	ProfileImageURL string    `json:"profile_image_url"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	Phone           string    `json:"phone"`
	CreatedAt       time.Time `json:"created_at"`
	GoogleID        string    `gorm:"unique;not null" json:"google_id"`
	Email           string    `gorm:"unique;not null" json:"email"`
	Habits          []Habit   `json:"habits,omitempty"`
}

// Habit represents a habit that a user wants to track
type Habit struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Frequency   int       `json:"frequency"`
	CreatedAt   time.Time `json:"created_at"`
	User        User      `json:"-"`
	Logs        []Log     `json:"logs,omitempty"`
	Tags        []Tag     `json:"tags,omitempty"`
}

// Log represents a habit tracking entry
type Log struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	HabitID   uint      `json:"habit_id"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	Habit     Habit     `json:"-"`
}

// Tag represents a label attached to a habit
type Tag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	HabitID   uint      `json:"habit_id"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Habit     Habit     `json:"-"`
}

// HabitWithStreak represents a habit with its current streak information
type HabitWithStreak struct {
	Habit
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	LastLogDate   *time.Time `json:"last_log_date,omitempty"`
}
