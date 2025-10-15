# Code Refactoring: Separation of Concerns

This document describes the refactoring of `file-to-test.go` into multiple packages with proper separation of concerns.

## Original Structure
The original file contained all code in a single file with mixed responsibilities:
- Data models
- Constants
- Database operations
- Business logic
- HTTP handlers
- Utility functions
- Context management

## Refactored Structure

### 📁 `models/`
Contains all data structures and domain models:
- **`user.go`**: `BUser` and `Claims` structs
- **`unit.go`**: `Unit` struct  
- **`access_code.go`**: `GateAccessCode` and `GateAccessCodes` with validation logic

### 📁 `constants/`
Contains all application constants:
- **`access_codes.go`**: Access code states, lock states, and validation messages
- **`context.go`**: Context keys for storing values in request context

### 📁 `repository/`
Contains data access layer:
- **`bank.go`**: Database operations, command center client interface and implementation

### 📁 `services/`
Contains business logic layer:
- **`access_code.go`**: Access code validation, user validation, unit validation, and business rules

### 📁 `handlers/`
Contains HTTP request handling:
- **`access_code.go`**: HTTP handlers for access code operations, request parsing, and response handling

### 📁 `utils/`
Contains utility functions:
- **`slice.go`**: Helper functions for slice operations and data conversion

### 📁 `context/`
Contains context management:
- **`context.go`**: Context utilities for storing and retrieving values from request context

### 📁 `config/`
Contains configuration management:
- **`config.go`**: Application configuration structure and loading logic
- **`env.go`**: Environment file (.env) parsing and loading

### 📁 `middleware/`
Contains HTTP middleware for request processing:
- **`context.go`**: Bank context middleware for dependency injection
- **`logging.go`**: Request logging middleware with timing and status codes
- **`cors.go`**: CORS middleware with configurable origins, methods, and headers

### 📁 `examples/`
Contains example implementations:
- **`server_with_middleware.go`**: Complete server setup with full middleware chain

## Benefits of Refactoring

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Maintainability**: Code is easier to understand, modify, and extend
3. **Testability**: Individual components can be unit tested in isolation
4. **Reusability**: Components can be reused across different parts of the application
5. **Dependency Management**: Clear dependency hierarchy and interfaces
6. **Code Organization**: Related functionality is grouped together

## Dependency Flow

```
handlers/ → services/ → repository/
    ↓         ↓           ↓
  models/   models/    models/
    ↓         ↓           ↓
constants/ constants/ constants/
    ↓         ↓           ↓
  utils/    utils/     utils/
    ↓         ↓           ↓
 context/  context/   context/
```

## Usage

### Environment Configuration

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your actual values:**
   - Set `JWT_SECRET` and `ENCRYPTION_KEY` (required)
   - Configure database connection details
   - Set external service URLs and API keys

### Running the Server

#### Basic Setup with Middleware
The refactored code automatically includes bank context middleware:

```go
// Setup server with environment configuration and middleware
cfg, err := SetupServer()
if err != nil {
    log.Fatalf("Failed to setup server: %v", err)
}

// Start the server (middleware is already configured)
log.Printf("Server listening on http://%s", cfg.GetServerAddress())
http.ListenAndServe(":"+cfg.Port, nil)
```

#### Full Middleware Chain Example
For a complete middleware setup with logging and CORS:

```go
// Run the example with full middleware chain
go run examples/server_with_middleware.go
```

#### Available Middleware
- **BankMiddleware**: Injects Bank instance into request context
- **LoggingMiddleware**: Logs requests with method, path, status, and timing
- **CORSMiddleware**: Handles cross-origin requests with configurable options

Or simply run:
```bash
go run file-to-test.go          # Basic server with bank middleware
go run examples/server_with_middleware.go  # Full middleware chain
```

## Files Created

- `go.mod` - Go module definition
- `.env.example` - Example environment configuration file
- `models/user.go` - User and Claims models
- `models/unit.go` - Unit model
- `models/access_code.go` - Access code models and validation
- `constants/access_codes.go` - Access code related constants
- `constants/context.go` - Context key constants
- `config/config.go` - Application configuration management
- `config/env.go` - Environment file parsing and loading
- `middleware/context.go` - Bank context middleware for dependency injection
- `middleware/logging.go` - Request logging middleware
- `middleware/cors.go` - CORS middleware with configurable options
- `examples/server_with_middleware.go` - Complete middleware chain example
- `repository/bank.go` - Database operations and command center client
- `services/access_code.go` - Business logic for access code operations
- `handlers/access_code.go` - HTTP request handlers
- `utils/slice.go` - Utility functions
- `context/context.go` - Context management utilities
- `file-to-test.go` - Updated to use refactored packages with environment configuration

The original functionality is preserved while providing a much cleaner, more maintainable architecture.
