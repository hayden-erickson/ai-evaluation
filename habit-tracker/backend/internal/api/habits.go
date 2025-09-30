package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hayden-erickson/habit-tracker/internal/auth"
	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/hayden-erickson/habit-tracker/internal/models"
	"github.com/jinzhu/gorm"
)

// HabitInput represents the input for creating/updating a habit
type HabitInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Frequency   int    `json:"frequency" binding:"required,min=1"`
}

// GetHabits returns all habits for the current user
func GetHabits(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	var habits []models.Habit

	if err := db.DB.Where("user_id = ?", user.ID).Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch habits"})
		return
	}

	c.JSON(http.StatusOK, habits)
}

// GetHabit returns a specific habit by ID
func GetHabit(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	id := c.Param("id")

	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).
		Preload("Logs").
		Preload("Tags").
		First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	c.JSON(http.StatusOK, habit)
}

// CreateHabit creates a new habit
func CreateHabit(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	var input HabitInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	habit := models.Habit{
		UserID:      user.ID,
		Name:        input.Name,
		Description: input.Description,
		Frequency:   input.Frequency,
	}

	if err := db.DB.Create(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create habit"})
		return
	}

	c.JSON(http.StatusCreated, habit)
}

// UpdateHabit updates an existing habit
func UpdateHabit(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	id := c.Param("id")

	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var input HabitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := models.Habit{
		Name:        input.Name,
		Description: input.Description,
		Frequency:   input.Frequency,
	}

	if err := db.DB.Model(&habit).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update habit"})
		return
	}

	c.JSON(http.StatusOK, habit)
}

// DeleteHabit deletes a habit
func DeleteHabit(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	id := c.Param("id")

	result := db.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.Habit{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete habit"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHabitWithStreak returns a habit with its streak information
func GetHabitWithStreak(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	id := c.Param("id")

	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).
		Preload("Logs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	currentStreak, longestStreak := calculateStreaks(habit)
	var lastLogDate *time.Time
	if len(habit.Logs) > 0 {
		lastLogDate = &habit.Logs[0].CreatedAt
	}

	habitWithStreak := models.HabitWithStreak{
		Habit:         habit,
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastLogDate:   lastLogDate,
	}

	c.JSON(http.StatusOK, habitWithStreak)
}

// calculateStreaks calculates the current and longest streaks for a habit
func calculateStreaks(habit models.Habit) (currentStreak, longestStreak int) {
	if len(habit.Logs) == 0 {
		return 0, 0
	}

	currentStreak = 1
	longestStreak = 1
	currentStreakCount := 1

	// Convert all log dates to the user's timezone and get just the dates
	var dates []time.Time
	for _, log := range habit.Logs {
		dates = append(dates, log.CreatedAt)
	}

	// Sort dates in descending order
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})

	// Calculate streaks
	for i := 1; i < len(dates); i++ {
		diff := dates[i-1].Sub(dates[i]).Hours() / 24

		if diff <= float64(habit.Frequency) {
			currentStreakCount++
			if currentStreakCount > longestStreak {
				longestStreak = currentStreakCount
			}
		} else {
			currentStreakCount = 1
		}
	}

	// Check if the current streak is still active
	if len(dates) > 0 {
		daysSinceLastLog := time.Since(dates[0]).Hours() / 24
		if daysSinceLastLog > float64(habit.Frequency) {
			currentStreak = 0
		} else {
			currentStreak = currentStreakCount
		}
	}

	return currentStreak, longestStreak
}
