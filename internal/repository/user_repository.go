package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	UpdateSelf(ctx context.Context, u *models.User) error
	DeleteSelf(ctx context.Context, id int64) error
}

type userRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(ctx context.Context, u *models.User) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO users (email, password_hash, profile_image_url, name, time_zone, phone_number, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email, u.PasswordHash, u.ProfileImageURL, u.Name, u.TimeZone, u.Phone, u.Role, time.Now().UTC())
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, profile_image_url, name, time_zone, phone_number, role, created_at FROM users WHERE email = ?`, email)
	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.ProfileImageURL, &u.Name, &u.TimeZone, &u.Phone, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, profile_image_url, name, time_zone, phone_number, role, created_at FROM users WHERE id = ?`, id)
	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.ProfileImageURL, &u.Name, &u.TimeZone, &u.Phone, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return u, nil
}

func (r *userRepository) UpdateSelf(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET profile_image_url = ?, name = ?, time_zone = ?, phone_number = ? WHERE id = ?`,
		u.ProfileImageURL, u.Name, u.TimeZone, u.Phone, u.ID)
	return err
}

func (r *userRepository) DeleteSelf(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}
