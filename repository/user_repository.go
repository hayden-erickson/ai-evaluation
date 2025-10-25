package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.CreateUserRequest, passwordHash string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	Update(id int64, req *models.UpdateUserRequest, passwordHash *string) (*models.User, error)
	Delete(id int64) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(user *models.CreateUserRequest, passwordHash string) (*models.User, error) {
	// Insert the user
	result, err := r.db.Exec(
		"INSERT INTO users (profile_image_url, name, time_zone, phone_number, password_hash) VALUES (?, ?, ?, ?, ?)",
		user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber, passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted user ID: %w", err)
	}

	// Retrieve and return the created user
	return r.GetByID(id)
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	var createdAt string

	// Query the user
	err := r.db.QueryRow(
		"SELECT id, profile_image_url, name, time_zone, phone_number, password_hash, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.ProfileImageURL, &user.Name, &user.TimeZone, &user.PhoneNumber, &user.PasswordHash, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse created_at timestamp
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Try alternative format
		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
	}

	return user, nil
}

// GetByPhoneNumber retrieves a user by phone number
func (r *userRepository) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := &models.User{}
	var createdAt string

	// Query the user
	err := r.db.QueryRow(
		"SELECT id, profile_image_url, name, time_zone, phone_number, password_hash, created_at FROM users WHERE phone_number = ?",
		phoneNumber,
	).Scan(&user.ID, &user.ProfileImageURL, &user.Name, &user.TimeZone, &user.PhoneNumber, &user.PasswordHash, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse created_at timestamp
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Try alternative format
		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
	}

	return user, nil
}

// Update updates a user in the database
func (r *userRepository) Update(id int64, req *models.UpdateUserRequest, passwordHash *string) (*models.User, error) {
	// Build dynamic update query
	query := "UPDATE users SET "
	args := []interface{}{}
	updates := []string{}

	if req.ProfileImageURL != nil {
		updates = append(updates, "profile_image_url = ?")
		args = append(args, *req.ProfileImageURL)
	}
	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.TimeZone != nil {
		updates = append(updates, "time_zone = ?")
		args = append(args, *req.TimeZone)
	}
	if req.PhoneNumber != nil {
		updates = append(updates, "phone_number = ?")
		args = append(args, *req.PhoneNumber)
	}
	if passwordHash != nil {
		updates = append(updates, "password_hash = ?")
		args = append(args, *passwordHash)
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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Retrieve and return the updated user
	return r.GetByID(id)
}

// Delete deletes a user from the database
func (r *userRepository) Delete(id int64) error {
	// Delete the user
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
