-- Add duration_seconds column to habits table
ALTER TABLE habits ADD COLUMN duration_seconds INT;

-- Add duration_seconds column to logs table
ALTER TABLE logs ADD COLUMN duration_seconds INT;
