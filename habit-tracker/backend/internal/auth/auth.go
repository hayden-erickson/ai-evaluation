package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hayden-erickson/habit-tracker/internal/db"
	"github.com/hayden-erickson/habit-tracker/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	config *oauth2.Config
	// Define JWT secret key from environment variable
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// Initialize sets up the OAuth2 config
func Initialize() {
	config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback", // Update this for production
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// Claims represents the JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GoogleUser represents the user info received from Google
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// LoginHandler initiates the Google OAuth flow
func LoginHandler(c *gin.Context) {
	url := config.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// CallbackHandler handles the Google OAuth callback
func CallbackHandler(c *gin.Context) {
	var requestBody struct {
		AccessToken string `json:"accessToken"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+requestBody.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Find or create user
	var user models.User
	if err := db.DB.Where("google_id = ?", googleUser.ID).First(&user).Error; err != nil {
		// Create new user
		user = models.User{
			GoogleID:        googleUser.ID,
			Email:           googleUser.Email,
			Name:            googleUser.Name,
			ProfileImageURL: googleUser.Picture,
			TimeZone:        "UTC", // Default timezone, can be updated later
		}
		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate JWT
	tokenJWT, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// In production, redirect to frontend with token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenJWT,
		"user":  user,
	})
}

// AuthMiddleware validates the JWT token and sets the user in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			return
		}

		tokenString := auth[7:] // Remove "Bearer " prefix
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var user models.User
		if err := db.DB.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// generateToken creates a new JWT token for a user
func generateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GetCurrentUser returns the current authenticated user from context
func GetCurrentUser(c *gin.Context) models.User {
	user, _ := c.Get("user")
	return user.(models.User)
}
