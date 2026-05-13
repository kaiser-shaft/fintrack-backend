package domain

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	IncomeCategoryType  CategoryType = "income"
	ExpenseCategoryType CategoryType = "expense"
)

type Category struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      CategoryType
	CreatedAt time.Time
}
