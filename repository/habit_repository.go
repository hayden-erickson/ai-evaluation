package repository

import (
	"database/sql"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// HabitRepository defines the interface for habit data operations
type HabitRepository interface {
	Create(habit *models.Habit) error
	GetByID(id int64) (*models.Habit, error)
	GetByUserID(userID int64) ([]*models.Habit, error)
	Update(habit *models.Habit) error
	Delete(id int64) error
}

type habitRepository struct {
	db *sql.DB
}

// NewHabitRepository creates a new HabitRepository instance
func NewHabitRepository(db *sql.DB) HabitRepository {
	return &habitRepository{db: db}
}

// Create inserts a new habit into the database
func (r *habitRepository) Create(habit *models.Habit) error {
	query := `
		INSERT INTO habits (user_id, name, description, duration_seconds)
		VALUES (?, ?, ?, ?)
	`
	var duration interface{}
	if habit.DurationSeconds != nil {
		duration = *habit.DurationSeconds
	} else {
		duration = nil
	}
	result, err := r.db.Exec(query, habit.UserID, habit.Name, habit.Description, duration)
	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	habit.ID = id

	// Retrieve the created_at timestamp
	createdHabit, err := r.GetByID(id)
	if err != nil {
		return err
	}
	habit.CreatedAt = createdHabit.CreatedAt

	return nil
}

// GetByID retrieves a habit by its ID
func (r *habitRepository) GetByID(id int64) (*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, duration_seconds, created_at
		FROM habits
		WHERE id = ?
	`
	habit := &models.Habit{}
	var duration sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&duration,
		&habit.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("habit not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	if duration.Valid {
		v := int(duration.Int64)
		habit.DurationSeconds = &v
	}

	return habit, nil
}

// GetByUserID retrieves all habits for a specific user
func (r *habitRepository) GetByUserID(userID int64) ([]*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, duration_seconds, created_at
		FROM habits
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()

	var habits []*models.Habit
	for rows.Next() {
		habit := &models.Habit{}
		var duration sql.NullInt64
		if err := rows.Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &duration, &habit.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		if duration.Valid {
			v := int(duration.Int64)
			habit.DurationSeconds = &v
		}
		habits = append(habits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

// Update updates an existing habit in the database
func (r *habitRepository) Update(habit *models.Habit) error {
	query := `
		UPDATE habits
		SET name = ?, description = ?, duration_seconds = ?
		WHERE id = ?
	`
	var duration interface{}
	if habit.DurationSeconds != nil {
		duration = *habit.DurationSeconds
	} else {
		duration = nil
	}
	result, err := r.db.Exec(query, habit.Name, habit.Description, duration, habit.ID)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}

// Delete removes a habit from the database
func (r *habitRepository) Delete(id int64) error {
	query := `DELETE FROM habits WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}
