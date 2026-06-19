package models

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingList struct {
	Id          uuid.UUID    `db:"id" json:"id"`
	FamilyID    uuid.UUID    `db:"family_id" json:"family_id"`
	Name       string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt   time.Time    `db:"deleted_at" json:"deleted_at"`
}
