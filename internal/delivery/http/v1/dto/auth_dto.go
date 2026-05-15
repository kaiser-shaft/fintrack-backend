package dto

import (
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

func (req RegisterRequest) ToInput() usecase.RegisterInput {
	return usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req LoginRequest) ToInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewAuthResponse(user domain.User) AuthResponse {
	return AuthResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	}
}
