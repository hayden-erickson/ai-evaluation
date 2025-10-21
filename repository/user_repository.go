package repository

import (
	"database/sql"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByName(name string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (profile_image_url, name, time_zone, phone_number, password_hash)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, profile_image_url, name, time_zone, phone_number, password_hash, created_at
		FROM users
		WHERE id = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.ProfileImageURL,
		&user.Name,
		&user.TimeZone,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByName retrieves a user by their name (for authentication)
func (r *userRepository) GetByName(name string) (*models.User, error) {
	query := `
		SELECT id, profile_image_url, name, time_zone, phone_number, password_hash, created_at
		FROM users
		WHERE name = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, name).Scan(
		&user.ID,
		&user.ProfileImageURL,
		&user.Name,
		&user.TimeZone,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user in the database
func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET profile_image_url = ?, name = ?, time_zone = ?, phone_number = ?, password_hash = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete removes a user from the database
func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
