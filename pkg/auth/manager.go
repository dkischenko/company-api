package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type Manager struct {
	signinKey []byte
	tokenTTL  time.Duration
}

func NewManager(tokenTTL time.Duration) (*Manager, error) {
	var key string
	if key = os.Getenv("SIGNINKEY"); key == "" {
		return nil, errors.New("empty signin key passed")
	}

	return &Manager{signinKey: []byte(key), tokenTTL: tokenTTL}, nil
}

func (m *Manager) CreateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(m.tokenTTL).Unix()
	claims["iss_at"] = time.Now().Unix()
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.signinKey)
}
