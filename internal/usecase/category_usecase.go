package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type CreateCategoryInput struct {
	UserID uuid.UUID
	Name   string
	Type   domain.CategoryType
}

func (input CreateCategoryInput) ToDomain() domain.Category {
	return domain.Category{
		ID:     uuid.New(),
		UserID: input.UserID,
		Name:   input.Name,
		Type:   input.Type,
	}
}

type CategoryUsecase interface {
	Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error)
	List(ctx context.Context, userID uuid.UUID) ([]domain.Category, error)
}

type categoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(
	repo domain.CategoryRepository,
) CategoryUsecase {
	return &categoryUsecase{
		repo: repo,
	}
}

func (u *categoryUsecase) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	if input.Name == "" {
		return nil, domain.ErrCategoryNameRequired
	}
	if !input.Type.IsValid() {
		return nil, domain.ErrInvalidCategoryType
	}

	category := input.ToDomain()
	if err := u.repo.Create(ctx, &category); err != nil {
		return nil, fmt.Errorf("categoryUsecase.Create: %w", err)
	}

	return &category, nil
}
func (u *categoryUsecase) List(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	categories, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("categoryUsecase.GetByUserID: %w", err)
	}
	return categories, nil
}
