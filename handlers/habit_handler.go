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

// NewHabitHandler creates a new habit handler
func NewHabitHandler(service service.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

// Create handles habit creation
func (h *HabitHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	
	var req models.HabitCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	habit, err := h.service.Create(userID, &req)
	if err != nil {
		log.Printf("Error creating habit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// List handles listing all habits for the authenticated user
func (h *HabitHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	habits, err := h.service.GetByUserID(userID)
	if err != nil {
		log.Printf("Error getting habits: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

// HandleHabit handles GET, PUT, DELETE operations for a specific habit
func (h *HabitHandler) HandleHabit(w http.ResponseWriter, r *http.Request) {
	// Get habit ID from URL path
	habitIDStr := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	habitID, err := strconv.ParseInt(habitIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
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
		h.getHabit(w, r, habitID, userID)
	case http.MethodPut:
		h.updateHabit(w, r, habitID, userID)
	case http.MethodDelete:
		h.deleteHabit(w, r, habitID, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getHabit retrieves a habit by ID
func (h *HabitHandler) getHabit(w http.ResponseWriter, r *http.Request, habitID, userID int64) {
	habit, err := h.service.GetByID(habitID, userID)
	if err != nil {
		log.Printf("Error getting habit: %v", err)
		if err == service.ErrHabitNotFound {
			http.Error(w, "Habit not found", http.StatusNotFound)
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
	json.NewEncoder(w).Encode(habit)
}

// updateHabit updates a habit's information
func (h *HabitHandler) updateHabit(w http.ResponseWriter, r *http.Request, habitID, userID int64) {
	var req models.HabitUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.Update(habitID, userID, &req); err != nil {
		log.Printf("Error updating habit: %v", err)
		if err == service.ErrHabitNotFound {
			http.Error(w, "Habit not found", http.StatusNotFound)
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

// deleteHabit deletes a habit
func (h *HabitHandler) deleteHabit(w http.ResponseWriter, r *http.Request, habitID, userID int64) {
	if err := h.service.Delete(habitID, userID); err != nil {
		log.Printf("Error deleting habit: %v", err)
		if err == service.ErrHabitNotFound {
			http.Error(w, "Habit not found", http.StatusNotFound)
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
