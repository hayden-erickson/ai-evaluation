package repository

import (
	"database/sql"
	"new-api/internal/models"
)

// HabitRepository defines the interface for habit persistence
type HabitRepository interface {
	CreateHabit(habit *models.Habit) error
	GetHabitByID(id int64) (*models.Habit, error)
	GetHabitsByUserID(userID int64) ([]*models.Habit, error)
	UpdateHabit(habit *models.Habit) error
	DeleteHabit(id int64) error
}

type habitRepository struct {
	db *sql.DB
}

// NewHabitRepository creates a new HabitRepository
func NewHabitRepository(db *sql.DB) HabitRepository {
	return &habitRepository{db: db}
}

// CreateHabit creates a new habit in the database
func (r *habitRepository) CreateHabit(habit *models.Habit) error {
	query := `INSERT INTO habits (user_id, name, description) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, habit.UserID, habit.Name, habit.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	habit.ID = id
	return nil
}

// GetHabitByID retrieves a habit from the database by ID
func (r *habitRepository) GetHabitByID(id int64) (*models.Habit, error) {
	habit := &models.Habit{}
	query := `SELECT id, user_id, name, description, created_at FROM habits WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &habit.CreatedAt)
	if err != nil {
		return nil, err
	}
	return habit, nil
}

// GetHabitsByUserID retrieves all habits for a user
func (r *habitRepository) GetHabitsByUserID(userID int64) ([]*models.Habit, error) {
	query := `SELECT id, user_id, name, description, created_at FROM habits WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []*models.Habit
	for rows.Next() {
		habit := &models.Habit{}
		err := rows.Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &habit.CreatedAt)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	return habits, nil
}

// UpdateHabit updates a habit in the database
func (r *habitRepository) UpdateHabit(habit *models.Habit) error {
	query := `UPDATE habits SET name = ?, description = ? WHERE id = ?`
	_, err := r.db.Exec(query, habit.Name, habit.Description, habit.ID)
	return err
}

// DeleteHabit deletes a habit from the database
func (r *habitRepository) DeleteHabit(id int64) error {
	query := `DELETE FROM habits WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
