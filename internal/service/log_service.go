package service

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// LogService defines the interface for log business logic operations
type LogService interface {
	CreateLog(ctx context.Context, userID string, req *models.CreateLogRequest) (*models.Log, error)
	GetLog(ctx context.Context, logID, userID string) (*models.Log, error)
	ListLogs(ctx context.Context, userID, habitID string) ([]*models.Log, error)
	UpdateLog(ctx context.Context, logID, userID string, req *models.UpdateLogRequest) (*models.Log, error)
	DeleteLog(ctx context.Context, logID, userID string) error
}

// logService implements LogService interface
type logService struct {
	logRepo   repository.LogRepository
	habitRepo repository.HabitRepository
}

// NewLogService creates a new instance of LogService
func NewLogService(logRepo repository.LogRepository, habitRepo repository.HabitRepository) LogService {
	return &logService{
		logRepo:   logRepo,
		habitRepo: habitRepo,
	}
}

// CreateLog creates a new log entry for a habit
func (s *logService) CreateLog(ctx context.Context, userID string, req *models.CreateLogRequest) (*models.Log, error) {
	// Verify habit exists and belongs to user
	habit, err := s.habitRepo.GetByID(ctx, req.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}

	// Verify ownership - users can only create logs for their own habits
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: habit does not belong to user")
	}

	// Create log object
	log := &models.Log{
		HabitID: req.HabitID,
		Notes:   req.Notes,
	}

	// Save to database
	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return log, nil
}

// GetLog retrieves a specific log entry, ensuring it belongs to the requesting user
func (s *logService) GetLog(ctx context.Context, logID, userID string) (*models.Log, error) {
	log, err := s.logRepo.GetByID(ctx, logID)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit ownership
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: log does not belong to user")
	}

	return log, nil
}

// ListLogs retrieves all log entries for a specific habit
func (s *logService) ListLogs(ctx context.Context, userID, habitID string) ([]*models.Log, error) {
	// Verify habit exists and belongs to user
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: habit does not belong to user")
	}

	// Retrieve logs
	logs, err := s.logRepo.ListByHabitID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	return logs, nil
}

// UpdateLog updates an existing log entry, ensuring it belongs to the requesting user
func (s *logService) UpdateLog(ctx context.Context, logID, userID string, req *models.UpdateLogRequest) (*models.Log, error) {
	// Retrieve existing log
	log, err := s.logRepo.GetByID(ctx, logID)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit ownership
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: log does not belong to user")
	}

	// Apply updates only for fields that are provided
	if req.Notes != nil {
		log.Notes = *req.Notes
	}

	// Save updates to database
	if err := s.logRepo.Update(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to update log: %w", err)
	}

	return log, nil
}

// DeleteLog removes a log entry, ensuring it belongs to the requesting user
func (s *logService) DeleteLog(ctx context.Context, logID, userID string) error {
	// Retrieve existing log
	log, err := s.logRepo.GetByID(ctx, logID)
	if err != nil {
		return fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit ownership
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return fmt.Errorf("habit not found: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return fmt.Errorf("forbidden: log does not belong to user")
	}

	// Delete from database
	if err := s.logRepo.Delete(ctx, logID); err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	return nil
}
