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
	return &HabitHandler{
		service: service,
	}
}

// CreateHabit handles creating a new habit (POST /habits)
func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the habit
	habit, err := h.service.CreateHabit(userID, &req)
	if err != nil {
		log.Printf("Failed to create habit: %v", err)
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the created habit
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// GetHabit handles getting a habit by ID (GET /habits/{id})
func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract habit ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing habit ID in path")
		http.Error(w, "Habit ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Get the habit
	habit, err := h.service.GetHabit(id, userID)
	if err != nil {
		log.Printf("Failed to get habit: %v", err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Habit not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the habit
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
}

// GetUserHabits handles getting all habits for the authenticated user (GET /habits)
func (h *HabitHandler) GetUserHabits(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all habits for the user
	habits, err := h.service.GetUserHabits(userID)
	if err != nil {
		log.Printf("Failed to get habits: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the habits
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

// UpdateHabit handles updating a habit (PUT /habits/{id})
func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT requests
	if r.Method != http.MethodPut {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract habit ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing habit ID in path")
		http.Error(w, "Habit ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req models.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the habit
	habit, err := h.service.UpdateHabit(id, userID, &req)
	if err != nil {
		log.Printf("Failed to update habit: %v", err)
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Habit not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the updated habit
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
}

// DeleteHabit handles deleting a habit (DELETE /habits/{id})
func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE requests
	if r.Method != http.MethodDelete {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract habit ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing habit ID in path")
		http.Error(w, "Habit ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid habit ID: %v", err)
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Delete the habit
	if err := h.service.DeleteHabit(id, userID); err != nil {
		log.Printf("Failed to delete habit: %v", err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Habit not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
