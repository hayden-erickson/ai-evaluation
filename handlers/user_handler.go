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
	service   service.UserService
	jwtSecret string
}

// NewUserHandler creates a new user handler
func NewUserHandler(service service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	user, err := h.service.Register(&req)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		if err == service.ErrUserAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user authentication
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response, err := h.service.Login(&req, h.jwtSecret)
	if err != nil {
		log.Printf("Error logging in: %v", err)
		if err == service.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleUser handles GET, PUT, DELETE operations for users
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL path
	userIDStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	// Get authenticated user ID from context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Users can only access their own data
	if authUserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, userID)
	case http.MethodPut:
		h.updateUser(w, r, userID)
	case http.MethodDelete:
		h.deleteUser(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getUser retrieves a user by ID
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, userID int64) {
	user, err := h.service.GetByID(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		if err == service.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUser updates a user's information
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, userID int64) {
	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.Update(userID, &req); err != nil {
		log.Printf("Error updating user: %v", err)
		if err == service.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// deleteUser deletes a user
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, userID int64) {
	if err := h.service.Delete(userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		if err == service.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
