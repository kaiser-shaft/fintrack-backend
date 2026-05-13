package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error)
}
