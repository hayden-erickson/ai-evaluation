package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/config"
	"github.com/hayden-erickson/ai-evaluation/internal/db"
	httpserver "github.com/hayden-erickson/ai-evaluation/internal/http"
	"github.com/hayden-erickson/ai-evaluation/internal/http/middleware"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Wire repositories
	userRepo := repository.NewUserRepository(database)
	habitRepo := repository.NewHabitRepository(database)
	logRepo := repository.NewLogRepository(database)

	// Wire services
	authSvc := service.NewAuthService(userRepo, cfg)
	userSvc := service.NewUserService(userRepo)
	habitSvc := service.NewHabitService(habitRepo)
	logSvc := service.NewLogService(logRepo, habitRepo)

	// Build HTTP server
	mux := httpserver.NewRouter(authSvc, userSvc, habitSvc, logSvc, cfg.JWTSecret, cfg.JWTIssuer)

	// Secure headers middleware
	handler := middleware.SecureHeaders(middleware.RecoverPanic(mux))

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownErr <- srv.Shutdown(ctx)
	}()

	log.Printf("server starting on %s", cfg.ServerAddr)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Printf("server stopped")
}
