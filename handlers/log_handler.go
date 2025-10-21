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

// NewLogHandler creates a new LogHandler instance
func NewLogHandler(service service.LogService) *LogHandler {
	return &LogHandler{service: service}
}

// HandleLogs handles GET and POST /api/logs and /api/habits/{habit_id}/logs
func (h *LogHandler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Check if this is a request for logs by habit ID
		if strings.Contains(r.URL.Path, "/habits/") {
			h.listLogsByHabit(w, r, userID)
		} else {
			writeJSONError(w, http.StatusBadRequest, "habit_id is required")
		}
	case http.MethodPost:
		h.createLog(w, r, userID)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// listLogsByHabit handles GET request to list all logs for a habit
func (h *LogHandler) listLogsByHabit(w http.ResponseWriter, r *http.Request, userID int64) {
	// Parse habit ID from path /api/habits/{habit_id}/logs
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		log.Printf("Invalid path format: %s", r.URL.Path)
		writeJSONError(w, http.StatusBadRequest, "invalid path")
		return
	}

	habitID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid habit ID")
		return
	}

	logs, err := h.service.GetByHabitID(habitID, userID)
	if err != nil {
		log.Printf("Failed to get logs: %v", err)
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	// Return empty array if no logs found
	if logs == nil {
		logs = []*models.Log{}
	}

	writeJSONResponse(w, http.StatusOK, logs)
}

// createLog handles POST request to create a new log
func (h *LogHandler) createLog(w http.ResponseWriter, r *http.Request, userID int64) {
	var req models.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	logEntry, err := h.service.Create(userID, &req)
	if err != nil {
		log.Printf("Failed to create log: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, logEntry)
}

// HandleLog handles GET, PUT, DELETE /api/logs/{id}
func (h *LogHandler) HandleLog(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid log ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid log ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getLog(w, r, id, userID)
	case http.MethodPut:
		h.updateLog(w, r, id, userID)
	case http.MethodDelete:
		h.deleteLog(w, r, id, userID)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getLog handles GET request for a specific log
func (h *LogHandler) getLog(w http.ResponseWriter, r *http.Request, id, userID int64) {
	logEntry, err := h.service.GetByID(id, userID)
	if err != nil {
		log.Printf("Failed to get log: %v", err)
		writeJSONError(w, http.StatusNotFound, "log not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, logEntry)
}

// updateLog handles PUT request to update a log
func (h *LogHandler) updateLog(w http.ResponseWriter, r *http.Request, id, userID int64) {
	var req models.UpdateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	logEntry, err := h.service.Update(id, userID, &req)
	if err != nil {
		log.Printf("Failed to update log: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, logEntry)
}

// deleteLog handles DELETE request to delete a log
func (h *LogHandler) deleteLog(w http.ResponseWriter, r *http.Request, id, userID int64) {
	if err := h.service.Delete(id, userID); err != nil {
		log.Printf("Failed to delete log: %v", err)
		writeJSONError(w, http.StatusNotFound, "log not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, models.SuccessResponse{
		Message: "log deleted successfully",
	})
}
