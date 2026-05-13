package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	SecretKey     string        `envconfig:"JWT_SECRET_KEY" required:"true"`
	TokenDuration time.Duration `envconfig:"JWT_TOKEN_DURATION" default:"1h"`
	CookieSecure  bool          `envconfig:"JWT_COOKIE_SECURE" default:"false"`
}

type Manager struct {
	secretKey     string
	tokenDuration time.Duration
}

func New(c Config) *Manager {
	return &Manager{secretKey: c.SecretKey, tokenDuration: c.TokenDuration}
}

func (m *Manager) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(m.tokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	idStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid subject")
	}

	return uuid.Parse(idStr)
}
