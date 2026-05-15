package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
)

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"max=100"`
	Type string `json:"type" validate:"max=20"`
}

func (req CreateCategoryRequest) ToInput(userID uuid.UUID) usecase.CreateCategoryInput {
	return usecase.CreateCategoryInput{
		UserID: userID,
		Name:   req.Name,
		Type:   domain.CategoryType(req.Type),
	}
}

type CategoryResponse struct {
	ID        uuid.UUID           `json:"id"`
	UserID    uuid.UUID           `json:"user_id"`
	Name      string              `json:"name"`
	Type      domain.CategoryType `json:"type"`
	CreatedAt time.Time           `json:"created_at"`
}

func NewCategoryResponse(category domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		Type:      category.Type,
		CreatedAt: category.CreatedAt,
	}
}

func NewCategoryListResponse(categories []domain.Category) ListResponse[CategoryResponse] {
	return ListResponse[CategoryResponse]{
		Data: MapSlice(categories, NewCategoryResponse),
	}
}
