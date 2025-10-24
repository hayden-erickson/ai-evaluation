package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, req *models.LoginRequest) (*models.User, error)
}

// userService implements UserService interface
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user with validation
func (s *userService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.TimeZone == "" {
		return nil, fmt.Errorf("time zone is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		ID:              uuid.New().String(),
		ProfileImageURL: req.ProfileImageURL,
		Name:            req.Name,
		TimeZone:        req.TimeZone,
		PhoneNumber:     req.PhoneNumber,
		PasswordHash:    string(hashedPassword),
		CreatedAt:       time.Now().UTC(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user with validation
func (s *userService) Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		user.Name = *req.Name
	}
	if req.TimeZone != nil {
		if *req.TimeZone == "" {
			return nil, fmt.Errorf("time zone cannot be empty")
		}
		user.TimeZone = *req.TimeZone
	}
	if req.ProfileImageURL != nil {
		user.ProfileImageURL = *req.ProfileImageURL
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	// Save updates
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("user ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Login authenticates a user and returns their information
func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	// Validate credentials
	if req.PhoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Get user by phone number
	user, err := s.repo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

