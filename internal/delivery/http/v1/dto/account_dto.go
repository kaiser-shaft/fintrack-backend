package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
)

type CreateAccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

func (r CreateAccountRequest) MapToInput(userID uuid.UUID) usecase.CreateAccountInput {
	return usecase.CreateAccountInput{
		UserID:   userID,
		Name:     r.Name,
		Currency: r.Currency,
	}
}

type AccountResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAccountResponse(account domain.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Name:      account.Name,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

type GetAccountsResponse struct {
	Data []AccountResponse `json:"data"`
}

func NewGetAccountsResponse(accounts []domain.Account) GetAccountsResponse {
	data := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		data[i] = NewAccountResponse(account)
	}
	return GetAccountsResponse{Data: data}
}
