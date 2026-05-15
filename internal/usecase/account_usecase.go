package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type AccountUsecase interface {
	Create(ctx context.Context, input CreateAccountInput) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
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

type accountUsecase struct {
	repo domain.AccountRepository
}

func NewAccountUsecase(repo domain.AccountRepository) AccountUsecase {
	return &accountUsecase{repo: repo}
}

func (u *accountUsecase) Create(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
	if input.Name == "" {
		return nil, domain.ErrAccountNameRequired
	}

	account := input.ToDomain()
	if err := u.repo.Create(ctx, &account); err != nil {
		return nil, fmt.Errorf("accountUsecase.Create: %w", err)
	}

	return &account, nil
}

func (u *accountUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	accounts, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("accountUsecase.GetByUserID: %w", err)
	}
	return accounts, nil
}
