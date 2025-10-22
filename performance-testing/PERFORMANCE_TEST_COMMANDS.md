# Performance Test Examples
# Copy and paste these commands to run different test scenarios

# ============================================
# Quick Tests (for development)
# ============================================

# Minimal test - fastest
go run perftest.go -users=3 -habits=2 -logs=3 -workers=2 -duration=5s

# Quick validation
go run perftest.go -users=5 -habits=3 -logs=5 -workers=3 -duration=10s

# ============================================
# Standard Tests (for regular testing)
# ============================================

# Default test
go run perftest.go

# Standard with more users
go run perftest.go -users=20 -habits=5 -logs=10 -workers=5 -duration=30s

# ============================================
# Stress Tests (finding limits)
# ============================================

# Moderate stress
go run perftest.go -users=50 -habits=10 -logs=20 -workers=10 -duration=60s

# High stress
go run perftest.go -users=100 -habits=10 -logs=20 -workers=20 -duration=60s

# Extreme stress
go run perftest.go -users=200 -habits=15 -logs=30 -workers=30 -duration=120s

# ============================================
# Specific Endpoint Tests
# ============================================

# Test heavy habit creation
go run perftest.go -users=10 -habits=50 -logs=1 -workers=10 -duration=10s

# Test heavy log creation
go run perftest.go -users=10 -habits=5 -logs=100 -workers=15 -duration=10s

# Test many users, few habits
go run perftest.go -users=100 -habits=2 -logs=5 -workers=20 -duration=30s

# ============================================
# Load Test Focus (minimal setup, max load testing)
# ============================================

# Short setup, long load test
go run perftest.go -users=5 -habits=3 -logs=5 -workers=20 -duration=120s

# Medium setup, sustained load
go run perftest.go -users=20 -habits=5 -logs=10 -workers=15 -duration=300s

# ============================================
# Concurrency Tests
# ============================================

# Low concurrency
go run perftest.go -users=20 -habits=5 -logs=10 -workers=2 -duration=30s

# Medium concurrency
go run perftest.go -users=20 -habits=5 -logs=10 -workers=10 -duration=30s

# High concurrency
go run perftest.go -users=20 -habits=5 -logs=10 -workers=30 -duration=30s

# Very high concurrency (test database locks)
go run perftest.go -users=20 -habits=5 -logs=10 -workers=50 -duration=30s

# ============================================
# PowerShell Script Examples
# ============================================

# Quick mode
.\run_performance_test.ps1 -Mode quick

# Standard mode
.\run_performance_test.ps1 -Mode standard

# Stress mode
.\run_performance_test.ps1 -Mode stress

# Extreme mode
.\run_performance_test.ps1 -Mode extreme

# Custom via PowerShell
.\run_performance_test.ps1 -Users 25 -Habits 8 -Logs 15 -Workers 10 -DurationSeconds 45

# ============================================
# Production Simulation
# ============================================

# Simulate 1 week of light usage (1 user, 7 habits, daily logs)
go run perftest.go -users=1 -habits=7 -logs=7 -workers=1 -duration=5s

# Simulate 100 users over a month
go run perftest.go -users=100 -habits=10 -logs=30 -workers=10 -duration=60s

# ============================================
# CI/CD Pipeline Examples
# ============================================

# Fast CI test (< 30 seconds total)
go run perftest.go -users=5 -habits=3 -logs=5 -workers=3 -duration=5s

# Thorough CI test (< 2 minutes total)
go run perftest.go -users=15 -habits=5 -logs=10 -workers=5 -duration=30s

# Nightly build test (comprehensive)
go run perftest.go -users=50 -habits=10 -logs=20 -workers=10 -duration=120s

# ============================================
# Benchmarking Specific Operations
# ============================================

# Registration benchmark (no load test)
go run perftest.go -users=100 -habits=0 -logs=0 -workers=10 -duration=0s

# Login benchmark (requires users first)
# Run in two steps: create users, then test login by modifying code

# Read-heavy workload (GET operations during load test)
go run perftest.go -users=20 -habits=10 -logs=10 -workers=15 -duration=60s

# ============================================
# Database Stress Tests
# ============================================

# Test SQLite write locks
go run perftest.go -users=50 -habits=5 -logs=5 -workers=30 -duration=30s

# Test with many small transactions
go run perftest.go -users=100 -habits=2 -logs=5 -workers=20 -duration=30s

# Test with fewer large transactions
go run perftest.go -users=10 -habits=20 -logs=50 -workers=5 -duration=30s

# ============================================
# Before/After Comparison
# ============================================

# Baseline (run before changes)
go run perftest.go -users=20 -habits=5 -logs=10 -workers=5 -duration=30s > baseline.txt

# After changes (run after optimization)
go run perftest.go -users=20 -habits=5 -logs=10 -workers=5 -duration=30s > after.txt

# Compare results
# On Windows: fc baseline.txt after.txt
# On Unix: diff baseline.txt after.txt

# ============================================
# Notes
# ============================================

# - Start with smaller tests and gradually increase
# - Monitor system resources (Task Manager on Windows)
# - Check habits.db file size after tests
# - Clean database between major test runs
# - Results vary based on system performance
# - SQLite has concurrency limitations (~10-20 concurrent writers)
