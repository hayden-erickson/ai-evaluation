-- +migrate Up
ALTER TABLE logs ADD COLUMN completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +migrate Down
ALTER TABLE logs DROP COLUMN completed_at;