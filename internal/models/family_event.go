package models

import (
	"time"

	"github.com/google/uuid"
)

type FamilyEvent struct {
	Id             uuid.UUID `db:"id" json:"id"`
	FamilyID       uuid.UUID `db:"family_id" json:"family_id"`
	Title          string    `db:"title" json:"title"`
	Description    string    `db:"description" json:"description"`
	StartTime      time.Time `db:"start_time" json:"start_time"`
	EndTime        time.Time `db:"end_time" json:"end_time"`
	OrganizerId    uuid.UUID `db:"organizer_id" json:"organizer_id"`
	Color          string    `db:"color" json:"color"`
	IsRecurring    bool      `db:"is_recurring" json:"is_recurring"`
	RecurrenceRule string    `db:"recurrence_rule" json:"recurrence_rule"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt      time.Time `db:"deleted_at" json:"deleted_at"`
}
