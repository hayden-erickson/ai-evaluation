package service

import (
	"context"
	"errors"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

type LogService interface {
	Create(ctx context.Context, requesterID int64, l *models.LogEntry) (int64, error)
	GetByID(ctx context.Context, requesterID, id int64) (*models.LogEntry, error)
	ListByHabit(ctx context.Context, requesterID, habitID int64) ([]*models.LogEntry, error)
	Update(ctx context.Context, requesterID int64, l *models.LogEntry) error
	Delete(ctx context.Context, requesterID, id int64) error
}

type logService struct {
	logs   repository.LogRepository
	habits repository.HabitRepository
}

func NewLogService(logs repository.LogRepository, habits repository.HabitRepository) LogService {
	return &logService{logs: logs, habits: habits}
}

func (s *logService) Create(ctx context.Context, requesterID int64, l *models.LogEntry) (int64, error) {
	h, err := s.habits.GetByID(ctx, l.HabitID)
	if err != nil { return 0, err }
	if h == nil { return 0, errors.New("habit not found") }
	if h.UserID != requesterID { return 0, errors.New("forbidden") }
	return s.logs.Create(ctx, l)
}

func (s *logService) GetByID(ctx context.Context, requesterID, id int64) (*models.LogEntry, error) {
	l, err := s.logs.GetByID(ctx, id)
	if err != nil || l == nil { return l, err }
	h, err := s.habits.GetByID(ctx, l.HabitID)
	if err != nil { return nil, err }
	if h == nil || h.UserID != requesterID { return nil, errors.New("forbidden") }
	return l, nil
}

func (s *logService) ListByHabit(ctx context.Context, requesterID, habitID int64) ([]*models.LogEntry, error) {
	h, err := s.habits.GetByID(ctx, habitID)
	if err != nil { return nil, err }
	if h == nil { return nil, errors.New("habit not found") }
	if h.UserID != requesterID { return nil, errors.New("forbidden") }
	return s.logs.ListByHabit(ctx, habitID)
}

func (s *logService) Update(ctx context.Context, requesterID int64, l *models.LogEntry) error {
	existing, err := s.logs.GetByID(ctx, l.ID)
	if err != nil { return err }
	if existing == nil { return errors.New("not found") }
	h, err := s.habits.GetByID(ctx, existing.HabitID)
	if err != nil { return err }
	if h == nil || h.UserID != requesterID { return errors.New("forbidden") }
	return s.logs.Update(ctx, l)
}

func (s *logService) Delete(ctx context.Context, requesterID, id int64) error {
	l, err := s.logs.GetByID(ctx, id)
	if err != nil { return err }
	if l == nil { return errors.New("not found") }
	h, err := s.habits.GetByID(ctx, l.HabitID)
	if err != nil { return err }
	if h == nil || h.UserID != requesterID { return errors.New("forbidden") }
	return s.logs.Delete(ctx, id)
}
