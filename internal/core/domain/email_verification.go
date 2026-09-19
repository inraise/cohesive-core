package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Email     string
	TokenHash string

	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func NewEmailVerification(
	userID uuid.UUID,
	email string,
	tokenHash string,
	expiresAt time.Time,
) EmailVerification {
	return EmailVerification{
		ID:        UninitializedID,
		UserID:    userID,
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (e *EmailVerification) IsValid() bool {
	return e.UsedAt == nil && time.Now().Before(e.ExpiresAt)
}
