package main

import (
	"context"

	"github.com/kaiser-shaft/fintrack-backend/config"
	"github.com/kaiser-shaft/fintrack-backend/internal/app"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	if err := app.Run(ctx, cfg); err != nil {
		panic(err)
	}
}
