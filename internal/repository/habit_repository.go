package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

type HabitRepository interface {
	Create(ctx context.Context, h *models.Habit) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Habit, error)
	ListByUser(ctx context.Context, userID int64) ([]*models.Habit, error)
	Update(ctx context.Context, h *models.Habit) error
	Delete(ctx context.Context, id int64) error
}

type habitRepository struct{ db *sql.DB }

func NewHabitRepository(db *sql.DB) HabitRepository { return &habitRepository{db: db} }

func (r *habitRepository) Create(ctx context.Context, h *models.Habit) (int64, error) {
	var dur any
	if h.DurationSeconds != nil { dur = *h.DurationSeconds } else { dur = nil }
	res, err := r.db.ExecContext(ctx, `INSERT INTO habits (user_id, name, description, duration_seconds, created_at) VALUES (?, ?, ?, ?, ?)`,
		h.UserID, h.Name, h.Description, dur, time.Now().UTC())
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (r *habitRepository) GetByID(ctx context.Context, id int64) (*models.Habit, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, description, duration_seconds, created_at FROM habits WHERE id = ?`, id)
	h := &models.Habit{}
	var dur sql.NullInt64
	if err := row.Scan(&h.ID, &h.UserID, &h.Name, &h.Description, &dur, &h.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	if dur.Valid { v := dur.Int64; h.DurationSeconds = &v }
	return h, nil
}

func (r *habitRepository) ListByUser(ctx context.Context, userID int64) ([]*models.Habit, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, description, duration_seconds, created_at FROM habits WHERE user_id = ? ORDER BY id DESC`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.Habit
	for rows.Next() {
		h := &models.Habit{}
		var dur sql.NullInt64
		if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &h.Description, &dur, &h.CreatedAt); err != nil { return nil, err }
		if dur.Valid { v := dur.Int64; h.DurationSeconds = &v }
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *habitRepository) Update(ctx context.Context, h *models.Habit) error {
	var dur any
	if h.DurationSeconds != nil { dur = *h.DurationSeconds } else { dur = nil }
	_, err := r.db.ExecContext(ctx, `UPDATE habits SET name = ?, description = ?, duration_seconds = ? WHERE id = ? AND user_id = ?`, h.Name, h.Description, dur, h.ID, h.UserID)
	return err
}

func (r *habitRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM habits WHERE id = ?`, id)
	return err
}
