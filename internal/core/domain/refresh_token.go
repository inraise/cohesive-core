package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID     uuid.UUID
	UserID uuid.UUID

	TokenHash string

	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func NewRefreshTokenUninitialized(
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
) RefreshToken {
	return RefreshToken{
		ID:        UninitializedID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (t *RefreshToken) IsValid() bool {
	return t.RevokedAt == nil && time.Now().Before(t.ExpiresAt)
}