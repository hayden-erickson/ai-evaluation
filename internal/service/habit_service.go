package service

import (
	"errors"
	"new-api/internal/models"
	"new-api/internal/repository"
)

// HabitService defines the interface for habit business logic
type HabitService interface {
	CreateHabit(habit *models.Habit) error
	GetHabit(id, userID int64) (*models.Habit, error)
	GetUserHabits(userID int64) ([]*models.Habit, error)
	UpdateHabit(habit *models.Habit, userID int64) error
	DeleteHabit(id, userID int64) error
}

type habitService struct {
	habitRepo repository.HabitRepository
}

// NewHabitService creates a new HabitService
func NewHabitService(habitRepo repository.HabitRepository) HabitService {
	return &habitService{habitRepo: habitRepo}
}

// CreateHabit creates a new habit
func (s *habitService) CreateHabit(habit *models.Habit) error {
	return s.habitRepo.CreateHabit(habit)
}

// GetHabit retrieves a habit by ID and ensures it belongs to the user
func (s *habitService) GetHabit(id, userID int64) (*models.Habit, error) {
	habit, err := s.habitRepo.GetHabitByID(id)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, errors.New("user does not have permission to view this habit")
	}
	return habit, nil
}

// GetUserHabits retrieves all habits for a user
func (s *habitService) GetUserHabits(userID int64) ([]*models.Habit, error) {
	return s.habitRepo.GetHabitsByUserID(userID)
}

// UpdateHabit updates a habit and ensures it belongs to the user
func (s *habitService) UpdateHabit(habit *models.Habit, userID int64) error {
	existingHabit, err := s.habitRepo.GetHabitByID(habit.ID)
	if err != nil {
		return err
	}
	if existingHabit.UserID != userID {
		return errors.New("user does not have permission to update this habit")
	}
	return s.habitRepo.UpdateHabit(habit)
}

// DeleteHabit deletes a habit and ensures it belongs to the user
func (s *habitService) DeleteHabit(id, userID int64) error {
	habit, err := s.habitRepo.GetHabitByID(id)
	if err != nil {
		return err
	}
	if habit.UserID != userID {
		return errors.New("user does not have permission to delete this habit")
	}
	return s.habitRepo.DeleteHabit(id)
}
