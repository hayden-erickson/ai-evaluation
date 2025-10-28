package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string
}

type User struct {
	ID          int64
	Name        string
	PhoneNumber string
	TimeZone    string
}

type LogSummary struct {
	UserID     int64
	Day1HasLog bool
	Day2HasLog bool
}

func main() {
	log.Println("Starting habit notification cron job...")

	// Load configuration from environment variables
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MySQL database
	db, err := connectDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Twilio client
	twilioClient := NewTwilioClient(config.TwilioAccountSID, config.TwilioAuthToken, config.TwilioFromNumber)

	// Process notifications
	if err := processNotifications(db, twilioClient); err != nil {
		log.Fatalf("Failed to process notifications: %v", err)
	}

	log.Println("Habit notification cron job completed successfully")
}

func loadConfig() (*Config, error) {
	config := &Config{
		DBHost:           getEnv("DB_HOST", ""),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", ""),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", ""),
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
	}

	// Validate required fields
	if config.DBHost == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if config.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if config.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if config.TwilioAccountSID == "" {
		return nil, fmt.Errorf("TWILIO_ACCOUNT_SID is required")
	}
	if config.TwilioAuthToken == "" {
		return nil, fmt.Errorf("TWILIO_AUTH_TOKEN is required")
	}
	if config.TwilioFromNumber == "" {
		return nil, fmt.Errorf("TWILIO_FROM_NUMBER is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDatabase(config *Config) (*sql.DB, error) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}

func processNotifications(db *sql.DB, twilioClient *TwilioClient) error {
	log.Println("Fetching users and their log summaries...")

	// Get all users
	users, err := getAllUsers(db)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	log.Printf("Found %d users to process", len(users))

	// Get log summaries for all users
	logSummaries, err := getLogSummaries(db)
	if err != nil {
		return fmt.Errorf("failed to get log summaries: %w", err)
	}

	// Create a map for quick lookup
	summaryMap := make(map[int64]*LogSummary)
	for i := range logSummaries {
		summaryMap[logSummaries[i].UserID] = &logSummaries[i]
	}

	// Process each user
	notificationsSent := 0
	notificationsFailed := 0

	for _, user := range users {
		summary, exists := summaryMap[user.ID]

		// Determine if notification should be sent
		shouldNotify := false
		message := ""

		if !exists || (!summary.Day1HasLog && !summary.Day2HasLog) {
			// Case 1: No logs over the last 2 days
			shouldNotify = true
			message = fmt.Sprintf("Hi %s! We noticed you haven't logged your habit in the past 2 days. Keep up your streak!", user.Name)
		} else if summary.Day1HasLog && !summary.Day2HasLog {
			// Case 2: Log on day 1, but no log on day 2
			shouldNotify = true
			message = fmt.Sprintf("Hi %s! You logged yesterday but not today. Don't break your streak!", user.Name)
		}

		if shouldNotify {
			log.Printf("Sending notification to user %d (%s) at %s", user.ID, user.Name, user.PhoneNumber)

			err := twilioClient.SendSMS(user.PhoneNumber, message)
			if err != nil {
				log.Printf("Failed to send notification to user %d: %v", user.ID, err)
				notificationsFailed++
			} else {
				log.Printf("Successfully sent notification to user %d", user.ID)
				notificationsSent++
			}
		} else {
			log.Printf("User %d (%s) has logs for both days, no notification needed", user.ID, user.Name)
		}
	}

	log.Printf("Notification summary: %d sent, %d failed", notificationsSent, notificationsFailed)

	if notificationsFailed > 0 {
		return fmt.Errorf("%d notifications failed to send", notificationsFailed)
	}

	return nil
}

func getAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT id, name, phone_number, time_zone
		FROM users
		ORDER BY id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.TimeZone)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func getLogSummaries(db *sql.DB) ([]LogSummary, error) {
	// Calculate date boundaries
	now := time.Now().UTC()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfYesterday := startOfToday.Add(-24 * time.Hour)
	startOfTwoDaysAgo := startOfToday.Add(-48 * time.Hour)

	log.Printf("Checking logs from %s to %s", startOfTwoDaysAgo.Format(time.RFC3339), startOfToday.Format(time.RFC3339))

	// Query to get log summaries
	// Day 1 = yesterday (24-48 hours ago)
	// Day 2 = today (0-24 hours ago)
	query := `
		SELECT 
			h.user_id,
			MAX(CASE 
				WHEN l.created_at >= ? AND l.created_at < ? 
				THEN 1 ELSE 0 
			END) AS day1_has_log,
			MAX(CASE 
				WHEN l.created_at >= ? AND l.created_at < ? 
				THEN 1 ELSE 0 
			END) AS day2_has_log
		FROM habits h
		LEFT JOIN logs l ON h.id = l.habit_id
		WHERE l.created_at >= ?
		GROUP BY h.user_id
	`

	rows, err := db.Query(query,
		startOfYesterday, startOfToday, // day1 range
		startOfToday, now, // day2 range
		startOfTwoDaysAgo, // overall filter
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query log summaries: %w", err)
	}
	defer rows.Close()

	var summaries []LogSummary
	for rows.Next() {
		var summary LogSummary
		var day1HasLog, day2HasLog int
		err := rows.Scan(&summary.UserID, &day1HasLog, &day2HasLog)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log summary row: %w", err)
		}
		summary.Day1HasLog = day1HasLog == 1
		summary.Day2HasLog = day2HasLog == 1
		summaries = append(summaries, summary)

		log.Printf("User %d: Day1=%v, Day2=%v", summary.UserID, summary.Day1HasLog, summary.Day2HasLog)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating log summary rows: %w", err)
	}

	return summaries, nil
}
