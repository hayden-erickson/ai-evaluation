-- Add optional duration in seconds to habits and logs
ALTER TABLE habits ADD COLUMN duration_seconds INTEGER;
ALTER TABLE logs ADD COLUMN duration_seconds INTEGER;

-- For data integrity: nothing enforced at DB layer per requirement;
-- application code will enforce: if habit has duration, corresponding log must as well.


