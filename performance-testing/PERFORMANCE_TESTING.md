# API Performance Testing

This directory contains a comprehensive performance testing suite for the Habit Tracker REST API. The tests are designed to measure throughput, latency, and system behavior under various load conditions.

## Overview

The performance test suite includes:

- **User Registration Testing**: Tests the `/api/register` endpoint
- **Authentication Testing**: Tests the `/api/login` endpoint
- **Habit CRUD Operations**: Tests habit creation, retrieval, update, and deletion
- **Log CRUD Operations**: Tests log entry operations
- **Concurrent Load Testing**: Sustained load testing with multiple concurrent workers
- **Statistical Analysis**: Detailed latency metrics and success rates

## Prerequisites

- Go 1.16 or higher
- The API server must be running before starting the tests
- PowerShell 5.1 or higher (for Windows script)

## Quick Start

### Method 1: Using PowerShell Script (Recommended for Windows)

The easiest way to run performance tests is using the PowerShell script:

```powershell
# Quick test (5 users, short duration)
.\run_performance_test.ps1 -Mode quick

# Standard test (default settings)
.\run_performance_test.ps1 -Mode standard

# Stress test (50 users, higher load)
.\run_performance_test.ps1 -Mode stress

# Extreme test (100 users, maximum load)
.\run_performance_test.ps1 -Mode extreme

# Custom test with specific parameters
.\run_performance_test.ps1 -Users 25 -Habits 8 -Logs 15 -Workers 10 -DurationSeconds 45
```

### Method 2: Using Go Directly

Build and run the test manually:

```powershell
# Build the test
go build -o perftest.exe perftest.go

# Run with default settings
.\perftest.exe

# Run with custom settings
.\perftest.exe -url=http://localhost:8080 -users=20 -habits=10 -logs=15 -workers=10 -duration=60s
```

## Test Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-url` | `http://localhost:8080` | Base URL of the API server |
| `-users` | `10` | Number of test users to create |
| `-habits` | `5` | Number of habits to create per user |
| `-logs` | `10` | Number of log entries per habit |
| `-workers` | `5` | Number of concurrent workers for operations |
| `-duration` | `30s` | Duration of the sustained load test phase |

## Test Modes (PowerShell Script)

### Quick Mode
- **Users**: 5
- **Habits per User**: 3
- **Logs per Habit**: 5
- **Workers**: 3
- **Duration**: 10 seconds
- **Use Case**: Quick validation, development testing

### Standard Mode (Default)
- **Users**: 10
- **Habits per User**: 5
- **Logs per Habit**: 10
- **Workers**: 5
- **Duration**: 30 seconds
- **Use Case**: Regular performance testing, CI/CD

### Stress Mode
- **Users**: 50
- **Habits per User**: 10
- **Logs per Habit**: 20
- **Workers**: 20
- **Duration**: 60 seconds
- **Use Case**: Finding bottlenecks, capacity planning

### Extreme Mode
- **Users**: 100
- **Habits per User**: 15
- **Logs per Habit**: 30
- **Workers**: 50
- **Duration**: 120 seconds
- **Use Case**: Maximum load testing, breaking point analysis

## Test Phases

The performance test runs through the following phases:

1. **Health Check**: Verifies the API is responding
2. **User Registration**: Creates test users concurrently
3. **User Login**: Authenticates all users and obtains JWT tokens
4. **Habit Creation**: Creates habits for each user
5. **Log Creation**: Creates log entries for each habit
6. **Read Operations**: Tests GET endpoints for habits and logs
7. **Load Testing**: Sustained concurrent requests for the specified duration

## Output and Results

### Console Output

The test provides real-time progress updates and a detailed summary:

```
========================================
API Performance Test
========================================
Base URL: http://localhost:8080
Users: 10
Habits per User: 5
Logs per Habit: 10
Concurrent Workers: 5
Load Test Duration: 30s
========================================

✓ Health check passed

Phase 1: User Registration
✓ Registered 10 users

Phase 2: User Login
✓ Logged in 10 users

...

========================================
Performance Test Results
========================================
Total Requests: 1234
Successful Requests: 1230
Failed Requests: 4
Success Rate: 99.68%

Average Latency: 45ms
Min Latency: 12ms
Max Latency: 234ms

User Registration (10 requests):
  Avg: 120ms, Min: 98ms, Max: 156ms

User Login (10 requests):
  Avg: 45ms, Min: 32ms, Max: 78ms
...
```

### JSON Results File

A detailed JSON file is generated with the format: `performance_test_results_YYYYMMDD_HHMMSS.json`

Example content:
```json
{
  "timestamp": "2025-10-22T14:30:00Z",
  "total_requests": 1234,
  "successful_requests": 1230,
  "failed_requests": 4,
  "success_rate": 99.68,
  "average_latency_ms": 45,
  "min_latency_ms": 12,
  "max_latency_ms": 234
}
```

## Performance Benchmarks

Expected performance on a typical development machine:

| Operation | Expected Latency | Target Success Rate |
|-----------|------------------|---------------------|
| User Registration | < 200ms | > 99% |
| User Login | < 100ms | > 99% |
| Create Habit | < 50ms | > 99.5% |
| Create Log | < 50ms | > 99.5% |
| Get Habits | < 30ms | > 99.9% |
| Get Logs | < 30ms | > 99.9% |

## Interpreting Results

### Success Rate
- **> 99%**: Excellent
- **95-99%**: Good, investigate occasional failures
- **< 95%**: Poor, requires investigation

### Latency
- **Average Latency**: Should be < 100ms for most operations
- **Max Latency**: Spikes > 1s indicate potential issues
- **Min Latency**: Useful for best-case performance baseline

### High Failure Rate Causes
1. Database lock contention (SQLite limitation)
2. Connection pool exhaustion
3. Memory pressure
4. CPU throttling
5. Network issues

## Performance Tuning Tips

### For SQLite
1. Enable WAL mode: `PRAGMA journal_mode=WAL`
2. Increase cache size: `PRAGMA cache_size=10000`
3. Use connection pooling appropriately
4. Consider moving to PostgreSQL/MySQL for production

### For the API
1. Adjust worker pool size based on CPU cores
2. Implement request rate limiting
3. Add caching for frequently accessed data
4. Use database connection pooling
5. Profile using `pprof` to find bottlenecks

### For Testing
1. Start with small loads and gradually increase
2. Monitor system resources (CPU, memory, disk I/O)
3. Run multiple test iterations for consistency
4. Compare results across code changes

## Troubleshooting

### API Server Not Responding
```
✗ API server is not responding at http://localhost:8080
```
**Solution**: Start the API server first: `go run main.go`

### High Failure Rate
```
Failed Requests: 450 (45%)
```
**Possible Causes**:
- Database locked (SQLite concurrency limits)
- API server overwhelmed
- Network issues

**Solutions**:
- Reduce concurrent workers
- Reduce number of users/habits/logs
- Check server logs for specific errors

### Build Errors
```
✗ Failed to build performance test
```
**Solution**: Ensure you're in the correct directory and Go is installed:
```powershell
go version
go mod download
```

## Advanced Usage

### Custom Test Scenarios

Create a custom test by editing the parameters:

```go
// Test extreme habit creation
.\performance_test.exe -users=5 -habits=100 -logs=1 -workers=10 -duration=0s

// Test heavy log creation
.\performance_test.exe -users=5 -habits=5 -logs=1000 -workers=20 -duration=0s

// Test read-heavy workload (modify code to increase read operations)
```

### Integration with CI/CD

Add to your CI pipeline:

```yaml
# GitHub Actions example
- name: Run Performance Tests
  run: |
    go run main.go &
    SERVER_PID=$!
    sleep 5
    go run perftest.go -users=10 -habits=5 -logs=10 -duration=30s
    kill $SERVER_PID
```

### Comparing Results

Run tests before and after changes to measure impact:

```powershell
# Baseline test
.\run_performance_test.ps1 -Mode standard | Out-File baseline_results.txt

# Make code changes...

# Comparison test
.\run_performance_test.ps1 -Mode standard | Out-File new_results.txt

# Compare files
Compare-Object (Get-Content baseline_results.txt) (Get-Content new_results.txt)
```

## Known Limitations

1. **SQLite Concurrency**: SQLite has limited concurrent write support. High worker counts may cause database lock errors.
2. **Local Testing**: Tests run on localhost, which doesn't account for network latency.
3. **Single Machine**: Cannot test distributed load scenarios.
4. **State Cleanup**: Each test creates new users/data. You may want to clean the database between runs.

## Cleanup

To reset the test database:

```powershell
# Stop the server if running
# Delete the database file
Remove-Item habits.db

# Restart the server (migrations will recreate tables)
go run main.go
```

## Contributing

When adding new endpoints:

1. Add corresponding test functions in `perftest.go`
2. Update test phases to include the new operations
3. Add latency tracking for the new operations
4. Update this README with expected performance metrics

## License

Same as the main project.
