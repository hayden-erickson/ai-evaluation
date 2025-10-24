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

// HabitHandler handles habit-related HTTP requests
type HabitHandler struct {
	habitService service.HabitService
}

// NewHabitHandler creates a new habit handler
func NewHabitHandler(habitService service.HabitService) *HabitHandler {
	return &HabitHandler{habitService: habitService}
}

// CreateHabit handles POST /habits - creates a new habit
func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Method not allowed: %s on /habits", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request body
	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create habit
	habit, err := h.habitService.Create(r.Context(), userID, &req)
	if err != nil {
		log.Printf("ERROR: Failed to create habit: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// GetHabits handles GET /habits - retrieves all habits for the authenticated user
func (h *HabitHandler) GetHabits(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		log.Printf("ERROR: Method not allowed: %s on /habits", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get all habits for user
	habits, err := h.habitService.GetByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to get habits: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to get habits")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habits)
}

// GetHabit handles GET /habits/{id} - retrieves a specific habit
func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		log.Printf("ERROR: Method not allowed: %s on /habits/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get habit ID from URL path
	habitID := strings.TrimPrefix(r.URL.Path, "/habits/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get habit
	habit, err := h.habitService.GetByID(r.Context(), habitID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, "habit not found")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habit)
}

// UpdateHabit handles PUT /habits/{id} - updates a habit
func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT method
	if r.Method != http.MethodPut {
		log.Printf("ERROR: Method not allowed: %s on /habits/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get habit ID from URL path
	habitID := strings.TrimPrefix(r.URL.Path, "/habits/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request body
	var req models.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update habit
	habit, err := h.habitService.Update(r.Context(), habitID, userID, &req)
	if err != nil {
		log.Printf("ERROR: Failed to update habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habit)
}

// DeleteHabit handles DELETE /habits/{id} - deletes a habit
func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	// Only accept DELETE method
	if r.Method != http.MethodDelete {
		log.Printf("ERROR: Method not allowed: %s on /habits/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get habit ID from URL path
	habitID := strings.TrimPrefix(r.URL.Path, "/habits/")

	// Get authenticated user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Delete habit
	if err := h.habitService.Delete(r.Context(), habitID, userID); err != nil {
		log.Printf("ERROR: Failed to delete habit: %v", err)
		if strings.Contains(err.Error(), "unauthorized") {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to delete habit")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
