package service

import (
	"errors"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// LogService defines the interface for log business logic
type LogService interface {
	CreateLog(log *models.Log, userID int64) error
	GetLog(id, userID int64) (*models.Log, error)
	GetHabitLogs(habitID, userID int64) ([]*models.Log, error)
	UpdateLog(log *models.Log, userID int64) error
	DeleteLog(id, userID int64) error
}

type logService struct {
	logRepo   repository.LogRepository
	habitRepo repository.HabitRepository
}

// NewLogService creates a new LogService
func NewLogService(logRepo repository.LogRepository, habitRepo repository.HabitRepository) LogService {
	return &logService{logRepo: logRepo, habitRepo: habitRepo}
}

// CreateLog creates a new log
func (s *logService) CreateLog(log *models.Log, userID int64) error {
	habit, err := s.habitRepo.GetHabitByID(log.HabitID)
	if err != nil {
		return err
	}
	if habit.UserID != userID {
		return errors.New("user does not have permission to create logs for this habit")
	}
	return s.logRepo.CreateLog(log)
}

// GetLog retrieves a log by ID and ensures the user has permission
func (s *logService) GetLog(id, userID int64) (*models.Log, error) {
	log, err := s.logRepo.GetLogByID(id)
	if err != nil {
		return nil, err
	}
	habit, err := s.habitRepo.GetHabitByID(log.HabitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, errors.New("user does not have permission to view this log")
	}
	return log, nil
}

// GetHabitLogs retrieves all logs for a habit and ensures the user has permission
func (s *logService) GetHabitLogs(habitID, userID int64) ([]*models.Log, error) {
	habit, err := s.habitRepo.GetHabitByID(habitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, errors.New("user does not have permission to view logs for this habit")
	}
	return s.logRepo.GetLogsByHabitID(habitID)
}

// UpdateLog updates a log and ensures the user has permission
func (s *logService) UpdateLog(log *models.Log, userID int64) error {
	existingLog, err := s.logRepo.GetLogByID(log.ID)
	if err != nil {
		return err
	}
	habit, err := s.habitRepo.GetHabitByID(existingLog.HabitID)
	if err != nil {
		return err
	}
	if habit.UserID != userID {
		return errors.New("user does not have permission to update this log")
	}
	return s.logRepo.UpdateLog(log)
}

// DeleteLog deletes a log and ensures the user has permission
func (s *logService) DeleteLog(id, userID int64) error {
	log, err := s.logRepo.GetLogByID(id)
	if err != nil {
		return err
	}
	habit, err := s.habitRepo.GetHabitByID(log.HabitID)
	if err != nil {
		return err
	}
	if habit.UserID != userID {
		return errors.New("user does not have permission to delete this log")
	}
	return s.logRepo.DeleteLog(id)
}
