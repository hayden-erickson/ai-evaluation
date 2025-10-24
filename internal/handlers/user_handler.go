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

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser handles POST /users - creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Method not allowed: %s on /users", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse request body
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create user
	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: Failed to create user: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Failed to generate token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Return response
	resp := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Login handles POST /login - authenticates a user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Method not allowed: %s on /login", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode login request: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Authenticate user
	user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: Login failed: %v", err)
		respondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Failed to generate token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Return response
	resp := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetUser handles GET /users/{id} - retrieves a user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		log.Printf("ERROR: Method not allowed: %s on /users/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")

	// Get authenticated user ID from context
	authUserID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Ensure user can only access their own data
	if userID != authUserID {
		log.Printf("ERROR: User %s attempted to access user %s data", authUserID, userID)
		respondWithError(w, http.StatusForbidden, "forbidden: cannot access other users' data")
		return
	}

	// Get user
	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to get user: %v", err)
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT /users/{id} - updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT method
	if r.Method != http.MethodPut {
		log.Printf("ERROR: Method not allowed: %s on /users/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")

	// Get authenticated user ID from context
	authUserID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Ensure user can only update their own data
	if userID != authUserID {
		log.Printf("ERROR: User %s attempted to update user %s data", authUserID, userID)
		respondWithError(w, http.StatusForbidden, "forbidden: cannot update other users' data")
		return
	}

	// Parse request body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update user
	user, err := h.userService.Update(r.Context(), userID, &req)
	if err != nil {
		log.Printf("ERROR: Failed to update user: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /users/{id} - deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Only accept DELETE method
	if r.Method != http.MethodDelete {
		log.Printf("ERROR: Method not allowed: %s on /users/{id}", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")

	// Get authenticated user ID from context
	authUserID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Ensure user can only delete their own account
	if userID != authUserID {
		log.Printf("ERROR: User %s attempted to delete user %s account", authUserID, userID)
		respondWithError(w, http.StatusForbidden, "forbidden: cannot delete other users' accounts")
		return
	}

	// Delete user
	if err := h.userService.Delete(r.Context(), userID); err != nil {
		log.Printf("ERROR: Failed to delete user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errResp := models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	}
	json.NewEncoder(w).Encode(errResp)
}
