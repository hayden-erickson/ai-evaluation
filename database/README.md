# Database Package

This package provides MySQL database operations for the AI Evaluation application.

## Features

- **Connection Management**: Robust MySQL connection with connection pooling
- **Transaction Support**: Context-aware transaction handling
- **Health Monitoring**: Connection health checks and statistics
- **Error Handling**: Comprehensive error handling with proper error wrapping
- **Security**: Parameterized queries to prevent SQL injection

## Configuration

### Environment Variables

```bash
DB_HOST=localhost          # Database host (default: localhost)
DB_PORT=3306              # Database port (default: 3306)
DB_USERNAME=root          # Database username (default: root)
DB_PASSWORD=your_password # Database password (required)
DB_DATABASE=ai_evaluation # Database name (default: ai_evaluation)
```

### Using Config Struct

```go
config := database.Config{
    Host:     "localhost",
    Port:     3306,
    Username: "myuser",
    Password: "mypassword",
    Database: "mydatabase",
}

bank, err := database.NewBank(config)
if err != nil {
    log.Fatal(err)
}
defer bank.Close()
```

### Using DSN String

```go
dsn := "user:password@tcp(localhost:3306)/dbname?parseTime=true"
bank, err := database.NewBankFromDSN(dsn)
if err != nil {
    log.Fatal(err)
}
defer bank.Close()
```

## Database Schema

Run the provided `schema.sql` file to create the required tables:

```bash
mysql -u username -p database_name < database/schema.sql
```

### Tables

- **business_users**: User accounts and company associations
- **sites**: Physical locations/sites
- **user_sites**: Many-to-many relationship between users and sites
- **units**: Individual rental units within sites
- **gate_access_codes**: Access codes for units

## Usage Examples

### Basic Operations

```go
// Get a user by ID
user, err := bank.GetBUserByID(123)
if err != nil {
    log.Printf("Error: %v", err)
}

// Get a unit
unit, err := bank.V2UnitGetById(456, 1)
if err != nil {
    log.Printf("Error: %v", err)
}

// Get access codes for units
codes, err := bank.GetCodesForUnits([]int{456, 789}, 1)
if err != nil {
    log.Printf("Error: %v", err)
}

// Update access codes
gacs := models.GateAccessCodes{
    {
        AccessCode: "1234",
        UnitID:     456,
        UserID:     123,
        SiteID:     1,
        State:      constants.AccessCodeStateSetup,
    },
}
err = bank.UpdateAccessCodes(gacs, 1)
if err != nil {
    log.Printf("Error: %v", err)
}
```

### Health Monitoring

```go
// Check connection health
if err := bank.Ping(); err != nil {
    log.Printf("Database connection unhealthy: %v", err)
}

// Get connection statistics
stats := bank.Stats()
log.Printf("Open connections: %d", stats.OpenConnections)
log.Printf("In use: %d", stats.InUse)
log.Printf("Idle: %d", stats.Idle)
```

### Transaction Handling

```go
ctx := context.Background()
tx, err := bank.BeginTx(ctx)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback()

// Perform operations with tx...

if err := tx.Commit(); err != nil {
    log.Fatal(err)
}
```

## Connection Pool Configuration

The database connection is configured with:
- **Max Open Connections**: 25
- **Max Idle Connections**: 5
- **Connection Max Lifetime**: 5 minutes
- **Character Set**: utf8mb4 with unicode collation

## Error Handling

The package uses wrapped errors for better debugging:

```go
user, err := bank.GetBUserByID(123)
if err != nil {
    if err.Error() == "no_ob_found" {
        // Handle user not found
    } else {
        // Handle other database errors
        log.Printf("Database error: %v", err)
    }
}
```

## Security Features

- **Parameterized Queries**: All queries use parameter placeholders to prevent SQL injection
- **Connection Encryption**: Supports SSL/TLS connections (configure in DSN)
- **Input Validation**: Proper validation of input parameters
- **Error Sanitization**: Database errors are wrapped to avoid exposing sensitive information

## Dependencies

- `github.com/go-sql-driver/mysql`: MySQL driver for Go
