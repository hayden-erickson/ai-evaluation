package service

import (
	"context"
	"errors"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

type HabitService interface {
	Create(ctx context.Context, h *models.Habit) (int64, error)
	GetByID(ctx context.Context, requesterID, id int64) (*models.Habit, error)
	ListByUser(ctx context.Context, userID int64) ([]*models.Habit, error)
	Update(ctx context.Context, requesterID int64, h *models.Habit) error
	Delete(ctx context.Context, requesterID, id int64) error
}

type habitService struct{ habits repository.HabitRepository }

func NewHabitService(habits repository.HabitRepository) HabitService { return &habitService{habits: habits} }

func (s *habitService) Create(ctx context.Context, h *models.Habit) (int64, error) { return s.habits.Create(ctx, h) }

func (s *habitService) GetByID(ctx context.Context, requesterID, id int64) (*models.Habit, error) {
	h, err := s.habits.GetByID(ctx, id)
	if err != nil || h == nil { return h, err }
	if h.UserID != requesterID { return nil, errors.New("forbidden") }
	return h, nil
}

func (s *habitService) ListByUser(ctx context.Context, userID int64) ([]*models.Habit, error) {
	return s.habits.ListByUser(ctx, userID)
}

func (s *habitService) Update(ctx context.Context, requesterID int64, h *models.Habit) error {
	existing, err := s.habits.GetByID(ctx, h.ID)
	if err != nil { return err }
	if existing == nil { return errors.New("not found") }
	if existing.UserID != requesterID { return errors.New("forbidden") }
	h.UserID = requesterID
	return s.habits.Update(ctx, h)
}

func (s *habitService) Delete(ctx context.Context, requesterID, id int64) error {
	existing, err := s.habits.GetByID(ctx, id)
	if err != nil { return err }
	if existing == nil { return errors.New("not found") }
	if existing.UserID != requesterID { return errors.New("forbidden") }
	return s.habits.Delete(ctx, id)
}
