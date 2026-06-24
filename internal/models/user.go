package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	FullName     string    `db:"full_name" json:"full_name"`
	AvatarUrl    string    `db:"avatar_url" json:"avatar_url"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
