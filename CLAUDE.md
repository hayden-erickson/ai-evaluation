# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go-based API for managing access codes for rental units. The application provides HTTP endpoints for editing access codes, integrates with a command center system for physical access control, and uses MySQL for data persistence.

## Build and Development Commands

### Setup
```bash
go mod download  # Download dependencies
```

### Running the Application
```bash
go run main.go   # Start the HTTP server on port 8080 (or PORT env var)
```

### Testing
```bash
go test ./...                    # Run all tests
go test ./commandcenter          # Run tests for a specific package
go test -v ./commandcenter       # Run with verbose output
```

Note: Some tests may fail due to missing `go.sum` entry. Run `go mod tidy` to resolve dependency issues.

### Database Setup
```bash
mysql -u username -p database_name < database/schema.sql
```

## Environment Variables

The application requires the following environment variables:

- `DB_HOST` - Database host (default: localhost)
- `DB_USERNAME` - Database username (default: root)
- `DB_PASSWORD` - Database password (required)
- `DB_DATABASE` - Database name (default: ai_evaluation)
- `PORT` - Server port (default: 8080)

## Architecture

### Core Design Pattern

The application follows a **dependency injection** pattern where the `Bank` (database client) is injected into HTTP request contexts via middleware. This allows handlers to access database operations without tight coupling.

### Package Structure

The codebase is organized with strict separation of concerns:

- **`/models`** - Domain models and validation logic
  - `BUser` - Business users with company and site associations
  - `Unit` - Rental units with lock states
  - `GateAccessCode` - Access codes with lifecycle states
  - `Claims` - Authentication claims from JWT/session
  - `GateAccessCodes` collection type with `Validate()` method

- **`/database`** - Data persistence layer
  - `Bank` struct - Main database client with MySQL connection pooling
  - Methods for CRUD operations on users, units, and access codes
  - `NewCommandCenterClient()` - Factory method for command center clients
  - Transaction support via `BeginTx()`

- **`/commandcenter`** - External system integration
  - Interface-based design (`AccessCodeManager`, `ClientFactory`)
  - `Client` - Handles access code operations with physical lock system
  - `Factory` - Creates clients with different configurations
  - Files organized by purpose: `interfaces.go`, `types.go`, `factory.go`, `client.go`

- **`/handlers`** - HTTP request handlers
  - `AccessCodeEditHandler` - Main endpoint at `/api/access-code/edit`
  - Validates user permissions, unit states, and access code changes
  - Coordinates database updates with command center operations

- **`/contextutil`** - Context management utilities
  - Helpers for injecting/extracting values from request contexts
  - Used for passing `Bank` and `Claims` through middleware chain

- **`/constants`** - Application-wide constants
  - Access code lifecycle states (active, setup, pending, remove, etc.)
  - Lock states (overlock, gatelock, prelet)
  - Validation message constants

- **`/utils`** - Shared utility functions
  - `UniqueIntSlice()` - Deduplicate integer slices
  - `ConvertToStringSlice()` - Type conversion helpers

### Key Architectural Patterns

1. **Context-based Dependency Injection**: The `Bank` is injected into request contexts via middleware in `main.go:SetupServer()`, allowing handlers to access it via `contextutil.BankFromContext()`.

2. **Interface Segregation**: The `commandcenter` package uses interfaces (`AccessCodeManager`, `ClientFactory`) to enable testing with mocks and future implementation swaps.

3. **Factory Pattern**: Command center clients are created via factory methods, supporting different configuration options (basic, with context, with full config).

4. **Transaction Safety**: The `Bank.UpdateAccessCodes()` method uses database transactions to ensure atomic updates.

5. **State Machine**: Access codes follow a lifecycle with states: setup → pending → active → remove → removing → removed. Units in overlock/gatelock/prelet states cannot have codes changed.

## Data Flow for Access Code Edits

1. Handler extracts claims (user identity, current site) from context
2. Validates user exists and belongs to correct company
3. Validates user has association with target site
4. For each unit:
   - Validates unit exists and belongs to site
   - Validates unit is not in overlock/gatelock/prelet state
   - Creates new `GateAccessCode` with state "setup"
   - Validates access code (checks for duplicates)
   - Retrieves existing codes for the unit
   - Marks old codes for removal (state → "remove")
   - Updates database: old codes to "remove", new code to "setup"
   - Calls command center to revoke old codes
   - Calls command center to set new codes
5. Records activity event for audit trail

## Database Schema

Key tables (see `database/schema.sql` for full schema):
- `business_users` - User accounts with company associations
- `sites` - Physical locations
- `user_sites` - Many-to-many relationship between users and sites
- `units` - Individual rental units with site association and rental state
- `gate_access_codes` - Access codes with unit/user/site associations and lifecycle state

## Important Constraints and Validation

1. **Company Isolation**: Users can only modify access codes for users in their company
2. **Site Association**: Both users and units must be associated with the current site
3. **Unit State Restrictions**: Cannot change codes for units in overlock, gatelock, or prelet states
4. **Duplicate Prevention**: New access codes are validated against existing codes
5. **Idempotency**: If the new code matches an existing active/setup/pending code, no changes are made

## Testing Strategy

The `commandcenter` package has comprehensive tests using the standard testing library. Tests verify:
- Factory creation methods
- Access code revocation and setting operations
- Edge cases (empty unit lists)
- Context propagation

Other packages currently lack test files but follow testable patterns (interfaces, dependency injection).
