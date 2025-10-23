package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// LogService defines the interface for log business logic
type LogService interface {
	Create(ctx context.Context, habitID string, userID string, req *models.CreateLogRequest) (*models.Log, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Log, error)
	GetByHabitID(ctx context.Context, habitID string, userID string) ([]*models.Log, error)
	Update(ctx context.Context, id string, userID string, req *models.UpdateLogRequest) (*models.Log, error)
	Delete(ctx context.Context, id string, userID string) error
}

// logService implements LogService interface
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

// Create creates a new log with validation
func (s *logService) Create(ctx context.Context, habitID string, userID string, req *models.CreateLogRequest) (*models.Log, error) {
	// Validate required fields
	if habitID == "" {
		return nil, fmt.Errorf("habit ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Verify habit exists and belongs to user
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	// Validate duration: if habit has duration, log must have it too
	if habit.DurationSeconds != nil && req.DurationSeconds == nil {
		return nil, fmt.Errorf("duration is required for this habit")
	}

	// Create log model
	log := &models.Log{
		ID:              uuid.New().String(),
		HabitID:         habitID,
		Notes:           req.Notes,
		DurationSeconds: req.DurationSeconds,
		CreatedAt:       time.Now().UTC(),
	}

	// Save to repository
	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return log, nil
}

// GetByID retrieves a log by ID, ensuring the user owns the associated habit
func (s *logService) GetByID(ctx context.Context, id string, userID string) (*models.Log, error) {
	// Validate IDs
	if id == "" {
		return nil, fmt.Errorf("log ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get log
	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify user owns the habit
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: log does not belong to user")
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a habit, ensuring the user owns the habit
func (s *logService) GetByHabitID(ctx context.Context, habitID string, userID string) ([]*models.Log, error) {
	// Validate IDs
	if habitID == "" {
		return nil, fmt.Errorf("habit ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Verify user owns the habit
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	// Get logs
	logs, err := s.logRepo.GetByHabitID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	return logs, nil
}

// Update updates a log with validation, ensuring the user owns the associated habit
func (s *logService) Update(ctx context.Context, id string, userID string, req *models.UpdateLogRequest) (*models.Log, error) {
	// Validate IDs
	if id == "" {
		return nil, fmt.Errorf("log ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get existing log
	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Verify user owns the habit
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: log does not belong to user")
	}

	// Update fields if provided
	if req.Notes != nil {
		log.Notes = *req.Notes
	}
	if req.DurationSeconds != nil {
		log.DurationSeconds = req.DurationSeconds
	}

	// Validate duration: if habit has duration, log must have it too after update
	if habit.DurationSeconds != nil && log.DurationSeconds == nil {
		return nil, fmt.Errorf("duration is required for this habit")
	}

	// Save updates
	if err := s.logRepo.Update(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to update log: %w", err)
	}

	return log, nil
}

// Delete removes a log, ensuring the user owns the associated habit
func (s *logService) Delete(ctx context.Context, id string, userID string) error {
	// Validate IDs
	if id == "" {
		return fmt.Errorf("log ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Get log
	log, err := s.logRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get log: %w", err)
	}

	// Verify user owns the habit
	habit, err := s.habitRepo.GetByID(ctx, log.HabitID)
	if err != nil {
		return fmt.Errorf("habit not found")
	}
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized: log does not belong to user")
	}

	// Delete log
	if err := s.logRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	return nil
}
