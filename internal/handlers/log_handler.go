package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler creates a new instance of LogHandler
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// CreateLog handles POST requests to create a new log entry
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	var req models.CreateLogRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode create log request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Create log request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create log
	logEntry, err := h.logService.CreateLog(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Failed to create log: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: habit does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create log")
		return
	}

	respondJSON(w, http.StatusCreated, logEntry)
}

// GetLog handles GET requests to retrieve a specific log entry
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract log ID from path
	logID := r.PathValue("id")
	if logID == "" {
		respondError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	// Retrieve log
	logEntry, err := h.logService.GetLog(r.Context(), logID, userID)
	if err != nil {
		log.Printf("Failed to get log: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: log does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusNotFound, "Log not found")
		return
	}

	respondJSON(w, http.StatusOK, logEntry)
}

// ListLogs handles GET requests to retrieve all log entries for a habit
func (h *LogHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract habit ID from query parameter
	habitID := r.URL.Query().Get("habit_id")
	if habitID == "" {
		respondError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	// Retrieve logs
	logs, err := h.logService.ListLogs(r.Context(), userID, habitID)
	if err != nil {
		log.Printf("Failed to list logs: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: habit does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve logs")
		return
	}

	respondJSON(w, http.StatusOK, logs)
}

// UpdateLog handles PUT requests to update a log entry
func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract log ID from path
	logID := r.PathValue("id")
	if logID == "" {
		respondError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	var req models.UpdateLogRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode update log request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Update log request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update log
	logEntry, err := h.logService.UpdateLog(r.Context(), logID, userID, &req)
	if err != nil {
		log.Printf("Failed to update log: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: log does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update log")
		return
	}

	respondJSON(w, http.StatusOK, logEntry)
}

// DeleteLog handles DELETE requests to remove a log entry
func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract log ID from path
	logID := r.PathValue("id")
	if logID == "" {
		respondError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	// Delete log
	if err := h.logService.DeleteLog(r.Context(), logID, userID); err != nil {
		log.Printf("Failed to delete log: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: log does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete log")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Log deleted successfully"})
}
