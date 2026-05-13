package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/v1"
)

func NewRouter(authHandler *v1.AuthHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})
	})

	return r
}
