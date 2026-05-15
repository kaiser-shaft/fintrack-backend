package http

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/middleware"
	v1 "github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/v1"
	"github.com/kaiser-shaft/fintrack-backend/pkg/jwt"
)

func NewRouter(
	authHandler *v1.AuthHandler,
	accountHandler *v1.AccountHandler,
	categoryHandler *v1.CategoryHandler,
	jwtManager *jwt.Manager,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)

	authMW := middleware.NewAuthMiddleware(jwtManager)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMW.Handler)

			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", accountHandler.Create)
				r.Get("/", accountHandler.List)
			})

			r.Route("/categories", func(r chi.Router) {
				r.Post("/", categoryHandler.Create)
				r.Get("/", categoryHandler.List)
			})
		})
	})

	return r
}
