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
	return &LogHandler{service: service}
}

// Create handles log creation
func (h *LogHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var req models.LogCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	logEntry, err := h.service.Create(userID, &req)
	if err != nil {
		log.Printf("Error creating log: %v", err)
		if err == service.ErrUnauthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

// ListByHabit handles listing all logs for a specific habit
func (h *LogHandler) ListByHabit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get habit ID from query parameter
	habitIDStr := r.URL.Query().Get("habit_id")
	if habitIDStr == "" {
		http.Error(w, "habit_id query parameter is required", http.StatusBadRequest)
		return
	}
	
	habitID, err := strconv.ParseInt(habitIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid habit_id", http.StatusBadRequest)
		return
	}
	
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	logs, err := h.service.GetByHabitID(habitID, userID)
	if err != nil {
		log.Printf("Error getting logs: %v", err)
		if err == service.ErrUnauthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// HandleLog handles GET, PUT, DELETE operations for a specific log
func (h *LogHandler) HandleLog(w http.ResponseWriter, r *http.Request) {
	// Get log ID from URL path
	logIDStr := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}
	
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		h.getLog(w, r, logID, userID)
	case http.MethodPut:
		h.updateLog(w, r, logID, userID)
	case http.MethodDelete:
		h.deleteLog(w, r, logID, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getLog retrieves a log by ID
func (h *LogHandler) getLog(w http.ResponseWriter, r *http.Request, logID, userID int64) {
	logEntry, err := h.service.GetByID(logID, userID)
	if err != nil {
		log.Printf("Error getting log: %v", err)
		if err == service.ErrLogNotFound {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		}
		if err == service.ErrUnauthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

// updateLog updates a log's information
func (h *LogHandler) updateLog(w http.ResponseWriter, r *http.Request, logID, userID int64) {
	var req models.LogUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.Update(logID, userID, &req); err != nil {
		log.Printf("Error updating log: %v", err)
		if err == service.ErrLogNotFound {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		}
		if err == service.ErrUnauthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// deleteLog deletes a log
func (h *LogHandler) deleteLog(w http.ResponseWriter, r *http.Request, logID, userID int64) {
	if err := h.service.Delete(logID, userID); err != nil {
		log.Printf("Error deleting log: %v", err)
		if err == service.ErrLogNotFound {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		}
		if err == service.ErrUnauthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
