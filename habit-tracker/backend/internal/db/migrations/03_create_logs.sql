-- +migrate Up
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX idx_logs_habit_id ON logs(habit_id);
CREATE INDEX idx_logs_created_at ON logs(created_at);

-- +migrate Down
DROP TABLE logs;