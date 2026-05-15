package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountNameRequired = errors.New("account name is required")
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Balance   float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error
}
