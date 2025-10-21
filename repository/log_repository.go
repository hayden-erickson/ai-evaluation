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
	Update(id int64, update *models.LogUpdateRequest) error
	Delete(id int64) error
}

// logRepository implements LogRepository interface
type logRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new log repository instance
func NewLogRepository(db *sql.DB) LogRepository {
	return &logRepository{db: db}
}

// Create inserts a new log into the database
func (r *logRepository) Create(log *models.Log) error {
	query := `
		INSERT INTO logs (habit_id, notes, created_at)
		VALUES (?, ?, ?)
	`
	
	result, err := r.db.Exec(query, log.HabitID, log.Notes, log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	log.ID = id
	return nil
}

// GetByID retrieves a log by its ID
func (r *logRepository) GetByID(id int64) (*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, created_at
		FROM logs
		WHERE id = ?
	`
	
	var log models.Log
	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.HabitID,
		&log.Notes,
		&log.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}
	
	return &log, nil
}

// GetByHabitID retrieves all logs for a specific habit
func (r *logRepository) GetByHabitID(habitID int64) ([]*models.Log, error) {
	query := `
		SELECT id, habit_id, notes, created_at
		FROM logs
		WHERE habit_id = ?
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()
	
	logs := []*models.Log{}
	for rows.Next() {
		var log models.Log
		err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.Notes,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, &log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}
	
	return logs, nil
}

// Update updates a log's information
func (r *logRepository) Update(id int64, update *models.LogUpdateRequest) error {
	// Build dynamic update query based on provided fields
	query := "UPDATE logs SET "
	args := []interface{}{}
	updates := []string{}
	
	if update.Notes != nil {
		updates = append(updates, "notes = ?")
		args = append(args, *update.Notes)
	}
	
	// If no fields to update, return early
	if len(updates) == 0 {
		return nil
	}
	
	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, id)
	
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}
	
	return nil
}

// Delete removes a log from the database
func (r *logRepository) Delete(id int64) error {
	query := "DELETE FROM logs WHERE id = ?"
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}
	
	return nil
}
