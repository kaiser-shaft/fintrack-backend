package pgpool

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	Name     string `envconfig:"DB_NAME" default:"fintrack"`
}

type Pool struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, c Config) (*Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()
}
