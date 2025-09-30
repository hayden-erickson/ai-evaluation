-- +migrate Up
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX idx_tags_habit_id ON tags(habit_id);
CREATE UNIQUE INDEX idx_tags_habit_value ON tags(habit_id, value);

-- +migrate Down
DROP TABLE tags;