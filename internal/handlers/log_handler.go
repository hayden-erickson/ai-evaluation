package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/hayden-erickson/ai-evaluation/internal/auth"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
	"strconv"
)

// LogHandler handles log-related requests
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// CreateLog handles the creation of a new log
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var log models.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.logService.CreateLog(&log, userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create log")
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "Log created successfully",
		Data:    log,
	})
}

// GetLog handles retrieving a log by ID
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid log ID")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	log, err := h.logService.GetLog(id, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Log not found")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    log,
	})
}

// GetHabitLogs handles retrieving all logs for a habit
func (h *LogHandler) GetHabitLogs(w http.ResponseWriter, r *http.Request) {
	habitID, err := strconv.ParseInt(r.URL.Query().Get("habit_id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logs, err := h.logService.GetHabitLogs(habitID, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve logs")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    logs,
	})
}

// UpdateLog handles updating a log
func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	var log models.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.logService.UpdateLog(&log, userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update log")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Log updated successfully",
	})
}

// DeleteLog handles deleting a log
func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid log ID")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.logService.DeleteLog(id, userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete log")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Log deleted successfully",
	})
}

// LogRoutes registers the log routes
func (h *LogHandler) LogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateLog(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				h.GetLog(w, r)
			} else if r.URL.Query().Get("habit_id") != "" {
				h.GetHabitLogs(w, r)
			} else {
				sendErrorResponse(w, http.StatusBadRequest, "Missing 'id' or 'habit_id' query parameter")
			}
		case http.MethodPut:
			h.UpdateLog(w, r)
		case http.MethodDelete:
			h.DeleteLog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
