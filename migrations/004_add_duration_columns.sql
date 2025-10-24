-- Add duration column to habits table
ALTER TABLE habits ADD COLUMN duration INTEGER;

-- Add duration column to logs table
ALTER TABLE logs ADD COLUMN duration INTEGER;
