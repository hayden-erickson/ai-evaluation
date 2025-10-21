package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(req *models.CreateUserRequest) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id int64) error
	Authenticate(name, password string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user with validation
func (s *userService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ProfileImageURL: req.ProfileImageURL,
		Name:            req.Name,
		TimeZone:        req.TimeZone,
		PhoneNumber:     req.PhoneNumber,
		PasswordHash:    string(hashedPassword),
		CreatedAt:       time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

// Update updates a user with validation
func (s *userService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.ProfileImageURL != nil {
		user.ProfileImageURL = *req.ProfileImageURL
	}
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		user.Name = *req.Name
	}
	if req.TimeZone != nil {
		if strings.TrimSpace(*req.TimeZone) == "" {
			return nil, fmt.Errorf("time zone cannot be empty")
		}
		user.TimeZone = *req.TimeZone
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Password != nil {
		if len(*req.Password) < 6 {
			return nil, fmt.Errorf("password must be at least 6 characters")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user by ID
func (s *userService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Authenticate verifies user credentials and returns the user if valid
func (s *userService) Authenticate(name, password string) (*models.User, error) {
	user, err := s.repo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// validateCreateRequest validates the create user request
func (s *userService) validateCreateRequest(req *models.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.TimeZone) == "" {
		return fmt.Errorf("time zone is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}
