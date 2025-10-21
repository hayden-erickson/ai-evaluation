package models

import "time"

// User represents a user in the system
type User struct {
	ID              int64     `json:"id"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	PhoneNumber     string    `json:"phone_number,omitempty"`
	PasswordHash    string    `json:"-"` // Never expose password hash in JSON
	CreatedAt       time.Time `json:"created_at"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Name            string `json:"name"`
	TimeZone        string `json:"time_zone"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	Password        string `json:"password"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	Name            *string `json:"name,omitempty"`
	TimeZone        *string `json:"time_zone,omitempty"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
	Password        *string `json:"password,omitempty"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
