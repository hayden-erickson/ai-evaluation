package handlers

import (
	"encoding/json"
	"net/http"
	"new-api/internal/auth"
	"new-api/internal/models"
	"new-api/internal/service"
	"strconv"
)

// HabitHandler handles habit-related requests
type HabitHandler struct {
	habitService service.HabitService
}

// NewHabitHandler creates a new HabitHandler
func NewHabitHandler(habitService service.HabitService) *HabitHandler {
	return &HabitHandler{habitService: habitService}
}

// CreateHabit handles the creation of a new habit
func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	habit.UserID = userID

	if err := h.habitService.CreateHabit(&habit); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create habit")
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "Habit created successfully",
		Data:    habit,
	})
}

// GetHabit handles retrieving a habit by ID
func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	habit, err := h.habitService.GetHabit(id, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Habit not found")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    habit,
	})
}

// GetUserHabits handles retrieving all habits for a user
func (h *HabitHandler) GetUserHabits(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	habits, err := h.habitService.GetUserHabits(userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve habits")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    habits,
	})
}

// UpdateHabit handles updating a habit
func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.habitService.UpdateHabit(&habit, userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update habit")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Habit updated successfully",
	})
}

// DeleteHabit handles deleting a habit
func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.habitService.DeleteHabit(id, userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete habit")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Habit deleted successfully",
	})
}

// HabitRoutes registers the habit routes
func (h *HabitHandler) HabitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/habits", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateHabit(w, r)
		case http.MethodGet:
			// Differentiate between getting a single habit and all user habits
			if r.URL.Query().Get("id") != "" {
				h.GetHabit(w, r)
			} else {
				h.GetUserHabits(w, r)
			}
		case http.MethodPut:
			h.UpdateHabit(w, r)
		case http.MethodDelete:
			h.DeleteHabit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
