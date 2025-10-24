package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/auth"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService service.UserService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token for the new user
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}

// Login handles user login (simplified - just generates token based on phone_number)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Since there's no password in the schema, we'll use phone_number to look up the user
	user, err := h.userService.GetUserByPhoneNumber(loginReq.PhoneNumber)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}

// AuthRoutes registers the authentication routes
func (h *AuthHandler) AuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Register(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Login(w, r)
	})
}
