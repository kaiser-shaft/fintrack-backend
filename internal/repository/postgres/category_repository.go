package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) domain.CategoryRepository {
	return &categoryRepository{pool: pool}
}

const (
	createCategoryQuery = `
		INSERT INTO categories (id, user_id, name, type)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	getCategoriesByUserIDQuery = `
		SELECT id, user_id, name, type, created_at
		FROM categories
		WHERE user_id = $1`
)

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	row := r.pool.QueryRow(
		ctx,
		createCategoryQuery,
		category.ID,
		category.UserID,
		category.Name,
		category.Type,
	)
	err := row.Scan(&category.CreatedAt)
	if err != nil {
		return fmt.Errorf("categoryRepository.Create.Scan: %w", err)
	}
	return nil
}
func (r *categoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	rows, err := r.pool.Query(ctx, getCategoriesByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("categoryRepository.GetByUserID.Query: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("categoryRepository.GetByUserID.Scan: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("categoryRepository.GetByUserID.Err: %w", err)
	}

	return categories, nil
}
