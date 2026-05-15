package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/pkg/pgpool"
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
	getAccountByIDQuery = `
		SELECT id, user_id, name, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
		LIMIT 1`
	getAccountByIDForUpdateQuery = `
		SELECT id, user_id, name, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE`
	getAccountsByUserIDQuery = `
		SELECT id, user_id, name, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1`
	updateAccountBalanceQuery = `
		UPDATE accounts
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2`
)

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	row := pgpool.GetRunner(ctx, r.pool).QueryRow(
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

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	row := pgpool.GetRunner(ctx, r.pool).QueryRow(ctx, getAccountByIDQuery, id)
	var account domain.Account
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("accountRepository.GetByID.Scan: %w", err)
	}
	return &account, nil
}

func (r *accountRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	row := pgpool.GetRunner(ctx, r.pool).QueryRow(ctx, getAccountByIDForUpdateQuery, id)
	var account domain.Account
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("accountRepository.GetByID.Scan: %w", err)
	}
	return &account, nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	rows, err := pgpool.GetRunner(ctx, r.pool).Query(ctx, getAccountsByUserIDQuery, userID)
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

func (r *accountRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	_, err := pgpool.GetRunner(ctx, r.pool).Exec(ctx, updateAccountBalanceQuery, amount, id)
	if err != nil {
		return fmt.Errorf("accountRepository.UpdateBalance.Exec: %w", err)
	}
	return nil
}
