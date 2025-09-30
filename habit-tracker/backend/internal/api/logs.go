package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hayden-erickson/habit-tracker/internal/auth"
	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/hayden-erickson/habit-tracker/internal/models"
)

// LogInput represents the input for creating/updating a log
type LogInput struct {
	Notes string `json:"notes"`
}

// GetLogs returns all logs for a specific habit
func GetLogs(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("habitId")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var logs []models.Log
	if err := db.DB.Where("habit_id = ?", habitID).Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// CreateLog creates a new log for a habit
func CreateLog(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("habitId")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var input LogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log := models.Log{
		HabitID: habit.ID,
		Notes:   input.Notes,
	}

	if err := db.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create log"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// UpdateLog updates an existing log
func UpdateLog(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("habitId")
	logID := c.Param("logId")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var log models.Log
	if err := db.DB.Where("id = ? AND habit_id = ?", logID, habitID).First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	var input LogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&log).Update("notes", input.Notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// DeleteLog deletes a log
func DeleteLog(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("habitId")
	logID := c.Param("logId")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	result := db.DB.Where("id = ? AND habit_id = ?", logID, habitID).Delete(&models.Log{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete log"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	c.Status(http.StatusNoContent)
}