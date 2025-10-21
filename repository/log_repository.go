package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// LogRepository defines the interface for log data operations
type LogRepository interface {
	Create(habitID int64, log *models.CreateLogRequest) (*models.Log, error)
	GetByID(id int64) (*models.Log, error)
	GetByHabitID(habitID int64) ([]*models.Log, error)
	Update(id int64, req *models.UpdateLogRequest) (*models.Log, error)
	Delete(id int64) error
}

// logRepository implements LogRepository
type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new log repository
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// Create creates a new log in the database
func (r *logRepository) Create(habitID int64, log *models.CreateLogRequest) (*models.Log, error) {
	// Insert the log
	result, err := r.db.Exec(
		"INSERT INTO logs (habit_id, notes) VALUES (?, ?)",
		habitID, log.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted log ID: %w", err)
	}

	// Retrieve and return the created log
	return r.GetByID(id)
}

// GetByID retrieves a log by ID
func (r *logRepository) GetByID(id int64) (*models.Log, error) {
	log := &models.Log{}
	var createdAt string

	// Query the log
	err := r.db.QueryRow(
		"SELECT id, habit_id, notes, created_at FROM logs WHERE id = ?",
		id,
	).Scan(&log.ID, &log.HabitID, &log.Notes, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Parse created_at timestamp
	log.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Try alternative format
		log.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
	}

	return log, nil
}

// GetByHabitID retrieves all logs for a habit
func (r *logRepository) GetByHabitID(habitID int64) ([]*models.Log, error) {
	// Query all logs for the habit
	rows, err := r.db.Query(
		"SELECT id, habit_id, notes, created_at FROM logs WHERE habit_id = ? ORDER BY created_at DESC",
		habitID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	logs := []*models.Log{}
	for rows.Next() {
		log := &models.Log{}
		var createdAt string

		err := rows.Scan(&log.ID, &log.HabitID, &log.Notes, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		// Parse created_at timestamp
		log.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			// Try alternative format
			log.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse created_at: %w", err)
			}
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// Update updates a log in the database
func (r *logRepository) Update(id int64, req *models.UpdateLogRequest) (*models.Log, error) {
	// Build dynamic update query
	query := "UPDATE logs SET "
	args := []interface{}{}
	updates := []string{}

	if req.Notes != nil {
		updates = append(updates, "notes = ?")
		args = append(args, *req.Notes)
	}

	// No fields to update
	if len(updates) == 0 {
		return r.GetByID(id)
	}

	// Build the complete query
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	args = append(args, id)

	// Execute the update
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update log: %w", err)
	}

	// Retrieve and return the updated log
	return r.GetByID(id)
}

// Delete deletes a log from the database
func (r *logRepository) Delete(id int64) error {
	// Delete the log
	result, err := r.db.Exec("DELETE FROM logs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("log not found")
	}

	return nil
}
