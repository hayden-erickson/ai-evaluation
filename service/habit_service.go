package service

import (
	"fmt"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// HabitService defines the interface for habit business logic
type HabitService interface {
	Create(userID int64, req *models.CreateHabitRequest) (*models.Habit, error)
	GetByID(id, userID int64) (*models.Habit, error)
	GetByUserID(userID int64) ([]*models.Habit, error)
	Update(id, userID int64, req *models.UpdateHabitRequest) (*models.Habit, error)
	Delete(id, userID int64) error
}

type habitService struct {
	repo repository.HabitRepository
}

// NewHabitService creates a new HabitService instance
func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{repo: repo}
}

// Create creates a new habit with validation
func (s *habitService) Create(userID int64, req *models.CreateHabitRequest) (*models.Habit, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create habit
	habit := &models.Habit{
		UserID:          userID,
		Name:            req.Name,
		Description:     req.Description,
		DurationSeconds: req.DurationSeconds,
	}

	if err := s.repo.Create(habit); err != nil {
		return nil, err
	}

	return habit, nil
}

// GetByID retrieves a habit by ID and ensures it belongs to the user
func (s *habitService) GetByID(id, userID int64) (*models.Habit, error) {
	habit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure habit belongs to the user
	if habit.UserID != userID {
		return nil, fmt.Errorf("habit not found")
	}

	return habit, nil
}

// GetByUserID retrieves all habits for a user
func (s *habitService) GetByUserID(userID int64) ([]*models.Habit, error) {
	return s.repo.GetByUserID(userID)
}

// Update updates a habit with validation
func (s *habitService) Update(id, userID int64, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Get existing habit and verify ownership
	habit, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		habit.Name = *req.Name
	}
	if req.Description != nil {
		habit.Description = *req.Description
	}
	if req.DurationSeconds != nil {
		if *req.DurationSeconds <= 0 {
			return nil, fmt.Errorf("duration_seconds must be positive")
		}
		habit.DurationSeconds = req.DurationSeconds
	}

	if err := s.repo.Update(habit); err != nil {
		return nil, err
	}

	return habit, nil
}

// Delete deletes a habit and ensures it belongs to the user
func (s *habitService) Delete(id, userID int64) error {
	// Verify ownership before deleting
	habit, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	// Verify ownership
	if habit.UserID != userID {
		return fmt.Errorf("habit not found")
	}

	return s.repo.Delete(id)
}

// validateCreateRequest validates the create habit request
func (s *habitService) validateCreateRequest(req *models.CreateHabitRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if req.DurationSeconds != nil && *req.DurationSeconds <= 0 {
		return fmt.Errorf("duration_seconds must be positive")
	}
	return nil
}
