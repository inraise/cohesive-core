package models

import (
	"time"

	"github.com/google/uuid"
)

type StatusTask string
type PriorityTask string

const (
	Pending   StatusTask = "pending"
	Completed StatusTask = "completed"
	Overdue   StatusTask = "overdue"
)

const (
	Low    PriorityTask = "low"
	Medium PriorityTask = "medium"
	High   PriorityTask = "high"
)

type Task struct {
	Id          uuid.UUID    `db:"id" json:"id"`
	FamilyID    uuid.UUID    `db:"family_id" json:"family_id"`
	Title       string       `db:"title" json:"title"`
	Description string       `db:"description" json:"description"`
	AssigneeId  uuid.UUID    `db:"assignee_id" json:"assignee_id"`
	AssignerId  uuid.UUID    `db:"assigner_id" json:"assigner_id"`
	Priority    PriorityTask `db:"priority" json:"priority"`
	Deadline    time.Time    `db:"deadline" json:"deadline"`
	Status      StatusTask   `db:"status" json:"status"`
	CompletedAt time.Time    `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt   time.Time    `db:"deleted_at" json:"deleted_at"`
}
