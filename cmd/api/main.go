package main

import (
	"log"
	"net/http"
	"new-api/internal/auth"
	"new-api/internal/database"
	"new-api/internal/handlers"
	"new-api/internal/repository"
	"new-api/internal/service"
)

func main() {
	database.InitDB()
	db := database.DB

	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	logRepo := repository.NewLogRepository(db)

	userService := service.NewUserService(userRepo)
	habitService := service.NewHabitService(habitRepo)
	logService := service.NewLogService(logRepo, habitRepo)

	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)
	logHandler := handlers.NewLogHandler(logService)

	mux := http.NewServeMux()

	userHandler.UserRoutes(mux)

	// Protected routes
	habitMux := http.NewServeMux()
	habitHandler.HabitRoutes(habitMux)
	mux.Handle("/habits/", auth.JWTMiddleware(habitMux))

	logMux := http.NewServeMux()
	logHandler.LogRoutes(logMux)
	mux.Handle("/logs/", auth.JWTMiddleware(logMux))

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
