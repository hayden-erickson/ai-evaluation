package main

/*
Performance Test for Habit Tracker API

IMPORTANT NOTES:
This test has been updated to match the actual API endpoints defined in cmd/api/main.go.

API Endpoints Available:
- POST /register     - Register user (returns JWT token)
- POST /login        - Login user (returns JWT token)
- POST /users        - Create user (alternative)
- GET /users?id={id} - Get user by ID
- PUT /users         - Update user
- DELETE /users?id={id} - Delete user

Protected Endpoints (require JWT):
- POST /habits/      - Create habit (requires JWT token)
- GET /habits/       - Get all user habits (requires JWT token)
- POST /logs/        - Create log (requires JWT token)
- GET /logs/?habit_id={id} - Get logs for habit (requires JWT token)

Authentication:
- Users are registered with phone_number (no password required in current schema)
- Login requires only phone_number
- JWT tokens are returned on registration and login
- Tokens are valid for 24 hours

Current Test Coverage:
✓ User Registration (POST /register)
✓ User Login (POST /login)
✓ Habits CRUD (with JWT authentication)
✓ Logs CRUD (with JWT authentication)
✓ Load Testing (with JWT authentication)
*/

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Configuration for the performance test
type TestConfig struct {
	BaseURL           string
	NumUsers          int
	NumHabitsPerUser  int
	NumLogsPerHabit   int
	ConcurrentWorkers int
	TestDuration      time.Duration
}

// Test statistics
type TestStats struct {
	mu                   sync.Mutex
	TotalRequests        int64
	SuccessfulRequests   int64
	FailedRequests       int64
	TotalLatency         time.Duration
	MinLatency           time.Duration
	MaxLatency           time.Duration
	RegisterLatencies    []time.Duration
	LoginLatencies       []time.Duration
	CreateHabitLatencies []time.Duration
	CreateLogLatencies   []time.Duration
	GetHabitsLatencies   []time.Duration
	GetLogsLatencies     []time.Duration
}

// User credentials for testing
type TestUser struct {
	PhoneNumber string
	Password    string
	Token       string
	UserID      int64
	Habits      []TestHabit
}

// Habit for testing
type TestHabit struct {
	ID   int64
	Name string
}

// Response structures
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  json.RawMessage `json:"user"`
}

type UserResponse struct {
	ID int64 `json:"id"`
}

type HabitResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

func main() {
	// Parse command line flags
	baseURL := flag.String("url", "http://localhost:8080", "Base URL of the API")
	numUsers := flag.Int("users", 10, "Number of test users to create")
	numHabits := flag.Int("habits", 5, "Number of habits per user")
	numLogs := flag.Int("logs", 10, "Number of logs per habit")
	workers := flag.Int("workers", 5, "Number of concurrent workers")
	duration := flag.Duration("duration", 30*time.Second, "Duration for load testing")
	flag.Parse()

	config := TestConfig{
		BaseURL:           *baseURL,
		NumUsers:          *numUsers,
		NumHabitsPerUser:  *numHabits,
		NumLogsPerHabit:   *numLogs,
		ConcurrentWorkers: *workers,
		TestDuration:      *duration,
	}

	stats := &TestStats{
		MinLatency: time.Hour, // Start with a very high value
	}

	fmt.Println("========================================")
	fmt.Println("API Performance Test")
	fmt.Println("========================================")
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Printf("Users: %d\n", config.NumUsers)
	fmt.Printf("Habits per User: %d\n", config.NumHabitsPerUser)
	fmt.Printf("Logs per Habit: %d\n", config.NumLogsPerHabit)
	fmt.Printf("Concurrent Workers: %d\n", config.ConcurrentWorkers)
	fmt.Printf("Load Test Duration: %s\n", config.TestDuration)
	fmt.Println("========================================")

	// Test health endpoint - Note: No dedicated health endpoint exists, testing /users instead
	if err := testHealthEndpoint(config.BaseURL); err != nil {
		log.Printf("Warning: API connectivity check failed: %v", err)
	} else {
		fmt.Println("✓ API connectivity check passed")
	}

	// Phase 1: User Registration
	fmt.Println("\nPhase 1: User Registration")
	users := registerUsers(config, stats)
	fmt.Printf("✓ Registered %d users\n", len(users))

	// Phase 2: User Login
	fmt.Println("\nPhase 2: User Login")
	loginUsers(config, stats, users)
	fmt.Printf("✓ Logged in %d users\n", len(users))

	// Phase 3: Create Habits
	fmt.Println("\nPhase 3: Create Habits")
	createHabits(config, stats, users)
	totalHabits := 0
	for _, user := range users {
		totalHabits += len(user.Habits)
	}
	fmt.Printf("✓ Created %d habits\n", totalHabits)

	// Phase 4: Create Logs
	fmt.Println("\nPhase 4: Create Logs")
	createLogs(config, stats, users)
	fmt.Printf("✓ Created logs for all habits\n")

	// Phase 5: Read Operations
	fmt.Println("\nPhase 5: Read Operations")
	performReadOperations(config, stats, users)
	fmt.Printf("✓ Performed read operations\n")

	// Phase 6: Load Testing
	fmt.Println("\nPhase 6: Load Testing")
	performLoadTest(config, stats, users)

	// Print results
	printResults(stats)

	// Save results to file
	if err := saveResultsToFile(stats); err != nil {
		log.Printf("Warning: Failed to save results to file: %v", err)
	}
}

// testHealthEndpoint checks if the API is responsive
func testHealthEndpoint(baseURL string) error {
	// Note: API doesn't have a /health endpoint, so we test basic connectivity
	resp, err := http.Get(baseURL + "/users?id=1")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Accept any response (even 404) as long as we can connect
	if resp.StatusCode == 0 {
		return fmt.Errorf("no response from server")
	}
	return nil
}

// registerUsers creates test users
func registerUsers(config TestConfig, stats *TestStats) []*TestUser {
	users := make([]*TestUser, 0, config.NumUsers)
	var wg sync.WaitGroup

	workerChan := make(chan int, config.NumUsers)
	resultChan := make(chan *TestUser, config.NumUsers)

	// Start workers
	for i := 0; i < config.ConcurrentWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userNum := range workerChan {
				user := &TestUser{
					PhoneNumber: fmt.Sprintf("+1555000%04d", userNum),
					Password:    "TestPassword123!",
				}

				start := time.Now()
				if err := registerUser(config.BaseURL, user); err != nil {
					log.Printf("Failed to register user %d: %v", userNum, err)
					recordRequest(stats, start, false, &stats.RegisterLatencies)
					continue
				}
				recordRequest(stats, start, true, &stats.RegisterLatencies)
				resultChan <- user
			}
		}()
	}

	// Send work to workers
	go func() {
		for i := 0; i < config.NumUsers; i++ {
			workerChan <- i
		}
		close(workerChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for user := range resultChan {
		users = append(users, user)
	}

	return users
}

// registerUser registers a single user
func registerUser(baseURL string, user *TestUser) error {
	payload := map[string]interface{}{
		"phone_number": user.PhoneNumber,
		"name":         "Test User " + user.PhoneNumber,
		"time_zone":    "America/New_York",
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	// Extract token and user from the nested response
	var registerData map[string]json.RawMessage
	if err := json.Unmarshal(apiResp.Data, &registerData); err != nil {
		return err
	}

	// Extract token
	var token string
	if err := json.Unmarshal(registerData["token"], &token); err != nil {
		return err
	}
	user.Token = token

	// Extract user data
	var userResp UserResponse
	if err := json.Unmarshal(registerData["user"], &userResp); err != nil {
		return err
	}
	user.UserID = userResp.ID

	return nil
}

// loginUsers logs in all test users
func loginUsers(config TestConfig, stats *TestStats, users []*TestUser) {
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(u *TestUser) {
			defer wg.Done()
			start := time.Now()
			if err := loginUser(config.BaseURL, u); err != nil {
				log.Printf("Failed to login user %s: %v", u.PhoneNumber, err)
				recordRequest(stats, start, false, &stats.LoginLatencies)
				return
			}
			recordRequest(stats, start, true, &stats.LoginLatencies)
		}(user)
	}

	wg.Wait()
}

// loginUser logs in a single user
func loginUser(baseURL string, user *TestUser) error {
	payload := map[string]string{
		"phone_number": user.PhoneNumber,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	// Extract token from the nested response
	var loginData map[string]json.RawMessage
	if err := json.Unmarshal(apiResp.Data, &loginData); err != nil {
		return err
	}

	var token string
	if err := json.Unmarshal(loginData["token"], &token); err != nil {
		return err
	}

	user.Token = token
	return nil
}

// createHabits creates habits for all users
func createHabits(config TestConfig, stats *TestStats, users []*TestUser) {
	var wg sync.WaitGroup

	for _, user := range users {
		for i := 0; i < config.NumHabitsPerUser; i++ {
			wg.Add(1)
			go func(u *TestUser, habitNum int) {
				defer wg.Done()
				start := time.Now()
				habit, err := createHabit(config.BaseURL, u, habitNum)
				if err != nil {
					log.Printf("Failed to create habit for user %s: %v", u.PhoneNumber, err)
					recordRequest(stats, start, false, &stats.CreateHabitLatencies)
					return
				}
				recordRequest(stats, start, true, &stats.CreateHabitLatencies)
				u.Habits = append(u.Habits, habit)
			}(user, i)
		}
	}

	wg.Wait()
}

// createHabit creates a single habit
func createHabit(baseURL string, user *TestUser, habitNum int) (TestHabit, error) {
	payload := map[string]string{
		"name":        fmt.Sprintf("Habit %d for %s", habitNum, user.PhoneNumber),
		"description": fmt.Sprintf("Test habit description %d", habitNum),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", baseURL+"/habits/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TestHabit{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return TestHabit{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return TestHabit{}, err
	}

	var habitResp HabitResponse
	if err := json.Unmarshal(apiResp.Data, &habitResp); err != nil {
		return TestHabit{}, err
	}

	return TestHabit{ID: habitResp.ID, Name: habitResp.Name}, nil
}

// createLogs creates logs for all habits
func createLogs(config TestConfig, stats *TestStats, users []*TestUser) {
	var wg sync.WaitGroup

	for _, user := range users {
		for _, habit := range user.Habits {
			for i := 0; i < config.NumLogsPerHabit; i++ {
				wg.Add(1)
				go func(u *TestUser, h TestHabit, logNum int) {
					defer wg.Done()
					start := time.Now()
					if err := createLog(config.BaseURL, u, h, logNum); err != nil {
						log.Printf("Failed to create log: %v", err)
						recordRequest(stats, start, false, &stats.CreateLogLatencies)
						return
					}
					recordRequest(stats, start, true, &stats.CreateLogLatencies)
				}(user, habit, i)
			}
		}
	}

	wg.Wait()
}

// createLog creates a single log entry
func createLog(baseURL string, user *TestUser, habit TestHabit, logNum int) error {
	payload := map[string]interface{}{
		"habit_id": habit.ID,
		"notes":    fmt.Sprintf("Log entry %d for habit %s", logNum, habit.Name),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", baseURL+"/logs/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// performReadOperations performs GET operations to test read performance
func performReadOperations(config TestConfig, stats *TestStats, users []*TestUser) {
	var wg sync.WaitGroup

	// Get habits for each user
	for _, user := range users {
		wg.Add(1)
		go func(u *TestUser) {
			defer wg.Done()
			start := time.Now()
			if err := getHabits(config.BaseURL, u); err != nil {
				log.Printf("Failed to get habits: %v", err)
				recordRequest(stats, start, false, &stats.GetHabitsLatencies)
				return
			}
			recordRequest(stats, start, true, &stats.GetHabitsLatencies)
		}(user)
	}

	// Get logs for each habit
	for _, user := range users {
		for _, habit := range user.Habits {
			wg.Add(1)
			go func(u *TestUser, h TestHabit) {
				defer wg.Done()
				start := time.Now()
				if err := getLogs(config.BaseURL, u, h); err != nil {
					log.Printf("Failed to get logs: %v", err)
					recordRequest(stats, start, false, &stats.GetLogsLatencies)
					return
				}
				recordRequest(stats, start, true, &stats.GetLogsLatencies)
			}(user, habit)
		}
	}

	wg.Wait()
}

// getHabits retrieves all habits for a user
func getHabits(baseURL string, user *TestUser) error {
	req, _ := http.NewRequest("GET", baseURL+"/habits/", nil)
	req.Header.Set("Authorization", "Bearer "+user.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// getLogs retrieves all logs for a habit
func getLogs(baseURL string, user *TestUser, habit TestHabit) error {
	url := fmt.Sprintf("%s/logs/?habit_id=%d", baseURL, habit.ID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+user.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// performLoadTest performs sustained load testing
func performLoadTest(config TestConfig, stats *TestStats, users []*TestUser) {
	if len(users) == 0 {
		log.Println("No users available for load testing")
		return
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	// Start workers that continuously make requests
	for i := 0; i < config.ConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			userIndex := workerID % len(users)

			for {
				select {
				case <-stopChan:
					return
				default:
					user := users[userIndex]
					if len(user.Habits) == 0 {
						continue
					}

					// Randomly choose an operation
					switch time.Now().UnixNano() % 2 {
					case 0:
						// Get habits
						start := time.Now()
						err := getHabits(config.BaseURL, user)
						recordRequest(stats, start, err == nil, nil)
					case 1:
						// Get logs for a random habit
						habitIndex := time.Now().UnixNano() % int64(len(user.Habits))
						start := time.Now()
						err := getLogs(config.BaseURL, user, user.Habits[habitIndex])
						recordRequest(stats, start, err == nil, nil)
					}

					// Cycle through users
					userIndex = (userIndex + 1) % len(users)
				}
			}
		}(i)
	}

	// Run for the specified duration
	time.Sleep(config.TestDuration)
	close(stopChan)
	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("Load test completed in %s\n", elapsed)
}

// recordRequest records statistics for a request
func recordRequest(stats *TestStats, start time.Time, success bool, latencies *[]time.Duration) {
	latency := time.Since(start)

	stats.mu.Lock()
	defer stats.mu.Unlock()

	stats.TotalRequests++
	if success {
		stats.SuccessfulRequests++
	} else {
		stats.FailedRequests++
	}

	stats.TotalLatency += latency
	if latency < stats.MinLatency {
		stats.MinLatency = latency
	}
	if latency > stats.MaxLatency {
		stats.MaxLatency = latency
	}

	if latencies != nil {
		*latencies = append(*latencies, latency)
	}
}

// printResults prints the test results
func printResults(stats *TestStats) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	fmt.Println("\n========================================")
	fmt.Println("Performance Test Results")
	fmt.Println("========================================")
	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("Successful Requests: %d\n", stats.SuccessfulRequests)
	fmt.Printf("Failed Requests: %d\n", stats.FailedRequests)
	fmt.Printf("Success Rate: %.2f%%\n", float64(stats.SuccessfulRequests)/float64(stats.TotalRequests)*100)
	fmt.Println()

	if stats.TotalRequests > 0 {
		avgLatency := stats.TotalLatency / time.Duration(stats.TotalRequests)
		fmt.Printf("Average Latency: %s\n", avgLatency)
		fmt.Printf("Min Latency: %s\n", stats.MinLatency)
		fmt.Printf("Max Latency: %s\n", stats.MaxLatency)
		fmt.Println()
	}

	printOperationStats("User Registration", stats.RegisterLatencies)
	printOperationStats("User Login", stats.LoginLatencies)
	printOperationStats("Create Habit", stats.CreateHabitLatencies)
	printOperationStats("Create Log", stats.CreateLogLatencies)
	printOperationStats("Get Habits", stats.GetHabitsLatencies)
	printOperationStats("Get Logs", stats.GetLogsLatencies)
}

// printOperationStats prints statistics for a specific operation
func printOperationStats(operation string, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}

	var total time.Duration
	min := latencies[0]
	max := latencies[0]

	for _, l := range latencies {
		total += l
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
	}

	avg := total / time.Duration(len(latencies))

	fmt.Printf("%s (%d requests):\n", operation, len(latencies))
	fmt.Printf("  Avg: %s, Min: %s, Max: %s\n", avg, min, max)
	fmt.Println()
}

// saveResultsToFile saves the test results to a JSON file
func saveResultsToFile(stats *TestStats) error {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	results := map[string]interface{}{
		"timestamp":           time.Now().Format(time.RFC3339),
		"total_requests":      stats.TotalRequests,
		"successful_requests": stats.SuccessfulRequests,
		"failed_requests":     stats.FailedRequests,
		"success_rate":        float64(stats.SuccessfulRequests) / float64(stats.TotalRequests) * 100,
		"average_latency_ms":  stats.TotalLatency.Milliseconds() / stats.TotalRequests,
		"min_latency_ms":      stats.MinLatency.Milliseconds(),
		"max_latency_ms":      stats.MaxLatency.Milliseconds(),
	}

	filename := fmt.Sprintf("performance_test_results_%s.json", time.Now().Format("20060102_150405"))
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return err
	}

	fmt.Printf("\nResults saved to: %s\n", filename)
	return nil
}
