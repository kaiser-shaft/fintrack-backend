package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
)

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

func (u *authUsecase) Register(ctx context.Context, email, password string) error {
	existUser, err := u.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("authUsecase.Register.GetByEmail: %w", err)
	}
	if existUser != nil {
		return domain.ErrEmailExists
	}

	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("authUsecase.Register.Hash: %w", err)
	}

	user := domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = u.repo.Create(ctx, &user)
	if err != nil {
		return fmt.Errorf("authUsecase.Register.Create: %w", err)
	}

	return nil
}

func (u *authUsecase) authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("authUsecase.authenticate.GetByEmail: %w", err)
	}

	if !u.hasher.Compare(password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := u.authenticate(ctx, email, password)
	if err != nil {
		return nil, err
	}
	token, err := u.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("authUsecase.Login.GenerateToken: %w", err)
	}

	return &LoginResult{
		User:  *user,
		Token: token,
	}, nil
}
