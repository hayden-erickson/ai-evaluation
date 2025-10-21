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

// HabitHandler handles habit-related HTTP requests
type HabitHandler struct {
	service service.HabitService
}

// NewHabitHandler creates a new HabitHandler instance
func NewHabitHandler(service service.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

// HandleHabits handles GET and POST /api/habits
func (h *HabitHandler) HandleHabits(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.listHabits(w, r, userID)
	case http.MethodPost:
		h.createHabit(w, r, userID)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// listHabits handles GET request to list all habits for a user
func (h *HabitHandler) listHabits(w http.ResponseWriter, r *http.Request, userID int64) {
	habits, err := h.service.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get habits: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to get habits")
		return
	}

	// Return empty array if no habits found
	if habits == nil {
		habits = []*models.Habit{}
	}

	writeJSONResponse(w, http.StatusOK, habits)
}

// createHabit handles POST request to create a new habit
func (h *HabitHandler) createHabit(w http.ResponseWriter, r *http.Request, userID int64) {
	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	habit, err := h.service.Create(userID, &req)
	if err != nil {
		log.Printf("Failed to create habit: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, habit)
}

// HandleHabit handles GET, PUT, DELETE /api/habits/{id}
func (h *HabitHandler) HandleHabit(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid habit ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getHabit(w, r, id, userID)
	case http.MethodPut:
		h.updateHabit(w, r, id, userID)
	case http.MethodDelete:
		h.deleteHabit(w, r, id, userID)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getHabit handles GET request for a specific habit
func (h *HabitHandler) getHabit(w http.ResponseWriter, r *http.Request, id, userID int64) {
	habit, err := h.service.GetByID(id, userID)
	if err != nil {
		log.Printf("Failed to get habit: %v", err)
		writeJSONError(w, http.StatusNotFound, "habit not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, habit)
}

// updateHabit handles PUT request to update a habit
func (h *HabitHandler) updateHabit(w http.ResponseWriter, r *http.Request, id, userID int64) {
	var req models.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	habit, err := h.service.Update(id, userID, &req)
	if err != nil {
		log.Printf("Failed to update habit: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, habit)
}

// deleteHabit handles DELETE request to delete a habit
func (h *HabitHandler) deleteHabit(w http.ResponseWriter, r *http.Request, id, userID int64) {
	if err := h.service.Delete(id, userID); err != nil {
		log.Printf("Failed to delete habit: %v", err)
		writeJSONError(w, http.StatusNotFound, "habit not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, models.SuccessResponse{
		Message: "habit deleted successfully",
	})
}
