package core_jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	secretKey []byte
	ttl       time.Duration
}

func NewTokenManager(secretKey string, ttl time.Duration) (*TokenManager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("empty secret key")
	}
	return &TokenManager{
		secretKey: []byte(secretKey),
		ttl:       ttl,
	}, nil
}

func (m *TokenManager) GenerateAccessToken(userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.ttl)

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (m *TokenManager) ValidateToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
