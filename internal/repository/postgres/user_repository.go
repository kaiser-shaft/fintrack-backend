package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

const (
	createUserQuery = `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`
	getUserByEmailQuery = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1`
)

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	row := r.pool.QueryRow(
		ctx,
		createUserQuery,
		user.ID,
		user.Email,
		user.PasswordHash,
	)

	err := row.Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("userRepository.Create.Scan: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, getUserByEmailQuery, email)
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.GetByEmail.Scan: %w", err)
	}
	return &user, nil
}
