package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// LogRepository defines the interface for log data access operations
type LogRepository interface {
	Create(ctx context.Context, log *models.Log) error
	GetByID(ctx context.Context, id string) (*models.Log, error)
	ListByHabitID(ctx context.Context, habitID string) ([]*models.Log, error)
	Update(ctx context.Context, log *models.Log) error
	Delete(ctx context.Context, id string) error
}

// logRepository implements LogRepository interface
type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new instance of LogRepository
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// Create inserts a new log entry into the database
func (r *logRepository) Create(ctx context.Context, log *models.Log) error {
	// Generate UUID if not provided
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	// Set created timestamp
	log.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO logs (id, habit_id, created_at, notes)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.HabitID,
		log.CreatedAt,
		log.Notes,
	)

	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

// GetByID retrieves a log entry by its ID
func (r *logRepository) GetByID(ctx context.Context, id string) (*models.Log, error) {
	query := `
		SELECT id, habit_id, created_at, notes
		FROM logs
		WHERE id = ?
	`

	log := &models.Log{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.HabitID,
		&log.CreatedAt,
		&log.Notes,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	return log, nil
}

// ListByHabitID retrieves all log entries for a specific habit
func (r *logRepository) ListByHabitID(ctx context.Context, habitID string) ([]*models.Log, error) {
	query := `
		SELECT id, habit_id, created_at, notes
		FROM logs
		WHERE habit_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.Log, 0)

	// Iterate through result set
	for rows.Next() {
		log := &models.Log{}
		err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.CreatedAt,
			&log.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// Update modifies an existing log entry in the database
func (r *logRepository) Update(ctx context.Context, log *models.Log) error {
	query := `
		UPDATE logs
		SET notes = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		log.Notes,
		log.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Check if log exists
	if rows == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}

// Delete removes a log entry from the database
func (r *logRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM logs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Check if log exists
	if rows == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}
