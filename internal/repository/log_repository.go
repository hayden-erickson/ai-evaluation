package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// LogRepository defines the interface for log data operations
type LogRepository interface {
	Create(ctx context.Context, log *models.Log) error
	GetByID(ctx context.Context, id string) (*models.Log, error)
	GetByHabitID(ctx context.Context, habitID string) ([]*models.Log, error)
	Update(ctx context.Context, log *models.Log) error
	Delete(ctx context.Context, id string) error
}

// logRepository implements LogRepository interface
type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new log repository
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// Create inserts a new log into the database
func (r *logRepository) Create(ctx context.Context, log *models.Log) error {
	query := `
		INSERT INTO logs (id, habit_id, notes, duration_seconds, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.HabitID,
		log.Notes,
		log.DurationSeconds,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

// GetByID retrieves a log by its ID
func (r *logRepository) GetByID(ctx context.Context, id string) (*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, duration_seconds, created_at
		FROM logs
		WHERE id = ?
	`

	log := &models.Log{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.HabitID,
		&log.Notes,
		&log.DurationSeconds,
		&log.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a specific habit
func (r *logRepository) GetByHabitID(ctx context.Context, habitID string) ([]*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, duration_seconds, created_at
		FROM logs
		WHERE habit_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.Log, 0)
	for rows.Next() {
		log := &models.Log{}
		err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.Notes,
			&log.DurationSeconds,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// Update modifies an existing log in the database
func (r *logRepository) Update(ctx context.Context, log *models.Log) error {
	query := `
		UPDATE logs
		SET notes = ?, duration_seconds = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		log.Notes,
		log.DurationSeconds,
		log.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}

// Delete removes a log from the database
func (r *logRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM logs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}
