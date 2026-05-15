package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type CreateTransactionInput struct {
	AccountID   uuid.UUID
	CategoryID  uuid.UUID
	Amount      float64
	Description string
	Date        time.Time
}

func (input CreateTransactionInput) ToDomain() domain.Transaction {
	return domain.Transaction{
		ID:          uuid.New(),
		AccountID:   input.AccountID,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        input.Date,
	}
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionUsecase interface {
	Create(ctx context.Context, input CreateTransactionInput) error
}

type transactionUsecase struct {
	transRepo   domain.TransactionRepository
	accountRepo domain.AccountRepository
	txManager   TransactionManager
}

func NewTransactionUsecase(
	transRepo domain.TransactionRepository,
	accountRepo domain.AccountRepository,
	txManager TransactionManager,
) TransactionUsecase {
	return &transactionUsecase{
		transRepo:   transRepo,
		accountRepo: accountRepo,
		txManager:   txManager,
	}
}

func (u *transactionUsecase) Create(ctx context.Context, input CreateTransactionInput) error {
	return u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		account, err := u.accountRepo.GetByIDForUpdate(txCtx, input.AccountID)
		if err != nil {
			return fmt.Errorf("transactionUsecase.Create.GetByIDForUpdate: %w", err)
		}

		if account.Balance+input.Amount < 0 {
			return domain.ErrInsufficientFunds
		}

		transaction := input.ToDomain()
		if err := u.transRepo.Create(txCtx, &transaction); err != nil {
			return fmt.Errorf("transactionUsecase.Create.TransRepo: %w", err)
		}

		if err := u.accountRepo.UpdateBalance(txCtx, input.AccountID, input.Amount); err != nil {
			return fmt.Errorf("transactionUsecase.Create.UpdateBalance: %w", err)
		}

		return nil
	})
}
