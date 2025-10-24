package service

import (
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// LogService defines the interface for log business logic
type LogService interface {
	Create(userID int64, req *models.CreateLogRequest) (*models.Log, error)
	GetByID(id, userID int64) (*models.Log, error)
	GetByHabitID(habitID, userID int64) ([]*models.Log, error)
	Update(id, userID int64, req *models.UpdateLogRequest) (*models.Log, error)
	Delete(id, userID int64) error
}

type logService struct {
	logRepo   repository.LogRepository
	habitRepo repository.HabitRepository
}

// NewLogService creates a new LogService instance
func NewLogService(logRepo repository.LogRepository, habitRepo repository.HabitRepository) LogService {
	return &logService{
		logRepo:   logRepo,
		habitRepo: habitRepo,
	}
}

// Create creates a new log with validation
func (s *logService) Create(userID int64, req *models.CreateLogRequest) (*models.Log, error) {
	// Verify habit exists and belongs to user
	habit, err := s.habitRepo.GetByID(req.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("habit not found")
	}

	// Enforce duration requirement: if habit has a duration, log must include duration
	if habit.DurationSeconds != nil {
		if req.DurationSeconds == nil {
			return nil, fmt.Errorf("duration_seconds is required for logs of this habit")
		}
		if *req.DurationSeconds <= 0 {
			return nil, fmt.Errorf("duration_seconds must be positive")
		}
	} else {
		// If provided, still validate positivity
		if req.DurationSeconds != nil && *req.DurationSeconds <= 0 {
			return nil, fmt.Errorf("duration_seconds must be positive")
		}
	}

	// Create log
	log := &models.Log{
		HabitID:         req.HabitID,
		Notes:           req.Notes,
		DurationSeconds: req.DurationSeconds,
	}

	if err := s.logRepo.Create(log); err != nil {
		return nil, err
	}

	return log, nil
}

// GetByID retrieves a log by ID and ensures it belongs to the user's habit
func (s *logService) GetByID(id, userID int64) (*models.Log, error) {
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("log not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("log not found")
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a habit
func (s *logService) GetByHabitID(habitID, userID int64) ([]*models.Log, error) {
	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("habit not found")
	}

	return s.logRepo.GetByHabitID(habitID)
}

// Update updates a log with validation
func (s *logService) Update(id, userID int64, req *models.UpdateLogRequest) (*models.Log, error) {
	// Get existing log and verify ownership
	log, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Notes != nil {
		log.Notes = *req.Notes
	}
	if req.DurationSeconds != nil {
		if *req.DurationSeconds <= 0 {
			return nil, fmt.Errorf("duration_seconds must be positive")
		}
		log.DurationSeconds = req.DurationSeconds
	}

	if err := s.logRepo.Update(log); err != nil {
		return nil, err
	}

	return log, nil
}

// Delete deletes a log and ensures it belongs to the user's habit
func (s *logService) Delete(id, userID int64) error {
	// Verify ownership before deleting
	log, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	// Verify the log's habit belongs to the user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return fmt.Errorf("log not found")
	}
	if habit.UserID != userID {
		return fmt.Errorf("log not found")
	}

	return s.logRepo.Delete(id)
}
