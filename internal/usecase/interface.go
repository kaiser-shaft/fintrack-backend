package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type LoginResult struct {
	User  domain.User
	Token string
}

type AuthUsecase interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*LoginResult, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type JWTManager interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}
