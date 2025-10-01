package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hayden-erickson/habit-tracker/internal/api"
	"github.com/hayden-erickson/habit-tracker/internal/auth"
	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db.Initialize()
	defer db.Close()

	// Initialize authentication
	auth.Initialize()

	// Start reminder scheduler
	api.StartReminderScheduler()

	// Initialize router
	r := gin.Default()

	// Initialize middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Authentication routes
	r.GET("/auth/google/login", auth.LoginHandler)
	r.POST("/auth/google/callback", auth.CallbackHandler)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		// Habits routes
		habits := protected.Group("/habits")
		{
			habits.GET("", api.GetHabits)
			habits.POST("", api.CreateHabit)
			habits.GET("/:id", api.GetHabit)
			habits.PUT("/:id", api.UpdateHabit)
			habits.DELETE("/:id", api.DeleteHabit)
			habits.GET("/:id/streak", api.GetHabitWithStreak)
			habits.GET("/by-tag", api.GetHabitsByTag)

			// Logs routes
			logs := habits.Group("/:id/logs")
			{
				logs.GET("", api.GetLogs)
				logs.POST("", api.CreateLog)
				logs.PUT("/:logId", api.UpdateLog)
				logs.DELETE("/:logId", api.DeleteLog)
			}

			// Tags routes
			tags := habits.Group("/:id/tags")
			{
				tags.GET("", api.GetTags)
				tags.POST("", api.CreateTag)
				tags.DELETE("/:tagId", api.DeleteTag)
			}
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
