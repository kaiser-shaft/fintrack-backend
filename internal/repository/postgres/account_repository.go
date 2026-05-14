package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type accountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) domain.AccountRepository {
	return &accountRepository{pool: pool}
}

const (
	createAccountQuery = `
		INSERT INTO accounts (id, user_id, name, balance, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`
	getAccountsByUserIDQuery = `
		SELECT id, user_id, name, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1`
)

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	row := r.pool.QueryRow(
		ctx,
		createAccountQuery,
		account.ID,
		account.UserID,
		account.Name,
		account.Balance,
		account.Currency,
	)

	err := row.Scan(&account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return fmt.Errorf("accountRepository.Create.Scan: %w", err)
	}

	return nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx, getAccountsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("accountRepository.GetByUserID.Query: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("accountRepository.GetByUserID.Scan: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("accountRepository.GetByUserID.Err: %w", err)
	}

	return accounts, nil
}
