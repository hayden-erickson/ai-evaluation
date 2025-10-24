package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

var (
	// ErrLogNotFound is returned when a log is not found
	ErrLogNotFound = errors.New("log not found")
)

// LogService defines the interface for log business logic
type LogService interface {
	Create(userID int64, req *models.LogCreateRequest) (*models.Log, error)
	GetByID(id, userID int64) (*models.Log, error)
	GetByHabitID(habitID, userID int64) ([]*models.Log, error)
	Update(id, userID int64, req *models.LogUpdateRequest) error
	Delete(id, userID int64) error
}

// logService implements LogService interface
type logService struct {
	logRepo   repository.LogRepository
	habitRepo repository.HabitRepository
}

// NewLogService creates a new log service instance
func NewLogService(logRepo repository.LogRepository, habitRepo repository.HabitRepository) LogService {
	return &logService{
		logRepo:   logRepo,
		habitRepo: habitRepo,
	}
}

// Create creates a new log entry for a habit
func (s *logService) Create(userID int64, req *models.LogCreateRequest) (*models.Log, error) {
	// Validate required fields
	if req.HabitID == 0 {
		return nil, errors.New("habit_id is required")
	}

	// Verify the habit exists and belongs to the user
	habit, err := s.habitRepo.GetByID(req.HabitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}
	if habit.UserID != userID {
		return nil, ErrUnauthorized
	}

	if habit.Duration != nil && req.Duration == nil {
		return nil, errors.New("log must have a duration if the habit has a duration")
	}

	if habit.Duration == nil && req.Duration != nil {
		return nil, errors.New("log cannot have a duration if the habit does not have a duration")
	}

	// Create log
	log := &models.Log{
		HabitID:   req.HabitID,
		Notes:     utils.SanitizeString(req.Notes),
		Duration:  req.Duration,
		CreatedAt: time.Now(),
	}

	if err := s.logRepo.Create(log); err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return log, nil
}

// GetByID retrieves a log by ID, ensuring the associated habit belongs to the user
func (s *logService) GetByID(id, userID int64) (*models.Log, error) {
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}
	if log == nil {
		return nil, ErrLogNotFound
	}

	// Verify the associated habit belongs to the user
	habit, err := s.habitRepo.GetByID(log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit == nil || habit.UserID != userID {
		return nil, ErrUnauthorized
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a habit, ensuring it belongs to the user
func (s *logService) GetByHabitID(habitID, userID int64) ([]*models.Log, error) {
	// Verify the habit exists and belongs to the user
	habit, err := s.habitRepo.GetByID(habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}
	if habit.UserID != userID {
		return nil, ErrUnauthorized
	}

	logs, err := s.logRepo.GetByHabitID(habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, nil
}

// Update updates a log, ensuring the associated habit belongs to the user
func (s *logService) Update(id, userID int64, req *models.LogUpdateRequest) error {
	// Verify log exists and belongs to user (through habit)
	log, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	if log == nil {
		return ErrLogNotFound
	}

	// Sanitize fields
	if req.Notes != nil {
		notes := utils.SanitizeString(*req.Notes)
		req.Notes = &notes
	}

	return s.logRepo.Update(id, req)
}

// Delete removes a log, ensuring the associated habit belongs to the user
func (s *logService) Delete(id, userID int64) error {
	// Verify log exists and belongs to user (through habit)
	log, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	if log == nil {
		return ErrLogNotFound
	}

	return s.logRepo.Delete(id)
}
