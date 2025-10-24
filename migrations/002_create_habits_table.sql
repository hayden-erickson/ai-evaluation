-- Create habits table
CREATE TABLE IF NOT EXISTS habits (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on user_id for efficient querying
CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);

-- Create index on created_at for efficient querying
CREATE INDEX IF NOT EXISTS idx_habits_created_at ON habits(created_at);
