-- Create logs table
CREATE TABLE IF NOT EXISTS logs (
    id TEXT PRIMARY KEY,
    habit_id TEXT NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

-- Create index on habit_id for efficient querying
CREATE INDEX IF NOT EXISTS idx_logs_habit_id ON logs(habit_id);

-- Create index on created_at for efficient querying
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
