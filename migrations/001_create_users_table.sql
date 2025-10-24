-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    profile_image_url TEXT,
    name TEXT NOT NULL,
    time_zone TEXT NOT NULL,
    phone_number TEXT,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on created_at for efficient querying
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
