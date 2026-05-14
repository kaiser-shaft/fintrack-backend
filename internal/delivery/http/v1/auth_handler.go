package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
	"github.com/kaiser-shaft/fintrack-backend/pkg/render"
	"github.com/kaiser-shaft/fintrack-backend/pkg/validator"
)

type AuthHandler struct {
	usecase       usecase.AuthUsecase
	validate      *validator.Validator
	logger        *slog.Logger
	cookieSecure  bool
	tokenDuration time.Duration
}

func NewAuthHandler(
	usecase usecase.AuthUsecase,
	validate *validator.Validator,
	logger *slog.Logger,
	cookieSecure bool,
	tokenDuration time.Duration,
) *AuthHandler {
	return &AuthHandler{
		usecase:       usecase,
		validate:      validate,
		logger:        logger,
		cookieSecure:  cookieSecure,
		tokenDuration: tokenDuration,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, "invalid request body", http.StatusBadRequest, nil)
		return
	}

	if res := h.validate.Validate(req); res.HasError {
		render.Error(w, "validation failed", res.StatusCode(), res.Fields)
		return
	}

	err := h.usecase.Register(r.Context(), req.MapToInput())
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			render.Error(w, err.Error(), http.StatusConflict, nil)
			return
		}
		h.logger.Error("failed to register user", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, "invalid request body", http.StatusBadRequest, nil)
		return
	}

	if res := h.validate.Validate(req); res.HasError {
		render.Error(w, "validation failed", res.StatusCode(), res.Fields)
		return
	}

	out, err := h.usecase.Login(r.Context(), req.MapToInput())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			render.Error(w, err.Error(), http.StatusUnauthorized, nil)
			return
		}
		h.logger.Error("failed to login user", slog.Any("error", err))
		render.Error(w, "internal server error", http.StatusInternalServerError, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    out.Token,
		Path:     "/",
		Secure:   h.cookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.tokenDuration.Seconds()),
	})

	render.JSON(w, NewAuthResponse(out.User), http.StatusOK)
}
