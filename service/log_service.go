package service

import (
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// LogService defines the interface for log business logic
type LogService interface {
	CreateLog(habitID int64, userID int64, req *models.CreateLogRequest) (*models.Log, error)
	GetLog(id int64, userID int64) (*models.Log, error)
	GetHabitLogs(habitID int64, userID int64) ([]*models.Log, error)
	UpdateLog(id int64, userID int64, req *models.UpdateLogRequest) (*models.Log, error)
	DeleteLog(id int64, userID int64) error
}

// logService implements LogService
type logService struct {
	logRepo   repository.LogRepository
	habitRepo repository.HabitRepository
}

// NewLogService creates a new log service
func NewLogService(logRepo repository.LogRepository, habitRepo repository.HabitRepository) LogService {
	return &logService{
		logRepo:   logRepo,
		habitRepo: habitRepo,
	}
}

// CreateLog creates a new log
func (s *logService) CreateLog(habitID int64, userID int64, req *models.CreateLogRequest) (*models.Log, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify that the habit exists and belongs to the user
	habit, err := s.habitRepo.GetByID(habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to habit")
	}

	// Create the log
	log, err := s.logRepo.Create(habitID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return log, nil
}

// GetLog retrieves a log by ID
func (s *logService) GetLog(id int64, userID int64) (*models.Log, error) {
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify that the log's habit belongs to the user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to log")
	}

	return log, nil
}

// GetHabitLogs retrieves all logs for a habit
func (s *logService) GetHabitLogs(habitID int64, userID int64) ([]*models.Log, error) {
	// Verify that the habit exists and belongs to the user
	habit, err := s.habitRepo.GetByID(habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to habit")
	}

	// Get all logs for the habit
	logs, err := s.logRepo.GetByHabitID(habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	return logs, nil
}

// UpdateLog updates a log
func (s *logService) UpdateLog(id int64, userID int64, req *models.UpdateLogRequest) (*models.Log, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the log to check ownership
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify that the log's habit belongs to the user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to log")
	}

	// Update the log
	updatedLog, err := s.logRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update log: %w", err)
	}

	return updatedLog, nil
}

// DeleteLog deletes a log
func (s *logService) DeleteLog(id int64, userID int64) error {
	// Get the log to check ownership
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get log: %w", err)
	}

	// Verify that the log's habit belongs to the user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return fmt.Errorf("failed to get habit: %w", err)
	}
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized access to log")
	}

	// Delete the log
	if err := s.logRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	return nil
}
