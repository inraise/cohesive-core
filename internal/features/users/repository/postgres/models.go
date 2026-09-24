package users_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID      uuid.UUID
	Version int

	Email        string
	PasswordHash string
	IsVerified   bool

	FirstName string
	LastName  *string
	Age       *int

	CreatedAt time.Time
	UpdatedAt time.Time
}
