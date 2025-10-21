package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// HabitService defines the interface for habit business logic
type HabitService interface {
	Create(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Habit, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Habit, error)
	Update(ctx context.Context, id string, userID string, req *models.UpdateHabitRequest) (*models.Habit, error)
	Delete(ctx context.Context, id string, userID string) error
}

// habitService implements HabitService interface
type habitService struct {
	repo repository.HabitRepository
}

// NewHabitService creates a new habit service
func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{repo: repo}
}

// Create creates a new habit with validation
func (s *habitService) Create(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Create habit model
	habit := &models.Habit{
		ID:          uuid.New().String(),
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

// GetByID retrieves a habit by ID, ensuring it belongs to the user
func (s *habitService) GetByID(ctx context.Context, id string, userID string) (*models.Habit, error) {
	// Validate IDs
	if id == "" {
		return nil, fmt.Errorf("habit ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get habit
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if habit belongs to user (RBAC)
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized: habit does not belong to user")
	}

	return habit, nil
}

// GetByUserID retrieves all habits for a user
func (s *habitService) GetByUserID(ctx context.Context, userID string) ([]*models.Habit, error) {
	// Validate user ID
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	habits, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}

	return habits, nil
}

// Update updates a habit with validation, ensuring it belongs to the user
func (s *habitService) Update(ctx context.Context, id string, userID string, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Validate IDs
	if id == "" {
		return nil, fmt.Errorf("habit ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get existing habit
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if habit belongs to user (RBAC)
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

	// Save updates
	if err := s.repo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return habit, nil
}

// Delete removes a habit, ensuring it belongs to the user
func (s *habitService) Delete(ctx context.Context, id string, userID string) error {
	// Validate IDs
	if id == "" {
		return fmt.Errorf("habit ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Get habit to verify ownership
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if habit belongs to user (RBAC)
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized: habit does not belong to user")
	}

	// Delete habit
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}

