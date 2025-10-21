package models

import (
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID              int64     `json:"id"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"time_zone"`
	PhoneNumber     string    `json:"phone_number"`
	PasswordHash    string    `json:"-"` // Never expose password hash in JSON
	CreatedAt       time.Time `json:"created_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Name            string `json:"name"`
	TimeZone        string `json:"time_zone"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	Name            *string `json:"name,omitempty"`
	TimeZone        *string `json:"time_zone,omitempty"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
	Password        *string `json:"password,omitempty"`
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	// Validate required fields
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.TimeZone == "" {
		return errors.New("time_zone is required")
	}
	if r.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	// Validate password strength
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	// At least one field must be provided
	if r.ProfileImageURL == nil && r.Name == nil && r.TimeZone == nil && r.PhoneNumber == nil && r.Password == nil {
		return errors.New("at least one field must be provided for update")
	}
	// Validate password strength if provided
	if r.Password != nil && len(*r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// LoginRequest represents the request to login
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// Validate validates the LoginRequest
func (r *LoginRequest) Validate() error {
	if r.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// LoginResponse represents the response to a successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
