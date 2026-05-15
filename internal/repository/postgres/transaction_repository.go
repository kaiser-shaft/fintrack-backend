package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/pkg/pgpool"
)

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) domain.TransactionRepository {
	return &transactionRepository{pool: pool}
}

const (
	createTransactionQuery = `
		INSERT INTO transactions (id, account_id, category_id, amount, description, transaction_date
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`
)

func (r *transactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	row := pgpool.GetRunner(ctx, r.pool).QueryRow(
		ctx,
		createTransactionQuery,
		transaction.ID,
		transaction.AccountID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Description,
		transaction.Date,
	)

	if err := row.Scan(&transaction.CreatedAt); err != nil {
		return fmt.Errorf("transactionRepository.Create.Scan: %w", err)
	}

	return nil
}
