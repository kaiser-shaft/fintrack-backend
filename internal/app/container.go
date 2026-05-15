package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/kaiser-shaft/fintrack-backend/config"
	"github.com/kaiser-shaft/fintrack-backend/internal/delivery/http"
	v1 "github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/v1"
	"github.com/kaiser-shaft/fintrack-backend/internal/domain"
	"github.com/kaiser-shaft/fintrack-backend/internal/repository/postgres"
	"github.com/kaiser-shaft/fintrack-backend/internal/usecase"
	"github.com/kaiser-shaft/fintrack-backend/pkg/hasher"
	"github.com/kaiser-shaft/fintrack-backend/pkg/httpserver"
	"github.com/kaiser-shaft/fintrack-backend/pkg/jwt"
	"github.com/kaiser-shaft/fintrack-backend/pkg/logger"
	"github.com/kaiser-shaft/fintrack-backend/pkg/pgpool"
	"github.com/kaiser-shaft/fintrack-backend/pkg/validator"
)

type Container struct {
	ctx context.Context
	cfg *config.Config
	mu  sync.Mutex

	pgPool       *pgpool.Pool
	logger       *slog.Logger
	passHasher   *hasher.PasswordHasher
	reqValidator *validator.Validator
	jwtManager   *jwt.Manager

	userRepository     domain.UserRepository
	accountRepository  domain.AccountRepository
	categoryRepository domain.CategoryRepository

	authUsecase     usecase.AuthUsecase
	accountUsecase  usecase.AccountUsecase
	categoryUsecase usecase.CategoryUsecase

	httpServer *httpserver.Server
}

func NewContainer(ctx context.Context, cfg *config.Config) *Container {
	return &Container{
		ctx: ctx,
		cfg: cfg,
	}
}

func (c *Container) PgPool() (*pgpool.Pool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getPgPool()
}

func (c *Container) Logger() *slog.Logger {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getLogger()
}

func (c *Container) PassHasher() *hasher.PasswordHasher {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getPassHasher()
}

func (c *Container) ReqValidator() *validator.Validator {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getReqValidator()
}

func (c *Container) JWTManager() *jwt.Manager {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getJWTManager()
}

func (c *Container) UserRepository() (domain.UserRepository, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getUserRepository()
}

func (c *Container) AccountRepository() (domain.AccountRepository, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getAccountRepository()
}

func (c *Container) CategoryRepository() (domain.CategoryRepository, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getCategoryRepository()
}

func (c *Container) AuthUsecase() (usecase.AuthUsecase, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getAuthUsecase()
}

func (c *Container) HTTPServer() (*httpserver.Server, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpServer == nil {
		reqValidator := c.getReqValidator()
		logger := c.getLogger()

		authUC, err := c.getAuthUsecase()
		if err != nil {
			return nil, fmt.Errorf("container.HTTPServer: %w", err)
		}
		authH := v1.NewAuthHandler(
			authUC,
			reqValidator,
			logger,
			c.cfg.JWT.CookieSecure,
			c.cfg.JWT.TokenDuration,
		)

		accountUC, err := c.getAccountUsecase()
		if err != nil {
			return nil, fmt.Errorf("container.HTTPServer: %w", err)
		}
		accountH := v1.NewAccountHandler(
			accountUC,
			reqValidator,
			logger,
		)

		categoryUC, err := c.getCategoryUsecase()
		if err != nil {
			return nil, fmt.Errorf("container.HTTPServer: %w", err)
		}
		categoryH := v1.NewCategoryHandler(
			categoryUC,
			reqValidator,
			logger,
		)

		jwtManager := c.getJWTManager()
		router := http.NewRouter(
			authH,
			accountH,
			categoryH,
			jwtManager,
		)

		c.httpServer = httpserver.New(router, c.cfg.HTTP)
	}

	return c.httpServer, nil
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpServer != nil {
		c.httpServer.Close()
	}

	if c.pgPool != nil {
		c.pgPool.Close()
	}
}

// Вспомогательные методы чтобы не вызывать дважды блокировку

func (c *Container) getPgPool() (*pgpool.Pool, error) {
	if c.pgPool == nil {
		pgPool, err := pgpool.New(c.ctx, c.cfg.Postgres)
		if err != nil {
			return nil, fmt.Errorf("container.getPgPool: %w", err)
		}
		c.pgPool = pgPool
	}

	return c.pgPool, nil
}

func (c *Container) getLogger() *slog.Logger {
	if c.logger == nil {
		c.logger = logger.Init(c.cfg.Log)
	}

	return c.logger
}

func (c *Container) getPassHasher() *hasher.PasswordHasher {
	if c.passHasher == nil {
		c.passHasher = hasher.New()
	}

	return c.passHasher
}

func (c *Container) getReqValidator() *validator.Validator {
	if c.reqValidator == nil {
		c.reqValidator = validator.New()
	}

	return c.reqValidator
}

func (c *Container) getJWTManager() *jwt.Manager {
	if c.jwtManager == nil {
		c.jwtManager = jwt.New(c.cfg.JWT)
	}

	return c.jwtManager
}

func (c *Container) getUserRepository() (domain.UserRepository, error) {
	if c.userRepository == nil {
		pool, err := c.getPgPool()
		if err != nil {
			return nil, fmt.Errorf("container.getUserRepository: %w", err)
		}
		c.userRepository = postgres.NewUserRepository(pool.Pool)
	}

	return c.userRepository, nil
}

func (c *Container) getAccountRepository() (domain.AccountRepository, error) {
	if c.accountRepository == nil {
		pool, err := c.getPgPool()
		if err != nil {
			return nil, fmt.Errorf("container.getAccountRepository: %w", err)
		}
		c.accountRepository = postgres.NewAccountRepository(pool.Pool)
	}

	return c.accountRepository, nil
}

func (c *Container) getCategoryRepository() (domain.CategoryRepository, error) {
	if c.categoryRepository == nil {
		pool, err := c.getPgPool()
		if err != nil {
			return nil, fmt.Errorf("container.getCategoryRepository: %w", err)
		}
		c.categoryRepository = postgres.NewCategoryRepository(pool.Pool)
	}

	return c.categoryRepository, nil
}

func (c *Container) getAuthUsecase() (usecase.AuthUsecase, error) {
	if c.authUsecase == nil {
		userRepo, err := c.getUserRepository()
		if err != nil {
			return nil, fmt.Errorf("container.getAuthUsecase: %w", err)
		}
		passHasher := c.getPassHasher()
		jwtManager := c.getJWTManager()

		c.authUsecase = usecase.NewAuthUsecase(
			userRepo,
			passHasher,
			jwtManager,
		)
	}

	return c.authUsecase, nil
}

func (c *Container) getAccountUsecase() (usecase.AccountUsecase, error) {
	if c.accountUsecase == nil {
		accountRepo, err := c.getAccountRepository()
		if err != nil {
			return nil, fmt.Errorf("container.getAccountUsecase: %w", err)
		}
		c.accountUsecase = usecase.NewAccountUsecase(accountRepo)
	}

	return c.accountUsecase, nil
}

func (c *Container) getCategoryUsecase() (usecase.CategoryUsecase, error) {
	if c.categoryUsecase == nil {
		categoryRepo, err := c.getCategoryRepository()
		if err != nil {
			return nil, fmt.Errorf("container.getCategoryUsecase: %w", err)
		}
		c.categoryUsecase = usecase.NewCategoryUsecase(categoryRepo)
	}

	return c.categoryUsecase, nil
}
