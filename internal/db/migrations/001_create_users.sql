CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    profile_image_url TEXT,
    name TEXT,
    time_zone TEXT,
    phone_number TEXT,
    role TEXT NOT NULL,
    created_at DATETIME NOT NULL
);
