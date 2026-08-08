package auth_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID      uuid.UUID
	Version int

	Email        string
	PasswordHash string

	FirstName string
	LastName  *string
	Age       *int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshTokenModel struct {
	ID     uuid.UUID
	UserID uuid.UUID

	TokenHash string

	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}