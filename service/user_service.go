package service

import (
	"fmt"
	"time"

	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

// UserService defines the interface for user business logic
type UserService interface {
	Register(req *models.CreateUserRequest) (*models.User, error)
	RegisterAndLogin(req *models.CreateUserRequest) (*models.LoginResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetUser(id int64) (*models.User, error)
	UpdateUser(id int64, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id int64) error
}

// userService implements UserService
type userService struct {
	repo       repository.UserRepository
	hasher     *utils.PasswordHasher
	jwtManager *utils.JWTManager
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, jwtManager *utils.JWTManager) UserService {
	return &userService{
		repo:       repo,
		hasher:     utils.NewPasswordHasher(),
		jwtManager: jwtManager,
	}
}

// Register registers a new user
func (s *userService) Register(req *models.CreateUserRequest) (*models.User, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if existingUser != nil {
		return nil, fmt.Errorf("user with phone number %s already exists", req.PhoneNumber)
	}

	// Hash the password
	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user, err := s.repo.Create(req, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// RegisterAndLogin registers a new user and returns a login response with JWT token
func (s *userService) RegisterAndLogin(req *models.CreateUserRequest) (*models.LoginResponse, error) {
	// Register the user
	user, err := s.Register(req)
	if err != nil {
		return nil, err
	}

	// Generate a JWT token (24 hours expiration)
	token, err := s.jwtManager.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the user by phone number
	user, err := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify the password
	valid, err := s.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate a JWT token (24 hours expiration)
	token, err := s.jwtManager.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(id int64) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Hash the password if provided
	var passwordHash *string
	if req.Password != nil {
		hash, err := s.hasher.Hash(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		passwordHash = &hash
	}

	// Update the user
	user, err := s.repo.Update(id, req, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
