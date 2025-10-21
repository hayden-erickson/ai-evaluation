package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/middleware"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler creates a new log handler
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// CreateLog handles POST /habits/{habit_id}/logs - creates a new log
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Method not allowed: %s on /habits/{habit_id}/logs", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get habit ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/habits/")
	parts := strings.Split(path, "/logs")
	if len(parts) != 2 {
		log.Printf("ERROR: Invalid URL path: %s", r.URL.Path)
		respondWithError(w, http.StatusBadRequest, "invalid URL path")
		return
	}
	habitID := parts[0]

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request body
	var req models.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create log
	logEntry, err := h.logService.Create(r.Context(), habitID, userID, &req)
	if err != nil {
		log.Printf("ERROR: Failed to create log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

// GetLogs handles GET /habits/{habit_id}/logs - retrieves all logs for a habit
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		log.Printf("ERROR: Method not allowed: %s on /habits/{habit_id}/logs", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get habit ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/habits/")
	parts := strings.Split(path, "/logs")
	if len(parts) != 2 {
		log.Printf("ERROR: Invalid URL path: %s", r.URL.Path)
		respondWithError(w, http.StatusBadRequest, "invalid URL path")
		return
	}
	habitID := parts[0]

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get logs
	logs, err := h.logService.GetByHabitID(r.Context(), habitID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get logs: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to get logs")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// GetLog handles GET /logs/{id} - retrieves a specific log
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		log.Printf("ERROR: Method not allowed: %s on /logs/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get log ID from URL path
	logID := strings.TrimPrefix(r.URL.Path, "/logs/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get log
	logEntry, err := h.logService.GetByID(r.Context(), logID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, "log not found")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logEntry)
}

// UpdateLog handles PUT /logs/{id} - updates a log
func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT method
	if r.Method != http.MethodPut {
		log.Printf("ERROR: Method not allowed: %s on /logs/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get log ID from URL path
	logID := strings.TrimPrefix(r.URL.Path, "/logs/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request body
	var req models.UpdateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update log
	logEntry, err := h.logService.Update(r.Context(), logID, userID, &req)
	if err != nil {
		log.Printf("ERROR: Failed to update log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logEntry)
}

// DeleteLog handles DELETE /logs/{id} - deletes a log
func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	// Only accept DELETE method
	if r.Method != http.MethodDelete {
		log.Printf("ERROR: Method not allowed: %s on /logs/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get log ID from URL path
	logID := strings.TrimPrefix(r.URL.Path, "/logs/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		return
	}

	// Delete log
	if err := h.logService.Delete(r.Context(), logID, userID); err != nil {
		log.Printf("ERROR: Failed to delete log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to delete log")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

