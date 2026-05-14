package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type JWTManager interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User  domain.User
	Token string
}

type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) error
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type CreateAccountInput struct {
	UserID   uuid.UUID
	Name     string
	Currency string
}

type AccountUsecase interface {
	Create(ctx context.Context, input CreateAccountInput) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
}
