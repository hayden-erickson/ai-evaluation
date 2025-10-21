package service

import (
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// HabitService defines the interface for habit business logic
type HabitService interface {
	CreateHabit(userID int64, req *models.CreateHabitRequest) (*models.Habit, error)
	GetHabit(id int64, userID int64) (*models.Habit, error)
	GetUserHabits(userID int64) ([]*models.Habit, error)
	UpdateHabit(id int64, userID int64, req *models.UpdateHabitRequest) (*models.Habit, error)
	DeleteHabit(id int64, userID int64) error
}

// habitService implements HabitService
type habitService struct {
	repo repository.HabitRepository
}

// NewHabitService creates a new habit service
func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{
		repo: repo,
	}
}

// CreateHabit creates a new habit
func (s *habitService) CreateHabit(userID int64, req *models.CreateHabitRequest) (*models.Habit, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create the habit
	habit, err := s.repo.Create(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return habit, nil
}

// GetHabit retrieves a habit by ID
func (s *habitService) GetHabit(id int64, userID int64) (*models.Habit, error) {
	habit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if the habit belongs to the user
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to habit")
	}

	return habit, nil
}

// GetUserHabits retrieves all habits for a user
func (s *habitService) GetUserHabits(userID int64) ([]*models.Habit, error) {
	habits, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}

	return habits, nil
}

// UpdateHabit updates a habit
func (s *habitService) UpdateHabit(id int64, userID int64, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the habit to check ownership
	habit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if the habit belongs to the user
	if habit.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to habit")
	}

	// Update the habit
	updatedHabit, err := s.repo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return updatedHabit, nil
}

// DeleteHabit deletes a habit
func (s *habitService) DeleteHabit(id int64, userID int64) error {
	// Get the habit to check ownership
	habit, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get habit: %w", err)
	}

	// Check if the habit belongs to the user
	if habit.UserID != userID {
		return fmt.Errorf("unauthorized access to habit")
	}

	// Delete the habit
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}
