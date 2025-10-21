package service

import (
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(user *models.User) error
	GetUser(id int64) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser creates a new user
func (s *userService) CreateUser(user *models.User) error {
	// Add any business logic here, e.g., validation
	return s.userRepo.CreateUser(user)
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(id int64) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// UpdateUser updates a user
func (s *userService) UpdateUser(user *models.User) error {
	// Add any business logic here
	return s.userRepo.UpdateUser(user)
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id int64) error {
	return s.userRepo.DeleteUser(id)
}
