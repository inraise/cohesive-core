package models

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingItem struct {
	Id          uuid.UUID `db:"id" json:"id"`
	ListId      uuid.UUID `db:"list_id" json:"list_id"`
	Text        string    `db:"text" json:"text"`
	Category    string    `db:"category" json:"category"`
	Completed   bool      `db:"completed" json:"completed"`
	CompletedBy uuid.UUID `db:"completed_by" json:"completed_by"`
	CompletedAt time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at" json:"deleted_at"`
}
