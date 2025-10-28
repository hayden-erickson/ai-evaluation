# Duration Feature Implementation Summary

## Overview
Added optional duration tracking (in seconds) to both habits and logs with validation logic that enforces:
- If a habit has a duration, corresponding logs must also have a duration
- If a habit has no duration, log durations are optional
- All durations must be non-negative

## Changes Made

### 1. Models Updated
**Files Modified:**
- `internal/models/habit.go`
- `internal/models/log.go`

**Changes:**
- Added `Duration *int` field to `Habit`, `CreateHabitRequest`, and `UpdateHabitRequest` structs
- Added `Duration *int` field to `Log`, `CreateLogRequest`, and `UpdateLogRequest` structs
- Duration is stored as a pointer to allow for optional values (nil = no duration)

### 2. Database Migration
**New File:**
- `migrations/004_add_duration_columns.sql`

**Changes:**
- Added `duration INTEGER` column to `habits` table
- Added `duration INTEGER` column to `logs` table

### 3. Validation Logic
**New File:**
- `pkg/validator/duration_validator.go`
  - Contains `ValidateLogDuration()` function that enforces the business rule

**Modified File:**
- `pkg/validator/validator.go`
  - Updated `ValidateCreateHabitRequest()` to validate non-negative duration

### 4. Service Layer
**Files Modified:**
- `internal/service/habit_service.go`
  - Updated `Create()` to include duration field
  - Updated `Update()` to handle duration updates with validation
  
- `internal/service/log_service.go`
  - Updated `Create()` to validate log duration against habit requirements
  - Updated `Update()` to validate log duration against habit requirements

### 5. Repository Layer
**Files Modified:**
- `internal/repository/habit_repository.go`
  - Updated all SQL queries to include `duration` column
  - Modified `Create()`, `GetByID()`, `Update()`, and `ListByUserID()` methods

- `internal/repository/log_repository.go`
  - Updated all SQL queries to include `duration` column
  - Modified `Create()`, `GetByID()`, `Update()`, and `ListByHabitID()` methods

## API Usage Examples

### Creating a Habit with Duration
```json
POST /api/habits
{
  "name": "Morning Run",
  "description": "Daily morning exercise",
  "duration": 1800
}
```

### Creating a Habit without Duration
```json
POST /api/habits
{
  "name": "Read a Book",
  "description": "Daily reading habit"
}
```

### Creating a Log for a Habit with Duration (Required)
```json
POST /api/logs
{
  "habit_id": "habit-123",
  "notes": "Felt great today",
  "duration": 1800
}
```

### Creating a Log for a Habit without Duration (Optional)
```json
POST /api/logs
{
  "habit_id": "habit-456",
  "notes": "Read 20 pages"
}
```

## Validation Rules

1. **Habit Duration:**
   - Optional field
   - Must be non-negative if provided
   - Can be set during creation or update

2. **Log Duration:**
   - Required if the associated habit has a duration
   - Optional if the associated habit has no duration
   - Must be non-negative if provided

3. **Error Messages:**
   - "log must have a duration because the habit requires one" - when log is missing required duration
   - "duration must be non-negative" - when duration is negative

## Migration Instructions

To apply the database changes:
1. Run the migration: `migrations/004_add_duration_columns.sql`
2. Existing habits and logs will have `NULL` duration values (which maps to `nil` in Go)
3. No data migration needed - the feature is backward compatible

## Minimal Code Changes

This implementation follows the requirement of minimal modifications to existing code:
- **No changes** to handler layer (automatically handles new fields via JSON unmarshaling)
- **Only additive changes** to models (new optional fields)
- **Minimal modifications** to service and repository layers (added field handling)
- **New validation logic** in separate file to avoid modifying existing validation functions
- **Single migration file** to add database columns

All existing functionality remains intact, and the new duration feature is fully optional and backward compatible.
