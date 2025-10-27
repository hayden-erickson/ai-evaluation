package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// User represents a user with their phone number
type User struct {
	ID          int64
	PhoneNumber string
	TimeZone    string
	Name        string
}

// LogEntry represents a habit log entry
type LogEntry struct {
	UserID    int64
	CreatedAt time.Time
}

func main() {
	log.Println("Starting Habit Log Notification CronJob")

	// Get configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := os.Getenv("DB_PASSWORD") // Required
	dbName := getEnv("DB_NAME", "habits")

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")   // Required
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")     // Required
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER") // Required

	// Validate required environment variables
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}
	if twilioAccountSID == "" {
		log.Fatal("TWILIO_ACCOUNT_SID environment variable is required")
	}
	if twilioAuthToken == "" {
		log.Fatal("TWILIO_AUTH_TOKEN environment variable is required")
	}
	if twilioPhoneNumber == "" {
		log.Fatal("TWILIO_PHONE_NUMBER environment variable is required")
	}

	// Connect to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Twilio client
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSID,
		Password: twilioAuthToken,
	})

	// Process notifications
	if err := processNotifications(db, twilioClient, twilioPhoneNumber); err != nil {
		log.Fatalf("Failed to process notifications: %v", err)
	}

	log.Println("Habit Log Notification CronJob completed successfully")
}

// processNotifications queries the database and sends notifications
func processNotifications(db *sql.DB, twilioClient *twilio.RestClient, fromPhoneNumber string) error {
	// Get all users
	users, err := getAllUsers(db)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	log.Printf("Found %d users to check", len(users))

	// Calculate time windows for the past 2 days
	now := time.Now()
	twoDaysAgo := now.AddDate(0, 0, -2)
	oneDayAgo := now.AddDate(0, 0, -1)

	notificationsSent := 0
	notificationsFailed := 0

	for _, user := range users {
		// Get logs for this user over the past 2 days
		logs, err := getUserLogsInTimeRange(db, user.ID, twoDaysAgo)
		if err != nil {
			log.Printf("Error getting logs for user %d: %v", user.ID, err)
			continue
		}

		shouldNotify := false
		var reason string

		if len(logs) == 0 {
			// No logs over the last 2 days
			shouldNotify = true
			reason = "no logs in the past 2 days"
		} else {
			// Check if there's a log on day 2 (yesterday)
			hasLogDay2 := false
			hasLogDay1 := false

			for _, logEntry := range logs {
				if logEntry.CreatedAt.After(oneDayAgo) {
					hasLogDay2 = true
				} else if logEntry.CreatedAt.After(twoDaysAgo) && logEntry.CreatedAt.Before(oneDayAgo) {
					hasLogDay1 = true
				}
			}

			if hasLogDay1 && !hasLogDay2 {
				// Has log on day 1, but not on day 2
				shouldNotify = true
				reason = "logged on day 1 but not on day 2"
			}
		}

		if shouldNotify {
			log.Printf("Sending notification to user %d (%s) - Reason: %s", user.ID, user.PhoneNumber, reason)
			if err := sendNotification(twilioClient, fromPhoneNumber, user, reason); err != nil {
				log.Printf("Failed to send notification to user %d: %v", user.ID, err)
				notificationsFailed++
			} else {
				notificationsSent++
			}
		}
	}

	log.Printf("Notifications sent: %d, failed: %d", notificationsSent, notificationsFailed)
	return nil
}

// getAllUsers retrieves all users from the database
func getAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, phone_number, time_zone, name FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.PhoneNumber, &user.TimeZone, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// getUserLogsInTimeRange retrieves all logs for a user since the given time
func getUserLogsInTimeRange(db *sql.DB, userID int64, since time.Time) ([]LogEntry, error) {
	query := `
		SELECT l.created_at, h.user_id
		FROM logs l
		INNER JOIN habits h ON l.habit_id = h.id
		WHERE h.user_id = ? AND l.created_at >= ?
		ORDER BY l.created_at DESC
	`

	rows, err := db.Query(query, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var logEntry LogEntry
		if err := rows.Scan(&logEntry.CreatedAt, &logEntry.UserID); err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	return logs, rows.Err()
}

// sendNotification sends an SMS notification via Twilio
func sendNotification(client *twilio.RestClient, fromPhone string, user User, reason string) error {
	message := fmt.Sprintf("Hi %s! We noticed you haven't logged your habits recently. Keep up your streak! 💪", user.Name)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(user.PhoneNumber)
	params.SetFrom(fromPhone)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if resp.ErrorCode != nil {
		return fmt.Errorf("twilio error: %d - %s", *resp.ErrorCode, *resp.ErrorMessage)
	}

	log.Printf("Message sent successfully to %s (SID: %s)", user.PhoneNumber, *resp.Sid)
	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
