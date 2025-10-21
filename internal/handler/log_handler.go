package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	service service.LogService
}

// NewLogHandler creates a new log handler instance
func NewLogHandler(service service.LogService) *LogHandler {
	return &LogHandler{service: service}
}

// HandleLogs routes log requests based on HTTP method
func (h *LogHandler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("Error: User ID not found in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.createLog(w, r, userID)
	case http.MethodGet:
		// Check if this is a list by habit or get by ID
		path := strings.TrimPrefix(r.URL.Path, "/api/logs")
		if path == "" || path == "/" {
			// List by habit ID from query parameter
			h.listLogsByHabit(w, r, userID)
		} else {
			h.getLog(w, r, userID)
		}
	case http.MethodPut:
		h.updateLog(w, r, userID)
	case http.MethodDelete:
		h.deleteLog(w, r, userID)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// createLog handles POST /api/logs
func (h *LogHandler) createLog(w http.ResponseWriter, r *http.Request, userID string) {
	var req models.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding create log request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	logEntry, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Error creating log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, logEntry)
}

// getLog handles GET /api/logs/{id}
func (h *LogHandler) getLog(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	logEntry, err := h.service.GetByID(r.Context(), userID, id)
	if err != nil {
		log.Printf("Error getting log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, logEntry)
}

// updateLog handles PUT /api/logs/{id}
func (h *LogHandler) updateLog(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	var req models.UpdateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding update log request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	logEntry, err := h.service.Update(r.Context(), userID, id, &req)
	if err != nil {
		log.Printf("Error updating log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, logEntry)
}

// deleteLog handles DELETE /api/logs/{id}
func (h *LogHandler) deleteLog(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	if err := h.service.Delete(r.Context(), userID, id); err != nil {
		log.Printf("Error deleting log: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listLogsByHabit handles GET /api/logs?habit_id={habitID}
func (h *LogHandler) listLogsByHabit(w http.ResponseWriter, r *http.Request, userID string) {
	habitID := r.URL.Query().Get("habit_id")
	if habitID == "" {
		respondWithError(w, http.StatusBadRequest, "habit_id query parameter is required")
		return
	}

	logs, err := h.service.ListByHabitID(r.Context(), userID, habitID)
	if err != nil {
		log.Printf("Error listing logs: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}
