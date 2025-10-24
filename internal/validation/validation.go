package validation

import (
	"errors"
	"regexp"
)

var (
	emailRx = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

func ValidateEmail(s string) error {
	if len(s) < 3 || len(s) > 254 || !emailRx.MatchString(s) {
		return errors.New("invalid email")
	}
	return nil
}

func ValidatePassword(s string) error {
	if len(s) < 8 { return errors.New("password must be at least 8 characters") }
	return nil
}

func ValidateNonEmpty(s string, field string) error {
	if len(s) == 0 { return errors.New(field + " is required") }
	return nil
}

func ValidateMaxLen(s string, n int, field string) error {
	if len(s) > n { return errors.New(field + " too long") }
	return nil
}
