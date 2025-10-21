package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists is returned when trying to create a duplicate user
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService defines the interface for user business logic
type UserService interface {
	Register(req *models.UserCreateRequest) (*models.User, error)
	Login(req *models.LoginRequest, jwtSecret string) (*models.LoginResponse, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, req *models.UserUpdateRequest) error
	Delete(id int64) error
}

// userService implements UserService interface
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register creates a new user account
func (s *userService) Register(req *models.UserCreateRequest) (*models.User, error) {
	// Validate required fields
	if err := utils.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}
	if err := utils.ValidateRequired(req.Password, "password"); err != nil {
		return nil, err
	}
	if err := utils.ValidateMinLength(req.Password, 8, "password"); err != nil {
		return nil, err
	}
	if err := utils.ValidateRequired(req.PhoneNumber, "phone_number"); err != nil {
		return nil, err
	}
	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return nil, err
	}
	
	// Check if user already exists
	existing, err := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}
	
	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &models.User{
		ProfileImageURL: utils.SanitizeString(req.ProfileImageURL),
		Name:            utils.SanitizeString(req.Name),
		TimeZone:        utils.SanitizeString(req.TimeZone),
		PhoneNumber:     utils.SanitizeString(req.PhoneNumber),
		PasswordHash:    hashedPassword,
		CreatedAt:       time.Now(),
	}
	
	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(req *models.LoginRequest, jwtSecret string) (*models.LoginResponse, error) {
	// Validate required fields
	if err := utils.ValidateRequired(req.PhoneNumber, "phone_number"); err != nil {
		return nil, err
	}
	if err := utils.ValidateRequired(req.Password, "password"); err != nil {
		return nil, err
	}
	
	// Get user by phone number
	user, err := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	
	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	
	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(id int64) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Update updates a user's information
func (s *userService) Update(id int64, req *models.UserUpdateRequest) error {
	// Verify user exists
	user, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	
	// Validate and sanitize fields
	if req.Name != nil {
		name := utils.SanitizeString(*req.Name)
		req.Name = &name
		if err := utils.ValidateRequired(*req.Name, "name"); err != nil {
			return err
		}
	}
	if req.PhoneNumber != nil {
		phoneNumber := utils.SanitizeString(*req.PhoneNumber)
		req.PhoneNumber = &phoneNumber
		if err := utils.ValidatePhoneNumber(*req.PhoneNumber); err != nil {
			return err
		}
	}
	if req.Password != nil {
		if err := utils.ValidateMinLength(*req.Password, 8, "password"); err != nil {
			return err
		}
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		req.Password = &hashedPassword
	}
	if req.ProfileImageURL != nil {
		profileImageURL := utils.SanitizeString(*req.ProfileImageURL)
		req.ProfileImageURL = &profileImageURL
	}
	if req.TimeZone != nil {
		timeZone := utils.SanitizeString(*req.TimeZone)
		req.TimeZone = &timeZone
	}
	
	return s.repo.Update(id, req)
}

// Delete removes a user
func (s *userService) Delete(id int64) error {
	// Verify user exists
	user, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	
	return s.repo.Delete(id)
}
