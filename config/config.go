package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaiser-shaft/fintrack-backend/pkg/httpserver"
	"github.com/kaiser-shaft/fintrack-backend/pkg/jwt"
	"github.com/kaiser-shaft/fintrack-backend/pkg/logger"
	"github.com/kaiser-shaft/fintrack-backend/pkg/pgpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTP     httpserver.Config
	Postgres pgpool.Config
	Log      logger.Config
	JWT      jwt.Config
}

func New() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if _, err := os.Stat(".env"); err == nil {
			configPath = ".env"
		}
	}

	var cfg Config
	if _, err := os.Stat(configPath); err == nil {
		if err := godotenv.Load(configPath); err != nil {
			return nil, fmt.Errorf("error loading %s file: %w", configPath, err)
		}
	}

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := New()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}
