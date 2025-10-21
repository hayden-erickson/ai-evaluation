package repository

import (
	"database/sql"
	"new-api/internal/models"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (profile_image_url, name, time_zone, phone_number) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// GetUserByID retrieves a user from the database by ID
func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, profile_image_url, name, time_zone, phone_number, created_at FROM users WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.ProfileImageURL, &user.Name, &user.TimeZone, &user.PhoneNumber, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user in the database
func (r *userRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET profile_image_url = ?, name = ?, time_zone = ?, phone_number = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.ProfileImageURL, user.Name, user.TimeZone, user.PhoneNumber, user.ID)
	return err
}

// DeleteUser deletes a user from the database
func (r *userRepository) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
