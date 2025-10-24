package validator

import (
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// ValidateLogDuration validates that a log's duration matches the habit's duration requirement
// If habit has a duration, the log must also have a duration
// If habit has no duration, the log's duration is optional
func ValidateLogDuration(habit *models.Habit, logDuration *int) error {
	if habit.Duration != nil && logDuration == nil {
		return fmt.Errorf("log must have a duration because the habit requires one")
	}
	
	// Validate duration is non-negative if provided
	if logDuration != nil && *logDuration < 0 {
		return fmt.Errorf("duration must be non-negative")
	}
	
	return nil
}
