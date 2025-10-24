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

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register handles user registration (POST /users/register)
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the user
	loginResp, err := h.service.Register(&req)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		// Check if it's a validation error or duplicate user
		if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the login response (with token and user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loginResp)
}

// Login handles user login (POST /users/login)
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Login the user
	resp, err := h.service.Login(&req)
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return the login response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUser handles getting user details (GET /users/{id})
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing user ID in path")
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Users can only view their own profile
	if authUserID != id {
		log.Printf("User %d attempted to access user %d's profile", authUserID, id)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get the user
	user, err := h.service.GetUser(id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles updating user details (PUT /users/{id})
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT requests
	if r.Method != http.MethodPut {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing user ID in path")
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Users can only update their own profile
	if authUserID != id {
		log.Printf("User %d attempted to update user %d's profile", authUserID, id)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse the request body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the user
	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles deleting a user (DELETE /users/{id})
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE requests
	if r.Method != http.MethodDelete {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		log.Println("Missing user ID in path")
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Users can only delete their own profile
	if authUserID != id {
		log.Printf("User %d attempted to delete user %d's profile", authUserID, id)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete the user
	if err := h.service.DeleteUser(id); err != nil {
		log.Printf("Failed to delete user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
