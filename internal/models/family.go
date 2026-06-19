package models

import (
	"time"

	"github.com/google/uuid"
)

type Family struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	InviteCode string    `db:"invite_code" json:"invite_code"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
