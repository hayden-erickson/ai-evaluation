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

// NotificationCandidate represents a user who should receive a notification
type NotificationCandidate struct {
	UserID      int64
	Name        string
	PhoneNumber string
	TimeZone    string
	Reason      string // Why they're receiving a notification
}

func main() {
	log.Println("Starting habit tracking notification cron job")

	// Get configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := os.Getenv("DB_PASSWORD") // Required, no default
	dbName := getEnv("DB_NAME", "habits")

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID") // Required
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")   // Required
	twilioFromNumber := os.Getenv("TWILIO_FROM_NUMBER") // Required

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
	if twilioFromNumber == "" {
		log.Fatal("TWILIO_FROM_NUMBER environment variable is required")
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)
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

	// Find users who need notifications
	candidates, err := findNotificationCandidates(db)
	if err != nil {
		log.Fatalf("Failed to find notification candidates: %v", err)
	}

	log.Printf("Found %d users who need notifications", len(candidates))

	// Send notifications
	successCount := 0
	failureCount := 0
	for _, candidate := range candidates {
		if err := sendNotification(twilioClient, twilioFromNumber, candidate); err != nil {
			log.Printf("Failed to send notification to user %d (%s): %v",
				candidate.UserID, candidate.PhoneNumber, err)
			failureCount++
		} else {
			log.Printf("Successfully sent notification to user %d (%s)",
				candidate.UserID, candidate.PhoneNumber)
			successCount++
		}
	}

	log.Printf("Notification job completed. Success: %d, Failures: %d", successCount, failureCount)
}

// findNotificationCandidates finds users who should receive notifications
// based on their habit logs over the past 2 days
func findNotificationCandidates(db *sql.DB) ([]NotificationCandidate, error) {
	// Calculate date boundaries
	now := time.Now().UTC()
	twoDaysAgo := now.Add(-48 * time.Hour)
	oneDayAgo := now.Add(-24 * time.Hour)

	// Optimized query using LEFT JOIN and aggregation to find users who need notifications
	// This is more efficient than multiple EXISTS subqueries
	query := `
		SELECT DISTINCT 
			u.id, 
			u.name, 
			u.phone_number, 
			u.time_zone,
			COALESCE(SUM(CASE WHEN l.created_at >= ? AND l.created_at < ? THEN 1 ELSE 0 END), 0) as day1_count,
			COALESCE(SUM(CASE WHEN l.created_at >= ? THEN 1 ELSE 0 END), 0) as day2_count
		FROM users u
		INNER JOIN habits h ON u.id = h.user_id
		LEFT JOIN logs l ON l.habit_id = h.id AND l.created_at >= ?
		GROUP BY u.id, u.name, u.phone_number, u.time_zone
		HAVING 
			-- Condition 1: No logs in the last 2 days
			(day1_count = 0 AND day2_count = 0)
			OR
			-- Condition 2: Has logs between 2 days ago and 1 day ago, but none in the last day
			(day1_count > 0 AND day2_count = 0)
	`

	rows, err := db.Query(query, twoDaysAgo, oneDayAgo, oneDayAgo, twoDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to query notification candidates: %w", err)
	}
	defer rows.Close()

	var candidates []NotificationCandidate
	for rows.Next() {
		var c NotificationCandidate
		var day1Count, day2Count int
		if err := rows.Scan(&c.UserID, &c.Name, &c.PhoneNumber, &c.TimeZone, &day1Count, &day2Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Determine the reason for notification based on log counts
		if day1Count == 0 && day2Count == 0 {
			c.Reason = "no_logs_2_days"
		} else if day1Count > 0 && day2Count == 0 {
			c.Reason = "logged_yesterday_not_today"
		} else {
			c.Reason = "missed_logging"
		}
		
		candidates = append(candidates, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return candidates, nil
}

// sendNotification sends a notification via Twilio
func sendNotification(client *twilio.RestClient, fromNumber string, candidate NotificationCandidate) error {
	// Create a personalized message based on the reason
	var message string
	switch candidate.Reason {
	case "no_logs_2_days":
		message = fmt.Sprintf("Hi %s! We noticed you haven't logged any habits in the past 2 days. Don't break your streak! Log your habits now.", candidate.Name)
	case "logged_yesterday_not_today":
		message = fmt.Sprintf("Hi %s! You logged your habits yesterday but not today. Keep up the momentum and log your habits now!", candidate.Name)
	default:
		message = fmt.Sprintf("Hi %s! Don't forget to log your habits today. Stay on track with your goals!", candidate.Name)
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(candidate.PhoneNumber)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("twilio API error: %w", err)
	}

	if resp.Sid != nil {
		log.Printf("Message sent successfully. SID: %s", *resp.Sid)
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
