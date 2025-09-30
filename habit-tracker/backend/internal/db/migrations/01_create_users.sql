-- +migrate Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_image_url TEXT,
    name TEXT NOT NULL,
    time_zone TEXT NOT NULL,
    phone TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    google_id TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL
);

-- +migrate Down
DROP TABLE users;