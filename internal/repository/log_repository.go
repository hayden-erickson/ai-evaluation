package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

type LogRepository interface {
	Create(ctx context.Context, l *models.LogEntry) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.LogEntry, error)
	ListByHabit(ctx context.Context, habitID int64) ([]*models.LogEntry, error)
	Update(ctx context.Context, l *models.LogEntry) error
	Delete(ctx context.Context, id int64) error
}

type logRepository struct{ db *sql.DB }

func NewLogRepository(db *sql.DB) LogRepository { return &logRepository{db: db} }

func (r *logRepository) Create(ctx context.Context, l *models.LogEntry) (int64, error) {
	var dur any
	if l.DurationSeconds != nil { dur = *l.DurationSeconds } else { dur = nil }
	res, err := r.db.ExecContext(ctx, `INSERT INTO logs (habit_id, notes, duration_seconds, created_at) VALUES (?, ?, ?, ?)`,
		l.HabitID, l.Notes, dur, time.Now().UTC())
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (r *logRepository) GetByID(ctx context.Context, id int64) (*models.LogEntry, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, habit_id, notes, duration_seconds, created_at FROM logs WHERE id = ?`, id)
	l := &models.LogEntry{}
	var dur sql.NullInt64
	if err := row.Scan(&l.ID, &l.HabitID, &l.Notes, &dur, &l.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	if dur.Valid { v := dur.Int64; l.DurationSeconds = &v }
	return l, nil
}

func (r *logRepository) ListByHabit(ctx context.Context, habitID int64) ([]*models.LogEntry, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, habit_id, notes, duration_seconds, created_at FROM logs WHERE habit_id = ? ORDER BY id DESC`, habitID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.LogEntry
	for rows.Next() {
		l := &models.LogEntry{}
		var dur sql.NullInt64
		if err := rows.Scan(&l.ID, &l.HabitID, &l.Notes, &dur, &l.CreatedAt); err != nil { return nil, err }
		if dur.Valid { v := dur.Int64; l.DurationSeconds = &v }
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *logRepository) Update(ctx context.Context, l *models.LogEntry) error {
	var dur any
	if l.DurationSeconds != nil { dur = *l.DurationSeconds } else { dur = nil }
	_, err := r.db.ExecContext(ctx, `UPDATE logs SET notes = ?, duration_seconds = ? WHERE id = ?`, l.Notes, dur, l.ID)
	return err
}

func (r *logRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM logs WHERE id = ?`, id)
	return err
}
