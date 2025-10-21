-- Create logs table
-- Stores log entries for habits
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

-- Create index on habit_id for faster queries when fetching habit's logs
CREATE INDEX IF NOT EXISTS idx_logs_habit_id ON logs(habit_id);
