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

type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) error
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type AccountUsecase interface {
	Create(ctx context.Context, input CreateAccountInput) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
}

type CategoryUsecase interface {
	Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Category, error)
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

type CreateAccountInput struct {
	UserID   uuid.UUID
	Name     string
	Currency string
}

func (i CreateAccountInput) ToDomain() domain.Account {
	return domain.Account{
		ID:       uuid.New(),
		UserID:   i.UserID,
		Name:     i.Name,
		Balance:  0,
		Currency: i.Currency,
	}
}

type CreateCategoryInput struct {
	UserID uuid.UUID
	Name   string
	Type   domain.CategoryType
}

func (i CreateCategoryInput) ToDomain() domain.Category {
	return domain.Category{
		ID:     uuid.New(),
		UserID: i.UserID,
		Name:   i.Name,
		Type:   i.Type,
	}
}
