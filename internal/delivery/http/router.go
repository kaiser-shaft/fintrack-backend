package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpMiddleware "github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/middleware"
	v1 "github.com/kaiser-shaft/fintrack-backend/internal/delivery/http/v1"
	"github.com/kaiser-shaft/fintrack-backend/pkg/jwt"
)

func NewRouter(
	authHandler *v1.AuthHandler,
	jwtManager *jwt.Manager,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	authMW := httpMiddleware.NewAuthMiddleware(jwtManager)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMW.Handler)

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				userID, _ := httpMiddleware.GetUserID(r.Context())
				w.Write([]byte("Hello, user " + userID.String()))
			})
		})
	})

	return r
}
