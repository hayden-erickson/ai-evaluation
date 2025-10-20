package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// HabitHandler handles habit-related HTTP requests
type HabitHandler struct {
	habitService service.HabitService
}

// NewHabitHandler creates a new instance of HabitHandler
func NewHabitHandler(habitService service.HabitService) *HabitHandler {
	return &HabitHandler{
		habitService: habitService,
	}
}

// CreateHabit handles POST requests to create a new habit
func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	var req models.CreateHabitRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode create habit request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Create habit request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create habit
	habit, err := h.habitService.CreateHabit(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Failed to create habit: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create habit")
		return
	}

	respondJSON(w, http.StatusCreated, habit)
}

// GetHabit handles GET requests to retrieve a specific habit
func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract habit ID from path
	habitID := r.PathValue("id")
	if habitID == "" {
		respondError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	// Retrieve habit
	habit, err := h.habitService.GetHabit(r.Context(), habitID, userID)
	if err != nil {
		log.Printf("Failed to get habit: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: habit does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusNotFound, "Habit not found")
		return
	}

	respondJSON(w, http.StatusOK, habit)
}

// ListHabits handles GET requests to retrieve all habits for a user
func (h *HabitHandler) ListHabits(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Retrieve habits
	habits, err := h.habitService.ListHabits(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to list habits: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve habits")
		return
	}

	respondJSON(w, http.StatusOK, habits)
}

// UpdateHabit handles PUT requests to update a habit
func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract habit ID from path
	habitID := r.PathValue("id")
	if habitID == "" {
		respondError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	var req models.UpdateHabitRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode update habit request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Update habit request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update habit
	habit, err := h.habitService.UpdateHabit(r.Context(), habitID, userID, &req)
	if err != nil {
		log.Printf("Failed to update habit: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: habit does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update habit")
		return
	}

	respondJSON(w, http.StatusOK, habit)
}

// DeleteHabit handles DELETE requests to remove a habit
func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value("user_id").(string)

	// Extract habit ID from path
	habitID := r.PathValue("id")
	if habitID == "" {
		respondError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	// Delete habit
	if err := h.habitService.DeleteHabit(r.Context(), habitID, userID); err != nil {
		log.Printf("Failed to delete habit: %v", err)
		// Check if error is due to forbidden access
		if err.Error() == "forbidden: habit does not belong to user" {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete habit")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Habit deleted successfully"})
}
