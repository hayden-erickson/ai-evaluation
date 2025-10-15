# AI Evaluation - Refactored Structure

This project has been refactored to follow proper separation of concerns with the following package structure:

## Package Structure

### `/models`
Contains all data structures and domain models:
- `BUser` - Business user model
- `Unit` - Rental unit model  
- `GateAccessCode` - Access code model
- `Claims` - Authentication claims model

### `/constants`
Contains all application constants:
- Access code states
- Lock states
- Validation messages

### `/database`
Contains database operations and the Bank struct:
- `Bank` - Main database client
- Database query methods
- Command center client factory

### `/commandcenter`
Contains command center client operations:
- `Client` - Command center client
- Access code management operations

### `/contextutil`
Contains context utilities:
- Context key management
- Context value extraction helpers

### `/handlers`
Contains HTTP request handlers:
- `AccessCodeEditHandler` - Main access code edit endpoint
- Request validation and processing

### `/utils`
Contains utility functions:
- `UniqueIntSlice` - Remove duplicates from int slice
- `ConvertToStringSlice` - Convert access codes to strings

## Usage

The main application is in `main.go` which sets up the HTTP server and routes. The refactored structure provides:

1. **Better maintainability** - Each package has a single responsibility
2. **Improved testability** - Packages can be tested in isolation
3. **Clear dependencies** - Import structure shows relationships between components
4. **Easier navigation** - Related code is grouped together

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 with the access code edit endpoint available at `/api/access-code/edit`.
