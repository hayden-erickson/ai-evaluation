# Performance Testing Suite

This folder contains all the tools and documentation for performance testing the Habit Tracker REST API.

## 📁 Files in This Folder

| File | Purpose |
|------|---------|
| `perftest.go` | Main performance testing application |
| `run_performance_test.ps1` | PowerShell runner with predefined test modes |
| `run_performance_test.bat` | Simple batch file for quick testing |
| `QUICKSTART_PERFORMANCE.md` | Quick start guide - **START HERE** |
| `PERFORMANCE_TESTING.md` | Complete documentation and reference |
| `PERFORMANCE_TEST_COMMANDS.md` | 40+ ready-to-use test commands |

## 🚀 Quick Start

### 1. Start the API Server
```powershell
# From the project root directory
go run main.go
```

### 2. Run Performance Tests
```powershell
# From this directory (performance-testing)
.\run_performance_test.ps1 -Mode quick
```

## 📚 Documentation

- **New to performance testing?** → Start with `QUICKSTART_PERFORMANCE.md`
- **Need detailed info?** → Check `PERFORMANCE_TESTING.md`
- **Looking for examples?** → Browse `PERFORMANCE_TEST_COMMANDS.md`

## 🎯 Test Modes

```powershell
# Quick test (5 users, 10 seconds)
.\run_performance_test.ps1 -Mode quick

# Standard test (10 users, 30 seconds)
.\run_performance_test.ps1 -Mode standard

# Stress test (50 users, 60 seconds)
.\run_performance_test.ps1 -Mode stress

# Extreme test (100 users, 120 seconds)
.\run_performance_test.ps1 -Mode extreme
```

## 💡 Direct Usage

You can also run the test directly with custom parameters:

```powershell
go run perftest.go -users=20 -habits=10 -logs=15 -workers=10 -duration=60s
```

## 📊 What Gets Tested

- User Registration
- User Authentication
- Habit Creation & Retrieval
- Log Creation & Retrieval
- Concurrent Load Handling

## 🔧 Requirements

- Go 1.16 or higher
- API server running on http://localhost:8080 (or custom URL)
- PowerShell 5.1+ (for .ps1 script)

## 📈 Output

Results are displayed in the console and saved to:
- `performance_test_results_YYYYMMDD_HHMMSS.json`

---

**Ready to start?** Open `QUICKSTART_PERFORMANCE.md` for step-by-step instructions!
