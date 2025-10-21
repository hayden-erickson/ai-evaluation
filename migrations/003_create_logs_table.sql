-- Migration: Create logs table
-- Description: Stores log entries for habit tracking
-- Created: 2025-10-20

CREATE TABLE IF NOT EXISTS logs (
    id TEXT PRIMARY KEY,
    habit_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_logs_habit_id ON logs(habit_id);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
