package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// HabitHandler handles habit-related HTTP requests
type HabitHandler struct {
	service service.HabitService
}

// NewHabitHandler creates a new habit handler instance
func NewHabitHandler(service service.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

// HandleHabits routes habit requests based on HTTP method
func (h *HabitHandler) HandleHabits(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("Error: User ID not found in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.createHabit(w, r, userID)
	case http.MethodGet:
		// Check if this is a list or get by ID
		path := strings.TrimPrefix(r.URL.Path, "/api/habits")
		if path == "" || path == "/" {
			h.listHabits(w, r, userID)
		} else {
			h.getHabit(w, r, userID)
		}
	case http.MethodPut:
		h.updateHabit(w, r, userID)
	case http.MethodDelete:
		h.deleteHabit(w, r, userID)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// createHabit handles POST /api/habits
func (h *HabitHandler) createHabit(w http.ResponseWriter, r *http.Request, userID string) {
	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding create habit request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	habit, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Error creating habit: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, habit)
}

// getHabit handles GET /api/habits/{id}
func (h *HabitHandler) getHabit(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	habit, err := h.service.GetByID(r.Context(), userID, id)
	if err != nil {
		log.Printf("Error getting habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, habit)
}

// updateHabit handles PUT /api/habits/{id}
func (h *HabitHandler) updateHabit(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	var req models.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding update habit request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	habit, err := h.service.Update(r.Context(), userID, id, &req)
	if err != nil {
		log.Printf("Error updating habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, habit)
}

// deleteHabit handles DELETE /api/habits/{id}
func (h *HabitHandler) deleteHabit(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	if err := h.service.Delete(r.Context(), userID, id); err != nil {
		log.Printf("Error deleting habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listHabits handles GET /api/habits
func (h *HabitHandler) listHabits(w http.ResponseWriter, r *http.Request, userID string) {
	habits, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("Error listing habits: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, habits)
}
