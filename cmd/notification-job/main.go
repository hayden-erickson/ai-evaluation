package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// UserNotificationStatus represents a user and their log status
type UserNotificationStatus struct {
	UserID      int64
	Name        string
	PhoneNumber string
	TimeZone    string
	HasDay1Log  bool
	HasDay2Log  bool
}

func main() {
	log.Println("Starting habit log notification job...")

	// Get environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := getEnv("DB_NAME", "habits")

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

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

	// Query users and their log status
	users, err := getUsersWithLogStatus(db)
	if err != nil {
		log.Fatalf("Failed to get users with log status: %v", err)
	}

	log.Printf("Found %d users to process", len(users))

	// Send notifications
	sentCount := 0
	for _, user := range users {
		shouldNotify := false
		var message string

		// Determine if notification should be sent
		if !user.HasDay1Log && !user.HasDay2Log {
			// No logs over the last 2 days
			shouldNotify = true
			message = fmt.Sprintf("Hi %s! We noticed you haven't logged any habits in the past 2 days. Don't break your streak! Log your habits today.", user.Name)
		} else if user.HasDay1Log && !user.HasDay2Log {
			// Has log on day 1 but not day 2
			shouldNotify = true
			message = fmt.Sprintf("Hi %s! You logged habits yesterday but not today. Keep your momentum going and log your habits!", user.Name)
		}

		if shouldNotify {
			err := sendNotification(twilioClient, twilioPhoneNumber, user.PhoneNumber, message)
			if err != nil {
				log.Printf("Failed to send notification to user %d (%s): %v", user.UserID, user.PhoneNumber, err)
			} else {
				log.Printf("Successfully sent notification to user %d (%s)", user.UserID, user.PhoneNumber)
				sentCount++
			}
		} else {
			log.Printf("User %d (%s) has logs for both days, skipping notification", user.UserID, user.Name)
		}
	}

	log.Printf("Notification job completed. Sent %d notifications out of %d users", sentCount, len(users))
}

// getUsersWithLogStatus queries the database for all users and their log status over the past 2 days
func getUsersWithLogStatus(db *sql.DB) ([]UserNotificationStatus, error) {
	// Calculate date boundaries for the past 2 days
	now := time.Now().UTC()
	day1Start := now.AddDate(0, 0, -1).Truncate(24 * time.Hour) // Start of yesterday
	day1End := day1Start.Add(24 * time.Hour)                    // End of yesterday
	day2Start := now.Truncate(24 * time.Hour)                   // Start of today
	day2End := day2Start.Add(24 * time.Hour)                    // End of today

	log.Printf("Querying logs for Day 1: %s to %s", day1Start.Format(time.RFC3339), day1End.Format(time.RFC3339))
	log.Printf("Querying logs for Day 2: %s to %s", day2Start.Format(time.RFC3339), day2End.Format(time.RFC3339))

	query := `
		SELECT 
			u.id,
			u.name,
			u.phone_number,
			u.time_zone,
			EXISTS(
				SELECT 1 FROM logs l
				INNER JOIN habits h ON l.habit_id = h.id
				WHERE h.user_id = u.id
				AND l.created_at >= ?
				AND l.created_at < ?
			) AS has_day1_log,
			EXISTS(
				SELECT 1 FROM logs l
				INNER JOIN habits h ON l.habit_id = h.id
				WHERE h.user_id = u.id
				AND l.created_at >= ?
				AND l.created_at < ?
			) AS has_day2_log
		FROM users u
		ORDER BY u.id
	`

	rows, err := db.Query(query, day1Start, day1End, day2Start, day2End)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []UserNotificationStatus
	for rows.Next() {
		var user UserNotificationStatus
		var hasDay1Log, hasDay2Log int

		err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.PhoneNumber,
			&user.TimeZone,
			&hasDay1Log,
			&hasDay2Log,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.HasDay1Log = hasDay1Log == 1
		user.HasDay2Log = hasDay2Log == 1

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// sendNotification sends an SMS notification via Twilio
func sendNotification(client *twilio.RestClient, from, to, message string) error {
	params := &api.CreateMessageParams{}
	params.SetFrom(from)
	params.SetTo(to)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("twilio API error: %w", err)
	}

	if resp.ErrorCode != nil {
		return fmt.Errorf("twilio error code %d: %s", *resp.ErrorCode, *resp.ErrorMessage)
	}

	log.Printf("Message sent successfully. SID: %s", *resp.Sid)
	return nil
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
