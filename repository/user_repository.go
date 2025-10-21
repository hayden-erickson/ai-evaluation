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
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	Update(id int64, update *models.UserUpdateRequest) error
	Delete(id int64) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (profile_image_url, name, time_zone, phone_number, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.Exec(query, user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber, user.PasswordHash, user.CreatedAt)
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
	
	var user models.User
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
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

// GetByPhoneNumber retrieves a user by their phone number
func (r *userRepository) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, profile_image_url, name, time_zone, phone_number, password_hash, created_at
		FROM users
		WHERE phone_number = ?
	`
	
	var user models.User
	err := r.db.QueryRow(query, phoneNumber).Scan(
		&user.ID,
		&user.ProfileImageURL,
		&user.Name,
		&user.TimeZone,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone number: %w", err)
	}
	
	return &user, nil
}

// Update updates a user's information
func (r *userRepository) Update(id int64, update *models.UserUpdateRequest) error {
	// Build dynamic update query based on provided fields
	query := "UPDATE users SET "
	args := []interface{}{}
	updates := []string{}
	
	if update.ProfileImageURL != nil {
		updates = append(updates, "profile_image_url = ?")
		args = append(args, *update.ProfileImageURL)
	}
	if update.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *update.Name)
	}
	if update.TimeZone != nil {
		updates = append(updates, "time_zone = ?")
		args = append(args, *update.TimeZone)
	}
	if update.PhoneNumber != nil {
		updates = append(updates, "phone_number = ?")
		args = append(args, *update.PhoneNumber)
	}
	if update.Password != nil {
		updates = append(updates, "password_hash = ?")
		args = append(args, *update.Password)
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
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

// Delete removes a user from the database
func (r *userRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}
