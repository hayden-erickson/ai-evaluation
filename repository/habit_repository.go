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
	Update(id int64, update *models.HabitUpdateRequest) error
	Delete(id int64) error
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
func (r *habitRepository) Create(habit *models.Habit) error {
	query := `
		INSERT INTO habits (user_id, name, description, created_at)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.Exec(query, habit.UserID, habit.Name, habit.Description, habit.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	habit.ID = id
	return nil
}

// GetByID retrieves a habit by its ID
func (r *habitRepository) GetByID(id int64) (*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, created_at
		FROM habits
		WHERE id = ?
	`
	
	var habit models.Habit
	err := r.db.QueryRow(query, id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	
	return &habit, nil
}

// GetByUserID retrieves all habits for a specific user
func (r *habitRepository) GetByUserID(userID int64) ([]*models.Habit, error) {
	query := `
		SELECT id, user_id, name, description, created_at
		FROM habits
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()
	
	habits := []*models.Habit{}
	for rows.Next() {
		var habit models.Habit
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Name,
			&habit.Description,
			&habit.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habits = append(habits, &habit)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}
	
	return habits, nil
}

// Update updates a habit's information
func (r *habitRepository) Update(id int64, update *models.HabitUpdateRequest) error {
	// Build dynamic update query based on provided fields
	query := "UPDATE habits SET "
	args := []interface{}{}
	updates := []string{}
	
	if update.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *update.Name)
	}
	if update.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *update.Description)
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
		return fmt.Errorf("failed to update habit: %w", err)
	}
	
	return nil
}

// Delete removes a habit from the database
func (r *habitRepository) Delete(id int64) error {
	query := "DELETE FROM habits WHERE id = ?"
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	
	return nil
}
