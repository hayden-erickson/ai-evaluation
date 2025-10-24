package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// HandleUsers routes user requests based on HTTP method
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodGet:
		// Check if this is a list or get by ID
		path := strings.TrimPrefix(r.URL.Path, "/api/users")
		if path == "" || path == "/" {
			h.listUsers(w, r)
		} else {
			h.getUser(w, r)
		}
	case http.MethodPut:
		h.updateUser(w, r)
	case http.MethodDelete:
		h.deleteUser(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// createUser handles POST /api/users
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding create user request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Create(r.Context(), &req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// getUser handles GET /api/users/{id}
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// updateUser handles PUT /api/users/{id}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding update user request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// deleteUser handles DELETE /api/users/{id}
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		log.Printf("Error deleting user: %v", err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listUsers handles GET /api/users
func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		log.Printf("Error listing users: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// HandleLogin handles POST /api/login
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding login request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.Login(r.Context(), &req)
	if err != nil {
		log.Printf("Error during login: %v", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
