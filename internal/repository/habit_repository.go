package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// HabitRepository defines the interface for habit data operations
type HabitRepository interface {
	Create(ctx context.Context, habit *models.Habit) error
	GetByID(ctx context.Context, id string) (*models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string) ([]*models.Habit, error)
}

// habitRepository implements HabitRepository interface
type habitRepository struct {
	db *sql.DB
}

// NewHabitRepository creates a new habit repository instance
func NewHabitRepository(db *sql.DB) HabitRepository {
	return &habitRepository{db: db}
}

// Create inserts a new habit into the database
func (r *habitRepository) Create(ctx context.Context, habit *models.Habit) error {
	query := `
		INSERT INTO habits (id, user_id, name, description, duration, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		habit.ID,
		habit.UserID,
		habit.Name,
		habit.Description,
		habit.Duration,
		habit.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}
	return nil
}

// GetByID retrieves a habit by its ID
func (r *habitRepository) GetByID(ctx context.Context, id string) (*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, duration, created_at
		FROM habits
		WHERE id = ?
	`
	habit := &models.Habit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.Duration,
		&habit.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("habit not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	return habit, nil
}

// Update modifies an existing habit in the database
func (r *habitRepository) Update(ctx context.Context, habit *models.Habit) error {
	query := `
		UPDATE habits
		SET name = ?, description = ?, duration = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		habit.Name,
		habit.Description,
		habit.Duration,
		habit.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}
	
	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}
	
	return nil
}

// Delete removes a habit from the database
func (r *habitRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM habits WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	
	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}
	
	return nil
}

// ListByUserID retrieves all habits for a specific user
func (r *habitRepository) ListByUserID(ctx context.Context, userID string) ([]*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, duration, created_at
		FROM habits
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}
	defer rows.Close()

	var habits []*models.Habit
	for rows.Next() {
		habit := &models.Habit{}
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Name,
			&habit.Description,
			&habit.Duration,
			&habit.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habits = append(habits, habit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}
