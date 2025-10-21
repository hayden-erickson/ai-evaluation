package repository

import (
	"database/sql"
	"new-api/internal/models"
)

// LogRepository defines the interface for log persistence
type LogRepository interface {
	CreateLog(log *models.Log) error
	GetLogByID(id int64) (*models.Log, error)
	GetLogsByHabitID(habitID int64) ([]*models.Log, error)
	UpdateLog(log *models.Log) error
	DeleteLog(id int64) error
}

type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new LogRepository
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// CreateLog creates a new log in the database
func (r *logRepository) CreateLog(log *models.Log) error {
	query := `INSERT INTO logs (habit_id, notes) VALUES (?, ?)`
	res, err := r.db.Exec(query, log.HabitID, log.Notes)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	log.ID = id
	return nil
}

// GetLogByID retrieves a log from the database by ID
func (r *logRepository) GetLogByID(id int64) (*models.Log, error) {
	log := &models.Log{}
	query := `SELECT id, habit_id, notes, created_at FROM logs WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&log.ID, &log.HabitID, &log.Notes, &log.CreatedAt)
	if err != nil {
		return nil, err
	}
	return log, nil
}

// GetLogsByHabitID retrieves all logs for a habit
func (r *logRepository) GetLogsByHabitID(habitID int64) ([]*models.Log, error) {
	query := `SELECT id, habit_id, notes, created_at FROM logs WHERE habit_id = ?`
	rows, err := r.db.Query(query, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		log := &models.Log{}
		err := rows.Scan(&log.ID, &log.HabitID, &log.Notes, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// UpdateLog updates a log in the database
func (r *logRepository) UpdateLog(log *models.Log) error {
	query := `UPDATE logs SET notes = ? WHERE id = ?`
	_, err := r.db.Exec(query, log.Notes, log.ID)
	return err
}

// DeleteLog deletes a log from the database
func (r *logRepository) DeleteLog(id int64) error {
	query := `DELETE FROM logs WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
