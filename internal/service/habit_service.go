package service

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// HabitService defines the interface for habit business logic operations
type HabitService interface {
	CreateHabit(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error)
	GetHabit(ctx context.Context, habitID, userID string) (*models.Habit, error)
	ListHabits(ctx context.Context, userID string) ([]*models.Habit, error)
	UpdateHabit(ctx context.Context, habitID, userID string, req *models.UpdateHabitRequest) (*models.Habit, error)
	DeleteHabit(ctx context.Context, habitID, userID string) error
}

// habitService implements HabitService interface
type habitService struct {
	habitRepo repository.HabitRepository
	userRepo  repository.UserRepository
}

// NewHabitService creates a new instance of HabitService
func NewHabitService(habitRepo repository.HabitRepository, userRepo repository.UserRepository) HabitService {
	return &habitService{
		habitRepo: habitRepo,
		userRepo:  userRepo,
	}
}

// CreateHabit creates a new habit for a user
func (s *habitService) CreateHabit(ctx context.Context, userID string, req *models.CreateHabitRequest) (*models.Habit, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Create habit object
	habit := &models.Habit{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	// Save to database
	if err := s.habitRepo.Create(ctx, habit); err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return habit, nil
}

// GetHabit retrieves a specific habit, ensuring it belongs to the requesting user
func (s *habitService) GetHabit(ctx context.Context, habitID, userID string) (*models.Habit, error) {
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Verify ownership - users can only access their own habits
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: habit does not belong to user")
	}

	return habit, nil
}

// ListHabits retrieves all habits for a specific user
func (s *habitService) ListHabits(ctx context.Context, userID string) ([]*models.Habit, error) {
	habits, err := s.habitRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}

	return habits, nil
}

// UpdateHabit updates an existing habit, ensuring it belongs to the requesting user
func (s *habitService) UpdateHabit(ctx context.Context, habitID, userID string, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Retrieve existing habit
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return nil, fmt.Errorf("forbidden: habit does not belong to user")
	}

	// Apply updates only for fields that are provided
	if req.Name != nil {
		habit.Name = *req.Name
	}
	if req.Description != nil {
		habit.Description = *req.Description
	}

	// Save updates to database
	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return habit, nil
}

// DeleteHabit removes a habit, ensuring it belongs to the requesting user
func (s *habitService) DeleteHabit(ctx context.Context, habitID, userID string) error {
	// Retrieve existing habit
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return fmt.Errorf("failed to get habit: %w", err)
	}

	// Verify ownership
	if habit.UserID != userID {
		return fmt.Errorf("forbidden: habit does not belong to user")
	}

	// Delete from database
	if err := s.habitRepo.Delete(ctx, habitID); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}
