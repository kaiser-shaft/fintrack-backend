package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
)

type CreateTransactionRequest struct {
	AccountID   uuid.UUID `json:"account_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

func (req CreateTransactionRequest) ToInput() usecase.CreateTransactionInput {
	return usecase.CreateTransactionInput{
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	}
}
