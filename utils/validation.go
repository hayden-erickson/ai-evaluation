package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
)

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length requirement
func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return errors.New(fieldName + " must be at least " + string(rune(minLength)) + " characters")
	}
	return nil
}

// ValidatePhoneNumber validates phone number format (basic validation)
func ValidatePhoneNumber(phoneNumber string) error {
	// Remove common separators
	cleaned := strings.ReplaceAll(phoneNumber, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	// Check if it contains only digits and optional + at start
	matched, err := regexp.MatchString(`^\+?\d{10,15}$`, cleaned)
	if err != nil || !matched {
		return errors.New("invalid phone number format")
	}
	return nil
}

// SanitizeString removes potentially harmful characters from input
func SanitizeString(input string) string {
	// Trim whitespace
	return strings.TrimSpace(input)
}
