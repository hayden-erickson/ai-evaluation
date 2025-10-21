package validator

import (
	"fmt"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// ValidateCreateUserRequest validates user creation request
func ValidateCreateUserRequest(req *models.CreateUserRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}

	// Validate time zone
	if strings.TrimSpace(req.TimeZone) == "" {
		return fmt.Errorf("time zone is required")
	}

	// Validate password
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}

// ValidateCreateHabitRequest validates habit creation request
func ValidateCreateHabitRequest(req *models.CreateHabitRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

// ValidateCreateLogRequest validates log creation request
func ValidateCreateLogRequest(req *models.CreateLogRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate habit ID
	if strings.TrimSpace(req.HabitID) == "" {
		return fmt.Errorf("habit ID is required")
	}

	return nil
}
