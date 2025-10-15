# Command Center Package

This package provides a clean, well-organized interface for managing access codes through the command center system.

## Package Structure

### Files

- **`interfaces.go`** - Core interfaces for dependency injection and testing
  - `AccessCodeManager` - Main interface for access code operations
  - `ClientFactory` - Factory interface for creating clients

- **`types.go`** - All data types and structures
  - `Client` - Main client struct
  - `ClientConfig` - Configuration options
  - `AccessCodeOptions` - Type-safe options
  - Request/Response types for future API integration

- **`factory.go`** - Factory implementation for creating clients
  - `Factory` - Concrete factory implementation
  - Multiple creation methods for different use cases

- **`client.go`** - Core client implementation
  - Implements `AccessCodeManager` interface
  - Handles access code revocation and setting
  - Includes logging and error handling

- **`client_test.go`** - Comprehensive test suite
- **`doc.go`** - Package documentation
- **`README.md`** - This file

## Usage Examples

### Basic Usage
```go
factory := commandcenter.NewFactory()
client := factory.NewClient(siteID)
err := client.SetAccessCodes(units, options)
```

### With Context
```go
factory := commandcenter.NewFactory()
client := factory.(*commandcenter.Factory).NewClientWithContext(siteID, ctx)
err := client.RevokeAccessCodes(units, options)
```

### With Full Configuration
```go
config := commandcenter.ClientConfig{
    SiteID:  123,
    Context: ctx,
    BaseURL: "https://api.commandcenter.com",
    APIKey:  "your-api-key",
}
factory := commandcenter.NewFactory()
client := factory.(*commandcenter.Factory).NewClientWithConfig(config)
```

## Benefits of This Structure

1. **Interface Segregation** - Clean interfaces for easy testing and mocking
2. **Single Responsibility** - Each file has a specific purpose
3. **Dependency Injection** - Easy to swap implementations
4. **Testability** - Comprehensive test coverage
5. **Extensibility** - Easy to add new features and configurations
6. **Type Safety** - Strong typing for all operations

## Testing

Run tests with:
```bash
go test ./commandcenter
```
