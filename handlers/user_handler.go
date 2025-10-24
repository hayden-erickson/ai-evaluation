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
	service    service.UserService
	jwtManager *middleware.JWTManager
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service service.UserService, jwtManager *middleware.JWTManager) *UserHandler {
	return &UserHandler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// Register handles POST /api/users (user registration)
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse request body
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create user
	user, err := h.service.Create(&req)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Return response with token
	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

// Login handles POST /api/login (user authentication)
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Authenticate user
	user, err := h.service.Authenticate(req.Name, req.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Return response with token
	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// HandleUser handles GET, PUT, DELETE /api/users/{id}
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Printf("Failed to get user ID from context: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Users can only access their own data
	if id != userID {
		log.Printf("User %d attempted to access user %d", userID, id)
		writeJSONError(w, http.StatusForbidden, "forbidden")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getUser handles GET request for a specific user
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}

// updateUser handles PUT request to update a user
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int64) {
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Update(id, &req)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}

// deleteUser handles DELETE request to delete a user
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.service.Delete(id); err != nil {
		log.Printf("Failed to delete user: %v", err)
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, models.SuccessResponse{
		Message: "user deleted successfully",
	})
}

// writeJSONError writes a JSON error response
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

// writeJSONResponse writes a JSON success response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
