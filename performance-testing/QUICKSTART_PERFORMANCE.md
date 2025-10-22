# Quick Start: Performance Testing

## 🚀 Fastest Way to Test

### Option 1: PowerShell Script (Recommended)
```powershell
# Make sure API is running first
go run main.go

# In another terminal, run:
.\run_performance_test.ps1 -Mode quick
```

### Option 2: Batch File (Simplest)
```cmd
REM Double-click or run:
run_performance_test.bat
```

### Option 3: Direct Go Command
```powershell
go run perftest.go -users=5 -habits=3 -logs=5 -workers=3 -duration=10s
```

## 📊 Test Modes Available

| Mode | Users | Habits | Logs | Duration | Best For |
|------|-------|--------|------|----------|----------|
| **quick** | 5 | 3 | 5 | 10s | Development testing |
| **standard** | 10 | 5 | 10 | 30s | Regular validation |
| **stress** | 50 | 10 | 20 | 60s | Finding bottlenecks |
| **extreme** | 100 | 15 | 30 | 120s | Maximum load testing |

## 📝 Example Commands

```powershell
# Quick test
.\run_performance_test.ps1 -Mode quick

# Custom test
.\run_performance_test.ps1 -Users 20 -Habits 10 -DurationSeconds 60

# Using Go directly
go run perftest.go -users=10 -habits=5 -logs=10 -duration=30s
```

## 📈 What Gets Tested

1. ✅ User Registration (`/api/register`)
2. ✅ User Login (`/api/login`)
3. ✅ Create Habits (`POST /api/habits`)
4. ✅ Create Logs (`POST /api/logs`)
5. ✅ Get Habits (`GET /api/habits`)
6. ✅ Get Logs (`GET /api/logs?habit_id=X`)
7. ✅ Sustained Load Testing (continuous requests)

## 📊 Understanding Results

### Success Rates
- **> 99%** = Excellent ✅
- **95-99%** = Good ⚠️
- **< 95%** = Needs investigation ❌

### Latency Expectations
- **Average < 100ms** = Good performance ✅
- **Average 100-500ms** = Acceptable ⚠️
- **Average > 500ms** = Performance issue ❌

## 🔧 Troubleshooting

### "API server is not responding"
```powershell
# Start the API server first
go run main.go

# Wait for: "Starting server on :8080"
```

### High failure rate (> 5%)
- Try reducing `-workers` parameter
- SQLite has concurrency limits (~10-20 writers)
- Check server logs for specific errors

### Build errors
```powershell
# Make sure you're in the right directory
cd c:\path\to\new-api

# Update dependencies
go mod download
```

## 📁 Files Generated

- `performance_test_results_YYYYMMDD_HHMMSS.json` - Detailed metrics
- Check console output for summary statistics

## 📚 More Information

- Full documentation: `PERFORMANCE_TESTING.md`
- All test commands: `PERFORMANCE_TEST_COMMANDS.md`
- Source code: `perftest.go`

## 💡 Tips

1. **Always start with quick mode** to verify everything works
2. **Monitor system resources** during stress tests
3. **Clean database** between major test runs
4. **Compare results** before/after code changes
5. **Save results** for performance regression tracking

---

**Ready to test?** Run: `.\run_performance_test.ps1 -Mode quick`
