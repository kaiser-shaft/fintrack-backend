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

type AccountHandler struct {
	usecase  usecase.AccountUsecase
	validate *validator.Validator
	logger   *slog.Logger
}

func NewAccountHandler(
	uc usecase.AccountUsecase,
	validate *validator.Validator,
	logger *slog.Logger,
) *AccountHandler {
	return &AccountHandler{
		usecase:  uc,
		validate: validate,
		logger:   logger,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		render.Error(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, "invalid request body", http.StatusBadRequest, nil)
		return
	}

	if res := h.validate.Validate(req); res.HasError {
		render.Error(w, "validation failed", res.StatusCode(), res.Fields)
		return
	}

	account, err := h.usecase.Create(r.Context(), req.MapToInput(userID))
	if err != nil {
		if errors.Is(err, domain.ErrAccountNameRequired) {
			render.Error(w, err.Error(), http.StatusBadRequest, nil)
			return
		}
		h.logger.Error("failed to create account", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	render.JSON(w, dto.NewAccountResponse(*account), http.StatusCreated)
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		render.Error(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	accounts, err := h.usecase.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get accounts", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	render.JSON(w, dto.NewGetAccountsResponse(accounts), http.StatusOK)
}
