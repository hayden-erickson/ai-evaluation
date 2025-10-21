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

// HabitService defines the interface for habit business logic
type HabitService interface {
	Create(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error)
	GetByID(ctx context.Context, userID, id string) (*models.Habit, error)
	Update(ctx context.Context, userID, id string, req *models.UpdateHabitRequest) (*models.Habit, error)
	Delete(ctx context.Context, userID, id string) error
	ListByUserID(ctx context.Context, userID string) ([]*models.Habit, error)
}

// habitService implements HabitService interface
type habitService struct {
	repo repository.HabitRepository
}

// NewHabitService creates a new habit service instance
func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{repo: repo}
}

// Create creates a new habit with validation
func (s *habitService) Create(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error) {
	// Validate input
	if err := validator.ValidateCreateHabitRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Generate habit ID
	id := auth.GenerateID()

	// Create habit model
	habit := &models.Habit{
		ID:          id,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, habit); err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return habit, nil
}

// GetByID retrieves a habit by ID with user ownership check
func (s *habitService) GetByID(ctx context.Context, userID, id string) (*models.Habit, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("habit ID is required")
	}

	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	return habit, nil
}

// Update updates an existing habit with user ownership check
func (s *habitService) Update(ctx context.Context, userID, id string, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("habit ID is required")
	}

	// Get existing habit
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		habit.Name = *req.Name
	}
	if req.Description != nil {
		habit.Description = *req.Description
	}

	// Save to repository
	if err := s.repo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return habit, nil
}

// Delete deletes a habit by ID with user ownership check
func (s *habitService) Delete(ctx context.Context, userID, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("habit ID is required")
	}

	// Get existing habit to check ownership
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get habit: %w", err)
	}

	// Check ownership
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized: habit does not belong to user")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}

// ListByUserID retrieves all habits for a specific user
func (s *habitService) ListByUserID(ctx context.Context, userID string) ([]*models.Habit, error) {
	habits, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}

	return habits, nil
}
