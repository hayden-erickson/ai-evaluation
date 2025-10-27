package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user to be notified
type User struct {
	ID          int
	PhoneNumber string
	Message     string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		fmt.Println("OK")
		os.Exit(0)
	}

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer db.Close()

	users, err := getUsersToNotify(db)
	if err != nil {
		log.Fatalf("Failed to get users to notify: %v", err)
	}

	for _, user := range users {
		sendTwilioNotification(user.PhoneNumber, user.Message)
	}
}

func setupDatabase() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection successful")
	return db, nil
}

func getUsersToNotify(db *sql.DB) ([]User, error) {
	query := `
        SELECT u.id, u.phone_number
        FROM users u
        WHERE u.id NOT IN (
            SELECT h.user_id
            FROM logs l
            JOIN habits h ON l.habit_id = h.id
            WHERE l.created_at >= ?
        )
    `
	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	rows, err := db.Query(query, twoDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("query for users with no logs failed: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.PhoneNumber); err != nil {
			log.Printf("Error scanning user with no logs: %v", err)
			continue
		}
		user.Message = "You haven't logged any habits in the last 2 days. Keep going!"
		users = append(users, user)
	}

	query = `
        SELECT u.id, u.phone_number
        FROM users u
        JOIN habits h ON u.id = h.user_id
        JOIN logs l ON h.id = l.habit_id
        GROUP BY u.id, u.phone_number
        HAVING
            MAX(CASE WHEN l.created_at BETWEEN ? AND ? THEN 1 ELSE 0 END) = 1
            AND MAX(CASE WHEN l.created_at > ? THEN 1 ELSE 0 END) = 0
    `
	day1Start := time.Now().Add(-48 * time.Hour)
	day1End := time.Now().Add(-24 * time.Hour)
	day2Start := time.Now().Add(-24 * time.Hour)

	rows, err = db.Query(query, day1Start, day1End, day2Start)
	if err != nil {
		return nil, fmt.Errorf("query for users with logs on day 1 but not day 2 failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.PhoneNumber); err != nil {
			log.Printf("Error scanning user with logs on day 1 but not day 2: %v", err)
			continue
		}
		user.Message = "You logged a habit yesterday, but not today. Don't break your streak!"
		users = append(users, user)
	}

	return users, nil
}

func sendTwilioNotification(toPhoneNumber string, message string) {
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if twilioAccountSid == "" || twilioAuthToken == "" || twilioPhoneNumber == "" {
		log.Println("Twilio environment variables not set. Skipping notification.")
		return
	}

	if toPhoneNumber == "" {
		log.Println("User phone number is empty. Skipping notification.")
		return
	}

	urlString := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", twilioAccountSid)
	msgData := url.Values{}
	msgData.Set("To", toPhoneNumber)
	msgData.Set("From", twilioPhoneNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", urlString, &msgDataReader)
	req.SetBasicAuth(twilioAccountSid, twilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Twilio notification to %s: %v", toPhoneNumber, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Printf("Twilio notification sent successfully to %s: %s", toPhoneNumber, data["sid"])
		}
	} else {
		log.Printf("Twilio API error for %s: %s", toPhoneNumber, resp.Status)
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Println(data)
		}
	}
}
