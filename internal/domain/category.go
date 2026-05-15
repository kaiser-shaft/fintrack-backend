package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrInvalidCategoryType  = errors.New("invalid category type")
)

type CategoryType string

const (
	IncomeCategoryType  CategoryType = "income"
	ExpenseCategoryType CategoryType = "expense"
)

func (t CategoryType) IsValid() bool {
	switch t {
	case IncomeCategoryType, ExpenseCategoryType:
		return true
	default:
		return false
	}
}

type Category struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      CategoryType
	CreatedAt time.Time
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error)
}
