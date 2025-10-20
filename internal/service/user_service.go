package service

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user business logic operations
type UserService interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user's information
func (s *userService) UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	// Retrieve existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Apply updates only for fields that are provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.ProfileImageURL != nil {
		user.ProfileImageURL = *req.ProfileImageURL
	}
	if req.TimeZone != nil {
		user.TimeZone = *req.TimeZone
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	// Save updates to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser removes a user from the system
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// hashPassword generates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// comparePasswords compares a plain text password with a hashed password
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
