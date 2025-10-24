package repository

import (
	"database/sql"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// LogRepository defines the interface for log data operations
type LogRepository interface {
	Create(log *models.Log) error
	GetByID(id int64) (*models.Log, error)
	GetByHabitID(habitID int64) ([]*models.Log, error)
	Update(log *models.Log) error
	Delete(id int64) error
}

type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new LogRepository instance
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// Create inserts a new log into the database
func (r *logRepository) Create(log *models.Log) error {
	query := `
		INSERT INTO logs (habit_id, notes, duration_seconds)
		VALUES (?, ?, ?)
	`
	var duration interface{}
	if log.DurationSeconds != nil {
		duration = *log.DurationSeconds
	} else {
		duration = nil
	}
	result, err := r.db.Exec(query, log.HabitID, log.Notes, duration)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.ID = id

	// Retrieve the created_at timestamp
	createdLog, err := r.GetByID(id)
	if err != nil {
		return err
	}
	log.CreatedAt = createdLog.CreatedAt

	return nil
}

// GetByID retrieves a log by its ID
func (r *logRepository) GetByID(id int64) (*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, duration_seconds, created_at
		FROM logs
		WHERE id = ?
	`
	log := &models.Log{}
	var duration sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.HabitID,
		&log.Notes,
		&duration,
		&log.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	if duration.Valid {
		v := int(duration.Int64)
		log.DurationSeconds = &v
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a specific habit
func (r *logRepository) GetByHabitID(habitID int64) ([]*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, duration_seconds, created_at
		FROM logs
		WHERE habit_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		log := &models.Log{}
		var duration sql.NullInt64
		if err := rows.Scan(&log.ID, &log.HabitID, &log.Notes, &duration, &log.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		if duration.Valid {
			v := int(duration.Int64)
			log.DurationSeconds = &v
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// Update updates an existing log in the database
func (r *logRepository) Update(log *models.Log) error {
	query := `
		UPDATE logs
		SET notes = ?, duration_seconds = ?
		WHERE id = ?
	`
	var duration interface{}
	if log.DurationSeconds != nil {
		duration = *log.DurationSeconds
	} else {
		duration = nil
	}
	result, err := r.db.Exec(query, log.Notes, duration, log.ID)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}

// Delete removes a log from the database
func (r *logRepository) Delete(id int64) error {
	query := `DELETE FROM logs WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}
