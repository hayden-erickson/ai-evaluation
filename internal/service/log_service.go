package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/pkg/auth"
	"github.com/hayden-erickson/ai-evaluation/pkg/validator"
)

// LogService defines the interface for log business logic
type LogService interface {
	Create(ctx context.Context, userID string, req *models.CreateLogRequest) (*models.Log, error)
	GetByID(ctx context.Context, userID, id string) (*models.Log, error)
	Update(ctx context.Context, userID, id string, req *models.UpdateLogRequest) (*models.Log, error)
	Delete(ctx context.Context, userID, id string) error
	ListByHabitID(ctx context.Context, userID, habitID string) ([]*models.Log, error)
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

// Create creates a new log with validation and ownership check
func (s *logService) Create(ctx context.Context, userID string, req *models.CreateLogRequest) (*models.Log, error) {
	// Validate input
	if err := validator.ValidateCreateLogRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Verify habit exists and belongs to user
	habit, err := s.habitRepo.GetByID(ctx, req.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	// Validate duration based on habit requirements
	if err := validator.ValidateLogDuration(habit, req.Duration); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Generate log ID
	id := auth.GenerateID()

	// Create log model
	log := &models.Log{
		ID:        id,
		HabitID:   req.HabitID,
		Notes:     req.Notes,
		Duration:  req.Duration,
		CreatedAt: time.Now().UTC(),
	}

	// Save to repository
	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return log, nil
}

// GetByID retrieves a log by ID with user ownership check
func (s *logService) GetByID(ctx context.Context, userID, id string) (*models.Log, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("log ID is required")
	}

	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: log does not belong to user")
	}

	return log, nil
}

// Update updates an existing log with user ownership check
func (s *logService) Update(ctx context.Context, userID, id string, req *models.UpdateLogRequest) (*models.Log, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("log ID is required")
	}

	// Get existing log
	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: log does not belong to user")
	}

	// Update fields if provided
	if req.Notes != nil {
		log.Notes = *req.Notes
	}
	if req.Duration != nil {
		// Validate duration based on habit requirements
		if err := validator.ValidateLogDuration(habit, req.Duration); err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
		log.Duration = req.Duration
	}

	// Save to repository
	if err := s.logRepo.Update(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to update log: %w", err)
	}

	return log, nil
}

// Delete deletes a log by ID with user ownership check
func (s *logService) Delete(ctx context.Context, userID, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("log ID is required")
	}

	// Get existing log to check ownership
	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get log: %w", err)
	}

	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return fmt.Errorf("habit not found: %w", err)
	}
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized: log does not belong to user")
	}

	if err := s.logRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	return nil
}

// ListByHabitID retrieves all logs for a specific habit with user ownership check
func (s *logService) ListByHabitID(ctx context.Context, userID, habitID string) ([]*models.Log, error) {
	// Validate input
	if habitID == "" {
		return nil, fmt.Errorf("habit ID is required")
	}

	// Verify habit belongs to user
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found: %w", err)
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	logs, err := s.logRepo.ListByHabitID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	return logs, nil
}
