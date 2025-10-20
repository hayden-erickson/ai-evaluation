package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode registration request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Registration request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Register user
	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode login request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validateRequest(&req); err != nil {
		log.Printf("Login request validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Authenticate user
	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	respondJSON(w, http.StatusOK, response)
}
