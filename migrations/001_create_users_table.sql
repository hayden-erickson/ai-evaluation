-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_image_url TEXT,
    name TEXT NOT NULL,
    time_zone TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create index on phone_number for faster lookups during authentication
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);
