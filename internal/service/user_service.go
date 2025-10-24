package service

import (
	"context"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	UpdateSelf(ctx context.Context, u *models.User) error
	DeleteSelf(ctx context.Context, id int64) error
}

type userService struct{ users repository.UserRepository }

func NewUserService(users repository.UserRepository) UserService { return &userService{users: users} }

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) { return s.users.GetByID(ctx, id) }
func (s *userService) UpdateSelf(ctx context.Context, u *models.User) error { return s.users.UpdateSelf(ctx, u) }
func (s *userService) DeleteSelf(ctx context.Context, id int64) error { return s.users.DeleteSelf(ctx, id) }
