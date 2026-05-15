package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) error
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type JWTManager interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User  domain.User
	Token string
}

type authUsecase struct {
	repo       domain.UserRepository
	hasher     Hasher
	jwtManager JWTManager
}

func NewAuthUsecase(
	repo domain.UserRepository,
	hasher Hasher,
	jwtManager JWTManager,
) AuthUsecase {
	return &authUsecase{
		repo:       repo,
		hasher:     hasher,
		jwtManager: jwtManager,
	}
}

func (u *authUsecase) Register(ctx context.Context, input RegisterInput) error {
	existUser, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("authUsecase.Register.GetByEmail: %w", err)
	}
	if existUser != nil {
		return domain.ErrEmailExists
	}

	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return fmt.Errorf("authUsecase.Register.Hash: %w", err)
	}

	user := domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	if err = u.repo.Create(ctx, &user); err != nil {
		return fmt.Errorf("authUsecase.Register.Create: %w", err)
	}

	return nil
}

func (u *authUsecase) authenticate(ctx context.Context, input LoginInput) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("authUsecase.authenticate.GetByEmail: %w", err)
	}

	if !u.hasher.Compare(input.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (u *authUsecase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := u.authenticate(ctx, input)
	if err != nil {
		return nil, err
	}
	token, err := u.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("authUsecase.Login.GenerateToken: %w", err)
	}

	return &LoginOutput{
		User:  *user,
		Token: token,
	}, nil
}
