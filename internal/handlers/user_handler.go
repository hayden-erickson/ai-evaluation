package handlers

import (
	"encoding/json"
	"net/http"
	"new-api/internal/models"
	"new-api/internal/service"
	"strconv"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUser handles retrieving a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.userService.UpdateUser(&user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "User updated successfully",
	})
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// UserRoutes registers the user routes
func (h *UserHandler) UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateUser(w, r)
		case http.MethodGet:
			h.GetUser(w, r)
		case http.MethodPut:
			h.UpdateUser(w, r)
		case http.MethodDelete:
			h.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
