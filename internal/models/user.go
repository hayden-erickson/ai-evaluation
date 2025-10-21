package models

import "time"

// User represents a user in the system
type User struct {
	ID              string    `json:"id"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	PhoneNumber     string    `json:"phone_number,omitempty"`
	PasswordHash    string    `json:"-"` // Never expose password hash in JSON
	CreatedAt       time.Time `json:"created_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Name            string `json:"name"`
	TimeZone        string `json:"time_zone"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	Password        string `json:"password"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	Name            *string `json:"name,omitempty"`
	TimeZone        *string `json:"time_zone,omitempty"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// LoginResponse represents the login response with JWT token
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

