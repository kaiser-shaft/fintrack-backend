package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/middleware"
	"github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/v1/dto"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
	"github.com/kaiser-shaft/fintrack-backend/pkg/render"
	"github.com/kaiser-shaft/fintrack-backend/pkg/validator"
)

type CategoryHandler struct {
	usecase  usecase.CategoryUsecase
	validate *validator.Validator
	logger   *slog.Logger
}

func NewCategoryHandler(
	uc usecase.CategoryUsecase,
	validate *validator.Validator,
	logger *slog.Logger,
) *CategoryHandler {
	return &CategoryHandler{
		usecase:  uc,
		validate: validate,
		logger:   logger,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		render.Error(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, "invalid request body", http.StatusBadRequest, nil)
		return
	}

	if res := h.validate.Validate(req); res.HasError {
		render.Error(w, "validation failed", res.StatusCode(), res.Fields)
		return
	}

	category, err := h.usecase.Create(r.Context(), req.ToInput(userID))
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNameRequired) || errors.Is(err, domain.ErrInvalidCategoryType) {
			render.Error(w, err.Error(), http.StatusBadRequest, nil)
			return
		}
		h.logger.Error("failed to create category", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	render.JSON(w, category, http.StatusCreated)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		render.Error(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	categories, err := h.usecase.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get categories", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	render.JSON(w, dto.NewCategoryListResponse(categories), http.StatusOK)
}
