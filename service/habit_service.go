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
	// ErrHabitNotFound is returned when a habit is not found
	ErrHabitNotFound = errors.New("habit not found")
	// ErrUnauthorized is returned when user doesn't have permission
	ErrUnauthorized = errors.New("unauthorized")
)

// HabitService defines the interface for habit business logic
type HabitService interface {
	Create(userID int64, req *models.HabitCreateRequest) (*models.Habit, error)
	GetByID(id, userID int64) (*models.Habit, error)
	GetByUserID(userID int64) ([]*models.Habit, error)
	Update(id, userID int64, req *models.HabitUpdateRequest) error
	Delete(id, userID int64) error
}

// habitService implements HabitService interface
type habitService struct {
	repo repository.HabitRepository
}

// NewHabitService creates a new habit service instance
func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{repo: repo}
}

// Create creates a new habit for a user
func (s *habitService) Create(userID int64, req *models.HabitCreateRequest) (*models.Habit, error) {
	// Validate required fields
	if err := utils.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}
	
	// Create habit
	habit := &models.Habit{
		UserID:      userID,
		Name:        utils.SanitizeString(req.Name),
		Description: utils.SanitizeString(req.Description),
		CreatedAt:   time.Now(),
	}
	
	if err := s.repo.Create(habit); err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}
	
	return habit, nil
}

// GetByID retrieves a habit by ID, ensuring it belongs to the user
func (s *habitService) GetByID(id, userID int64) (*models.Habit, error) {
	habit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}
	
	// Verify the habit belongs to the user
	if habit.UserID != userID {
		return nil, ErrUnauthorized
	}
	
	return habit, nil
}

// GetByUserID retrieves all habits for a user
func (s *habitService) GetByUserID(userID int64) ([]*models.Habit, error) {
	habits, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	return habits, nil
}

// Update updates a habit, ensuring it belongs to the user
func (s *habitService) Update(id, userID int64, req *models.HabitUpdateRequest) error {
	// Verify habit exists and belongs to user
	habit, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	if habit == nil {
		return ErrHabitNotFound
	}
	
	// Validate and sanitize fields
	if req.Name != nil {
		name := utils.SanitizeString(*req.Name)
		req.Name = &name
		if err := utils.ValidateRequired(*req.Name, "name"); err != nil {
			return err
		}
	}
	if req.Description != nil {
		description := utils.SanitizeString(*req.Description)
		req.Description = &description
	}
	
	return s.repo.Update(id, req)
}

// Delete removes a habit, ensuring it belongs to the user
func (s *habitService) Delete(id, userID int64) error {
	// Verify habit exists and belongs to user
	habit, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	if habit == nil {
		return ErrHabitNotFound
	}
	
	return s.repo.Delete(id)
}
