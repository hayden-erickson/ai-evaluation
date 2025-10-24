package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hayden-erickson/ai-evaluation/internal/config"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, u *models.User) (int64, error)
	Login(ctx context.Context, email, password string) (string, *models.User, error)
}

type authService struct {
	users repository.UserRepository
	cfg   *config.Config
}

func NewAuthService(users repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{users: users, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, u *models.User) (int64, error) {
	if u.Role == "" { u.Role = "user" }
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil { return 0, err }
	u.PasswordHash = string(hash)
	return s.users.Create(ctx, u)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil { return "", nil, err }
	if u == nil { return "", nil, errors.New("invalid credentials") }
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	// Issue JWT
	claims := jwt.MapClaims{
		"sub": u.ID,
		"email": u.Email,
		"role": u.Role,
		"iss": s.cfg.JWTIssuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil { return "", nil, err }
	// Clear sensitive fields before returning user
	u.PasswordHash = ""
	u.Password = ""
	// constant time no-op to keep timing consistent
	_ = subtle.ConstantTimeEq(1, 1)
	return signed, u, nil
}
