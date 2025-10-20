package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles GET requests to retrieve a user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path
	userID := r.PathValue("id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Get authenticated user ID from context
	authUserID := r.Context().Value("user_id").(string)

	// Verify user can only access their own data
	if userID != authUserID {
		log.Printf("Unauthorized access attempt: user %s tried to access user %s", authUserID, userID)
		respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Retrieve user
	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// UpdateUser handles PUT requests to update a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path
	userID := r.PathValue("id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Get authenticated user ID from context
	authUserID := r.Context().Value("user_id").(string)

	// Verify user can only update their own data
	if userID != authUserID {
		log.Printf("Unauthorized update attempt: user %s tried to update user %s", authUserID, userID)
		respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.UpdateUserRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode update user request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Update user request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path
	userID := r.PathValue("id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Get authenticated user ID from context
	authUserID := r.Context().Value("user_id").(string)

	// Verify user can only delete their own account
	if userID != authUserID {
		log.Printf("Unauthorized delete attempt: user %s tried to delete user %s", authUserID, userID)
		respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		log.Printf("Failed to delete user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
