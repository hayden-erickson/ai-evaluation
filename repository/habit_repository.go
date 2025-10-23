package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// HabitRepository defines the interface for habit data operations
type HabitRepository interface {
	Create(userID int64, habit *models.CreateHabitRequest) (*models.Habit, error)
	GetByID(id int64) (*models.Habit, error)
	GetByUserID(userID int64) ([]*models.Habit, error)
	Update(id int64, req *models.UpdateHabitRequest) (*models.Habit, error)
	Delete(id int64) error
}

// habitRepository implements HabitRepository
type habitRepository struct {
	db *sql.DB
}

// NewHabitRepository creates a new habit repository
func NewHabitRepository(db *sql.DB) HabitRepository {
	return &habitRepository{db: db}
}

// Create creates a new habit in the database
func (r *habitRepository) Create(userID int64, habit *models.CreateHabitRequest) (*models.Habit, error) {
	// Insert the habit
	result, err := r.db.Exec(
		"INSERT INTO habits (user_id, name, description, duration_seconds) VALUES (?, ?, ?, ?)",
		userID, habit.Name, habit.Description, habit.DurationSeconds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted habit ID: %w", err)
	}

	// Retrieve and return the created habit
	return r.GetByID(id)
}

// GetByID retrieves a habit by ID
func (r *habitRepository) GetByID(id int64) (*models.Habit, error) {
	habit := &models.Habit{}
	var createdAt string
	var durationSeconds sql.NullInt64

	// Query the habit
	err := r.db.QueryRow(
		"SELECT id, user_id, name, description, duration_seconds, created_at FROM habits WHERE id = ?",
		id,
	).Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &durationSeconds, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("habit not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	// Set duration_seconds if not null
	if durationSeconds.Valid {
		duration := int(durationSeconds.Int64)
		habit.DurationSeconds = &duration
	}

	// Parse created_at timestamp
	habit.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Try alternative format
		habit.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
	}

	return habit, nil
}

// GetByUserID retrieves all habits for a user
func (r *habitRepository) GetByUserID(userID int64) ([]*models.Habit, error) {
	// Query all habits for the user
	rows, err := r.db.Query(
		"SELECT id, user_id, name, description, duration_seconds, created_at FROM habits WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()

	habits := []*models.Habit{}
	for rows.Next() {
		habit := &models.Habit{}
		var createdAt string
		var durationSeconds sql.NullInt64

		err := rows.Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &durationSeconds, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		// Set duration_seconds if not null
		if durationSeconds.Valid {
			duration := int(durationSeconds.Int64)
			habit.DurationSeconds = &duration
		}

		// Parse created_at timestamp
		habit.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			// Try alternative format
			habit.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse created_at: %w", err)
			}
		}

		habits = append(habits, habit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

// Update updates a habit in the database
func (r *habitRepository) Update(id int64, req *models.UpdateHabitRequest) (*models.Habit, error) {
	// Build dynamic update query
	query := "UPDATE habits SET "
	args := []interface{}{}
	updates := []string{}

	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
	}
	if req.DurationSeconds != nil {
		updates = append(updates, "duration_seconds = ?")
		args = append(args, *req.DurationSeconds)
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
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	// Retrieve and return the updated habit
	return r.GetByID(id)
}

// Delete deletes a habit from the database
func (r *habitRepository) Delete(id int64) error {
	// Delete the habit
	result, err := r.db.Exec("DELETE FROM habits WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}
