package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hayden-erickson/habit-tracker/internal/auth"
	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/hayden-erickson/habit-tracker/internal/models"
)

// TagInput represents the input for creating a tag
type TagInput struct {
	Value string `json:"value" binding:"required"`
}

// GetTags returns all tags for a specific habit
func GetTags(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("id")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var tags []models.Tag
	if err := db.DB.Where("habit_id = ?", habitID).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// CreateTag adds a new tag to a habit
func CreateTag(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("id")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var input TagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if tag already exists
	var existingTag models.Tag
	if err := db.DB.Where("habit_id = ? AND value = ?", habitID, input.Value).First(&existingTag).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tag already exists"})
		return
	}

	tag := models.Tag{
		HabitID: habit.ID,
		Value:   input.Value,
	}

	if err := db.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// DeleteTag removes a tag from a habit
func DeleteTag(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	habitID := c.Param("id")
	tagID := c.Param("tagId")

	// Verify habit belongs to user
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, user.ID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	result := db.DB.Where("id = ? AND habit_id = ?", tagID, habitID).Delete(&models.Tag{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHabitsByTag returns all habits that have a specific tag
func GetHabitsByTag(c *gin.Context) {
	user := auth.GetCurrentUser(c)
	tagValue := c.Query("value")

	if tagValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag value is required"})
		return
	}

	var habits []models.Habit
	if err := db.DB.Joins("JOIN tags ON tags.habit_id = habits.id").
		Where("habits.user_id = ? AND tags.value = ?", user.ID, tagValue).
		Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch habits"})
		return
	}

	c.JSON(http.StatusOK, habits)
}
