-- Migration: Add duration_seconds to habits and logs
ALTER TABLE habits ADD COLUMN duration_seconds INTEGER;
ALTER TABLE logs ADD COLUMN duration_seconds INTEGER;
