package models

import "time"

// User represents a user in the system
type User struct {
	ID              string    `json:"id" db:"id"`
	ProfileImageURL string    `json:"profile_image_url" db:"profile_image_url"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	TimeZone        string    `json:"time_zone" db:"time_zone"`
	PhoneNumber     string    `json:"phone_number" db:"phone_number"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=100"`
	ProfileImageURL string `json:"profile_image_url" validate:"omitempty,url"`
	TimeZone        string `json:"time_zone" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"omitempty"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name            *string `json:"name" validate:"omitempty,min=1,max=100"`
	ProfileImageURL *string `json:"profile_image_url" validate:"omitempty,url"`
	TimeZone        *string `json:"time_zone" validate:"omitempty"`
	PhoneNumber     *string `json:"phone_number" validate:"omitempty"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
