package api

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/hayden-erickson/habit-tracker/internal/models"
	"github.com/jinzhu/gorm"
)

// CheckAndSendReminders checks all habits and sends reminders for those that need attention
func CheckAndSendReminders() {
	var habits []models.Habit
	if err := db.DB.Preload("User").Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Find(&habits).Error; err != nil {
		log.Printf("Error fetching habits: %v", err)
		return
	}

	for _, habit := range habits {
		if shouldSendReminder(habit) {
			sendReminder(habit)
		}
	}
}

// shouldSendReminder checks if a reminder should be sent for a habit
func shouldSendReminder(habit models.Habit) bool {
	if len(habit.Logs) == 0 {
		return true
	}

	lastLog := habit.Logs[0]
	daysSinceLastLog := time.Since(lastLog.CreatedAt).Hours() / 24
	return daysSinceLastLog > float64(habit.Frequency)
}

// sendReminder sends an SMS reminder for a habit
func sendReminder(habit models.Habit) {
	apiKey := os.Getenv("SMS_API_KEY")
	if apiKey == "" {
		log.Println("SMS_API_KEY not set")
		return
	}

	// Create the habit view URL (update this with your frontend URL)
	habitURL := fmt.Sprintf("http://localhost:3000/habits/%d", habit.ID)
	message := fmt.Sprintf(
		"Reminder: You haven't logged your habit '%s' in %d days. Keep up your streak! %s",
		habit.Name,
		habit.Frequency,
		habitURL,
	)

	// You'll need to replace this with your actual SMS provider's API
	// This is an example using Twilio-like API
	data := url.Values{}
	data.Set("To", habit.User.Phone)
	data.Set("From", os.Getenv("SMS_FROM_NUMBER"))
	data.Set("Body", message)

	// Make the API call to send the SMS
	// This is a placeholder - implement actual SMS sending logic
	log.Printf("Would send SMS to %s: %s", habit.User.Phone, message)
}

// StartReminderScheduler starts a goroutine that periodically checks for habits that need reminders
func StartReminderScheduler() {
	go func() {
		for {
			CheckAndSendReminders()
			// Check every hour
			time.Sleep(1 * time.Hour)
		}
	}()
}
