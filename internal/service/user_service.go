package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"github.com/hayden-erickson/ai-evaluation/pkg/auth"
	"github.com/hayden-erickson/ai-evaluation/pkg/validator"
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
}

// userService implements UserService interface
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user with validation and password hashing
func (s *userService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Validate input
	if err := validator.ValidateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Generate user ID
	id := auth.GenerateID()

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		ID:              id,
		ProfileImageURL: req.ProfileImageURL,
		Name:            req.Name,
		TimeZone:        req.TimeZone,
		PhoneNumber:     req.PhoneNumber,
		PasswordHash:    hashedPassword,
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
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.ProfileImageURL != nil {
		user.ProfileImageURL = *req.ProfileImageURL
	}
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
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Password != nil {
		// Hash new password
		hashedPassword, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	}

	// Save to repository
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user by ID
func (s *userService) Delete(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("user ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves all users
func (s *userService) List(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if req.ID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Get user
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
