package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAccountNameRequired = errors.New("account name is required")
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error)
}

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Balance   float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
