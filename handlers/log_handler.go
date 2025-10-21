package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/service"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	service service.LogService
}

// NewLogHandler creates a new log handler
func NewLogHandler(service service.LogService) *LogHandler {
	return &LogHandler{
		service: service,
	}
}

// CreateLog handles creating a new log (POST /habits/{habit_id}/logs)
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract habit ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		log.Println("Missing habit ID in path")
		http.Error(w, "Habit ID required", http.StatusBadRequest)
		return
	}

	habitID, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req models.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the log
	logEntry, err := h.service.CreateLog(habitID, userID, &req)
	if err != nil {
		log.Printf("Failed to create log: %v", err)
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Habit not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the created log
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

// GetLog handles getting a log by ID (GET /logs/{id})
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract log ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing log ID in path")
		http.Error(w, "Log ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid log ID: %v", err)
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	// Get the log
	logEntry, err := h.service.GetLog(id, userID)
	if err != nil {
		log.Printf("Failed to get log: %v", err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Log not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the log
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

// GetHabitLogs handles getting all logs for a habit (GET /habits/{habit_id}/logs)
func (h *LogHandler) GetHabitLogs(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract habit ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		log.Println("Missing habit ID in path")
		http.Error(w, "Habit ID required", http.StatusBadRequest)
		return
	}

	habitID, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Get all logs for the habit
	logs, err := h.service.GetHabitLogs(habitID, userID)
	if err != nil {
		log.Printf("Failed to get logs: %v", err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Habit not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the logs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// UpdateLog handles updating a log (PUT /logs/{id})
func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT requests
	if r.Method != http.MethodPut {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract log ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing log ID in path")
		http.Error(w, "Log ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid log ID: %v", err)
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req models.UpdateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the log
	logEntry, err := h.service.UpdateLog(id, userID, &req)
	if err != nil {
		log.Printf("Failed to update log: %v", err)
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Log not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the updated log
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

// DeleteLog handles deleting a log (DELETE /logs/{id})
func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE requests
	if r.Method != http.MethodDelete {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract log ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing log ID in path")
		http.Error(w, "Log ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid log ID: %v", err)
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	// Delete the log
	if err := h.service.DeleteLog(id, userID); err != nil {
		log.Printf("Failed to delete log: %v", err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Log not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
