package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	// FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]Transaction, error)
	// FindByAccountIDAndPeriod(ctx context.Context, accountID uuid.UUID, from, to time.Time) ([]Transaction, error)
	// CountByAccountID(ctx context.Context, accountID uuid.UUID) (int, error)
}

type Transaction struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	CategoryID  uuid.UUID
	Amount      float64
	Description string
	Date        time.Time
	CreatedAt   time.Time
}
