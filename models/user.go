package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              int64     `json:"id"`
	ProfileImageURL string    `json:"profile_image_url"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	PhoneNumber     string    `json:"phone_number"`
	PasswordHash    string    `json:"-"` // Never expose password hash in JSON
	CreatedAt       time.Time `json:"created_at"`
}

// UserCreateRequest represents the request to create a new user
type UserCreateRequest struct {
	ProfileImageURL string `json:"profile_image_url"`
	Name            string `json:"name" validate:"required"`
	TimeZone        string `json:"time_zone"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password" validate:"required,min=8"`
}

// UserUpdateRequest represents the request to update a user
type UserUpdateRequest struct {
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	Name            *string `json:"name,omitempty"`
	TimeZone        *string `json:"time_zone,omitempty"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
	Password        *string `json:"password,omitempty"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
