package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fintrack-backend/pkg/jwt"
	"github.com/kaiser-shaft/fintrack-backend/pkg/render"
)

type contextKey string

const userIDKey contextKey = "user_id"

type AuthMiddleware struct {
	jwtManager *jwt.Manager
}

func NewAuthMiddleware(jwtManager *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			render.Error(w, "unauthorized", http.StatusUnauthorized, nil)
			return
		}

		userID, err := m.jwtManager.ValidateToken(cookie.Value)
		if err != nil {
			render.Error(w, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	if userID, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return userID, true
	}
	return uuid.Nil, false
}
